package rest

import (
	"net/http"
	"strings"

	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authorizer auth.Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, agent, err := parse(c.Request)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		ctx := c.Request.Context()
		r := auth.RequestAddition{
			Token: token,
			Agent: agent,
		}
		ctx = auth.WithRequestCtx(ctx, r)
		ctx, err = authorizer.Verify(ctx)
		if err != nil {
			WrapGinError(c, errors.Newf(errors.Unauthenticated, err, ""))
			return
		}

		c.Request = c.Request.Clone(ctx)
		c.Next()
	}
}

func parse(r *http.Request) (token, agent string, err error) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "authorization" {
			token = cookie.Value
			break
		}
	}
	if token == "" {
		token = r.Header.Get("authorization")
	}
	if token == "" {
		return "", "", errors.New(errors.Unauthenticated, nil, "authorization not found")
	}

	splits := strings.SplitN(token, " ", 2)
	if len(splits) < 2 {
		return "", "", errors.New(errors.Unauthenticated, nil, "bad authorization")
	}
	if splits[0] != "Bearer" {
		return "", "", errors.New(errors.Unauthenticated, nil, "unsupported authorization type")
	}

	token = splits[1]
	agent = r.Header.Get("user-agent")

	return token, agent, nil
}
