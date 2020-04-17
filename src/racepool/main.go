package main

import (
	"flag"
	"github.com/go-redis/redis/v7"
	"github.com/qubard/claack-go/lib/microservice"
	"github.com/qubard/claack-go/websocket/socket"
)

func main() {
	var addr string
	flag.StringVar(&addr, "redis", "", "The address (ip:port) of a redis instance")
	flag.Parse()

	redis := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := redis.Ping().Result()

	edgeServer := socket.CreateSimpleEdgeServer(redis)

	if err == nil {
		pool := microservice.CreateRacePool(redis, edgeServer, "racepool", "enq", "deq")
		pool.Run()
	} else {
		panic(err)
	}
}
