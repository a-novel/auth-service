package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ResetPasswordHandler interface {
	Handle(c *gin.Context)
}

func NewResetPasswordHandler(service services.ResetPasswordService) ResetPasswordHandler {
	return &resetPasswordHandlerImpl{
		service: service,
	}
}

type resetPasswordHandlerImpl struct {
	service services.ResetPasswordService
}

func (h *resetPasswordHandlerImpl) Handle(c *gin.Context) {
	email := c.Query("email")

	deferred, err := h.service.ResetPassword(c, email, time.Now())
	if err != nil {
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidEntity, http.StatusBadRequest},
			{errors.ErrNotFound, http.StatusNotFound},
		})
		return
	}

	c.AbortWithStatus(http.StatusAccepted)

	if deferred != nil {
		if err := deferred(); err != nil {
			_ = c.Error(err)
		}
	}
}
