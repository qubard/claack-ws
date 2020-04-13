package socket

import (
	"github.com/go-redis/redis/v7"
)

// TODO: IsUserPresent and handling duplicate connections

type EdgeServer struct {
	Redis	   *redis.Client
	Id 		   string
	GlobalChan <-chan *redis.Message // Messages received by every hub node
	MainChan   <-chan *redis.Message // Messages specific to this hub only
}

// Find the channel the user maps to
func (server *EdgeServer) FindUser(id string) (string, error) {
	return server.Redis.Get(id).Result()
}

func (server *EdgeServer) RegisterUser(id string) {
	// Register user in redis to server.Id
	err := server.Redis.Set(id, server.Id, 0).Err()

	if err != nil {
		panic(err)
	}
}

func (server *EdgeServer) UnregisterUser(id string) {
	server.Redis.Del(id)
}

func (server *EdgeServer) IsUserPresent(id string) bool {
	_, err := server.Redis.Get(id).Result()
	return err != nil
}