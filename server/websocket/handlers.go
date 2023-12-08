package websocket

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/hub"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func newHandlers(env environment.Env, logger logger.Logger, user user.User, group group.Group, hub hub.Hub) (handlers, error) {
	return handlers{
		logger: logger,
		user:   user,
		group:  group,
		hub:    hub,
	}, nil
}

type handlers struct {
	logger logger.Logger

	user  user.User
	group group.Group
	hub   hub.Hub
}

func (h *handlers) Shack(authorizer auth.Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		token := strings.TrimSpace(c.Query("token"))
		if token == "" {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty subject"))
			return
		}

		ctx := c.Request.Context()
		r := auth.RequestAddition{
			Token: token,
			Agent: "",
		}
		ctx = auth.WithRequestCtx(ctx, r)
		ctx, err := authorizer.Verify(ctx)
		if err != nil {
			WrapGinError(c, errors.Newf(errors.Unauthenticated, err, ""))
			return
		}
		c.Request = c.Request.Clone(ctx)
		ui := auth.FromContext(ctx)
		subject := ui.Subject

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
		h.logger.Infof("[%s] login", u.Nickname)
		conn.WriteMessage(websocket.TextMessage, []byte("Welcome! "+u.Nickname))
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// 心跳检测
			if messageType == websocket.PingMessage || string(message) == "PING" || string(message) == "ping" {
				conn.WriteMessage(websocket.TextMessage, []byte("PONG"))
				continue
			}

			var m map[string]string
			if err = json.Unmarshal(message, &m); err != nil {
				conn.WriteJSON(KV{"error": err.Error()})
				continue
			}

			switch m["type"] {
			case "broadcast":
				content := m["content"]
				if err = h.hub.SendBroadcastMessage(ctx, subject, content); err != nil {
					h.logger.Errorf("%s send broadcast message failed: %w", subject, err)
					conn.WriteJSON(KV{"error": err.Error()})
					break
				}
			case "group":
				groupID, err := strconv.ParseInt(m["group_id"], 10, 64)
				if err != nil {
					conn.WriteJSON(KV{"error": err.Error()})
					break
				}
				content := m["content"]
				ok, err := h.group.IsMemberOfGroup(ctx, groupID, subject)
				if err != nil {
					conn.WriteJSON(KV{"error": err.Error()})
					break
				}
				if !ok {
					conn.WriteJSON(KV{"error": "你不是该群成员"})
					break
				}

				if err = h.hub.SendGroupMessage(ctx, subject, content, groupID); err != nil {
					h.logger.Errorf("%s send group message to %d failed: %w", subject, groupID, err)
					conn.WriteJSON(KV{"error": err.Error()})
					break
				}
			case "private":
				content := m["content"]
				receiver := m["receiver"]
				ok, err := h.user.IsFriendOfUser(ctx, subject, receiver)
				if err != nil {
					conn.WriteJSON(KV{"error": err.Error()})
					break
				}
				if !ok {
					conn.WriteJSON(KV{"error": "对方不是你的好友"})
					break
				}

				if err = h.hub.SendPrivateMessage(ctx, subject, content, receiver); err != nil {
					h.logger.Errorf("%s send private message to %d failed: %w", subject, receiver, err)
					conn.WriteJSON(KV{"error": err.Error()})
					break
				}
			default:
				conn.WriteJSON(KV{"error": "invalid message type"})
			}
		}
		h.hub.UnregisterClient(ctx, subject)
		h.logger.Infof("[%s] logout", u.Nickname)
	}
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

type KV map[string]any
