package group

import (
	"context"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/storage"
)

type InsertGroupInput struct {
	Name      string
	Type      entity.GroupType
	CreatedBy string
}

type Group interface {
	InsertGroup(ctx context.Context, input InsertGroupInput) (int64, error)
	GetGroupByID(ctx context.Context, id int64) (*entity.Group, error)
	ListGroupsByCreatedBy(ctx context.Context, createdBy string) ([]*entity.Group, error)
	DeleteGroup(ctx context.Context, id int64) error
}

func New(env environment.Env, storage storage.Storage) (Group, error) {
	return &group{storage: storage}, nil
}

type group struct {
	storage storage.Storage
}

func (g *group) InsertGroup(ctx context.Context, input InsertGroupInput) (int64, error) {}

func (g *group) GetGroupByID(ctx context.Context, id int64) (*entity.Group, error) {}

func (g *group) ListGroupsByCreatedBy(ctx context.Context, createdBy string) ([]*entity.Group, error) {
}

func (g *group) DeleteGroup(ctx context.Context, id int64) error {}
