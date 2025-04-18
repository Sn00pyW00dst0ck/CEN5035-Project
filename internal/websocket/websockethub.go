package websocket

import (
	"github.com/gorilla/websocket"
)

type WebSocketHub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

func (hub *WebSocketHub) Run() {
	for {
		select {
		case conn := <-hub.Register:
			hub.Clients[conn] = true
		case conn := <-hub.Unregister:
			if _, ok := hub.Clients[conn]; ok {
				delete(hub.Clients, conn)
				conn.Close()
			}
		case msg := <-hub.Broadcast:
			for conn := range hub.Clients {
				conn.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}
