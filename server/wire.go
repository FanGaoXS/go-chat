//go:build wireinject
// +build wireinject

package server

import (
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/internal/storage/postgres"

	"github.com/google/wire"
)

func initServer(env environment.Env, logger logger.Logger) (*Server, error) {
	panic(wire.Build(
		postgres.New,
		user.New,
		group.New,
		newServer,
	))
}
