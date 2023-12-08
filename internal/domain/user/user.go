package user

import (
	"context"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/storage"

	"github.com/google/uuid"
)

type RegisterInput struct {
	Nickname string
	Username string
	Password string
	Phone    string
}

type User interface {
	RegisterUser(ctx context.Context, input RegisterInput) (string, error)
	GetUserBySubject(ctx context.Context, subject string) (*entity.User, error)
	DeleteUser(ctx context.Context, subject string) error
	AllUsers(ctx context.Context) ([]*entity.User, error)

	IsFriendOfUser(ctx context.Context, userSubject, friendSubject string) (bool, error)
	AssignFriendsToUser(ctx context.Context, userSubject string, friendSubject ...string) error
	RemoveFriendsFromUser(ctx context.Context, userSubject string, friendSubject ...string) error
	ListFriendsOfUser(ctx context.Context, userSubject string) ([]*entity.User, error)
}

func New(env environment.Env, storage storage.Storage) (User, error) {
	return &user{
		env:     env,
		storage: storage,
	}, nil
}

type user struct {
	env     environment.Env
	storage storage.Storage
}

func (u *user) RegisterUser(ctx context.Context, input RegisterInput) (string, error) {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return "", err
	}

	i := &entity.User{
		Subject:  uuid.NewString(),
		Nickname: input.Nickname,
		Username: input.Username,
		Password: input.Password,
		Phone:    input.Phone,
	}
	err = u.storage.InsertUser(ses, i)
	if err != nil {
		return "", err
	}

	return i.Subject, nil
}

func (u *user) GetUserBySubject(ctx context.Context, subject string) (*entity.User, error) {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	i, err := u.storage.GetUserBySubject(ses, subject)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (u *user) DeleteUser(ctx context.Context, subject string) error {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	if err = u.storage.DeleteUser(ses, subject); err != nil {
		return err
	}

	return nil
}

func (u *user) AllUsers(ctx context.Context) ([]*entity.User, error) {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	users, err := u.storage.ListAllUsers(ses)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *user) IsFriendOfUser(ctx context.Context, userSubject, friendSubject string) (bool, error) {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return false, err
	}

	ok, err := u.storage.IsFriendOfUser(ses, userSubject, friendSubject)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (u *user) AssignFriendsToUser(ctx context.Context, userSubject string, friendSubject ...string) error {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	for _, fs := range friendSubject {
		if userSubject == fs {
			return errors.New(errors.InvalidArgument, nil, "不可以添加自己为好友")
		}

		uf := &entity.Friendship{
			UserSubject:   userSubject,
			FriendSubject: fs,
		}
		if err = u.storage.InsertFriendship(ses, uf); err != nil {
			return err
		}

		// TODO: 双向好友，好友申请

		uf = &entity.Friendship{
			UserSubject:   fs,
			FriendSubject: userSubject,
		}
		if err = u.storage.InsertFriendship(ses, uf); err != nil {
			return err
		}
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (u *user) RemoveFriendsFromUser(ctx context.Context, userSubject string, friendSubject ...string) error {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	for _, fs := range friendSubject {
		if err = u.storage.DeleteFriendship(ses, userSubject, fs); err != nil {
			return err
		}
		if err = u.storage.DeleteFriendship(ses, fs, userSubject); err != nil {
			return err
		}
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (u *user) ListFriendsOfUser(ctx context.Context, userSubject string) ([]*entity.User, error) {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	ufs, err := u.storage.ListFriendshipsByUserSubject(ses, userSubject)
	if err != nil {
		return nil, err
	}
	if len(ufs) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no friends with user: %s found", userSubject)
	}

	friends := make([]*entity.User, 0, len(ufs))
	for _, uf := range ufs {
		friend, err := u.storage.GetUserBySubject(ses, uf.FriendSubject)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}
