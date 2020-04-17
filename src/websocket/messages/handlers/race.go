package handlers

import (
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/websocket/socket"
)

func QueueRace(db *postgres.Database, client *socket.Client, msg interface{}) {
	// Simply publish an enqueue message to the racepool microservice
	// This message only requires the user id
	client.Hub.EdgeServer.Publish("enq", client.Credentials.Username)
}
