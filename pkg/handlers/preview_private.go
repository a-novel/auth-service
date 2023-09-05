package handlers

import (
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type PreviewPrivateHandler interface {
	Handle(c *gin.Context)
}

func NewPreviewPrivateHandler(service services.PreviewPrivateService) PreviewPrivateHandler {
	return &previewPrivateHandlerImpl{service: service}
}

type previewPrivateHandlerImpl struct {
	service services.PreviewPrivateService
}

func (h *previewPrivateHandlerImpl) Handle(c *gin.Context) {
	token := c.GetHeader("Authorization")

	preview, err := h.service.Preview(c, token, time.Now())
	if err != nil {
		errors.ErrorToHTTPCode(c, err, []errors.HTTPError{
			{errors.ErrInvalidCredentials, http.StatusForbidden},
		})
		return
	}

	c.JSON(http.StatusOK, preview)
}
