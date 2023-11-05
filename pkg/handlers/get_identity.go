package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GetIdentityHandler interface {
	Handle(c *gin.Context)
}

func NewGetIdentityHandler(service services.GetIdentityService) GetIdentityHandler {
	return &getIdentityHandlerImpl{service: service}
}

type getIdentityHandlerImpl struct {
	service services.GetIdentityService
}

func (h *getIdentityHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	identity, err := h.service.Get(c, token, time.Now())
	if err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
		}, false)
		return
	}

	c.JSON(http.StatusOK, identity)
}
