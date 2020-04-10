package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/lib/postgres/queries"
	"github.com/qubard/claack-go/lib/util"
	"github.com/qubard/claack-go/websocket/messages/types"
	"github.com/qubard/claack-go/websocket/socket"
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

type ProfileMessage struct {
	Type    int8
	Payload queries.FullProfileRow
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

	claack.SetUpgrader(&upgrader)
	claack.ParseFlags()

	claack.Hub.Bus.RegisterHandler(types.InitSocket, func(client *socket.Client, msg interface{}) {
		log.Println("Init Socket message", msg)
	})

	// TODO: global func for token extraction..
	claack.Hub.Bus.RegisterHandler(types.AuthUser, func(client *socket.Client, msg interface{}) {
		if token, ok := msg.(map[string]interface{})["token"]; ok && token != nil {
			if username, ok := util.ExtractField(token.(string), "username", []byte("key")); ok {
				// We have the user id (username), use it to find the user's profile
				// and send the necessary profile update back
				// Check that their last session is equal to this token
				lastSession, err := queries.FindSessionToken(claack.GetDB(), username.(string))

				if err == nil && lastSession == token {
					profile, err := queries.FindFullProfile(claack.GetDB(), username.(string))
					bytes, err := util.WritePackedMessage(ProfileMessage{
						Type:    types.ProfileMessage,
						Payload: *profile,
					})

					if err == nil {
						client.Send <- bytes
					}
				}
			}
		}
	})

	claack.StartHub()
	claack.HostEndpoint("/ws")

	defer claack.GetDB().Handle().Close()
}
