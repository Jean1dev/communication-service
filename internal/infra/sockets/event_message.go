package sockets

import "encoding/json"

type EventMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event EventMessage, c *ClientSocket) error

const (
	EventMyNotifications         = "my_notifications"
	EventSentNoficationsCaixinha = "notify_all_members_caixinha"
	ChatEvent                    = "chat"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}
