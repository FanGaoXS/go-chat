package rest

import (
	"fmt"
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/groupmember"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func New(
	env environment.Env,
	logger logger.Logger,
	user user.User,
	group group.Group,
	groupMember groupmember.GroupMember,
) (*Server, error) {
	handlers, err := NewHandlers(env, logger, user, group, groupMember)
	if err != nil {
		return nil, fmt.Errorf("create rest handles failed: %w", err)
	}

	router := gin.New()
	gin.ForceConsoleColor()
	router.Use(gin.Logger()) // middlewares

	v1 := router.Group("api/v1")
	{
		v1.POST("registerUser", handlers.RegisterUser())
		v1.GET("me", AuthMiddleware(user), handlers.Me())
		v1.GET("myGroups", AuthMiddleware(user), handlers.MyGroups())
	}

	g := v1.Group("group", AuthMiddleware(user))
	{
		g.POST("", handlers.InsertGroup())
		g.GET(":id", handlers.GetGroupByID())
		g.DELETE(":id", handlers.DeleteGroup())
		g.PUT("toPublic/:id", handlers.PublicGroup())
		g.PUT("toPrivate/:id", handlers.PrivateGroup())
		g.PUT("assignUser/:id", handlers.AssignUsersToGroup())
		g.GET("members/:id", handlers.GroupMembers())
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
