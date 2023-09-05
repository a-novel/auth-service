package handlers

import (
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
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
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
			{errors.ErrInvalidEntity, http.StatusForbidden},
			{errors.ErrUniqConstraintViolation, http.StatusConflict},
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
