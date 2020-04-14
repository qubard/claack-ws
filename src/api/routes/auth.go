package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/qubard/claack-go/lib/postgres/queries"
	"github.com/qubard/claack-go/lib/util"
)

// TODO: Throtle these endpoints
// Or use cloudflare
func (handler *RouteHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		dec := json.NewDecoder(r.Body)
		dst := LoginDetails{}
		err := dec.Decode(&dst)

		if err == nil {
			ok, err := recaptcha.Confirm("", dst.ReCaptcha)

			if !ok || err != nil {
				json.NewEncoder(w).Encode(FailedCaptcha)
			} else {
				username := strings.ToLower(util.RemoveWhitespace(dst.Username))
				password := strings.ToLower(util.RemoveWhitespace(dst.Password))

				if util.ValidFormString(username) && util.ValidFormString(password) {
					res, err := queries.FindAuthUser(handler.Db, username)
					if err == nil && util.VerifyPassword(res.Password, password) {
						token, _ := queries.GenerateSession(handler.Db, username, []byte(handler.Secret))

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
		} else {
			json.NewEncoder(w).Encode(MalformedInput)
		}
	}
}

func (handler *RouteHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		dec := json.NewDecoder(r.Body)
		dst := LoginDetails{}
		err := dec.Decode(&dst)

		if err == nil {
			username := strings.ToLower(util.RemoveWhitespace(dst.Username))
			password := strings.ToLower(util.RemoveWhitespace(dst.Password))

			if util.ValidFormString(username) && util.ValidFormString(password) {
				// Check that the user does not exist
				_, err := queries.FindAuthUser(handler.Db, username)

				if err == sql.ErrNoRows {
					// Register the user
					err := queries.InsertAuthUser(handler.Db, username, util.HashPassword(password))
					if err != nil {
						panic(err)
					}

					// Create a profile and rating row for the new user

					// This shouldn't be able to error, but how else can we
					// get the id of the inserted user?
					authRow, _ := queries.FindAuthUser(handler.Db, username)
					queries.CreateProfile(handler.Db, authRow.Id, "", "", "", "") // description, hash, location default empty
					queries.CreateRating(handler.Db, authRow.Id, 1200)            // Default rating is 1200
					queries.CreateRep(handler.Db, authRow.Id, 0)

					token, _ := queries.GenerateSession(handler.Db, username, []byte(handler.Secret))
					queries.InsertSession(handler.Db, authRow.Id, token)

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
		} else {
			json.NewEncoder(w).Encode(MalformedInput)
		}
	}
}
