//go:build wireinject
// +build wireinject

package server

import (
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/internal/storage/postgres"
	"github.com/gin-gonic/gin"

	"github.com/google/wire"
)

func initServer(env environment.Env, logger logger.Logger, httpServer *gin.Engine) (*Server, error) {
	panic(wire.Build(
		postgres.New,
		user.New,
		group.New,
		auth.NewAuthorizer,
		newServer,
	))
}
