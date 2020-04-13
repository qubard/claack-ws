package socket

import (
	"log"
	"net/http"
)

// serveWebsocket handles websocket requests from the peer.
func ServeWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request, bufferSize int) {
	conn, err := hub.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		Send: make(chan []byte, bufferSize),
	}

	// Register the client to the hub
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
