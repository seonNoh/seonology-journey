package ws

import "encoding/json"

// MessageType identifies the kind of client-to-server message.
type MessageType string

const (
	// MsgTypePing - keep-alive from client.
	MsgTypePing MessageType = "ping"
	// MsgTypeLocationUpdate - location sharing update.
	MsgTypeLocationUpdate MessageType = "location_update"
	// MsgTypeCursorMove - collaborative editing cursor.
	MsgTypeCursorMove MessageType = "cursor_move"
	// MsgTypeSubscribe - subscribe to additional events.
	MsgTypeSubscribe MessageType = "subscribe"
)

// IncomingMessage is the envelope for all client→server WebSocket messages.
type IncomingMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// LocationPayload is sent with MsgTypeLocationUpdate.
type LocationPayload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
}

// ServerMessageType identifies server→client messages.
type ServerMessageType string

const (
	// ServerMsgPong - response to client ping.
	ServerMsgPong ServerMessageType = "pong"
	// ServerMsgEvent - realtime event (schedule change, etc).
	ServerMsgEvent ServerMessageType = "event"
	// ServerMsgLocation - another user's location update.
	ServerMsgLocation ServerMessageType = "location"
	// ServerMsgError - error notification.
	ServerMsgError ServerMessageType = "error"
)

// OutgoingMessage is the envelope for all server→client WebSocket messages.
type OutgoingMessage struct {
	Type    ServerMessageType `json:"type"`
	Payload json.RawMessage   `json:"payload,omitempty"`
}
