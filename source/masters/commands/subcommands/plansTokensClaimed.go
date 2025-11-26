package subcommands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

// PlansTokensClaimed will show all the tokens a user has claimed that are claimed
var PlansTokensClaimed = PlansTokens.NewCommand(&commands.Command{
	Aliases: []string{"claimed"},
	Description: "shows all your claimed tokens",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.ADMIN},
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tokens, err := database.DB.GetClaimedTokens(session.User)
		if err != nil {
			return err
		}

		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "username.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "claimed.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "expires.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "plan.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens", "roles.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		for _, token := range tokens {
			user, err := database.DB.GetUserWithID(token.Owner)
			if err != nil {
				continue
			}

			roles, err := interactions.PopulateStringWithRoles(session.Terminal, user.Roles...)
			if err != nil {
				continue
			}

			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"username": user.Username}, "commands", "plans", "tokens", "value_username.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"claimed": user.Created}, "commands", "plans", "tokens", "value_claimed.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"expires": user.Created + user.Expiry}, "commands", "plans", "tokens", "value_expires.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"plan": token.Plan}, "commands", "plans", "tokens", "value_plan.tfx"),
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