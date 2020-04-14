package socket

import (
	"sync"
	"github.com/qubard/claack-go/lib/util"
	"github.com/qubard/claack-go/websocket/messages/types"
	"github.com/go-redis/redis/v7"
)

type EdgeServer struct {
	Redis	   *redis.Client
	Id 		   string // Id of this edge server for incoming relays
	GlobalChan <-chan *redis.Message // Messages received by every hub node
	RelayChan  <-chan *redis.Message // Messages specific to this hub only
	ClientTable map[string]*Client
	Mutex       sync.Mutex // We can either use a channel and a separate goroutine or a mutex to update ClientTable
}

// Find the channel the user maps to
func (server *EdgeServer) FindClient(client *Client) (string, error) {
	return server.FindClientById(client.Username)
}

func (server *EdgeServer) FindClientById(id string) (string, error) {
	return server.Redis.Get(id).Result()
}

// Sends a message to the user with id dstId
// If it currently exists on the server, otherwise relay it!
func (server* EdgeServer) RelayMessage(dstId string, msg interface{}) error {
	bytes, err := util.WritePackedMessage(msg)
	
	if err != nil {
		return err
	}

	// Generally avoid this case because it acquires the lock
	// But in practice this can't slow us down TOO much?
	if client := server.AcquireClient(dstId); client != nil {
		client.Send <- bytes
	} else {
		// Client is not present locally, construct a relay message
		// Do a lookup for their channelId, then relay a message there
		if channelId, err := server.FindClientById(dstId); err == nil {
			relayMsg := types.RelayMessage {
				DstId: dstId,
				Message: string(bytes),
			}
			// Relay the message to the other server
			if relayBytes, err := util.WritePackedMessage(relayMsg); err == nil {
				server.Redis.Publish(channelId, relayBytes)
			}
		}
	}

	return nil
}

func (server *EdgeServer) AcquireClient(id string) *Client {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	return server.ClientTable[id]
}

func (server *EdgeServer) RegisterClient(client *Client) {
	// Register user in redis to the current server's Id
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	err := server.Redis.Set(client.Username, server.Id, 0).Err()
	server.ClientTable[client.Username] = client

	if err != nil {
		panic(err)
	}
}

func (server *EdgeServer) UnregisterClient(client *Client) {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	server.Redis.Del(client.Username)
	delete(server.ClientTable, client.Username)
}

func (server *EdgeServer) IsClientPresent(client *Client) bool {
	_, err := server.Redis.Get(client.Username).Result()
	return err != nil
}