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

type ResendEmailValidationHandler interface {
	Handle(c *gin.Context)
}

func NewResendEmailValidationHandler(service services.ResendEmailValidationService) ResendEmailValidationHandler {
	return &resendEmailValidationHandlerImpl{
		service: service,
	}
}

type resendEmailValidationHandlerImpl struct {
	service services.ResendEmailValidationService
}

func (h *resendEmailValidationHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	deferred, err := h.service.ResendEmailValidation(c, token, time.Now())
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
