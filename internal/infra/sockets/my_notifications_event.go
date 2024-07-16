package sockets

import (
	"encoding/json"

	"github.com/Jean1dev/communication-service/internal/application"
)

func MyNotificationsEventHandler() EventHandler {
	return func(event EventMessage, c *ClientSocket) error {
		results := application.GetMyNotifications(c.Opt)
		data, err := json.Marshal(results)

		if err != nil {
			return err
		}

		var outgoingEvent EventMessage
		outgoingEvent.Payload = data
		outgoingEvent.Type = "MyNotificationsResponse"

		c.Egress <- outgoingEvent
		return nil
	}
}
