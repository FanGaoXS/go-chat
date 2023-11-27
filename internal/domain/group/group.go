package group

import (
	"context"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/storage"
)

type CreateGroupInput struct {
	Name      string
	Type      entity.GroupType
	CreatedBy string
}

type Group interface {
	CreateGroup(ctx context.Context, input CreateGroupInput) (int64, error)
	GetGroupByID(ctx context.Context, id int64) (*entity.Group, error)
	ListGroupsByCreatedBy(ctx context.Context, createdBy string) ([]*entity.Group, error)
	DeleteGroup(ctx context.Context, id int64) error
	PrivateGroup(ctx context.Context, id int64) error
	PublicGroup(ctx context.Context, id int64) error
}

func New(env environment.Env, storage storage.Storage) (Group, error) {
	return &group{storage: storage}, nil
}

type group struct {
	storage storage.Storage
}

func (g *group) CreateGroup(ctx context.Context, input CreateGroupInput) (int64, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return 0, err
	}
	ses, err = ses.Begin()
	if err != nil {
		return 0, err
	}

	i := &entity.Group{
		Name:      input.Name,
		Type:      input.Type,
		IsPublic:  false, // 默认不公开
		CreatedBy: input.CreatedBy,
	}
	gID, err := g.storage.InsertGroup(ses, i)
	if err != nil {
		return 0, err
	}

	if err = g.storage.InsertGroupMember(ses, input.CreatedBy, gID); err != nil {
		return 0, err
	}

	if err = ses.Commit(); err != nil {
		return 0, err
	}
	return gID, nil
}

func (g *group) GetGroupByID(ctx context.Context, id int64) (*entity.Group, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := g.storage.GetGroupByID(ses, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g *group) ListGroupsByCreatedBy(ctx context.Context, createdBy string) ([]*entity.Group, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := g.storage.ListGroupByCreatedBy(ses, createdBy)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g *group) DeleteGroup(ctx context.Context, id int64) error {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	err = g.storage.DeleteGroup(ses, id)
	if err != nil {
		return err
	}

	return nil
}

func (g *group) PrivateGroup(ctx context.Context, id int64) error {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	err = g.storage.UpdateGroupIsPublic(ses, id, false)
	if err != nil {
		return err
	}

	return nil
}

func (g *group) PublicGroup(ctx context.Context, id int64) error {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	err = g.storage.UpdateGroupIsPublic(ses, id, true)
	if err != nil {
		return err
	}

	return nil
}
