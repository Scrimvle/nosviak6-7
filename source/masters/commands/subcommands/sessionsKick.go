package subcommands

import (
	"time"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	reg "Nosviak4/source/masters/commands/commands"
)

var SessionsClose = reg.Sessions.NewCommand(&commands.Command{
	Aliases:     []string{"kick", "close"},
	Description: "close a users session",
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
			index.Reader.Reskin(index.ExecuteBrandingToStringNoErr(make(map[string]any), "session_closed.tfx"))
			time.Sleep(4 * time.Second)
			
			index.Cancel()
			index.Terminal.Channel.Close()
		})

		return session.ExecuteBranding(make(map[string]any), "commands", "sessions", "sessions_kicked.tfx")
	},

	/* params for the command */
	Args: []*commands.Arg{{
		Name: "id",
		Type: commands.STRING,
		OpenEnded: true,
		Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string { return s.Callback() },
	}},
})
