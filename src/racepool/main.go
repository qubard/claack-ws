package main

import (
	"flag"
	"github.com/go-redis/redis/v7"
	"github.com/qubard/claack-go/lib/microservice"
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

	if err == nil {
		pool := microservice.CreateRacePool(redis, "racepool", "enq", "deq")
		pool.Run()
	} else {
		panic(err)
	}
}
