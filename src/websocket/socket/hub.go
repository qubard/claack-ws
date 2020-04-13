package socket

import (
	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/lib/postgres"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	Bus        *HandlerBus
	EdgeServer *EdgeServer
	Upgrader   *websocket.Upgrader
}

func CreateHub(database *postgres.Database, edgeServer *EdgeServer) *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Bus:        CreateHandlerBus(database),
		EdgeServer: edgeServer,
		clientLookup: make(map[string]*Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				hub.EdgeServer.UnregisterUser(client.Username)
				delete(hub.clients, client)
				close(client.Send)
			}
		// Loop over all the clients and broadcast any messages
		case message := <-hub.EdgeServer.GlobalChan:
			for client := range hub.clients {
				select {
				case client.Send <- []byte(message.Payload):
				default:
					// Close the socket
					close(client.Send)
					// Remove the client from the map
					delete(hub.clients, client)
					hub.EdgeServer.UnregisterUser(client.Username)
				}
			}
		/*case message := <-hub.EdgeServer.MainChan:
			// Process messages only intended for specific users on this server
			// Find the user it belongs to
			message = nil*/
		}
	}
}
