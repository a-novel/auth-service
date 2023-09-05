package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"net/http"
)

type SlugExistsHandler interface {
	Handle(c *gin.Context)
}

func NewSlugExistsHandler(service services.SlugExistsService) SlugExistsHandler {
	return &slugExistsHandlerImpl{
		service: service,
	}
}

type slugExistsHandlerImpl struct {
	service services.SlugExistsService
}

func (h *slugExistsHandlerImpl) Handle(c *gin.Context) {
	slug := c.Query("slug")

	ok, err := h.service.SlugExists(c, slug)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(lo.Ternary(ok, http.StatusNoContent, http.StatusNotFound))
}
