package websocket

import (
	"fmt"
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func New(
	env environment.Env,
	logger logger.Logger,
	router *gin.Engine,
	authorizer auth.Authorizer,
) (*Server, error) {
	hdls, err := newHandlers(env, logger)
	if err != nil {
		return nil, fmt.Errorf("create websocket handlers failed: %w", err)
	}

	v1 := router.Group("ws/v1", AuthMiddleware(authorizer))
	{
		v1.POST("shack", hdls.Shack())
		v1.POST("broadcast", hdls.Broadcast())
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
	return s.server.Close()
}
