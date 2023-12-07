package websocket

import (
	"context"
	"fmt"
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/hub"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func New(
	env environment.Env,
	logger logger.Logger,
	router *gin.Engine,
	authorizer auth.Authorizer,
	user user.User,
	group group.Group,
	hub hub.Hub,
) (*Server, error) {
	hdls, err := newHandlers(env, logger, user, group, hub)
	if err != nil {
		return nil, fmt.Errorf("create websocket handlers failed: %w", err)
	}

	v1 := router.Group("ws/v1")
	{
		v1.GET("shack", hdls.Shack(authorizer))
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
