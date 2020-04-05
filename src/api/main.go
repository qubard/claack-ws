package main

import (
	"log"
	"net/http"
	"time"

	"github.com/qubard/claack-go/api/routes"
	"github.com/qubard/claack-go/lib/postgres"
)

type RaceResult struct {
	Id        int
	Text      string
	CreatedAt time.Time
}

func main() {
	// Connect to the database
	db, err := postgres.ConnectDB("user=postgres dbname=claack sslmode=disable")

	if err != nil {
		log.Println(err)
		panic(err)
	}

	handler := routes.AuthHandler{
		Db:     db,
		Secret: "key",
	}

	log.Println("Connected to database")

	http.HandleFunc("/auth/login/", handler.Login)
	http.HandleFunc("/auth/register/", handler.Register)

	err = http.ListenAndServe(":8000", nil)

	if err != nil {
		panic(err)
	}

	defer db.Handle().Close()
}
