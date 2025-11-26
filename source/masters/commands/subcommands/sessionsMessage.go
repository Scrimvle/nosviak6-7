package subcommands

import (
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/masters/terminal/interactions"

)

var SessionsMessage = reg.Sessions.NewCommand(&commands.Command{
	Aliases:     []string{"message", "msg"},
	Description: "send a message to a session",
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	CommandFunc: func(ac *commands.ArgContext, session *sessions.Session) error {
		recv := make([]*sessions.Session, 0)
		for _, value := range ac.Args[1].Values {
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
			index.Reader.PostAlert(&terminal.Alert{
				AlertCode:    terminal.MESSAGE,
				AlertMessage: index.ExecuteBrandingToStringNoErr(map[string]any{"sender": session.User.User(), "message": ac.Args[0].Values[0].ToString()}, "session_message.tfx"),
			})
		})

		return session.ExecuteBranding(make(map[string]any), "commands", "sessions", "message_sent.tfx")
	},

	/* params for the command */
	Args: []*commands.Arg{{
			Name: "message",
			Type: commands.STRING,	
		},
		{
			Name: "id",
			Type: commands.STRING,
			OpenEnded: true,
			Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
				return s.Callback()
			},
		},
	},
})
