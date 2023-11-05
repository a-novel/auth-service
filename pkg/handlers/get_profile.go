package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GetProfileHandler interface {
	Handle(c *gin.Context)
}

func NewGetProfileHandler(service services.GetProfileService) GetProfileHandler {
	return &getProfileHandlerImpl{service: service}
}

type getProfileHandlerImpl struct {
	service services.GetProfileService
}

func (h *getProfileHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	profile, err := h.service.Get(c, token, time.Now())
	if err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidCredentials, http.StatusForbidden},
		}, false)
		return
	}

	c.JSON(http.StatusOK, profile)
}
