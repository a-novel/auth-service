package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RotateSecretKeysHandler interface {
	Handle(c *gin.Context)
}

func NewRotateSecretKeysHandler(service services.RotateSecretKeysService) RotateSecretKeysHandler {
	return &rotateSecretKeysHandlerImpl{
		service: service,
	}
}

type rotateSecretKeysHandlerImpl struct {
	service services.RotateSecretKeysService
}

func (h *rotateSecretKeysHandlerImpl) Handle(c *gin.Context) {
	if err := h.service.RotateSecretKeys(c); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}
