package commands

import (
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"crypto/sha256"
	"encoding/hex"
)

// RegenerateAPIKey will allow the user to generate a brand new API Key
var RegenerateAPIKey = commands.ROOT.NewCommand(&commands.Command{
	Aliases: []string{"apikey"},
	Description: "generate a new API key",
	Permissions: make([]string, 0),
	CommandFunc: func(ac *commands.ArgContext, session *sessions.Session) error {
		if !session.User.API {
			return session.ExecuteBranding(make(map[string]any), "commands", "apiKey", "access_denied.tfx")
		}

		key := hex.EncodeToString(*database.NewSalt(8))
		session.User.APIKey = sha256.New().Sum([]byte(key))
		if err := database.DB.EditUser(session.User, session.User, commands.Conn.SendWebhook); err != nil {
			return session.ExecuteBranding(make(map[string]any), "commands", "apiKey", "error_occurred.tfx")
		}

		return session.ExecuteBranding(map[string]any{"key": key}, "commands", "apiKey", "generated.tfx")
	},
})