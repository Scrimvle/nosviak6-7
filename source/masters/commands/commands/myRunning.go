package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
)

var MyRunning = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"myrunning"},
	Permissions: make([]string, 0),
	Description: "view all your ongoing attacks",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "ongoing", "method.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "ongoing", "target.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "ongoing", "duration.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "ongoing", "created.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		ongoing, err := database.DB.GetOngoingUser(session.User)
		if err != nil {
			return err
		}

		for _, attack := range ongoing {
			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"method": attack.Method}, "commands", "ongoing", "value_method.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"target": attack.Target}, "commands", "ongoing", "value_target.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"duration": attack.Duration}, "commands", "ongoing", "value_duration.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"created": attack.Created}, "commands", "ongoing", "value_created.tfx"),
						Align: gotable2.AlignLeft,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})
