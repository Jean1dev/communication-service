package sockets

import "log"

func SimpleMessageEventHandler() EventHandler {
	return func(event EventMessage, c *ClientSocket) error {
		log.Println(event)
		return nil
	}
}
