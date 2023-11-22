package rest

import (
	"net/http"

	"fangaoxs.com/go-chat/environment"

	"github.com/gin-gonic/gin"
)

func New(env environment.Env) (*Server, error) {

	router := gin.New()
	router.Use(gin.Logger()) // middlewares

	v1 := router.Group("ap1/v1")
	user := v1.Group("user")
	{
		user.POST("/user")
		user.GET("/user")
		user.DELETE("/user")
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
