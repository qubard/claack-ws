package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/lib/util"
	"github.com/qubard/claack-go/websocket/messages/handlers"
	"github.com/qubard/claack-go/websocket/messages/types"
	"github.com/qubard/claack-go/websocket/socket"
)

type Application struct {
	Hub   *socket.Hub
	db    *postgres.Database
	redis *redis.Client
}

func (app *Application) InitRedis(addr string, password string) error {
	app.redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := app.redis.Ping().Result()
	return err
}

func (app *Application) RegisterHandlers() {
	app.Hub.Bus.RegisterHandler(types.QueueRace, socket.AuthMiddleware(handlers.QueueRace))
	app.Hub.Bus.RegisterHandler(types.AuthUser, handlers.AuthUser)
	app.Hub.Bus.RegisterHandler(types.AddMessage, handlers.AddMessage)
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

func (app *Application) HostEndpoint(endpoint string, ip string, port string, bufferSize int) {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWebsocket(app.Hub, w, r, bufferSize)
	})

	log.Println("Started running server on", endpoint, ip, port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic(err)
	}
}

func (app *Application) createEdgeServer(name string) *socket.EdgeServer {
	globalChan := util.CreateSubChannel(app.redis, "global")
	relayChan := util.CreateSubChannel(app.redis, name)

	if globalChan == nil || relayChan == nil {
		panic("Failed to create hub: invalid redis global channel or hub channel")
	}

	return &socket.EdgeServer{
		Redis:       app.redis,
		Id:          name,
		GlobalChan:  globalChan,
		RelayChan:   relayChan,
		ClientTable: make(map[string]*socket.Client),
		Mutex:       sync.Mutex{},
	}
}

func (app *Application) CreateHub(edgeName string) {
	// The hub needs some way of accessing the database
	// So we just pass it in like this (dependency injection) so its handlers
	// can use the database if they need to.
	app.Hub = socket.CreateHub(app.db, app.createEdgeServer(edgeName))
}

func CreateApp() *Application {
	return &Application{}
}
