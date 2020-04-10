package handlers

import (
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/websocket/socket"
	"log"
)

func QueueRace(db *postgres.Database, client *socket.Client, msg interface{}) {
	log.Println("queue race", msg)
}
