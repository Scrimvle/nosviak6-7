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
	"Nosviak4/source/swash/packages"
	"fmt"
	"strconv"
	"strings"
)

var UsersAddMins = reg.Users.NewCommand(&commands.Command{
	Aliases:     []string{"add_minutes", "add_mins"},
	Description: "add minutes to a users expiry",
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

		item, err := strconv.Atoi(ac.Header)
		if err != nil {
			return s.ExecuteBranding(make(map[string]any), "invalid_int.tfx")
		}

		user.Expiry += int64(item * 60)
		if err := database.DB.EditUser(user, s.User, commands.Conn.SendWebhook); err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error while trying to change %s's expiry: %v", user.Username, err)
			return s.ExecuteBranding(map[string]any{"target": user.User(), "err": err.Error()}, "commands", "users", "error_occurred.tfx")
		}

		sessions.WriteToSession(sessions.IndexSessions(user.Username), func(index *sessions.Session) {
			index.User = user
			index.Reader.PostAlert(&terminal.Alert{
				AlertMessage: index.ExecuteBrandingToStringNoErr(map[string]any{"promoter": s.User.User(), "field": strings.ReplaceAll(strings.ReplaceAll(ac.Command.Aliases[0], "add_", ""), "set_", ""), "value": packages.Until(int(user.Created + user.Expiry))}, "commands", "users", "alerts", "tunable_changed.tfx"),
				AlertCode:    terminal.MESSAGE,
			})
		})

		reg.UsersLog.WriteLog("%s has changed %s's expiry to expiry in %s", s.User.Username, user.Username, packages.Until(int(user.Expiry)))
		return s.ExecuteBranding(map[string]any{"target": user.User(), "field": strings.ReplaceAll(strings.ReplaceAll(ac.Command.Aliases[0], "add_", ""), "set_", ""), "value": packages.Until(int(user.Created + user.Expiry))}, "commands", "users", "alerts", "performer_tunable_set.tfx")
	},

	Args: []*commands.Arg{{
		Name:        "user",
		Type:        commands.STRING,
		OpenEnded:   false,
		Description: "user to be edited",
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
