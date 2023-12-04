package rest

import (
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
		return nil, fmt.Errorf("create rest handlers failed: %w", err)
	}

	v1 := router.Group("api/v1")
	{
		v1.POST("registerUser", hdls.RegisterUser())
		v1.GET("me", AuthMiddleware(authorizer), hdls.Me())
		v1.GET("myFriends", AuthMiddleware(authorizer), hdls.MyFriends())
		v1.PUT("assignFriends", AuthMiddleware(authorizer), hdls.AssignFriends())
		v1.DELETE("removeFriends", AuthMiddleware(authorizer), hdls.RemoveFriends())
		v1.GET("myGroups", AuthMiddleware(authorizer), hdls.MyGroups())
	}

	g := v1.Group("group", AuthMiddleware(authorizer))
	{
		g.POST("", hdls.InsertGroup())
		g.GET(":id", hdls.GetGroupByID())
		g.DELETE(":id", hdls.DeleteGroup())
		g.PUT("toPublic/:id", hdls.PublicGroup())
		g.PUT("toPrivate/:id", hdls.PrivateGroup())
		g.PUT("assignMembers/:id", hdls.AssignMembersToGroup())
		g.GET("members/:id", hdls.MembersOfGroup())
	}

	r := v1.Group("record", AuthMiddleware(authorizer))
	{
		r.POST("broadcast", hdls.BroadcastMessage())
		r.POST("group/:group_id", hdls.GroupMessage())
		r.POST("private", hdls.PrivateMessage())
		r.GET("broadcast")
		r.GET("group/:group_id")
		r.GET("private")
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
