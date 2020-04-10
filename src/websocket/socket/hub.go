package socket

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	Bus        *HandlerBus
	Upgrader   *websocket.Upgrader
}

func CreateHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Bus:        CreateHandlerBus(),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.Send)
			}
			// Loop over all the clients and broadcast any messages
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.Send <- message:
				default:
					// Close the socket
					close(client.Send)
					// Remove the client from the map
					delete(hub.clients, client)
				}
			}
		}
	}
}
