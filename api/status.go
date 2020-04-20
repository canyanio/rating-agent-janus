package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatusController contains status-related end-points
type StatusController struct{}

// NewStatusController returns a new StatusController
func NewStatusController() *StatusController {
	return new(StatusController)
}

// Status responds to GET /status
func (h StatusController) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
