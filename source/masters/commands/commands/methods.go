package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

var Methods = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"methods"},
	Permissions: make([]string, 0),
	Description: "list all available methods",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Align: gotable2.AlignCenter,
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "methods", "name.tfx"),
				},
				{
					Align: gotable2.AlignCenter,
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "methods", "description.tfx"),
				},
				{
					Align: gotable2.AlignCenter,
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "methods", "roles.tfx"),
				},
			},
		})

		for key, method := range source.Methods {
			access, err := interactions.PopulateStringWithRoles(session.Terminal, method.Access...)
			if err != nil {
				return err
			}

			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Align: gotable2.AlignLeft,
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"name": key}, "commands", "methods", "value_name.tfx"),
					},
					{
						Align: gotable2.AlignLeft,
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"description": method.Description}, "commands", "methods", "value_description.tfx"),
					},
					{
						Align: gotable2.AlignLeft,
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"roles": access}, "commands", "methods", "value_roles.tfx"),
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})
