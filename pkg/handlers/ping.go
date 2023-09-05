package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PingHandler interface {
	Handle(c *gin.Context)
}

func NewPingHandler() PingHandler {
	return &pingHandlerImpl{}
}

type pingHandlerImpl struct{}

func (h *pingHandlerImpl) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
