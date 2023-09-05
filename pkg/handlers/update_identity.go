package handlers

import (
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UpdateIdentityHandler interface {
	Handle(c *gin.Context)
}

func NewUpdateIdentityHandler(service services.UpdateIdentityService) UpdateIdentityHandler {
	return &updateIdentityHandlerImpl{
		service: service,
	}
}

type updateIdentityHandlerImpl struct {
	service services.UpdateIdentityService
}

func (h *updateIdentityHandlerImpl) Handle(c *gin.Context) {
	request := new(models.UpdateIdentityForm)
	token := c.GetHeader("Authorization")

	if err := c.BindJSON(request); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateIdentity(c, token, time.Now(), *request); err != nil {
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
			{errors.ErrInvalidEntity, http.StatusUnprocessableEntity},
		})
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}
