package websocket

import (
	"encoding/json"
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/domain/hub"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func newHandlers(env environment.Env, logger logger.Logger, user user.User, hub hub.Hub) (handlers, error) {
	return handlers{
		logger: logger,
		user:   user,
		hub:    hub,
	}, nil
}

type handlers struct {
	logger logger.Logger

	user user.User
	hub  hub.Hub
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

			var m map[string]string
			if err = json.Unmarshal(message, &m); err != nil {
				conn.WriteJSON(map[string]any{"error": err.Error()})
			}

			switch m["type"] {
			case "broadcast":
				content := m["content"]
				if err = h.hub.SendBroadcastMessage(ctx, subject, content); err != nil {
					h.logger.Errorf("%s send broadcast message failed: %w", subject, err)
					conn.WriteJSON(map[string]any{"error": err.Error()})
					break
				}
			case "group":
				groupID, err := strconv.ParseInt(m["group_id"], 10, 64)
				if err != nil {
					conn.WriteJSON(map[string]any{"error": err.Error()})
					break
				}
				content := m["content"]
				if err := h.hub.SendGroupMessage(ctx, subject, content, groupID); err != nil {
					h.logger.Errorf("%s send group message to %d failed: %w", subject, groupID, err)
					conn.WriteJSON(map[string]any{"error": err.Error()})
					break
				}
			case "private":
				content := m["content"]
				receiver := m["receiver"]
				if err := h.hub.SendPrivateMessage(ctx, subject, content, receiver); err != nil {
					h.logger.Errorf("%s send private message to %d failed: %w", subject, receiver, err)
					conn.WriteJSON(map[string]any{"error": err.Error()})
					break
				}
			default:
				conn.WriteJSON(map[string]any{"error": "invalid message type"})
			}
		}
		h.hub.UnregisterClient(ctx, subject)
		h.logger.Infof("%s logout", u.Nickname)
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
