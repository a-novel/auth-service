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

type RegisterHandler interface {
	Handle(c *gin.Context)
}

func NewRegisterHandler(service services.RegisterService) RegisterHandler {
	return &registerHandlerImpl{service: service}
}

type registerHandlerImpl struct {
	service services.RegisterService
}

func (h *registerHandlerImpl) Handle(c *gin.Context) {
	form := new(models.RegisterForm)
	if err := c.BindJSON(form); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, deferred, err := h.service.Register(c, *form, time.Now())
	if err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{services.ErrTaken, http.StatusConflict},
			{goframework.ErrInvalidEntity, http.StatusUnprocessableEntity},
		}, false)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token.TokenRaw})

	if deferred != nil {
		if err := deferred(); err != nil {
			_ = c.Error(err)
		}
	}
}
