package sockets

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     CheckOrigin,
	}
)

type ConnectionManager struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func CheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	log.Printf("Origin is %s", origin)
	return true
}

func NewManager() *ConnectionManager {
	m := &ConnectionManager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (c *ConnectionManager) setupEventHandlers() {
	c.handlers[EventSimpleMessage] = SimpleMessageEventHandler()
	c.handlers[EventMyNotifications] = MyNotificationsEventHandler()
}

func (c *ConnectionManager) routeEvent(event EventMessage, client *ClientSocket) error {
	if handler, ok := c.handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			log.Printf("error in event %s %v", event.Type, err.Error())
			return err
		}

		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (c *ConnectionManager) ServeWS(w http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Printf("New Connection %s", otp)
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error open connection %s", err.Error())
		return
	}

	client := NewSocketClient(conn, c, otp)
	c.addClient(client)

	go client.ReadMessages()
	go client.WriteMessages()
}

func (c *ConnectionManager) addClient(client *ClientSocket) {
	c.Lock()
	defer c.Unlock()

	c.clients[client] = true
}

func (c *ConnectionManager) removeClient(client *ClientSocket) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.clients[client]; ok {
		client.connection.Close()
		delete(c.clients, client)
	}
}
