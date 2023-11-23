package rest

import (
	"fmt"
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func New(env environment.Env, logger logger.Logger, user user.User) (*Server, error) {
	handlers, err := NewHandlers(env, logger, user)
	if err != nil {
		return nil, fmt.Errorf("create rest handles failed: %w", err)
	}

	router := gin.New()
	gin.ForceConsoleColor()
	router.Use(gin.Logger()) // middlewares

	v1 := router.Group("api/v1")
	ug := v1.Group("user")
	{
		ug.POST("", handlers.RegisterUser())
		ug.GET(":id", handlers.GetUserByID())
		ug.DELETE(":id", handlers.DeleteUser())
	}

	s := &http.Server{
		Addr:    env.RestListenAddr,
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
