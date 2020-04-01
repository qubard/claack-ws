package handlers

import (
    "github.com/qubard/claack-go/messages/types"
)

type Message struct {
    id types.MessageType
    data interface{}
}

// Concurrent read using the HandlerBus is OK
// The intended usage of handlers is as READ ONLY
type HandlerBus struct {
    handlers map[types.MessageType]func(interface{})
}

func CreateHandlerBus() *HandlerBus {
    return &HandlerBus{ 
        handlers: make(map[types.MessageType]func(interface{})),
    }
}

// Register a handler to the corresponding message id
func (handlerBus *HandlerBus) RegisterHandler(id types.MessageType, handler func(interface{})) {
    handlerBus.handlers[id] = handler
}

func (handlerBus *HandlerBus) GetHandler(id types.MessageType) func(interface{}) {
    return handlerBus.handlers[id]
}
