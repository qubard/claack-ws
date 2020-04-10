package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/websocket/socket"
)

type Application struct {
	cfg *Config
	Hub *socket.Hub
	db  *postgres.Database
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

func (app *Application) AddRaceText(text string) error {
	_, err := app.db.Handle().Exec(`INSERT INTO races(text) VALUES($1)`, text)
	return err
}

func (app *Application) InitDB(connStr string, schemaFile string) error {
	db, err := postgres.ConnectDB(connStr)

	if err != nil {
		log.Fatal("Could not connect to database!")
		return err
	}

	app.db = db
	log.Println("Connected to database.")

	// Try to create tables if they don't already exist
	err = app.db.CreateTables(schemaFile)

	if err != nil {
		log.Fatal("Could not create tables!")
		return err
	}

	log.Println("Created tables.")

	return nil
}

func (app *Application) GetDB() *postgres.Database {
	return app.db
}

func (app *Application) SetUpgrader(upgrader *websocket.Upgrader) {
	app.Hub.Upgrader = upgrader
}

func (app *Application) StartHub() {
	go app.Hub.Run()
}

func (app *Application) HostEndpoint(endpoint string) {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWebsocket(app.Hub, w, r, app.cfg.bufferSize)
	})

	log.Println("Started running server on", endpoint, app.cfg.ip, app.cfg.port)
	err := http.ListenAndServe(":"+app.cfg.port, nil)

	if err != nil {
		panic(err)
	}
}

func CreateApp() *Application {
	return &Application{
		Hub: socket.CreateHub(),
	}
}
