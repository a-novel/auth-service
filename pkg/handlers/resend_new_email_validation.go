package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ResendNewEmailValidationHandler interface {
	Handle(c *gin.Context)
}

func NewResendNewEmailValidationHandler(service services.ResendNewEmailValidationService) ResendNewEmailValidationHandler {
	return &resendNewEmailValidationHandlerImpl{
		service: service,
	}
}

type resendNewEmailValidationHandlerImpl struct {
	service services.ResendNewEmailValidationService
}

func (h *resendNewEmailValidationHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	deferred, err := h.service.ResendNewEmailValidation(c, token, time.Now())
	if err != nil {
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
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
