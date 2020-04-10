package queries

import (
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/lib/util"
)

// Inserts or updates a session for a user
func UpdateSession(database *postgres.Database, username string, token string) error {
	_, err := database.Handle().Exec(`INSERT INTO session VALUES((SELECT id from auth where username=$1), $2) ON CONFLICT (id) DO UPDATE SET token=$2`, username, token)
	return err
}

func FindSessionToken(database *postgres.Database, username string) (string, error) {
	row := database.Handle().QueryRow(`SELECT session.token FROM session INNER JOIN auth ON auth.username=$1`, username)
	var res string
	err := row.Scan(&res)
	return res, err
}

func GenerateSession(database *postgres.Database, username string, key []byte) (string, error) {
	token, _ := util.MakeSessionToken(username, key)
	err := UpdateSession(database, username, token)
	return token, err
}
