package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qubard/claack-go/lib/postgres/queries"
)

func (handler *RouteHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if username, ok := vars["username"]; ok {
		if profile, err := queries.FindFullProfile(handler.Db, username); err == nil {
			json.NewEncoder(w).Encode(profile)
		} else {
			json.NewEncoder(w).Encode(ProfileNotFound)
		}
	} else {
		json.NewEncoder(w).Encode(MalformedInput)
	}
}

// Using the authed username we can update that profile for that user
func (handler *RouteHandler) UpdateProfile(w http.ResponseWriter, r *http.Request, username string) {
	description := r.FormValue("Description")
	alias := r.FormValue("Alias")
	err := queries.UpdateProfileSimple(handler.Db, username, description, alias)
	if err == nil {
		// Return the new profile
		if profile, err := queries.FindFullProfile(handler.Db, username); err == nil {
			json.NewEncoder(w).Encode(profile)
		} else {
			json.NewEncoder(w).Encode(ProfileNotFound)
		}
		// Technically the err != nil case is impossible..
	} else {
		json.NewEncoder(w).Encode(InvalidProfileUpdate)
	}
}
