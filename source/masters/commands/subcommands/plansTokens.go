package subcommands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source"
	"Nosviak4/source/database"
	command "Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

// PlansTokens will show all the tokens made by this user
var PlansTokens = commands.Plans.NewCommand(&command.Command{
	Aliases: []string{"tokens"},
	Description: "shows all your tokens",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.ADMIN},
	CommandFunc: func(context *command.ArgContext, session *sessions.Session) error {
		tokens, err := database.DB.GetTokens(session.User)
		if err != nil {
			return err
		}

		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "token.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "plan.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "expires.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "length.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "roles.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		for _, token := range tokens {
			plan, ok := source.Presets[token.Plan]
			if !ok || plan == nil {
				continue
			}

			roles, err := interactions.PopulateStringWithRoles(session.Terminal, plan.Roles...)
			if err != nil {
				continue
			}

			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"token": token.Token}, "commands", "plans", "tokens", "value_token.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"plan": token.Plan}, "commands", "plans", "tokens", "value_plan.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"expires": token.Created + token.Expiry}, "commands", "plans", "tokens", "value_expires.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"length": plan.Days}, "commands", "plans", "tokens", "value_length.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"roles": roles}, "commands", "plans", "tokens", "value_roles.tfx"),
						Align: gotable2.AlignCenter,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})