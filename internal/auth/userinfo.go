package auth

import "context"

const keyInCtx = "context:userinfo"

type UserInfo struct {
	Subject  string
	Nickname string
	Phone    string

	Agent string
}

func WithContext(ctx context.Context, ui UserInfo) context.Context {
	return context.WithValue(ctx, keyInCtx, ui)
}

func FromContext(ctx context.Context) UserInfo {
	v := ctx.Value(keyInCtx)
	if v == nil {
		return UserInfo{}
	}
	return v.(UserInfo)
}
