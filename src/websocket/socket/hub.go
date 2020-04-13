package socket

import (
	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/lib/util"
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
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clients:       make(map[*Client]bool),
		Bus:           CreateHandlerBus(database),
		EdgeServer:    edgeServer,
	}
}

func (hub *Hub) UnregisterClient(client *Client) {
	hub.EdgeServer.UnregisterClient(client)
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
				// TODO: no more global keys (FLAG THIS ONCE AGAIN)
				username, ok := util.ExtractField(relay.DstToken, "username", []byte("key"))
				if ok {
					// Valid auth token for relay, lookup which client it belongs to
					if client, present := hub.EdgeServer.ClientTable[username.(string)]; present {
						// We have found a valid authorized client
						// Relay the message to the client
						// connected to.
						client.Send <- []byte(relay.Message)
					}
				}
			}
		}
	}
}
