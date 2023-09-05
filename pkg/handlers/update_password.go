package handlers

import (
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
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
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
			{errors.ErrInvalidEntity, http.StatusUnprocessableEntity},
			{errors.ErrNotFound, http.StatusForbidden},
		})
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}
