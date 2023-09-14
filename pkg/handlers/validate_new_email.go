package handlers

import (
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/bunovel"
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ValidateNewEmailHandler interface {
	Handle(c *gin.Context)
}

func NewValidateNewEmailHandler(service services.ValidateNewEmailService) ValidateNewEmailHandler {
	return &validateNewEmailHandlerImpl{
		service: service,
	}
}

type validateNewEmailHandlerImpl struct {
	service services.ValidateNewEmailService
}

func (h *validateNewEmailHandlerImpl) Handle(c *gin.Context) {
	query := new(models.ValidateEmailQuery)
	if err := c.BindQuery(query); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.service.ValidateNewEmail(c, query.ID.Value(), query.Code, time.Now()); err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{goframework.ErrInvalidEntity, http.StatusForbidden},
			{bunovel.ErrUniqConstraintViolation, http.StatusConflict},
		}, false)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
