package sockets

import (
	"encoding/json"
	"errors"
	"log"
)

type ChatInput struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

type SendMesageOutput struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

func messageTo(receiver string, input SendMesageOutput, c *ClientSocket) error {
	for wsclient := range c.manager.clients {
		if wsclient.Opt == receiver {
			var outgoingEvent EventMessage
			data, err := json.Marshal(input)

			if err != nil {
				return err
			}

			outgoingEvent.Payload = data
			outgoingEvent.Type = "ChatResponse"
			wsclient.Egress <- outgoingEvent
			return nil
		}
	}

	return errors.New("receiver not found")
}

func RealTimeChatHandler() EventHandler {
	return func(event EventMessage, c *ClientSocket) error {
		var input ChatInput
		err := json.Unmarshal(event.Payload, &input)
		if err != nil {
			log.Println("Error unmarshalling chat input: ", err)
			return err
		}

		if err := messageTo(input.Receiver, SendMesageOutput{Message: input.Message, Sender: input.Sender}, c); err != nil {
			log.Println("Error sending message to receiver: ", err)
			return err
		}

		return nil
	}
}
