package queries

import (
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/lib/util"
)

// Inserts or updates a session for a user
func UpdateSession(database *postgres.Database, username string, token string) error {
	_, err := database.Handle().Exec(`UPDATE session SET token=$2 FROM auth WHERE (auth.id=session.id and auth.username=$1)`, username, token)
	return err
}

func InsertSession(database *postgres.Database, id int, session string) error {
	_, err := database.Handle().Exec(
		`INSERT INTO session VALUES($1, $2)`, id, session)
	return err
}

func FindSessionToken(database *postgres.Database, username string) (string, error) {
	row := database.Handle().QueryRow(`SELECT session.token FROM session INNER JOIN auth ON (auth.id=session.id and auth.username=$1)`, username)
	var res string
	err := row.Scan(&res)
	return res, err
}

func GenerateSession(database *postgres.Database, username string, key []byte) (string, error) {
	token, _ := util.MakeSessionToken(username, key)
	err := UpdateSession(database, username, token)
	return token, err
}
