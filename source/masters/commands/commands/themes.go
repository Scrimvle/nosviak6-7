package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
)

var Themes = commands.ROOT.NewCommand(&commands.Command{
	Aliases: []string{"themes", "theme"},
	Description: "change the theme",
	Permissions: make([]string, 0),
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "themes", "name.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "themes", "description.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		for name, theme := range source.Themes {
			if theme.Hidden {
				continue
			}

			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"name": name}, "commands", "themes", "value_name.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"description": theme.Description}, "commands", "themes", "value_description.tfx"),
						Align: gotable2.AlignLeft,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})