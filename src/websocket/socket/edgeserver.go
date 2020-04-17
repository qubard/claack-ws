package socket

import (
	"github.com/go-redis/redis/v7"
	"github.com/qubard/claack-go/lib/util"
	"github.com/qubard/claack-go/websocket/messages/types"
	"sync"
)

type EdgeServer struct {
	redis       *redis.Client
	id          string                // Id of this edge server for incoming relays
	GlobalChan  <-chan *redis.Message // Messages received by every hub node
	RelayChan   <-chan *redis.Message // Messages specific to this hub only
	ClientTable map[string]*Client
	Mutex       sync.Mutex // We can either use a channel and a separate goroutine or a mutex to update ClientTable
}

func CreateSimpleEdgeServer(client *redis.Client) *EdgeServer {
	return (&EdgeServer{}).AttachRedis(client)
}

func CreateEdgeServer(client *redis.Client, name string) *EdgeServer {
	globalChan := util.CreateSubChannel(client, "global")
	relayChan := util.CreateSubChannel(client, name)

	if globalChan == nil || relayChan == nil {
		panic("Failed to create hub: invalid redis global channel or hub channel")
	}

	return (&EdgeServer{
		id:          name,
		GlobalChan:  globalChan,
		RelayChan:   relayChan,
		ClientTable: make(map[string]*Client),
		Mutex:       sync.Mutex{},
	}).AttachRedis(client)
}

// Find the channel the user maps to
func (server *EdgeServer) FindClient(client *Client) (string, error) {
	return server.FindClientById(client.Credentials.Username)
}

func (server *EdgeServer) FindClientById(id string) (string, error) {
	return server.redis.Get(id).Result()
}

func (server *EdgeServer) AttachRedis(redis *redis.Client) *EdgeServer {
	server.redis = redis
	return server
}

// Sends a message to the user with id dstId
// If it currently exists on the server, otherwise relay it!
func (server *EdgeServer) RelayMessage(dstId string, msg interface{}) error {
	bytes, err := util.WritePackedMessage(msg)

	if err != nil {
		return err
	}

	if client := server.AcquireClient(dstId); client != nil {
		client.Send <- bytes
	} else {
		// Client is not present locally, construct a relay message
		// Do a lookup for their channelId, then relay a message there
		server.SendLookupMessage(dstId, bytes)
	}

	return nil
}

func (server *EdgeServer) SendLookupMessage(dstId string, bytes []byte) {
	if channelId, err := server.FindClientById(dstId); err == nil {
		relayMsg := types.RelayMessage{
			DstId:   dstId,
			Message: string(bytes),
		}
		// Relay the message to the desired server
		if relayBytes, err := util.WritePackedMessage(relayMsg); err == nil {
			server.redis.Publish(channelId, relayBytes)
		}
	}
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
	err := server.redis.Set(client.Credentials.Username, server.id, 0).Err()
	server.ClientTable[client.Credentials.Username] = client

	if err != nil {
		panic(err)
	}
}

func (server *EdgeServer) Publish(channel string, msg interface{}) {
	server.redis.Publish(channel, msg)
}

func (server *EdgeServer) UnregisterClient(client *Client) {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	server.redis.Del(client.Credentials.Username)
	delete(server.ClientTable, client.Credentials.Username)
}

func (server *EdgeServer) IsClientPresent(client *Client) bool {
	_, err := server.redis.Get(client.Credentials.Username).Result()
	return err != nil
}
