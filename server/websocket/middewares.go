package websocket

import (
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/infras/errors"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authorizer auth.Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject, _ := c.Cookie("subject")
		if subject == "" {
			subject = c.GetHeader("subject")
		}
		if subject == "" {
			WrapGinError(c, errors.New(errors.Unauthenticated, nil, "empty subject in header or cookie"))
			return
		}
		agent := c.GetHeader("agent")

		ctx := c.Request.Context()
		r := auth.RequestAddition{
			Subject: subject,
			Agent:   agent,
		}
		ctx = auth.WithRequestCtx(ctx, r) // 将headers写入ctx，以便后续的ctx能够获取到
		ctx, err := authorizer.Verify(ctx)
		if err != nil {
			WrapGinError(c, errors.Newf(errors.Unauthenticated, err, ""))
			return
		}

		c.Request = c.Request.Clone(ctx)
		c.Next()
	}
}
