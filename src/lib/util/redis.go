package util

import (
	"github.com/go-redis/redis/v7"
)

func CreateSubChannel(redis *redis.Client, name string) <-chan *redis.Message {
	ch := redis.Subscribe(name)
	_, err := ch.Receive()

	if err != nil {
		return nil
	}

	return ch.Channel()
}
