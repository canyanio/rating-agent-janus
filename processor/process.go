package processor

import (
	"context"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/canyanio/rating-agent-janus/client/rabbitmq"
	"github.com/canyanio/rating-agent-janus/model"
	"github.com/canyanio/rating-agent-janus/state"
)

// JanusProcessorInterface is the interface for Server objects
type JanusProcessorInterface interface {
	Process(ctx context.Context, event *model.Event) error
}

// JanusProcessor is the Janus processor
type JanusProcessor struct {
	state  state.ManagerInterface
	client rabbitmq.ClientInterface
}

// NewJanusProcessor initializes a new Janus processor
func NewJanusProcessor(state state.ManagerInterface, client rabbitmq.ClientInterface) *JanusProcessor {
	return &JanusProcessor{
		state:  state,
		client: client,
	}
}

// Process raw bytes containing a Janus packet
func (s *JanusProcessor) Process(ctx context.Context, event *model.Event) error {
	if event.Type == model.JanusEventTypePlugin {
		eventPlugin := &model.EventPlugin{}
		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   eventPlugin,
			TagName:  "json",
		}
		decoder, _ := mapstructure.NewDecoder(cfg)
		decoder.Decode(event.Event)
		if eventPlugin == nil {
			return errors.New("unable to parse Janus Plugin event")
		}

		// Janus SIP plugin
		if eventPlugin.Plugin == model.JanusSIPPlugin {
			sipEvent := &model.EventPluginSIP{}
			cfg := &mapstructure.DecoderConfig{
				Metadata: nil,
				Result:   sipEvent,
				TagName:  "json",
			}
			decoder, _ := mapstructure.NewDecoder(cfg)
			decoder.Decode(eventPlugin.Data)
			if sipEvent == nil {
				return errors.New("unable to parse Janus SIP Plugin event")
			}

			if sipEvent.SIP != "" {
				sipMessage := model.SIPMessageFromString(sipEvent.SIP)
				sipMessage.Timestamp = time.Unix(event.Timestamp/1000000, event.Timestamp%1000000)
				err := s.processSIPMessage(ctx, sipMessage)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
