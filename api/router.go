package api

import (
	"github.com/gin-gonic/gin"

	"github.com/canyanio/rating-agent-janus/client/rabbitmq"
	"github.com/canyanio/rating-agent-janus/state"
)

// API URL used by the HTTP router
const (
	APIURLStatus = "/status"

	APIURLJanusWebhook = "/api/v1/janus-gateway"
)

// NewRouter returns the gin router
func NewRouter(state state.ManagerInterface, client rabbitmq.ClientInterface) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	status := NewStatusController()
	router.GET(APIURLStatus, status.Status)

	janus := NewJanusController(state, client)
	router.POST(APIURLJanusWebhook, janus.Webhook)

	return router
}
