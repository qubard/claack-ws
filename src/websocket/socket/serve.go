package socket

import (
	"log"
	"net/http"
	"time"
)

// serveWebsocket handles websocket requests from the peer.
func ServeWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request, bufferSize int) {
	conn, err := hub.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		Hub:  hub,
		conn: conn,
		Send: make(chan []byte, bufferSize),
		Limiter: &RateLimiter{
			Count:         0,
			LastRecNano:   time.Now().UnixNano(),
			WindowSizeMs:  100,
			ThrottleLimit: 10,
		},
	}

	// Register the client to the hub
	client.Hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
