package handlers

import (
	"github.com/qubard/claack-go/websocket/messages/types"
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/websocket/socket"
	"github.com/mitchellh/mapstructure"
)

type ChatMessagePayload struct {
	From string
	Id int
	Location string
	Text string
	Timestamp int
	To string
}

type ChatMessage struct {
	Type types.MessageType
	Payload ChatMessagePayload
}

func AddMessage(db *postgres.Database, client *socket.Client, msg interface{}) {
	var payload ChatMessagePayload
	if err := mapstructure.Decode(msg, &payload); err == nil {
		chatMsg := ChatMessage {
			Type: types.AddMessage,
			Payload: payload,
		}
		client.Hub.EdgeServer.RelayMessage(payload.To, chatMsg)
	}
}
