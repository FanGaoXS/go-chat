package auth

import (
	"context"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/errors"
)

type Authorizer interface {
	Verify(ctx context.Context) (context.Context, error)
}

func NewAuthorizer(env environment.Env, user user.User) (Authorizer, error) {
	if env.BypassAuth {
		return &noAuth{}, nil
	}

	return &auth{
		user: user,
	}, nil
}

type auth struct {
	user user.User
}

func (a *auth) Verify(ctx context.Context) (context.Context, error) {
	r := FromRequestCtx(ctx)

	// TODO: parse token to subject, email, phone ...

	subject := r.Token

	u, err := a.user.GetUserBySubject(ctx, subject)
	if err != nil {
		return nil, errors.New(errors.Unauthenticated, err, "")
	}

	ui := UserInfo{
		Subject:  subject,
		Nickname: u.Nickname,
		Phone:    u.Phone,
		Agent:    r.Agent,
	}
	return WithContext(ctx, ui), nil
}

type noAuth struct{}

func (n *noAuth) Verify(ctx context.Context) (context.Context, error) { return ctx, nil }
