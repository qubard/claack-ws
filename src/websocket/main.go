package main

import (
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

type PayloadP struct {
	Disconnect bool
}

type AuthUser struct {
	Id      int8
	Payload PayloadP
}

type ProfileDetails struct {
	Avatar   string
	Username string
}

type Profile struct {
	Profile ProfileDetails
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

	claack.CreateHub()

	claack.SetUpgrader(&upgrader)
	claack.ParseFlags()

	claack.Hub.Bus.RegisterHandler(types.QueueRace, handlers.QueueRace)
	claack.Hub.Bus.RegisterHandler(types.AuthUser, handlers.AuthUser)

	claack.StartHub()
	claack.HostEndpoint("/ws")

	defer claack.GetDB().Handle().Close()
}
