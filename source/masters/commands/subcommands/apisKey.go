package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	cmd "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// APIsKey will generate the user provided a key for api access
var APIsKey = cmd.APIs.NewCommand(&commands.Command{
	Aliases:     []string{"key"},
	Description: "generate a key for user",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.RESELLER},
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		user, err := database.DB.GetUserAsParentalFigure(fmt.Sprint(ac.Args[0].Values[0].Literal), s.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != s.User.ID {
			user, err := database.DB.GetUser(fmt.Sprint(ac.Args[0].Values[0].Literal))
			if err != nil || user == nil {
				return s.ExecuteBranding(map[string]any{"username": fmt.Sprint(ac.Args[0].Values[0].Literal)}, "commands", "users", "bad_user.tfx")
			}

			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "apis", "access_denied.tfx")
		}

		if !user.API {
			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "apis", "missing_api.tfx")
		}

		key := hex.EncodeToString(*database.NewSalt(8))
		user.APIKey = sha256.New().Sum([]byte(key))
		if err := database.DB.EditUser(user, s.User, commands.Conn.SendWebhook); err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error while trying to make %s's api key: %v", user.Username, err)
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "apis", "error_occurred.tfx")
		}

		return s.ExecuteBranding(map[string]any{"target": user.User(), "key": key}, "commands", "apis", "performer_apikey.tfx")
	},

	Args: []*commands.Arg{
		{
			Name:        "user",
			Type:        commands.STRING,
			Description: "user who will be modified",
			Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
				users, err := database.DB.GetUsersAsParent(s.User)
				if err != nil {
					return make([]string, 0)
				}

				appends := make([]string, 0)
				for _, user := range users {
					if !user.API {
						continue
					}

					appends = append(appends, user.Username)
				}

				return appends
			},
		},
	},
})
