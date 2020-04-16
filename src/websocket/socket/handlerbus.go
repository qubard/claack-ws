package socket

import (
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/websocket/messages/types"
)

type Message struct {
	id   types.MessageType
	data interface{}
}

type MessageHandler func(*postgres.Database, *Client, interface{})

// Concurrent read using the HandlerBus is OK
// The intended usage of handlers is as READ ONLY
type HandlerBus struct {
	handlers map[types.MessageType]MessageHandler
	db       *postgres.Database
}

func CreateHandlerBus(database *postgres.Database) *HandlerBus {
	return &HandlerBus{
		handlers: make(map[types.MessageType]MessageHandler),
		db:       database,
	}
}

// Register a handler to the corresponding message id
func (bus *HandlerBus) RegisterHandler(id types.MessageType, handler MessageHandler) {
	bus.handlers[id] = handler
}

func (bus *HandlerBus) GetHandler(id types.MessageType) MessageHandler {
	return bus.handlers[id]
}

func (bus *HandlerBus) InvokeHandler(client *Client, id types.MessageType, payload interface{}) {
	handler := bus.GetHandler(id)
	if handler != nil {
		handler(bus.db, client, payload)
	}
}

func (bus *HandlerBus) AttemptInvokeHandler(client *Client, msg interface{}) {
	if msgMap, ok := msg.(map[string]interface{}); ok {
		// NOTE: The type field is lowercase due to Redux
		// which doesn't let it be uppercase
		if msgType, ok := msgMap["type"].(types.MessageType); ok {
			if payload, ok := msgMap["Payload"]; ok {
				bus.InvokeHandler(client, msgType, payload)
			}
		}
	}
}
