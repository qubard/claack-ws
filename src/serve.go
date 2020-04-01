package main

import (
	"log"
	"net/http"
)

// serveWebsocket handles websocket requests from the peer.
func serveWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request, bufferSize int) {
	conn, err := hub.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, bufferSize),
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
