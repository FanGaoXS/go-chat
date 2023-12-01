package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/hub"
	"fangaoxs.com/go-chat/internal/domain/record"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func newHandlers(env environment.Env, logger logger.Logger, user user.User, group group.Group, hub hub.Hub, record record.Record) (handlers, error) {
	return handlers{
		logger: logger,
		user:   user,
		group:  group,
		hub:    hub,
		record: record,
	}, nil
}

type handlers struct {
	logger logger.Logger

	user   user.User
	group  group.Group
	hub    hub.Hub
	record record.Record
}

func (h *handlers) Shack() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		subject, ok := c.GetQuery("subject")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty subject"))
			return
		}
		ctx := c.Request.Context()
		u, err := h.user.GetUserBySubject(ctx, subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			WrapGinError(c, errors.New(errors.Internal, err, "create websocket connection failed"))
			return
		}
		defer conn.Close()

		h.hub.RegisterClient(ctx, subject, conn)
		h.logger.Infof("%s login", u.Nickname)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			var m hub.Message
			if err = json.Unmarshal(message, &m); err != nil {
				conn.WriteJSON(map[string]any{"error": err.Error()})
			}

			switch m.Type {
			case hub.MessageTypeBroadcast:
				if err = h.SendBroadcastMessage(ctx, subject, m.Content); err != nil {
					h.logger.Errorf("%s send broadcast message failed: %w", subject, err)
					conn.WriteJSON(map[string]any{"error": err.Error()})
				}
			case hub.MessageTypeGroup:
				groupID, err := strconv.ParseInt(m.Metadata, 10, 64)
				if err != nil {
					conn.WriteJSON(map[string]any{"error": err.Error()})
				}
				if err := h.SendGroupMessage(ctx, subject, m.Content, groupID); err != nil {
					h.logger.Errorf("%s send group message to %d failed: %w", subject, m.Metadata, err)
					conn.WriteJSON(map[string]any{"error": err.Error()})
				}
			case hub.MessageTypePrivate:
				if err := h.SendPrivateMessage(ctx, subject, m.Content, m.To); err != nil {
					h.logger.Errorf("%s send private message to %d failed: %w", subject, m.To, err)
					conn.WriteJSON(map[string]any{"error": err.Error()})
				}
			default:
				conn.WriteJSON(map[string]any{"error": "invalid message type"})
			}
		}
		h.hub.UnregisterClient(ctx, subject)
		h.logger.Infof("%s logout", u.Nickname)
	}
}

func (h *handlers) Broadcast() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		message, ok := c.GetPostForm("message")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty message"))
			return
		}

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		err := h.SendBroadcastMessage(ctx, ui.Subject, message)
		if err != nil {
			WrapGinError(c, err)
			return
		}
	}
}

func (h *handlers) Group() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		message, ok := c.GetPostForm("message")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty message"))
			return
		}
		groupID, err := strconv.ParseInt(c.Param("group_id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		err = h.SendGroupMessage(ctx, ui.Subject, message, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}
	}
}

func (h *handlers) Private() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		message, ok := c.GetPostForm("message")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty message"))
			return
		}
		to, ok := c.GetPostForm("to")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty to"))
			return
		}

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		err := h.SendPrivateMessage(ctx, ui.Subject, message, to)
		if err != nil {
			WrapGinError(c, err)
			return
		}
	}
}

func (h *handlers) SendBroadcastMessage(ctx context.Context, sender, content string) error {
	users, err := h.user.AllUsers(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		m := hub.Message{
			Type:    hub.MessageTypeBroadcast,
			From:    sender,
			To:      u.Subject,
			Content: content,
		}
		if err = h.hub.SendMessage(ctx, m); err != nil {
			// TODO: add to message queue
		}
	}

	return nil
}

func (h *handlers) SendGroupMessage(ctx context.Context, sender, content string, groupID int64) error {
	members, err := h.group.ListMembersOfGroup(ctx, groupID)
	if err != nil {
		return err
	}

	for _, u := range members {
		m := hub.Message{
			Type:      hub.MessageType,
			Metadata:  strconv.FormatInt(groupID, 10),
			From:      sender,
			To:        u.Subject,
			Content:   content,
			CreatedAt: time.Now(),
		}
		if err = h.hub.SendMessage(ctx, m); err != nil {
			// TODO: add to message queue
		}
	}

	return nil
}

func (h *handlers) SendPrivateMessage(ctx context.Context, sender, content, to string) error {
	if _, err := h.user.GetUserBySubject(ctx, to); err != nil {
		return err // check whether the user exists
	}

	m := hub.Message{
		Type:      hub.MessageTypePrivate,
		From:      sender,
		To:        to,
		Content:   content,
		CreatedAt: time.Now(),
	}
	if err := h.hub.SendMessage(ctx, m); err != nil {
		// TODO: add to message queue
	}

	return nil
}

// update http to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		// don't return errors to maintain backwards compatibility
	},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
