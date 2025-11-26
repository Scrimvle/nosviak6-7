package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
)

var Logins = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"logins"},
	Permissions: make([]string, 0),
	Description: "list all your logins",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "id.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "user.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "terminal.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "ip.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "created.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		logins, err := database.DB.GetLogins()
		if err != nil {
			return err
		}

		for _, login := range logins {
			if login.User != session.User.ID {
				continue
			}

			public := login.Login()
			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_id.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_user.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_terminal.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_ip.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_created.tfx"),
						Align: gotable2.AlignLeft,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})
