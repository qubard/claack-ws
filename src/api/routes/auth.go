package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/lib/util"
)

type AuthHandler struct {
	Db     *postgres.Database
	Secret string
}

// TODO: Throtle these endpoints
// Or use cloudflare
func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		dec := json.NewDecoder(r.Body)
		dst := LoginDetails{}
		err := dec.Decode(&dst)

		if err == nil {

			ok, err := recaptcha.Confirm(r.Header.Get("User-Agent"), dst.ReCaptcha)

			if !ok || err != nil {
				json.NewEncoder(w).Encode(FailedCaptcha)
			} else {
				username := strings.ToLower(util.RemoveWhitespace(dst.Username))
				password := strings.ToLower(util.RemoveWhitespace(dst.Password))

				if util.ValidFormString(username) && util.ValidFormString(password) {
					res, err := handler.Db.FindAuthUser(username)
					if err == nil && util.VerifyPassword(res.Password, password) {
						token, _ := handler.Db.GenerateSession(username, []byte(handler.Secret))

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

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		dec := json.NewDecoder(r.Body)
		dst := LoginDetails{}
		err := dec.Decode(&dst)

		if err == nil {
			username := strings.ToLower(util.RemoveWhitespace(dst.Username))
			password := strings.ToLower(util.RemoveWhitespace(dst.Password))

			if util.ValidFormString(username) && util.ValidFormString(password) {
				// Check that the user does not exist
				_, err := handler.Db.FindAuthUser(username)

				if err == sql.ErrNoRows {
					// Register the user
					err := handler.Db.InsertAuthUser(username, util.HashPassword(password))
					if err != nil {
						panic(err)
					}

					token, _ := handler.Db.GenerateSession(username, []byte(handler.Secret))

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
