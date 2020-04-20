package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mendersoftware/go-lib-micro/log"

	"github.com/canyanio/rating-agent-janus/model"
)

// JanusController container for end-points
type JanusController struct {
}

// NewJanusController returns a new StatusController
func NewJanusController() *JanusController {
	return &JanusController{}
}

// Webhook responds to POST /api/v1/metadata/workflows
func (h JanusController) Webhook(c *gin.Context) {
	l := log.FromContext(c.Request.Context())

	events := make([]model.Event, 100)

	err := c.BindJSON(&events)
	if err != nil {
		l.Errorf("unable to bind JSON: %v", err)

		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	l.Infof("events: %v", events)

	c.Status(http.StatusCreated)
}
