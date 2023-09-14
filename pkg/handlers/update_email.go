package handlers

import (
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UpdateEmailHandler interface {
	Handle(c *gin.Context)
}

func NewUpdateEmailHandler(service services.UpdateEmailService) UpdateEmailHandler {
	return &updateEmailHandlerImpl{
		service: service,
	}
}

type updateEmailHandlerImpl struct {
	service services.UpdateEmailService
}

func (h *updateEmailHandlerImpl) Handle(c *gin.Context) {
	request := new(models.UpdateEmailForm)
	token := c.GetHeader("Authorization")

	if err := c.BindJSON(request); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	deferred, err := h.service.UpdateEmail(c, token, request.NewEmail, time.Now())
	if err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{services.ErrTaken, http.StatusConflict},
			{goframework.ErrInvalidEntity, http.StatusUnprocessableEntity},
		}, false)
		return
	}

	c.AbortWithStatus(http.StatusAccepted)

	if deferred != nil {
		if err := deferred(); err != nil {
			_ = c.Error(err)
			return
		}
	}
}
