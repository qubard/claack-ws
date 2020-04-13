package socket

import (
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/lib/postgres"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	Bus        *HandlerBus
	globalChan <-chan *redis.Message // Messages received by every hub node
	mainChan    <-chan *redis.Message // Messages specific to this hub only
	id         string
	Upgrader   *websocket.Upgrader
}

func CreateHub(database *postgres.Database, globalChan <-chan *redis.Message, mainChan <-chan *redis.Message) *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Bus:        CreateHandlerBus(database),
		globalChan: globalChan, 
		mainChan:    mainChan,    
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
		case message := <-hub.globalChan:
			for client := range hub.clients {
				select {
				case client.Send <- []byte(message.Payload):
				default:
					// Close the socket
					close(client.Send)
					// Remove the client from the map
					delete(hub.clients, client)
				}
			}
		case message := <-hub.mainChan:
			// Process messages only intended for specific users on this server
		} 
	}
}
