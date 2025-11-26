package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"

	"golang.org/x/exp/slices"
)

var Plans = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"plans"},
	Description: "plan preset functionality",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.RESELLER},
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "name.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "maxtime.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "conns.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "cooldown.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "length.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "roles.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		for name, plan := range source.Presets {
			if !slices.Contains(plan.Roles, "member") {
				plan.Roles = append(plan.Roles, "member")
			}

			exec, err := interactions.PopulateStringWithRoles(session.Terminal, plan.Roles...)
			if err != nil {
				return err
			}

			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"name": name}, "commands", "plans", "value_name.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"maxtime": plan.Maxtime}, "commands", "plans", "value_maxtime.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"conns": plan.Concurrents}, "commands", "plans", "value_conns.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"cooldown": plan.Cooldown}, "commands", "plans", "value_cooldown.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"length": plan.Days}, "commands", "plans", "value_length.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"roles": exec}, "commands", "plans", "value_roles.tfx"),
						Align: gotable2.AlignLeft,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})
