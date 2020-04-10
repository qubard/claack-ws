package queries

import (
    "time"
    "github.com/qubard/claack-go/lib/postgres"
)

type ProfileRow struct {
    Id int
    Username string
    Alias string
    Description string
    Hash string // Avatar hash
    Location string // Location of avatar
}

type FullProfileRow struct {
    Id int
    Username string
    Alias string
    Description string
    Hash string // Avatar hash
    Location string // Location of avatar
    Elo int
    Rep int
    CreatedAt time.Time
}

func FindProfile(database *postgres.Database, username string) (*ProfileRow, error) {
    row := database.Handle().QueryRow(`SELECT auth.id, auth.username, profile.alias, profile.description, profile.hash, profile.location FROM auth INNER JOIN profile ON auth.username=$1 and profile.id=auth.id`, username)
    var res ProfileRow
    err := row.Scan(&res.Id, &res.Username, &res.Description, &res.Hash, &res.Location)
    return &res, err
}

func UpdateProfile(database *postgres.Database, username string, description string, hash string, location string) error {
    _, err := database.Handle().Exec(
        `UPDATE profile
        SET 
            description=$2, hash=$3, location=$4 
        FROM auth 
        WHERE
            auth.username=$1`, username, description, hash, location)
    return err
}

func UpdateProfileSimple(database *postgres.Database, username string, description string, alias string) error {
    _, err := database.Handle().Exec(`UPDATE profile SET description=$2, alias=$3 FROM auth WHERE auth.username=$1`, username, description, alias)
    return err
}

func CreateProfile(database *postgres.Database, id int, alias string, description string, hash string, location string) error {
    _, err := database.Handle().Exec(
        `INSERT INTO profile VALUES($1, $2, $3, $4, $5)`, id, alias, description, hash, location)
    return err
}

func FindFullProfile(database *postgres.Database, username string) (*FullProfileRow, error) {
    row := database.Handle().QueryRow(
        `SELECT auth.id, username, alias, description, hash, location, elo, rep.value, created_at 
        FROM auth 
        INNER JOIN profile ON auth.username=$1 and profile.id=auth.id 
        INNER JOIN rating ON rating.id=auth.id 
        INNER JOIN rep ON rep.id=auth.id`, username)
    var res FullProfileRow
    err := row.Scan(&res.Id, &res.Username, &res.Alias, &res.Description, &res.Hash, &res.Location, &res.Elo, &res.Rep, &res.CreatedAt)
    return &res, err
}
