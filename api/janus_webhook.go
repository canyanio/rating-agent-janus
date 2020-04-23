package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/pkg/errors"

	"github.com/canyanio/rating-agent-janus/client/rabbitmq"
	"github.com/canyanio/rating-agent-janus/model"
	"github.com/canyanio/rating-agent-janus/processor"
	"github.com/canyanio/rating-agent-janus/state"
)

// JanusController container for end-points
type JanusController struct {
	processor processor.JanusProcessorInterface
}

// NewJanusController returns a new StatusController
func NewJanusController(state state.ManagerInterface, client rabbitmq.ClientInterface) *JanusController {
	processor := processor.NewJanusProcessor(state, client)

	return &JanusController{
		processor: processor,
	}
}

// Webhook responds to POST /api/v1/metadata/workflows
func (h JanusController) Webhook(c *gin.Context) {
	ctx := c.Request.Context()
	l := log.FromContext(ctx)

	events := make([]*model.Event, 100)
	err := c.BindJSON(&events)
	if err != nil {
		l.Errorf("unable to bind JSON: %v", err)

		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var latestError error = nil
	for _, event := range events {
		err := h.processor.Process(ctx, event)
		if err != nil {
			l.Error(errors.Wrapf(err, "unable to process event: %v", event))
			latestError = err
		}
	}

	if latestError != nil {
		c.JSON(http.StatusBadRequest, latestError.Error())
	} else {
		c.Status(http.StatusCreated)
	}
}
