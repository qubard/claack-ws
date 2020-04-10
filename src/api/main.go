package main

import (
    "log"
    "net/http"
    "flag"
    
    "github.com/rs/cors"
    "github.com/dpapathanasiou/go-recaptcha"
    "github.com/qubard/claack-go/api/routes"
    "github.com/qubard/claack-go/lib/postgres"
    "github.com/gorilla/mux"
)

func main() {
    // Connect to the database
    // TODO: flag this
    db, err := postgres.ConnectDB("user=postgres dbname=claack sslmode=disable")

    if err != nil {
        log.Println(err)
        panic(err)
    }
    
    router := routes.RouteHandler{
        Db: db,
        Secret: "key",
    }
    
    var recaptchaSecret string
    flag.StringVar(&recaptchaSecret, "recaptchaSecret", "", "Your recaptcha secret key")
    flag.Parse()
    
    recaptcha.Init(recaptchaSecret)
    
    log.Println("Connected to database")
    
    muxer := mux.NewRouter()
    
    muxer.HandleFunc("/auth/login/", router.Login)
    muxer.HandleFunc("/auth/register/", router.Register)
    muxer.HandleFunc("/profile/{username}/", router.GetProfile)
    muxer.HandleFunc("/update/profile/", router.AuthMiddleware(router.UpdateProfile, []byte("key")))
    
    corsOpts := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowCredentials: true,
    })
    
    handler := corsOpts.Handler(muxer)
    
    err = http.ListenAndServe(":8000", handler)

    if err != nil {
        panic(err)
    }
    
    defer db.Handle().Close()
}
