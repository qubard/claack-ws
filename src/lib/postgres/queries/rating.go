package queries

import (
    "github.com/qubard/claack-go/lib/postgres"
)

type RatingRow struct {
    Id int
    Elo int
}

func CreateRating(database *postgres.Database, id int, elo int) error {
    _, err := database.Handle().Exec(
        `INSERT INTO rating VALUES($1, $2)`, id, elo)
    return err
}