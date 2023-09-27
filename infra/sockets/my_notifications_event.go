package sockets

import (
	"communication-service/application"
	"encoding/json"
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
		outgoingEvent.Type = "Response"

		c.Egress <- outgoingEvent
		return nil
	}
}
