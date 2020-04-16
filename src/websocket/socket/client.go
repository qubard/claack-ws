package socket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/lib/util"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 3 * 1024 // Max 3KiB
)

func (limiter *RateLimiter) DoLimit(client *Client) bool {
	currTimeNano := time.Now().UnixNano()
	if currTimeNano-client.Limiter.LastRecNano >= 1e6*limiter.WindowSizeMs {
		if client.Limiter.Count > client.Limiter.ThrottleLimit {
			return true
		}
		// Reset the 1 second window (1e6 ms)
		client.Limiter.Count = 0
		client.Limiter.LastRecNano = time.Now().UnixNano()
	} else {
		client.Limiter.Count = client.Limiter.Count + 1
	}
	return false
}

type RateLimiter struct {
	WindowSizeMs  int64  // Window size in milliseconds
	Count         uint64 // The number of packets since the last send
	LastRecNano   int64  // The time they last sent a packet in nanoseconds
	ThrottleLimit uint64 // The max value of `count` before throttling
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send     chan []byte
	Username string // Client's authed username
	Limiter  *RateLimiter
}

func (client *Client) SendMessage(msg interface{}) error {
	bytes, err := util.WritePackedMessage(msg)

	if err == nil {
		client.Send <- bytes
	}
	return err
}

// readPump pumps messages from the websocket connection to the hub.
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		msg, err := util.ReadPackedMessage(c.conn)

		if c.Limiter.DoLimit(c) {
			// If the client is rate limited, drop them immediately
			// And maybe ban them?
			return
		}

		if err == nil {
			// When we receive a message pass it to a registered handler
			c.Hub.Bus.AttemptInvokeHandler(c, msg)
		}

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
// Write to c.Send to send the client a message
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.conn.NextWriter(websocket.BinaryMessage)

			if err != nil {
				return
			}

			writer.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				writer.Write(<-c.Send)
			}

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
