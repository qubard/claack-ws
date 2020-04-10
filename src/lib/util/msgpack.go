package util

import (
	"bytes"
    "bufio"
    
    "github.com/vmihailenco/msgpack"
	"github.com/gorilla/websocket"
)

// Unfortunately we cannot define new methods on non-local types
func ReadPackedMessage(c *websocket.Conn) (interface{}, error) {
    _, message, _ := c.ReadMessage()
    var b bytes.Buffer
    b.Write(message)
    dec := msgpack.NewDecoder(bufio.NewReader(&b))
    return dec.DecodeInterface()
}

func WritePackedMessage(msg interface{}) ([]byte, error) {
    return msgpack.Marshal(msg)
}
