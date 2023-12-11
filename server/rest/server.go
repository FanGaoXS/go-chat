package rest

import (
	"context"
	"fmt"
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/application"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/hub"
	"fangaoxs.com/go-chat/internal/domain/records"
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
	record records.Records,
	application application.Application,
) (*Server, error) {
	hdls, err := newHandlers(env, logger, user, group, hub, record, application)
	if err != nil {
		return nil, fmt.Errorf("create rest handlers failed: %w", err)
	}

	v1 := router.Group("api/v1")
	v1.POST("registerUser", hdls.RegisterUser())

	p := v1.Group("personal", AuthMiddleware(authorizer))
	{
		p.GET("me", hdls.Me())
		p.POST("sendFriendRequest", hdls.SendFriendRequest())
		p.PUT("agreeFriendRequest/:id", hdls.AgreeFriendRequest())
		p.PUT("refuseFriendRequest/:id", hdls.RefuseFriendRequest())
		p.GET("friendRequestFromMe", hdls.FriendRequestFromMe())
		p.GET("friendRequestToMe", hdls.FriendRequestToMe())
		p.GET("myFriends", hdls.MyFriends())
		p.DELETE("removeFriends", hdls.RemoveFriends())
		p.GET("myGroups", hdls.MyGroups())
	}

	g := v1.Group("group", AuthMiddleware(authorizer))
	{
		g.POST("", hdls.CreateGroup())
		g.GET(":id", hdls.GetGroupByID())
		g.DELETE(":id", hdls.DeleteGroup())
		g.PUT("toPublic/:id", hdls.MakeGroupPublic())
		g.PUT("toPrivate/:id", hdls.MakeGroupPrivate())
		g.PUT("assignMembers/:id", hdls.AssignMembersToGroup())
		g.PUT("removeMembers/:id", hdls.RemoveMembersFromGroup())
		g.PUT("assignAdmins/:id", hdls.AssignAdminsToGroup())
		g.PUT("removeAdmins/:id", hdls.RemoveAdminsFromGroup())
		g.GET("members/:id", hdls.MembersOfGroup())
		g.GET("admins/:id", hdls.AdminsOfGroup())
	}

	r := v1.Group("record", AuthMiddleware(authorizer))
	{
		r.GET("broadcast", hdls.GetRecordBroadcast())
		r.GET("group/:group_id", hdls.GetRecordGroup())
		r.GET("private/:receiver", hdls.GetRecordPrivate())
		r.POST("broadcast", hdls.BroadcastMessage())
		r.POST("group/:group_id", hdls.GroupMessage())
		r.POST("private", hdls.PrivateMessage())
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
	return s.server.Shutdown(context.Background())
}
