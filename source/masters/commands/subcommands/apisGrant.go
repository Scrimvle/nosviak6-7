package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	cmd "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"
)

// APIsGrant will define the command for apis grant <user>
var APIsGrant = cmd.APIs.NewCommand(&commands.Command{
	Aliases: []string{"grant"},
	Description: "grant a user api access",
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

		/* checks if the user has API access already */
		if user.API || user == nil {
			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "apis", "already_has_access.tfx")
		}

		user.API = true
		if err := database.DB.EditUser(user, s.User, commands.Conn.SendWebhook); err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error while trying to set %s's api status: %v", user.Username, err)
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "apis", "error_occurred.tfx")
		}

		sessions.WriteToSession(sessions.IndexSessions(user.Username), func(index *sessions.Session) {
			index.User = user
			index.Reader.PostAlert(&terminal.Alert{
				AlertMessage: index.ExecuteBrandingToStringNoErr(map[string]any{"promoter": s.User.User()}, "commands", "apis", "promoted_user.tfx"),
				AlertCode:    terminal.MESSAGE,
			})
		})

		return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "apis", "performer_promoted_user.tfx")
	},

	Args: []*commands.Arg{
		{
			Name: "user",
			Type: commands.STRING,
			Description: "user who will be modified",
			Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
				users, err := database.DB.GetUsersAsParent(s.User)
				if err != nil {
					return make([]string, 0)
				}

				buf := make([]string, 0)
				for _, user := range users {
					buf = append(buf, user.Username)
				}

				return buf
			},
		},
	},
})