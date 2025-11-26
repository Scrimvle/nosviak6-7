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

var UsersBannedPurge = reg.Users.NewCommand(&commands.Command{
	Aliases: []string{"banned_purge"},
	Description: "remove all banned users",
	Permissions: []string{interactions.ADMIN},
	CommandFunc: func(ac *commands.ArgContext, session *sessions.Session) error {
		users, err := database.DB.GetBannedUsers()
		if err != nil {
			return session.ExecuteBranding(make(map[string]any), "commands", "users", "error_occurred.tfx")
		}

		usersRemoved := make([]string, 0)

		/* ranges through all the users which are banned. */
		for _, user := range users {
			err := database.DB.DeleteUser(user)
			if err != nil {
				source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error while trying to ban %s: %v", user.Username, err)
				return session.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "users", "error_occurred.tfx")
			}

			usersRemoved = append(usersRemoved, user.Username)
		}

		reg.UsersLog.WriteLog("%s purged all banned users", session.User.Username)
		return session.ExecuteBranding(map[string]any{"banned_deleted": len(usersRemoved)}, "commands", "users", "alerts", "performer_users_banned_purged.tfx")
	},
})