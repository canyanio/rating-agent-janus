package api

import (
	"github.com/gin-gonic/gin"
)

// API URL used by the HTTP router
const (
	APIURLStatus = "/status"

	APIURLJanusWebhook = "/api/v1/janus-gateway"
)

// NewRouter returns the gin router
func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	status := NewStatusController()
	router.GET(APIURLStatus, status.Status)

	janus := NewJanusController()
	router.POST(APIURLJanusWebhook, janus.Webhook)

	return router
}
