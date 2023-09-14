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

type UpdatePasswordHandler interface {
	Handle(c *gin.Context)
}

func NewUpdatePasswordHandler(service services.UpdatePasswordService) UpdatePasswordHandler {
	return &updatePasswordHandlerImpl{
		service: service,
	}
}

type updatePasswordHandlerImpl struct {
	service services.UpdatePasswordService
}

func (h *updatePasswordHandlerImpl) Handle(c *gin.Context) {
	request := new(models.UpdatePasswordForm)
	if err := c.BindJSON(request); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdatePassword(c, *request, time.Now()); err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{goframework.ErrInvalidEntity, http.StatusUnprocessableEntity},
			{bunovel.ErrNotFound, http.StatusForbidden},
		}, false)
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}
