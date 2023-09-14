package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/bunovel"
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
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
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{bunovel.ErrNotFound, http.StatusNotFound},
		}, false)
		return
	}

	c.AbortWithStatus(http.StatusAccepted)

	if deferred != nil {
		if err := deferred(); err != nil {
			_ = c.Error(err)
		}
	}
}
