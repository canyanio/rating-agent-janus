package model

// Event stores the event received from a Janus Gateway instance
type Event struct {
	Type      int         `json:"type"`
	Timestamp uint64      `json:"timestamp"`
	SessionID uint64      `json:"session_id"`
	HandleID  uint64      `json:"handle_id"`
	Event     interface{} `json:"event"`
}
