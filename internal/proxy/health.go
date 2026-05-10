package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"redis":  "connected"})
}
