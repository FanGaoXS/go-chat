package rest

import (
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/errors"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(user user.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject, _ := c.Cookie("subject")
		if subject == "" {
			subject = c.GetHeader("subject")
		}
		if subject == "" {
			WrapGinError(c, errors.New(errors.Unauthenticated, nil, "empty subject in header or cookie"))
			return
		}

		ctx := c.Request.Context()
		u, err := user.GetUserBySubject(ctx, subject)
		if err != nil {
			WrapGinError(c, errors.Newf(errors.Unauthenticated, err, ""))
			return
		}

		ui := auth.UserInfo{
			Subject:  u.Subject,
			Nickname: u.Nickname,
			Phone:    u.Phone,
		}

		ctx = auth.WithContext(ctx, ui)
		c.Request = c.Request.Clone(ctx)
		c.Next()
	}
}
