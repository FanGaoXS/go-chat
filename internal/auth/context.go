package auth

import "context"

const KeyInCtx = "context:userinfo"

type UserInfo struct {
	Subject  string
	Nickname string
	Phone    string
}

func WithContext(ctx context.Context, ui UserInfo) context.Context {
	return context.WithValue(ctx, KeyInCtx, ui)
}

func FromContext(ctx context.Context) UserInfo {
	v := ctx.Value(KeyInCtx)
	return v.(UserInfo)
}
