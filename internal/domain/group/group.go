package group

import (
	"context"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
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
	MakeGroupPrivate(ctx context.Context, id int64) error
	MakeGroupPublic(ctx context.Context, id int64) error
	ListGroupsOfUser(ctx context.Context, userSubject string) ([]*entity.Group, error)

	AssignMembersToGroup(ctx context.Context, groupID int64, userSubject ...string) error
	RemoveMembersFromGroup(ctx context.Context, groupID int64, userSubject ...string) error
	IsMemberOfGroup(ctx context.Context, groupID int64, memberSubject string) (bool, error)
	ListMembersOfGroup(ctx context.Context, groupID int64) ([]*entity.User, error)

	AssignAdminsToGroup(ctx context.Context, groupID int64, adminSubject ...string) error
	RemoveAdminsFromGroup(ctx context.Context, groupID int64, adminSubject ...string) error
	IsAdminOfGroup(ctx context.Context, groupID int64, memberSubject string) (bool, error)
	ListAdminsOfGroup(ctx context.Context, groupID int64) ([]*entity.User, error)
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

	e := &entity.GroupMember{
		UserSubject: input.CreatedBy,
		GroupID:     gID,
		IsAdmin:     true,
	}
	if err = g.storage.InsertGroupMember(ses, e); err != nil {
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

	res, err := g.storage.ListGroupsByCreatedBy(ses, createdBy)
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
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	err = g.storage.DeleteGroupMembersByGroupID(ses, id)
	if err != nil {
		return err
	}

	err = g.storage.DeleteGroup(ses, id)
	if err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (g *group) MakeGroupPrivate(ctx context.Context, id int64) error {
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

func (g *group) MakeGroupPublic(ctx context.Context, id int64) error {
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

func (g *group) ListGroupsOfUser(ctx context.Context, userSubject string) ([]*entity.Group, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	gms, err := g.storage.ListGroupMembersByUserSubject(ses, userSubject)
	if err != nil {
		return nil, err
	}

	groups := make([]*entity.Group, 0, len(gms))
	for _, gm := range gms {
		grp, err := g.storage.GetGroupByID(ses, gm.GroupID)
		if err != nil {
			return nil, err
		}

		groups = append(groups, grp)
	}

	if len(groups) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "empty groups of user: %s", userSubject)
	}

	return groups, nil
}

func (g *group) AssignMembersToGroup(ctx context.Context, groupID int64, userSubject ...string) error {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	for _, subject := range userSubject {
		i := &entity.GroupMember{
			UserSubject: subject,
			GroupID:     groupID,
			IsAdmin:     false,
		}
		err = g.storage.InsertGroupMember(ses, i)
		if err != nil {
			return err
		}
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (g *group) RemoveMembersFromGroup(ctx context.Context, groupID int64, userSubject ...string) error {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	for _, subject := range userSubject {
		err = g.storage.DeleteGroupMember(ses, subject, groupID)
		if err != nil {
			return err
		}
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (g *group) IsMemberOfGroup(ctx context.Context, groupID int64, memberSubject string) (bool, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return false, err
	}

	ok, err := g.storage.IsMemberOfGroup(ses, memberSubject, groupID)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (g *group) ListMembersOfGroup(ctx context.Context, groupID int64) ([]*entity.User, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	gms, err := g.storage.ListGroupMembersByGroupID(ses, groupID)
	if err != nil {
		return nil, err
	}

	members := make([]*entity.User, 0, len(gms))
	for _, gm := range gms {
		user, err := g.storage.GetUserBySubject(ses, gm.UserSubject)
		if err != nil {
			return nil, err
		}

		members = append(members, user)
	}

	if len(members) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "empty members of group: %d", groupID)
	}

	return members, nil
}

func (g *group) AssignAdminsToGroup(ctx context.Context, groupID int64, adminSubject ...string) error {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	for _, subject := range adminSubject {
		ok, err := g.storage.IsMemberOfGroup(ses, subject, groupID)
		if err != nil {
			return err
		}
		if !ok {
			return errors.Newf(errors.NotFound, nil, "[%s]不在[%d]群里", subject, groupID)
		}

		err = g.storage.UpdateGroupMemberIsAdmin(ses, subject, groupID, true)
		if err != nil {
			return err
		}
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (g *group) RemoveAdminsFromGroup(ctx context.Context, groupID int64, adminSubject ...string) error {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	for _, subject := range adminSubject {
		ok, err := g.storage.IsMemberOfGroup(ses, subject, groupID)
		if err != nil {
			return err
		}
		if !ok {
			return errors.Newf(errors.NotFound, nil, "[%s]不在[%d]群里", subject, groupID)
		}

		err = g.storage.UpdateGroupMemberIsAdmin(ses, subject, groupID, false)
		if err != nil {
			return err
		}
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (g *group) IsAdminOfGroup(ctx context.Context, groupID int64, memberSubject string) (bool, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return false, err
	}

	ok, err := g.storage.IsAdminOfGroup(ses, memberSubject, groupID)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (g *group) ListAdminsOfGroup(ctx context.Context, groupID int64) ([]*entity.User, error) {
	ses, err := g.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	gms, err := g.storage.ListGroupAdminsByGroupID(ses, groupID)
	if err != nil {
		return nil, err
	}

	members := make([]*entity.User, 0, len(gms))
	for _, gm := range gms {
		user, err := g.storage.GetUserBySubject(ses, gm.UserSubject)
		if err != nil {
			return nil, err
		}

		members = append(members, user)
	}

	if len(members) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "empty admins of group: %d", groupID)
	}

	return members, nil
}
