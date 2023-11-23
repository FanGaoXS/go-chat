package user

import (
	"context"
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/storage"
)

type RegisterInput struct {
	Nickname string
	Username string
	Password string
	Phone    string
}

type User interface {
	RegisterUser(ctx context.Context, input RegisterInput) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	DeleteUser(ctx context.Context, id int64) error
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

func (u *user) RegisterUser(ctx context.Context, input RegisterInput) (int64, error) {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return 0, err
	}

	user := &entity.User{
		Nickname: input.Nickname,
		Username: input.Username,
		Password: input.Password,
		Phone:    input.Phone,
	}
	id, err := u.storage.InsertUser(ses, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *user) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.storage.GetUserByID(ses, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *user) DeleteUser(ctx context.Context, id int64) error {
	ses, err := u.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	if err = u.storage.DeleteUser(ses, id); err != nil {
		return err
	}

	return nil
}
