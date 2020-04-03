package main

import (
    "log"
    
    "github.com/qubard/claack-go/lib/postgres"
)

func main() {
    // Connect to the database
    db, err := postgres.ConnectDB("user=postgres dbname=claack sslmode=disable")

    if err != nil {
        log.Fatal(err)
        panic("Could not initialize database!")
    } else {
        log.Println("Successfully connected to database!")
    }
    
    defer db.Handle().Close()
}
