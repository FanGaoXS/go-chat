package auth

import "context"

const requestAdditionInCtx = "context:request_addition"

type RequestAddition struct {
	Token string
	Agent string
}

func WithRequestCtx(ctx context.Context, add RequestAddition) context.Context {
	return context.WithValue(ctx, requestAdditionInCtx, add)
}

func FromRequestCtx(ctx context.Context) RequestAddition {
	v := ctx.Value(requestAdditionInCtx)
	if v == nil {
		return RequestAddition{}
	}
	return v.(RequestAddition)
}
