package model

import (
	"time"

	"github.com/mendersoftware/go-lib-micro/config"

	dconfig "github.com/canyanio/rating-agent-janus/config"
	"github.com/sipcapture/heplify-server/sipparser"
)

// SIPMessage represents a SIP message
type SIPMessage struct {
	*sipparser.SipMsg
	AccountTag            string
	DestinationAccountTag string
	Timestamp             time.Time
}

// SIPMessageFromString returns a SIPMessage from a string
func SIPMessageFromString(payload string) *SIPMessage {
	return parseSIPMessage(payload)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func parseSIPMessage(payload string) *SIPMessage {
	sipHeaderCaller := config.Config.GetString(dconfig.SettingSIPHeaderCaller)
	sipHeaderCallee := config.Config.GetString(dconfig.SettingSIPHeaderCallee)
	sipLocalDomains := config.Config.GetStringSlice(dconfig.SettingSIPLocalDomains)
	return parseSIPMessageWithSettings(payload, sipHeaderCaller, sipHeaderCallee, sipLocalDomains)
}

func parseSIPMessageWithSettings(payload string, sipHeaderCaller string, sipHeaderCallee string, sipLocalDomains []string) *SIPMessage {
	customHeaders := []string{}
	if sipHeaderCaller != "" {
		customHeaders = append(customHeaders, sipHeaderCaller)
	}
	if sipHeaderCallee != "" {
		customHeaders = append(customHeaders, sipHeaderCallee)
	}

	msg := sipparser.ParseMsg(payload, []string{}, customHeaders)

	accountTag := ""
	if sipHeaderCaller != "" && msg.CustomHeader[sipHeaderCaller] != "" {
		accountTag = msg.CustomHeader[sipHeaderCaller]
	} else if msg.PAssertedId != nil && (sipLocalDomains == nil || stringInSlice(msg.PaiHost, sipLocalDomains)) {
		accountTag = msg.PaiUser
	} else if sipLocalDomains != nil && stringInSlice(msg.FromHost, sipLocalDomains) {
		accountTag = msg.FromUser
	}

	destinationAccountTag := ""
	if sipHeaderCallee != "" && msg.CustomHeader[sipHeaderCallee] != "" {
		destinationAccountTag = msg.CustomHeader[sipHeaderCallee]
	} else if sipLocalDomains != nil && stringInSlice(msg.ToHost, sipLocalDomains) {
		destinationAccountTag = msg.ToUser
	}

	sipMessage := SIPMessage{
		SipMsg:                msg,
		AccountTag:            accountTag,
		DestinationAccountTag: destinationAccountTag,
	}
	return &sipMessage
}
