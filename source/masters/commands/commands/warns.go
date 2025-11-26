package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)


var Warns = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"warns"},
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	Description: "warn a user",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "warns", "id.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "warns", "user.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "warns", "issuer.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "warns", "reason.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "warns", "created.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "warns", "expires.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		warns, err := database.DB.GetWarns()
		if err != nil {
			return session.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "warns", "error_occurred.tfx")
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
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"id": warn.ID}, "commands", "warns", "value_id.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"receiver": receiver.User()}, "commands", "warns", "value_user.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"issuer": issuer.User()}, "commands", "warns", "value_issuer.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"reason": warn.Reason}, "commands", "warns", "value_reason.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"created": warn.Created}, "commands", "warns", "value_created.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"expires": warn.Created + warn.WeightedFor}, "commands", "warns", "value_expires.tfx"),
						Align: gotable2.AlignCenter,
					},
				},
			})
		}
		
		return session.Table(tablet, context.Command.Aliases[0])
	},
})
