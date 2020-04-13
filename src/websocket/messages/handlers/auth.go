package handlers

import (
	"github.com/qubard/claack-go/lib/postgres"
	"github.com/qubard/claack-go/lib/postgres/queries"
	"github.com/qubard/claack-go/lib/util"
	"github.com/qubard/claack-go/websocket/messages/types"
	"github.com/qubard/claack-go/websocket/socket"
)

type ProfileMessage struct {
	Type    int8
	Payload queries.FullProfileRow
}

func AuthUser(db *postgres.Database, client *socket.Client, msg interface{}) {
	if token, ok := msg.(map[string]interface{})["token"]; ok && token != nil {
		if username, ok := util.ExtractField(token.(string), "username", []byte("key")); ok {
			// We have the username, use it to find the user's profile
			// and send the necessary profile update back
			// Check that their last session is equal to this token
			lastSession, err := queries.FindSessionToken(db, username.(string))

			if err == nil && lastSession == token {
				profile, err := queries.FindFullProfile(db, username.(string))
				bytes, err := util.WritePackedMessage(ProfileMessage{
					Type:    types.ProfileUpdate,
					Payload: *profile,
				})

				// Send the user back their profile information
				if err == nil {
					client.Send <- bytes
				}
			}
		}
	}
}
