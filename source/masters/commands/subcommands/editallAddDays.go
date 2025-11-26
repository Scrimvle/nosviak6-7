package subcommands

import (
	"Nosviak4/source/database"
	cmd "Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"strconv"
	"strings"
)

// EditAllAddDays will forcefully add onto the days-expiry of all the users inside the scope.
var EditAllAddDays = commands.EditAll.NewCommand(&cmd.Command{
	Aliases: []string{"add_days"},
	Description: "increase everyone(s) expiry",
	Permissions: []string{interactions.ADMIN},
	CommandFunc: func(ac *cmd.ArgContext, s *sessions.Session) error {
		digit, err := strconv.Atoi(ac.Header)
		if err != nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "editall", "invalid_int.tfx")
		}

		users, err := database.DB.GetUsersAsParent(s.User)
		if err != nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "editall", "error_occurred.tfx")
		}

		/* ranges through all the users inside the scope */
		for _, user := range users {
			user.Expiry += int64(digit * 86400)

			/* EditUser will change the parameters */
			if err := database.DB.EditUser(user, s.User, cmd.Conn.SendWebhook); err != nil {
				return s.ExecuteBranding(make(map[string]any), "commands", "editall", "error_occurred.tfx")
			}

			/* our session won't get alerts. */
			if user.ID == s.User.ID {
				continue
			}
		}

		return s.ExecuteBranding(map[string]any{"field": strings.ReplaceAll(strings.ReplaceAll(ac.Command.Aliases[0], "add_", ""), "set_", ""), "value": digit}, "commands", "editall", "alerts", "performer_tunable_set.tfx")
	},
})