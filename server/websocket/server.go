package websocket

import (
	"fmt"
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/group"
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
) (*Server, error) {
	hub, err := NewHub(env, logger)
	if err != nil {
		return nil, fmt.Errorf("create websocket hub failed: %w", err)
	}

	hdls, err := newHandlers(env, logger, user, group, hub)
	if err != nil {
		return nil, fmt.Errorf("create websocket handlers failed: %w", err)
	}

	v1 := router.Group("ws/v1")
	{
		v1.GET("shack", hdls.Shack())
		v1.POST("private", AuthMiddleware(authorizer), hdls.Private())
		v1.POST("group/:group_id", AuthMiddleware(authorizer), hdls.Group())
		v1.POST("broadcast", AuthMiddleware(authorizer))
	}

	s := &http.Server{
		Addr:    env.WebsocketListenAddr,
		Handler: router,
	}
	return &Server{
		server: s,
		hub:    hub,
	}, nil
}

type Server struct {
	server *http.Server
	hub    *Hub
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	s.hub.Close()
	s.server.Close()
	return nil
}
