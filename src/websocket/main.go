package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/websocket/messages/handlers"
	"github.com/qubard/claack-go/websocket/messages/types"
)

var allowedHosts = map[string]bool{
	"http://localhost":  true,
	"http://claack.com": true,
}

func filterOrigin(r *http.Request) bool {
	host := r.Header.Get("Origin")
	_, allowed := allowedHosts[host]
	return allowed
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     filterOrigin,
}

func main() {
	var ip, port, addr string
	var bufferSize int
	flag.StringVar(&ip, "ip", "localhost", "The ip address to bind to")
	flag.StringVar(&port, "port", "4001", "The port to bind to")
	flag.IntVar(&bufferSize, "bufferSize", 1024*10, "The size of the client send buffer")
	flag.StringVar(&addr, "redis", "", "The address (ip:port) of a redis instance")
	flag.Parse()
	// Create and use our websocket app
	claack := CreateApp()
	err := claack.InitDB("user=postgres dbname=claack sslmode=disable", "claack.schema")

	if err != nil {
		log.Fatal(err)
		panic("Could not initialize database!")
	} else {
		log.Println("Successfully initialized database!")
	}

	if err := claack.InitRedis(addr, ""); err != nil {
		panic("Could not initialize redis connection!")
	}

	claack.CreateHub("hub1")
	log.Println("Successfully connected hub with edge server")

	claack.SetUpgrader(&upgrader)

	claack.Hub.Bus.RegisterHandler(types.QueueRace, handlers.QueueRace)
	claack.Hub.Bus.RegisterHandler(types.AuthUser, handlers.AuthUser)
	claack.Hub.Bus.RegisterHandler(types.AddMessage, handlers.AddMessage)

	claack.StartHub()
	claack.HostEndpoint("/ws", ip, port, bufferSize)

	defer claack.GetDB().Handle().Close()
}
