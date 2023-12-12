package sockets

import (
	"communication-service/application"
	"encoding/json"
	"fmt"
	"log"
)

type NotificationInput struct {
	User      string   `json:"user"`
	Desc      string   `json:"desc"`
	Caixinhas []string `json:"caixinhas"`
}

type NotificationOutput struct {
	Desc string `json:"desc"`
}

func broadcastoToAllUsers(userDestination string, input NotificationInput, c *ClientSocket) {
	for wsclient := range c.manager.clients {
		if wsclient.Opt == userDestination {
			var outgoingEvent EventMessage
			notificationOutput := &NotificationOutput{Desc: fmt.Sprintf("%s -> %s", input.User, input.Desc)}
			data, err := json.Marshal(notificationOutput)

			if err != nil {
				continue
			}

			outgoingEvent.Payload = data
			outgoingEvent.Type = "Notification"
			wsclient.Egress <- outgoingEvent
		}
	}
}

func SentNoficationsCaixinha() EventHandler {
	return func(event EventMessage, c *ClientSocket) error {
		var input NotificationInput
		err := json.Unmarshal(event.Payload, &input)

		if err != nil {
			log.Println("Erro ao fazer o unmarshal do JSON:", err)
			return err
		}

		list := application.NewNofiticationForCaixinha(input.Desc, input.User, input.Caixinhas)

		for _, it := range list {

			if it == input.User {
				continue
			}

			broadcastoToAllUsers(it, input, c)

		}

		return nil
	}
}
