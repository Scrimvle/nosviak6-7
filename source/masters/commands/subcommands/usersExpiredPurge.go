package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

var UsersExpiredPurge = reg.Users.NewCommand(&commands.Command{
	Aliases: []string{"expired_purge"},
	Description: "remove all expired users",
	Permissions: []string{interactions.ADMIN},
	CommandFunc: func(ac *commands.ArgContext, session *sessions.Session) error {
		users, err := database.DB.GetExpiredUsers()
		if err != nil {
			return session.ExecuteBranding(make(map[string]any), "commands", "users", "error_occurred.tfx")
		}

		usersRemoved := make([]string, 0)

		/* ranges through all the users which are expired. */
		for _, user := range users {
			err := database.DB.DeleteUser(user)
			if err != nil {
				source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error when trying to delete %s: %v", user.Username, err)
				return session.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "users", "error_occurred.tfx")
			}

			usersRemoved = append(usersRemoved, user.Username)
		}

		reg.UsersLog.WriteLog("%s purged all expired users", session.User.Username)
		return session.ExecuteBranding(map[string]any{"expired_deleted": len(usersRemoved)}, "commands", "users", "alerts", "performer_users_expired_purged.tfx")
	},
})