package routes

import (
    "net/http"
    "encoding/json"
 
    "github.com/qubard/claack-go/lib/postgres/queries"
    "github.com/qubard/claack-go/lib/postgres"
    "github.com/qubard/claack-go/lib/util"
)

// Not sure where we should put type definitions
type RouteHandler struct {
    Db *postgres.Database
    Secret string
}

// Interface might be cleaner here but then
// *RouteHandler has to be passed as an argument
type AuthHandler func(http.ResponseWriter, *http.Request, string)

func (handler *RouteHandler) AuthMiddleware(next AuthHandler, signKey []byte) func(http.ResponseWriter, *http.Request) {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify the session token
        token := r.FormValue("Token")
        
        // Session token has an id field specifically,
        if username, ok := util.ExtractField(token, "username", signKey); ok {
            // Check that the session is the last recorded session in the database.
            // TODO: Cache the last session token instead
            lastToken, err := queries.FindSessionToken(handler.Db, username.(string))
            if err == nil && lastToken == token {
                next(w, r, username.(string))
            } else {
                json.NewEncoder(w).Encode(ExpiredAuth)
            }
        } else {
            json.NewEncoder(w).Encode(InvalidAuth)
        }
    })
}

