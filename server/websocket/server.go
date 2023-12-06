package websocket

import (
	"context"
	"fmt"
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/domain/hub"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func New(
	env environment.Env,
	logger logger.Logger,
	router *gin.Engine,
	user user.User,
	hub hub.Hub,
) (*Server, error) {
	hdls, err := newHandlers(env, logger, user, hub)
	if err != nil {
		return nil, fmt.Errorf("create websocket handlers failed: %w", err)
	}

	v1 := router.Group("ws/v1")
	{
		v1.GET("shack", hdls.Shack())
	}

	s := &http.Server{
		Addr:    env.WebsocketListenAddr,
		Handler: router,
	}
	return &Server{
		server: s,
	}, nil
}

type Server struct {
	server *http.Server
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.server.Shutdown(context.Background())
}
