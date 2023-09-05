package handlers

import (
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
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
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
			{errors.ErrInvalidEntity, http.StatusForbidden},
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
