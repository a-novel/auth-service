package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type IntrospectTokenHandler interface {
	Handle(c *gin.Context)
}

func NewIntrospectTokenHandler(service services.IntrospectTokenService) IntrospectTokenHandler {
	return &introspectTokenHandlerImpl{service: service}
}

type introspectTokenHandlerImpl struct {
	service services.IntrospectTokenService
}

func (h *introspectTokenHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	result, err := h.service.IntrospectToken(c, token, time.Now(), true)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
