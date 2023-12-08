package websocket

import (
	"net/http"

	"fangaoxs.com/go-chat/internal/infras/errors"

	"github.com/gin-gonic/gin"
)

func WrapGinError(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	switch errors.Code(err) {
	case errors.NotFound:
		code = http.StatusNotFound
	case errors.InvalidArgument:
		code = http.StatusBadRequest
	case errors.Unimplemented:
		code = http.StatusNotImplemented
	case errors.PermissionDenied:
		code = http.StatusForbidden
	case errors.Unauthenticated:
		code = http.StatusUnauthorized
	case errors.Unavailable:
		code = http.StatusServiceUnavailable
	}
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": err.Error(),
	})
}
