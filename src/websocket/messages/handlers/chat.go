package handlers

import (
	"github.com/mitchellh/mapstructure"
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/websocket/messages/types"
	"github.com/qubard/claack-go/websocket/socket"
)

type ChatMessagePayload struct {
	From      string
	Id        int
	Location  string
	Text      string
	Timestamp int
	To        string
}

type MessageDeliveredPayload struct {
	To string
	Id int // The message identifier that is seen
}

type MessageDelivered struct {
	Type    types.MessageType
	Payload MessageDeliveredPayload
}

type ChatMessage struct {
	Type    types.MessageType
	Payload ChatMessagePayload
}

func AddMessage(db *postgres.Database, client *socket.Client, msg interface{}) {
	var payload ChatMessagePayload
	if err := mapstructure.Decode(msg, &payload); err == nil {
		// We don't use Payload.From since it's controlled by the client
		// instead use their auth credentials
		payload.From = client.Credentials.Username

		chatMsg := ChatMessage{
			Type:    types.AddMessage,
			Payload: payload,
		}

		// Send the actual message
		err = client.Hub.EdgeServer.RelayMessage(payload.To, chatMsg)

		// Let the client know their message has been delivered
		// Attach an Id so they know which message it was
		if err == nil {
			client.SendMessage(MessageDelivered{
				Type: types.MessageDelivered,
				Payload: MessageDeliveredPayload{
					To: payload.To,
					Id: payload.Id,
				},
			})
		}
		// TODO: alert the user their message could not be delivered
	}
}
