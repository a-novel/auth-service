package handlers

import (
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListHandler interface {
	Handle(c *gin.Context)
}

func NewListHandler(service services.ListService) ListHandler {
	return &listHandlerImpl{
		service: service,
	}
}

type listHandlerImpl struct {
	service services.ListService
}

func (l *listHandlerImpl) Handle(c *gin.Context) {
	query := new(models.ListQuery)
	if err := c.BindQuery(query); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	users, err := l.service.List(c, query.IDs.Value())
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
