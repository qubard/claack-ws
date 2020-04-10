package queries

import (
    "github.com/qubard/claack-go/lib/postgres"
)

type RepRow struct {
    Id int
    Reputation int
}

func CreateRep(database *postgres.Database, id int, rep int) error {
    _, err := database.Handle().Exec(
        `INSERT INTO rep VALUES($1, $2)`, id, rep)
    return err
}