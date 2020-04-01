package main

import (
	"github.com/gorilla/websocket"
    "github.com/qubard/claack-go/messages/handlers"
)

type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
    handlerBus *handlers.HandlerBus
    upgrader *websocket.Upgrader
}

func CreateHub() *Hub {
	return &Hub {
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
        handlerBus: handlers.CreateHandlerBus(),
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
        // Loop over all the clients and broadcast any messages 
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
                    case client.send <-message:
                    default:
                        // Close the socket
                        close(client.send)
                        // Remove the client from the map
                        delete(hub.clients, client)
				}
			}
		}
	}
}
