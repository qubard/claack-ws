package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Application struct {
	cfg *Config
	hub *Hub
}

type Config struct {
	ip         string
	port       string
	bufferSize int
}

func (app *Application) ParseFlags() {
	var cfg = Config{}
	flag.StringVar(&cfg.ip, "ip", "localhost", "The ip address to bind to")
	flag.StringVar(&cfg.port, "port", "4001", "The port to bind to")
	flag.IntVar(&cfg.bufferSize, "bufferSize", 1024, "The size of the client send buffer")
	flag.Parse()
	app.cfg = &cfg
}

func (app *Application) SetUpgrader(upgrader *websocket.Upgrader) {
	app.hub.upgrader = upgrader
}

func (app *Application) Start() {
	go app.hub.run()
}

func (app *Application) StartHTTP(endpoint string) {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(app.hub, w, r, app.cfg.bufferSize)
	})

	log.Println("Started running server on", endpoint, app.cfg.ip, app.cfg.port)
	err := http.ListenAndServe(":"+app.cfg.port, nil)

	if err != nil {
		panic(err)
	}
}

func CreateApp() *Application {
	return &Application{
		hub: CreateHub(),
	}
}
