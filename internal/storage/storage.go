package storage

import (
	"context"

	"fangaoxs.com/go-chat/internal/entity"
)

type Storage interface {
	Close() error
	NewSession(ctx context.Context) (Session, error)

	InsertUser(ses Session, i *entity.User) error
	GetUserBySubject(ses Session, subject string) (*entity.User, error)
	GetUserBySecret(ses Session, username, password string) (*entity.User, error)
	DeleteUser(ses Session, subject string) error

	InsertGroup(ses Session, i *entity.Group) (int64, error)
	GetGroupByID(ses Session, id int64) (*entity.Group, error)
	ListGroupByCreatedBy(ses Session, createdBy string) ([]*entity.Group, error)
	DeleteGroup(ses Session, id int64) error
}
