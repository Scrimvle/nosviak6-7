package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

var Sessions = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"sessions", "online"},
	Description: "moderate and manage sessions",
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "sessions", "id.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "sessions", "user.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "sessions", "connected.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "sessions", "ip.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "sessions", "idle.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "sessions", "roles.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		for _, session := range sessions.Sessions {
			roles, err := interactions.PopulateStringWithRoles(session.Terminal, session.User.Roles...)
			if err != nil {
				return err
			}

			if session.Reader == nil {
				session.Reader = session.Terminal.NewRead(session.ExecuteBrandingToStringNoErr(make(map[string]any), "prompt.tfx"))
			}

			tablet.Append(&gotable2.Row{Columns: []*gotable2.Column{
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"id": session.Opened}, "commands", "sessions", "value_id.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"username": session.User.Username}, "commands", "sessions", "value_user.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"connected": session.Terminal.ConnTime.Unix()}, "commands", "sessions", "value_connected.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"ip": session.ConnIP()}, "commands", "sessions", "value_ip.tfx"),
					Align: gotable2.AlignLeft,
				},

				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"idle": session.Reader.ReaderIdle.Unix()}, "commands", "sessions", "value_idle.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"roles": roles}, "commands", "sessions", "value_roles.tfx"),
					Align: gotable2.AlignLeft,
				},
			}})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})
