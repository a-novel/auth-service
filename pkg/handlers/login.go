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

type LoginHandler interface {
	Handle(c *gin.Context)
}

func NewLoginHandler(service services.LoginService) LoginHandler {
	return &loginHandlerImpl{service: service}
}

type loginHandlerImpl struct {
	service services.LoginService
}

func (h *loginHandlerImpl) Handle(c *gin.Context) {
	request := new(models.LoginForm)
	if err := c.BindJSON(request); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, err := h.service.Login(c, request.Email, request.Password, time.Now())
	if err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
			{bunovel.ErrNotFound, http.StatusNotFound},
			{goframework.ErrInvalidEntity, http.StatusUnprocessableEntity},
		}, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token.TokenRaw})
}
