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
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{goframework.ErrInvalidEntity, http.StatusUnprocessableEntity},
		}, false)
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}
