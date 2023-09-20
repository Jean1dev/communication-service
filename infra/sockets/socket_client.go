package sockets

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*ClientSocket]bool

type ClientSocket struct {
	Opt        string
	connection *websocket.Conn
	manager    *ConnectionManager
	egress     chan []byte
}

func NewSocketClient(conn *websocket.Conn, manager *ConnectionManager, opt string) *ClientSocket {
	return &ClientSocket{
		Opt:        opt,
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func (c *ClientSocket) ReadMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break
		}

		log.Println("messageType: ", messageType)
		log.Println("Payload ", string(payload))

		var request EventMessage
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message %v", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("error handeling message ", err.Error())
		}
	}
}

func (c *ClientSocket) WriteMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed ", err)
				}

				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
		}

		log.Println("message sent")
	}
}
