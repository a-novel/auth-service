package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PreviewHandler interface {
	Handle(c *gin.Context)
}

func NewPreviewHandler(service services.PreviewService) PreviewHandler {
	return &previewHandlerImpl{service: service}
}

type previewHandlerImpl struct {
	service services.PreviewService
}

func (h *previewHandlerImpl) Handle(c *gin.Context) {
	slug := c.Query("slug")

	preview, err := h.service.Preview(c, slug)
	if err != nil {
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrNotFound, http.StatusNotFound},
		})
		return
	}

	c.JSON(http.StatusOK, preview)
}
