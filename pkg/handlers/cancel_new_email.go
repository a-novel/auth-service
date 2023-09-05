package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CancelNewEmailHandler interface {
	Handle(c *gin.Context)
}

func NewCancelNewEmailHandler(service services.CancelNewEmailService) CancelNewEmailHandler {
	return &cancelNewEmailHandlerImpl{
		service: service,
	}
}

type cancelNewEmailHandlerImpl struct {
	service services.CancelNewEmailService
}

func (h *cancelNewEmailHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if err := h.service.CancelNewEmail(c, token, time.Now()); err != nil {
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
