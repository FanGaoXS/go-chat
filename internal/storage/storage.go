package storage

import (
	"context"

	"fangaoxs.com/go-chat/internal/entity"
)

type Storage interface {
	Close() error
	NewSession(ctx context.Context) (Session, error)

	InsertUser(ses Session, i *entity.User) (int64, error)
	GetUserByID(ses Session, id int64) (*entity.User, error)
	GetUserBySecret(ses Session, username, password string) (*entity.User, error)
}
