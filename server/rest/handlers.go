package rest

import (
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func NewHandlers(env environment.Env, logger logger.Logger, user user.User) (Handlers, error) {
	return Handlers{
		logger: logger,
		user:   user,
	}, nil
}

type Handlers struct {
	logger logger.Logger

	user user.User
}

func (h *Handlers) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		nickname := c.PostForm("nickname")
		username := c.PostForm("username")
		password := c.PostForm("password")
		phone := c.PostForm("phone")

		ctx := c.Request.Context()
		input := user.RegisterInput{
			Nickname: nickname,
			Username: username,
			Password: password,
			Phone:    phone,
		}
		subject, err := h.user.RegisterUser(ctx, input)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"subject": subject,
		})
	}
}

func (h *Handlers) GetUserBySubject() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		subject := c.Param("subject")

		ctx := c.Request.Context()
		u, err := h.user.GetUserBySubject(ctx, subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, u)
	}
}

func (h *Handlers) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// DELETE
		subject := c.Param("subject")

		ctx := c.Request.Context()
		err := h.user.DeleteUser(ctx, subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
