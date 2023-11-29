// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/internal/storage/postgres"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func initServer(env environment.Env, logger2 logger.Logger, httpServer *gin.Engine) (*Server, error) {
	storage, err := postgres.New(env)
	if err != nil {
		return nil, err
	}
	userUser, err := user.New(env, storage)
	if err != nil {
		return nil, err
	}
	authorizer, err := auth.NewAuthorizer(env, userUser)
	if err != nil {
		return nil, err
	}
	groupGroup, err := group.New(env, storage)
	if err != nil {
		return nil, err
	}
	server, err := newServer(env, logger2, httpServer, authorizer, userUser, groupGroup)
	if err != nil {
		return nil, err
	}
	return server, nil
}
