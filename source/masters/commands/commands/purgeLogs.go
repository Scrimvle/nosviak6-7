package commands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

var Clearlogs = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"clearlogs", "purgelogs"},
	Permissions: []string{interactions.ADMIN},
	Description: "clears the logs",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		switch context.Args[0].Values[0].Literal {

		case "logins", "warns", "attacks":
			if err := database.DB.Truncate(context.Args[0].Values[0].Literal.(string)); err != nil {
				source.LOGGER.AggregateTerminal().WriteLog(gologr.DEFAULT, "Error occurred while truncating logs: %v", err)
				return session.ExecuteBranding(map[string]any{"flusher": context.Args[0].Values[0].Literal, "err": err.Error()}, "commands", "clearlogs", "error_occurred.tfx")
			}

			return session.ExecuteBranding(map[string]any{"flusher": context.Args[0].Values[0].Literal}, "commands", "clearlogs", "success.tfx")

		default:
			return session.ExecuteBranding(map[string]any{"flusher": context.Args[0].Values[0].Literal}, "commands", "clearlogs", "unknown_flusher.tfx")
		}
	},

	Args: []*commands.Arg{
		{
			Name:        "flush",
			Type:        commands.STRING,
			OpenEnded:   false,
			Description: "what the clearlogs command targets",
			Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
				return []string{"logins", "warns", "attacks"}
			},
		},
	},
})
