package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
)

var Who = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"who", "whoami"},
	Permissions: make([]string, 0),
	Description: "unix like whoami command",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "who", "user.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "who", "connected.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "who", "ip.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		for _, session := range sessions.IndexSessions(session.User.Username) {
			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"row": session.User.User()}, "commands", "who", "value_user.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"connected": session.Terminal.ConnTime.Unix()}, "commands", "who", "value_connected.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"ip": session.ConnIP()}, "commands", "who", "value_ip.tfx"),
						Align: gotable2.AlignCenter,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})
