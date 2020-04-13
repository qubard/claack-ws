package socket

import (
	"github.com/go-redis/redis/v7"
)

type EdgeServer struct {
	Redis	   *redis.Client
	Id 		   string
	GlobalChan <-chan *redis.Message // Messages received by every hub node
	RelayChan  <-chan *redis.Message // Messages specific to this hub only
	ClientTable map[string]*Client
}

// Find the channel the user maps to
func (server *EdgeServer) FindClient(client *Client) (string, error) {
	return server.Redis.Get(client.Username).Result()
}

func (server *EdgeServer) RegisterClient(client *Client) {
	// Register user in redis to the current server's Id
	err := server.Redis.Set(client.Username, server.Id, 0).Err()
	server.ClientTable[client.Username] = client

	if err != nil {
		panic(err)
	}
}

func (server *EdgeServer) UnregisterClient(client *Client) {
	server.Redis.Del(client.Username)
	delete(server.ClientTable, client.Username)
}

func (server *EdgeServer) IsClientPresent(client *Client) bool {
	_, err := server.Redis.Get(client.Username).Result()
	return err != nil
}