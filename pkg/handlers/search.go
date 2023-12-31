package handlers

import (
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-apis"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SearchHandler interface {
	Handle(c *gin.Context)
}

func NewSearchHandler(service services.SearchService) SearchHandler {
	return &searchHandlerImpl{
		service: service,
	}
}

type searchHandlerImpl struct {
	service services.SearchService
}

func (s *searchHandlerImpl) Handle(c *gin.Context) {
	query := new(models.SearchQuery)
	if err := c.BindQuery(query); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	users, total, err := s.service.Search(c, query.Query, query.Limit, query.Offset)
	if err != nil {
		apis.ErrorToHTTPCode(c, err, []apis.HTTPError{
			{goframework.ErrInvalidEntity, http.StatusBadRequest},
		}, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"res":   users,
		"total": total,
	})
}
