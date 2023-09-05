package handlers

import (
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
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
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
			{services.ErrTaken, http.StatusConflict},
			{errors.ErrInvalidEntity, http.StatusUnprocessableEntity},
		})
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
