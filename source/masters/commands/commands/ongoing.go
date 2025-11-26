package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source"
	"Nosviak4/source/functions"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/attacks"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"strings"
)

var Ongoing = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"ongoing"},
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.RESELLER},
	Description: "view all the ongoing attacks",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "ongoing", "user.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "ongoing", "method.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "ongoing", "target.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "ongoing", "duration.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "ongoing", "created.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		ongoing, err := database.DB.GetOngoing()
		if err != nil {
			return err
		}

		for _, attack := range ongoing {
			user, err := database.DB.GetUserWithID(attack.User)
			if err != nil {
				continue
			}

			/* checks if we should enabled the masking. */
			if source.OPTIONS.Bool("attacks", "ongoing", "mask_enabled") && functions.CanAccessThemPermissions(session.User, source.OPTIONS.Strings("attacks", "ongoing", "mask_receivers")...) {
				attack.Target = attacks.NewTarget(attack.Target, source.Methods[attack.Method]).MaskTarget(source.OPTIONS.String("attacks", "ongoing", "mask_target"))
				user.Username = strings.Repeat(source.OPTIONS.String("attacks", "ongoing", "mask_user"), len(user.Username))
			}

			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": user.User()}, "commands", "ongoing", "value_user.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"method": attack.Method}, "commands", "ongoing", "value_method.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"target": attack.Target}, "commands", "ongoing", "value_target.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"duration": attack.Duration}, "commands", "ongoing", "value_duration.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"created": attack.Created}, "commands", "ongoing", "value_created.tfx"),
						Align: gotable2.AlignLeft,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})
