package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/websocket/messages/types"
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
	err := claack.InitDB("user=postgres dbname=claack sslmode=disable", "claack.schema")

	if err != nil {
		log.Fatal(err)
		panic("Could not initialize database!")
	} else {
		log.Println("Successfully initialized database!")
	}

	err = claack.AddRaceText("test")
	claack.SetUpgrader(&upgrader)
	claack.ParseFlags()

	claack.hub.bus.RegisterHandler(types.InitSocket, func(msg interface{}) {
		log.Println("Init Socket message", msg)
	})

	key := []byte("key")
	signed, err := SignUser("cub", key)

	log.Println(signed, err)

	token, err := ValidateUser(signed, key)

	log.Println(token.Valid, token.Claims)

	claack.StartHub()
	claack.HostEndpoint("/ws")

	defer claack.GetDB().Handle().Close()
}
