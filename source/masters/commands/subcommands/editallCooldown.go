package subcommands

import (
	"Nosviak4/source/database"
	cmd "Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/masters/terminal/interactions"
	"strconv"
	"strings"
)

// EditAllCooldown will forcefully change the cooldown of all users inside the scope.
var EditAllCooldown = commands.EditAll.NewCommand(&cmd.Command{
	Aliases: []string{"cooldown"},
	Description: "edit everyone(s) cooldown",
	Permissions: []string{interactions.ADMIN},
	CommandFunc: func(ac *cmd.ArgContext, s *sessions.Session) error {
		digit, err := strconv.Atoi(ac.Header)
		if err != nil  {
			return s.ExecuteBranding(make(map[string]any), "commands", "editall", "invalid_int.tfx")
		}

		users, err := database.DB.GetUsersAsParent(s.User)
		if err != nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "editall", "error_occurred.tfx")
		}

		/* ranges through all the users inside the scope */
		for _, user := range users {
			user.Cooldown = digit

			/* EditUser will change the parameters */
			if err := database.DB.EditUser(user, s.User, cmd.Conn.SendWebhook); err != nil {
				return s.ExecuteBranding(make(map[string]any), "commands", "editall", "error_occurred.tfx")
			}

			/* our session won't get alerts. */
			if user.ID == s.User.ID {
				continue
			}

			/* WriteToSession will change each sessions parameter directly and tell them */
			sessions.WriteToSession(sessions.IndexSessions(user.Username), func(index *sessions.Session) {
				index.User = user

				/* posts the alert */
				index.Reader.PostAlert(&terminal.Alert{
					AlertMessage: index.ExecuteBrandingToStringNoErr(map[string]any{"promoter": s.User.User(), "field": strings.ReplaceAll(strings.ReplaceAll(ac.Command.Aliases[0], "add_", ""), "set_", ""), "value": user.Conns}, "commands", "editall", "alerts", "tunable_set.tfx"),
					AlertCode:    terminal.MESSAGE,
				})
			})
		}

		return s.ExecuteBranding(map[string]any{"field": strings.ReplaceAll(strings.ReplaceAll(ac.Command.Aliases[0], "add_", ""), "set_", ""), "value": digit}, "commands", "editall", "alerts", "performer_tunable_set.tfx")
	},
})