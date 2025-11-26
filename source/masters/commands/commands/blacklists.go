package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

// Blacklists will show all the blacklisted targets and options to make a new blacklisted target.
var Blacklists = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"blacklist"},
	Description: "review your scope of blacklisted targets",
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "blacklist", "target.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "blacklist", "user.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text: session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "blacklist", "expires.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		blacklisted, err := database.DB.GetBlacklistedTargets()
		if err != nil {
			return session.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "blacklist", "error_occurred.tfx")
		}

		for _, blacklist := range blacklisted {
			user, err := database.DB.GetUserWithID(blacklist.User)
			if err != nil {
				return err
			}

			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"target": blacklist.Target}, "commands", "blacklist", "value_target.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"row": user.User()}, "commands", "blacklist", "value_user.tfx"),
						Align: gotable2.AlignCenter,
					},
					{
						Text: session.ExecuteBrandingToStringNoErr(map[string]any{"expires": int(blacklist.Created + blacklist.Expires)}, "commands", "blacklist", "value_expires.tfx"),
						Align: gotable2.AlignCenter,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})