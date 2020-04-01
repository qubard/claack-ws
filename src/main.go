package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/messages/types"
)

var allowedHosts = map[string]bool{
	"http://localhost:3000": true,
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
	// Create and use our websocket app
	claack := CreateApp()
	claack.SetUpgrader(&upgrader)
	claack.ParseFlags()

	claack.hub.handlerBus.RegisterHandler(types.InitSocket, func(msg interface{}) {
		log.Println("Init Socket message", msg)
	})

	claack.Start()
	claack.StartHTTP("/ws")
}
