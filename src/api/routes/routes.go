package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/lib/util"
)

type LoginDetails struct {
	Username string
	Password string
}

type AuthResponse struct {
	Error   string
	Message string
	Token   string
}

type AuthHandler struct {
	Db     *postgres.Database
	Secret string
}

var InvalidLogin AuthResponse = AuthResponse{
	Message: "Invalid login details",
}

var UserExists AuthResponse = AuthResponse{
	Message: "User already exists",
}

var InvalidRegister AuthResponse = AuthResponse{
	Message: "Invalid registration",
}

// Generate a new token for the user
// And update the last session from the database
func GenToken(username string, key []byte, db *postgres.Database) (string, error) {
	token, _ := util.MakeSessionToken(username, key)
	err := db.UpdateSession(username, token)
	return token, err
}

// TODO: Throtle these endpoints
// Or use cloudflare
func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := strings.ToLower(util.RemoveWhitespace(r.PostFormValue("username")))
		password := strings.ToLower(util.RemoveWhitespace(r.PostFormValue("password")))
		if util.ValidFormString(username) && util.ValidFormString(password) {
			res, err := handler.Db.FindAuthUser(username)
			if err == nil && util.VerifyPassword(res.Password, password) {
				token, sessionErr := GenToken(username, []byte(handler.Secret), handler.Db)

				if sessionErr != nil {
					panic(sessionErr)
				}

				response := AuthResponse{
					Message: "Valid password",
					Token:   token,
				}
				json.NewEncoder(w).Encode(response)
			} else {
				json.NewEncoder(w).Encode(InvalidLogin)
			}
		} else {
			json.NewEncoder(w).Encode(InvalidLogin)
		}
	}
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := strings.ToLower(util.RemoveWhitespace(r.PostFormValue("username")))
		password := strings.ToLower(util.RemoveWhitespace(r.PostFormValue("password")))
		if util.ValidFormString(username) && util.ValidFormString(password) {
			// Check that the user does not exist
			_, err := handler.Db.FindAuthUser(username)

			if err == sql.ErrNoRows {
				// Register the user
				err := handler.Db.InsertAuthUser(username, util.HashPassword(password))
				if err != nil {
					panic(err)
				}

				// TODO: change key
				token, sessionErr := GenToken(username, []byte(handler.Secret), handler.Db)

				if sessionErr != nil {
					panic(sessionErr)
				}

				response := AuthResponse{
					Message: "Successfully registered",
					Token:   token,
				}

				json.NewEncoder(w).Encode(response)
			} else {
				json.NewEncoder(w).Encode(UserExists)
			}
		} else {
			json.NewEncoder(w).Encode(InvalidRegister)
		}
	}
}
