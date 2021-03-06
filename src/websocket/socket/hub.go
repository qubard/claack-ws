package socket

import (
	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/lib/util"
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
	}
}

func (hub *Hub) UnregisterClient(client *Client) {
	// Don't unregister the user if they aren't even authed
	if client.Credentials != nil {
		hub.EdgeServer.UnregisterClient(client)
	}
	delete(hub.clients, client)
	close(client.Send)
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				hub.UnregisterClient(client)
			}
		// Loop over all the clients and broadcast any messages
		case message := <-hub.EdgeServer.GlobalChan:
			for client := range hub.clients {
				select {
				case client.Send <- []byte(message.Payload):
				default:
					hub.UnregisterClient(client)
				}
			}
		case message := <-hub.EdgeServer.RelayChan:
			// Process messages only intended for specific users on this server
			relay, err := util.ReadRelayMessage(message.Payload)
			if err == nil {
				// Valid auth token for relay, lookup which client it belongs to
				if client := hub.EdgeServer.AcquireClient(relay.DstId); client != nil {
					// We have found a valid authorized client
					// Relay the message to the client
					client.Send <- []byte(relay.Message)
				}
			}
		}
	}
}
