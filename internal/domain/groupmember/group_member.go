package groupmember

import (
	"context"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/storage"
)

type GroupMember interface {
	AssignUserToGroup(ctx context.Context, groupID int64, userSubject ...string) error
	ListGroupsOfUser(ctx context.Context, userSubject string) ([]*entity.Group, error)
	ListUsersOfGroup(ctx context.Context, groupID int64) ([]*entity.User, error)
}

func New(env environment.Env, storage storage.Storage) (GroupMember, error) {
	return &groupMember{
		storage: storage,
	}, nil
}

type groupMember struct {
	storage storage.Storage
}

func (g *groupMember) AssignUserToGroup(ctx context.Context, groupID int64, userSubject ...string) error {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	for _, subject := range userSubject {
		err = g.storage.InsertGroupMember(ses, subject, groupID)
		if err != nil {
			return err
		}
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (g *groupMember) ListGroupsOfUser(ctx context.Context, userSubject string) ([]*entity.Group, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	gms, err := g.storage.ListGroupMemberByUserSubject(ses, userSubject)
	if err != nil {
		return nil, err
	}

	groups := make([]*entity.Group, 0, len(gms))
	for _, gm := range gms {
		group, err := g.storage.GetGroupByID(ses, gm.GroupID)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (g *groupMember) ListUsersOfGroup(ctx context.Context, groupID int64) ([]*entity.User, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	gms, err := g.storage.ListGroupMemberByGroupID(ses, groupID)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, 0, len(gms))
	for _, gm := range gms {
		user, err := g.storage.GetUserBySubject(ses, gm.UserSubject)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
