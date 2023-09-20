package sockets

import "encoding/json"

type EventMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event EventMessage, c *ClientSocket) error

const (
	EventSimpleMessage   = "simple_message"
	EventMyNotifications = "my_notifications"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}
