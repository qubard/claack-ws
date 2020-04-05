package postgres

import (
	"time"
)

type AuthRow struct {
	Id        int
	Username  string
	Password  string
	CreatedAt time.Time
}

func (database *Database) FindAuthUser(username string) (*AuthRow, error) {
	row := database.Handle().QueryRow(`SELECT * from auth WHERE username=$1`, username)
	var res AuthRow
	err := row.Scan(&res.Id, &res.Username, &res.Password, &res.CreatedAt)
	return &res, err
}

func (database *Database) InsertAuthUser(username string, password string) error {
	_, err := database.Handle().Exec(`INSERT INTO auth (username, password) VALUES($1, $2)`, username, password)
	return err
}

// Inserts or updates a session for a user
func (database *Database) UpdateSession(username string, token string) error {
	_, err := database.Handle().Exec(`INSERT INTO sessions VALUES((SELECT id from auth where username=$1), $2) ON CONFLICT (id) DO UPDATE SET token=$2`, username, token)
	return err
}
