package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"
)

var UsersChangePW = reg.Users.NewCommand(&commands.Command{
	Aliases:     []string{"changepw"},
	Description: "changes the users password",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.RESELLER},
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		user, err := database.DB.GetUserAsParentalFigure(fmt.Sprint(ac.Args[0].Values[0].Literal), s.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != s.User.ID {
			user, err := database.DB.GetUser(fmt.Sprint(ac.Args[0].Values[0].Literal))
			if err != nil || user == nil {
				return s.ExecuteBranding(map[string]any{"username": fmt.Sprint(ac.Args[0].Values[0].Literal)}, "commands", "users", "bad_user.tfx")
			}

			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "access_denied.tfx")
		}

		if len(ac.Header) < source.OPTIONS.Ints("minimum_password_length") || len(ac.Header) > source.OPTIONS.Ints("maximum_password_length") {
			return s.ExecuteBranding(make(map[string]any), "commands", "users", "invalid_password.tfx")
		}

		user.Password = database.NewHash([]byte(ac.Header), &user.Key)
		if err := database.DB.EditUser(user, s.User, commands.Conn.SendWebhook); err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error when trying to promote %s to %s: %v", user.Username, ac.Header, err)
			return s.ExecuteBranding(map[string]any{"target": user.User(), "err": err.Error()}, "commands", "users", "error_occurred.tfx")
		}

		sessions.WriteToSession(sessions.IndexSessions(user.Username), func(index *sessions.Session) {
			index.User = user
			index.Reader.PostAlert(&terminal.Alert{
				AlertMessage: index.ExecuteBrandingToStringNoErr(map[string]any{"promoter": s.User.User()}, "commands", "users", "alerts", "password_changed.tfx"),
				AlertCode:    terminal.MESSAGE,
			})
		})

		reg.UsersLog.WriteLog("%s has changed %s's password", s.User.Username, user.Username)
		return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "alerts", "performer_password_changed.tfx")
	},

	Args: []*commands.Arg{{
		Name:        "user",
		Type:        commands.STRING,
		Description: "user to be promoted",
		Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
			child, err := database.DB.GetUsersAsParent(s.User)
			if err != nil {
				return make([]string, 0)
			}

			children := make([]string, 0)
			for _, child := range child {
				children = append(children, child.Username)
			}

			return children
		},
	}},
})
