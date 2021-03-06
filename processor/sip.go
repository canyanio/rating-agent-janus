package processor

import (
	"context"
	"github.com/pkg/errors"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/sirupsen/logrus"

	"github.com/canyanio/rating-agent-janus/client/rabbitmq"
	dconfig "github.com/canyanio/rating-agent-janus/config"
	"github.com/canyanio/rating-agent-janus/model"
)

// Server handler specific constants
const (
	MethodInvite          = "INVITE"
	MethodAck             = "ACK"
	MethodBye             = "BYE"
	MethodCancel          = "CANCEL"
	StateManagerTTLInvite = 600
	StateManagerTTLCall   = 3600 * 6
)

func (s *JanusProcessor) processSIPMessage(ctx context.Context, msg *model.SIPMessage) error {
	reqID := uuid.New()

	l := log.FromContext(ctx)

	requestMethod := msg.FirstMethod
	if requestMethod == "" {
		return errors.New("invalid SIP message")
	}

	callID := msg.CallID
	CSeqParts := strings.SplitN(msg.Cseq.Val, " ", 2)
	CSeqID := CSeqParts[0]

	l.WithFields(logrus.Fields{
		"req-id":        reqID,
		"requestMethod": requestMethod,
		"callID":        callID,
		"CSeqID":        CSeqID,
	}).Debug("received msg")

	var routingKey string
	var req interface{}
	if requestMethod == MethodInvite {
		call := &model.Call{
			Tenant:                config.Config.GetString(dconfig.SettingTenant),
			TransactionTag:        callID,
			AccountTag:            msg.FromUser,
			DestinationAccountTag: msg.ToUser,
			Source:                "sip:" + msg.FromUser + "@" + msg.FromHost,
			Destination:           "sip:" + msg.ToUser + "@" + msg.ToHost,
			TimestampInvite:       msg.Timestamp,
			CSeq:                  CSeqID,
		}
		err := s.state.Set(ctx, callID, call, StateManagerTTLInvite)
		if err != nil {
			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"method":  requestMethod,
				"call-id": callID,
				"ts":      msg.Timestamp,
				"err":     err.Error(),
			}).Error("unable to set the call status")
		} else {
			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"method":  requestMethod,
				"call-id": callID,
				"ts":      msg.Timestamp,
			}).Debug("call status set in the state manager, waiting for the ACK")
		}
	} else if requestMethod == MethodAck {
		var call model.Call
		err := s.state.Get(ctx, callID, &call)
		if err != nil || call.CSeq == "" {
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"method":  requestMethod,
				"call-id": callID,
				"ts":      msg.Timestamp,
				"err":     errStr,
			}).Error("unable to retrieve the call status, INVITE has not been received for this calls")
			return err
		}

		if CSeqID == call.CSeq && call.TimestampAck.IsZero() &&
			(call.AccountTag != "" || call.DestinationAccountTag != "") {
			call.TimestampAck = msg.Timestamp
			s.state.Set(ctx, call.TransactionTag, call, StateManagerTTLCall)

			routingKey = rabbitmq.QueueNameBeginTransaction
			req = &model.BeginTransaction{
				Request: model.BeginTransactionRequest{
					Tenant:                call.Tenant,
					TransactionTag:        call.TransactionTag,
					AccountTag:            call.AccountTag,
					DestinationAccountTag: call.DestinationAccountTag,
					Source:                call.Source,
					Destination:           call.Destination,
					TimestampBegin:        msg.Timestamp.UTC().Format(time.RFC3339),
				},
			}

			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"method":  requestMethod,
				"call-id": callID,
				"ts":      msg.Timestamp,
			}).Debug("call start detected: begin transaction")
		}
	} else if requestMethod == MethodBye || requestMethod == MethodCancel {
		s.state.Delete(ctx, callID)

		routingKey = rabbitmq.QueueNameEndTransaction
		req = &model.EndTransaction{
			Request: model.EndTransactionRequest{
				Tenant:                config.Config.GetString(dconfig.SettingTenant),
				TransactionTag:        callID,
				AccountTag:            msg.FromUser,
				DestinationAccountTag: msg.ToUser,
				TimestampEnd:          msg.Timestamp.UTC().Format(time.RFC3339),
			},
		}

		l.WithFields(logrus.Fields{
			"req-id":  reqID,
			"method":  requestMethod,
			"call-id": callID,
			"ts":      msg.Timestamp,
		}).Debug("call end detected: end transaction")
	}

	if req != nil {
		err := s.client.Publish(ctx, routingKey, req)
		if err != nil {
			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"method":  requestMethod,
				"call-id": callID,
				"err":     err.Error(),
			}).Error("unable to publish the request")

			return err
		}
	}

	return nil
}
