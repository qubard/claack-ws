package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/qubard/claack-go/api/routes"
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/rs/cors"
)

func main() {
	// Connect to the database
	// TODO: flag this
	db, err := postgres.ConnectDB("user=postgres dbname=claack sslmode=disable")

	if err != nil {
		log.Println(err)
		panic(err)
	}

	authHandler := routes.AuthHandler{
		Db:     db,
		Secret: "key",
	}

	var recaptchaSecret string
	flag.StringVar(&recaptchaSecret, "recaptchaSecret", "", "Your recaptcha secret key")
	flag.Parse()

	recaptcha.Init(recaptchaSecret)

	log.Println("Connected to database")

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/login/", authHandler.Login)
	mux.HandleFunc("/auth/register/", authHandler.Register)

	corsOpts := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := corsOpts.Handler(mux)

	err = http.ListenAndServe(":8000", handler)

	if err != nil {
		panic(err)
	}

	defer db.Handle().Close()
}
