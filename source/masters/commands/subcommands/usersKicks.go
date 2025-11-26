package subcommands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

var UsersKicks = reg.Users.NewCommand(&commands.Command{
	Aliases: []string{"kicks"},
	Description: "manage user kicks",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.RESELLER},
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "kicks", "id.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "kicks", "user.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "kicks", "issuer.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "kicks", "reason.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "kicks", "created.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "kicks", "expires.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		warns, err := database.DB.GetKicks()
		if err != nil {
			return session.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "users", "kicks", "error_occurred.tfx")
		}

		for _, warn := range warns {
			issuer, err := database.DB.GetUserWithID(warn.Issuer)
			if err != nil {
				continue
			}

			receiver, err := database.DB.GetUserWithID(warn.User)
			if err != nil {
				continue
			}

			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"id": warn.ID}, "commands", "users", "kicks", "value_id.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"receiver": receiver.User()}, "commands", "users", "kicks", "value_user.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"issuer": issuer.User()}, "commands", "users", "kicks", "value_issuer.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"reason": warn.Reason}, "commands", "users", "kicks", "value_reason.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"created": warn.Created}, "commands", "users", "kicks", "value_created.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"expires": warn.Created + warn.WeightedFor}, "commands", "users", "kicks", "value_expires.tfx"),
						Align: gotable2.AlignCenter,
					},
				},
			})
		}
		
		return session.Table(tablet, context.Command.Aliases[0])
	},
})