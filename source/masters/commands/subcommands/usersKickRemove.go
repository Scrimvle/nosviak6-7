package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/terminal/interactions"
)

var UsersKickRemove = UsersKicks.NewCommand(&commands.Command{
	Aliases:     []string{"remove"},
	Description: "remove a kick",
	Permissions: []string{interactions.ADMIN},
	CommandFunc: func(context *commands.ArgContext, s *sessions.Session) error {
		user, err := database.DB.GetUserAsParentalFigure(context.Args[0].Values[0].ToString(), s.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != s.User.ID {
			user, err := database.DB.GetUser(context.Args[0].Values[0].ToString())
			if err != nil || user == nil {
				return s.ExecuteBranding(map[string]any{"username": context.Args[1].Values[0].ToString()}, "commands", "users", "bad_user.tfx")
			}

			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "access_denied.tfx")
		}

		kicks, err := database.DB.GetOngoingKicks(user)
		if err != nil {
			return s.ExecuteBranding(map[string]any{"target": user.User(), "err": err.Error()}, "commands", "users", "error_occurred.tfx")
		}

		for _, kick := range kicks {
			if kick.User != user.ID {
				continue
			}

			if err := database.DB.RemoveKick(kick.ID); err != nil {
				source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error when trying to remove %s's kick issued by user(%d): %v", user.Username, kick.Issuer, err)
				return s.ExecuteBranding(map[string]any{"target": user.User(), "err": err.Error()}, "commands", "users", "error_occurred.tfx")
			}
		}

		reg.UsersLog.WriteLog("%s has removed all of %s's kicks", s.User.Username, user.Username)
		return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "alerts", "performer_kicks_removed.tfx")
	},

	Args: []*commands.Arg{{
		Name:        "user",
		Type:        commands.STRING,
		Description: "user to remove the kicks from",
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
