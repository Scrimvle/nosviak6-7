package commands

import (
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/masters/terminal/interactions"
	"strings"
)

var Broadcast = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"broadcast", "announce"},
	Permissions: []string{interactions.ADMIN},
	Description: "broadcast to all sessions",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		buf := make([]string, 0)
		for _, token := range context.Args[0].Values[:len(context.Args[0].Values) - 1] {
			buf = append(buf, token.ToString())
		}

		buf = append(buf, context.Args[0].Values[len(context.Args[0].Values) - 1].ToString())
		for _, recv := range sessions.Sessions {
			if recv.Opened == session.Opened {
				continue
			}

			recv.Reader.PostAlert(&terminal.Alert{AlertCode: terminal.MESSAGE, AlertMessage: recv.ExecuteBrandingToStringNoErr(map[string]any{"sender": session.User.User(), "message": strings.Join(buf, " ")}, "broadcast_recv.tfx")})
		}

		return session.ExecuteBranding(map[string]any{"sent": len(sessions.Sessions), "message": strings.Join(buf, " ")}, "broadcast_sent.tfx")
	},

	Args: []*commands.Arg{
		{
			Name: "message",
			Type: commands.STRING,
			OpenEnded: true,
			Description: "message to be broadcasted",
		},
	},
})
