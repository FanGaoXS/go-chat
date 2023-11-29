package auth

import "context"

const RequestAdditionInCtx = "context:request_addition"

type RequestAddition struct {
	Subject string

	Agent string
}

func WithRequestCtx(ctx context.Context, add RequestAddition) context.Context {
	return context.WithValue(ctx, RequestAdditionInCtx, add)
}

func FromRequestCtx(ctx context.Context) RequestAddition {
	v := ctx.Value(RequestAdditionInCtx)
	if v == nil {
		return RequestAddition{}
	}
	return v.(RequestAddition)
}
