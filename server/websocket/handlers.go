package websocket

import (
	"fangaoxs.com/go-chat/internal/auth"
	"fmt"
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func newHandlers(env environment.Env, logger logger.Logger) (handlers, error) {
	return handlers{
		logger: logger,
	}, nil
}

type handlers struct {
	logger logger.Logger
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			// don't return errors to maintain backwards compatibility
		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (h *handlers) Shack() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			WrapGinError(c, errors.New(errors.Internal, err, "create websocket connection failed"))
			return
		}
		defer ws.Close()
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)
		_ = ui

		// TODO: online
		for {
			// TODO: receive message
		}
		// TODO: offline
	}
}

func (h *handlers) Broadcast() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *handlers) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		defer ws.Close()

		for {
			// online
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				break
			}

			fmt.Println("messageType:", messageType)
			fmt.Println("p:", string(p))
		}

		// offline
	}
}
