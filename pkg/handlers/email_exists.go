package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"net/http"
)

type EmailExistsHandler interface {
	Handle(c *gin.Context)
}

func NewEmailExistsHandler(service services.EmailExistsService) EmailExistsHandler {
	return &emailExistsHandlerImpl{
		service: service,
	}
}

type emailExistsHandlerImpl struct {
	service services.EmailExistsService
}

func (h *emailExistsHandlerImpl) Handle(c *gin.Context) {
	email := c.Query("email")

	ok, err := h.service.EmailExists(c, email)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(lo.Ternary(ok, http.StatusNoContent, http.StatusNotFound))
}
