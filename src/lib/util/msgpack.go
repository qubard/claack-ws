package util

import (
	"bufio"
	"bytes"

	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack"
)

type RelayMessage struct {
	DstToken string // The token of the user we're sending the message to
	Message string // The packed message
}

// Unfortunately we cannot define new methods on non-local types
func ReadPackedMessage(c *websocket.Conn) (interface{}, error) {
	_, message, _ := c.ReadMessage()
	var b bytes.Buffer
	b.Write(message)
	dec := msgpack.NewDecoder(bufio.NewReader(&b))
	return dec.DecodeInterface()
}

// Read a (packed) relay message
func ReadRelayMessage(msg string) (*RelayMessage, error) {
	var relay RelayMessage
	err := msgpack.Unmarshal([]byte(msg), &relay)
	return &relay, err
}

func WritePackedMessage(msg interface{}) ([]byte, error) {
	return msgpack.Marshal(msg)
}
