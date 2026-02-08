package core

import (
	"net/http"

	"github.com/RadishXZ/fastgo/internal/pkg/errorsx"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Reason string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		errx := errorsx.FormError(err)
		c.JSON(errx.Code, ErrorResponse{
			Reason: errx.Reason,
			Message: errx.Message,
		})
		return
	}

	c.JSON(http.StatusOK, data)
}