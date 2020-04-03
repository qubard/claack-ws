package handlers

import (
	"github.com/qubard/claack-go/websocket/messages/types"
)

type Message struct {
	id   types.MessageType
	data interface{}
}

type MessageHandler func(interface{})

// Concurrent read using the HandlerBus is OK
// The intended usage of handlers is as READ ONLY
type HandlerBus struct {
	handlers map[types.MessageType]MessageHandler
}

func CreateHandlerBus() *HandlerBus {
	return &HandlerBus{
		handlers: make(map[types.MessageType]MessageHandler),
	}
}

// Register a handler to the corresponding message id
func (bus *HandlerBus) RegisterHandler(id types.MessageType, handler MessageHandler) {
	bus.handlers[id] = handler
}

func (bus *HandlerBus) GetHandler(id types.MessageType) func(interface{}) {
	return bus.handlers[id]
}
