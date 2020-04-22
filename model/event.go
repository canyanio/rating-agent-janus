package model

// Constants for the Janus events
const (
	JanusEventTypeSession int = 1 << iota
	JanusEventTypeHandle
	JanusEventTypeExternal
	JanusEventTypeJSep
	JanusEventTypeWebRTC
	JanusEventTypeMedia
	JanusEventTypePlugin
	JanusEventTypeTransport
	JanusEventTypeCore
)

// JanusSIPPlugin is the name of the SIP plugin in Janus
const JanusSIPPlugin = "janus.plugin.sip"

// Event stores the event received from a Janus Gateway instance
type Event struct {
	Emitter   string                 `json:"emitter,omitempty"`
	Type      int                    `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	SessionID int64                  `json:"session_id"`
	HandleID  int64                  `json:"handle_id"`
	OpaqueID  string                 `json:"opaque_id,omitempty"`
	Event     map[string]interface{} `json:"event"`
}

// EventPlugin stores the details of an event from a plugin
type EventPlugin struct {
	Plugin string                 `json:"plugin"`
	Data   map[string]interface{} `json:"data"`
}

// EventPluginSIP stores the data of a SIP event from the janus.plugin.sip plugin
type EventPluginSIP struct {
	Event string `json:"event"`
	SIP   string `json:"sip"`
}
