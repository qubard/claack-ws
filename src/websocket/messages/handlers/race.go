package handlers

import (
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/websocket/socket"
)

func QueueRace(db *postgres.Database, client *socket.Client, msg interface{}) {
	// How are we going to do this?
	// Need to register users to races and group them together
	// So their messages go to each other.
	client.Hub.EdgeServer.Redis.Publish("enq", client.Credentials.Username)
	//	client.Hub.EdgeServer.Redis.Publish("enq", "hello2")
}
