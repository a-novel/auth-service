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

type UpdateProfileHandler interface {
	Handle(c *gin.Context)
}

func NewUpdateProfileHandler(service services.UpdateProfileService) UpdateProfileHandler {
	return &updateProfileHandlerImpl{
		service: service,
	}
}

type updateProfileHandlerImpl struct {
	service services.UpdateProfileService
}

func (h *updateProfileHandlerImpl) Handle(c *gin.Context) {
	request := new(models.UpdateProfileForm)
	token := c.GetHeader("Authorization")

	if err := c.BindJSON(request); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateProfile(c, token, time.Now(), *request); err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{services.ErrTaken, http.StatusConflict},
			{goframework.ErrInvalidEntity, http.StatusUnprocessableEntity},
		}, false)
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}
