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

type ValidateEmailHandler interface {
	Handle(c *gin.Context)
}

func NewValidateEmailHandler(service services.ValidateEmailService) ValidateEmailHandler {
	return &validateEmailHandlerImpl{
		service: service,
	}
}

type validateEmailHandlerImpl struct {
	service services.ValidateEmailService
}

func (h *validateEmailHandlerImpl) Handle(c *gin.Context) {
	query := new(models.ValidateEmailQuery)
	if err := c.BindQuery(query); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.service.ValidateEmail(c, query.ID.Value(), query.Code, time.Now()); err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{goframework.ErrInvalidEntity, http.StatusForbidden},
		}, false)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
