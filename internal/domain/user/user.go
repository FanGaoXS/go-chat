package user

import (
	"context"
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
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
