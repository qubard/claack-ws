package queries

import (
    "time"
    
    "github.com/qubard/claack-go/lib/postgres"
)


type AuthRow struct {
    Id int
    Username string
    Password string
    CreatedAt time.Time
}

func FindAuthUser(database *postgres.Database, username string) (*AuthRow, error) {
    row := database.Handle().QueryRow(`SELECT * from auth WHERE username=$1`, username)
    var res AuthRow
    err := row.Scan(&res.Id, &res.Username, &res.Password, &res.CreatedAt)
    return &res, err
}

func InsertAuthUser(database *postgres.Database, username string, password string) error {
    _, err := database.Handle().Exec(`INSERT INTO auth (username, password) VALUES($1, $2)`, username, password)
    return err
}
