package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GetCredentialsHandler interface {
	Handle(c *gin.Context)
}

func NewGetCredentialsHandler(service services.GetCredentialsService) GetCredentialsHandler {
	return &getCredentialsHandlerImpl{service: service}
}

type getCredentialsHandlerImpl struct {
	service services.GetCredentialsService
}

func (h *getCredentialsHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	credentials, err := h.service.Get(c, token, time.Now())
	if err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
		}, false)
		return
	}

	c.JSON(http.StatusOK, credentials)
}
