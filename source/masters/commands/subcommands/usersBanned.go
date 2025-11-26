package subcommands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

var UsersBanned = reg.Users.NewCommand(&commands.Command{
	Aliases:     []string{"banned"},
	Description: "list all banned users",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.RESELLER},
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "id.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "username.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "maxtime.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "conns.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "cooldown.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "roles.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		users, err := database.DB.GetBannedUsers()
		if err != nil {
			return session.ExecuteBranding(make(map[string]any), "commands", "users", "error_occurred.tfx")
		}

		for _, user := range users {
			roles, err := interactions.PopulateStringWithRoles(session.Terminal, user.Roles...)
			if err != nil {
				return err
			}

			warns, err := database.DB.GetOngoingWarnings(user)
			if err != nil || warns == nil {
				warns = make([]*database.Warn, 0)
			}

			kicks, err := database.DB.GetOngoingKicks(user)
			if err != nil || kicks == nil {
				kicks = make([]*database.Kick, 0)
			}

			public := user.User()
			tablet.Append(&gotable2.Row{Columns: []*gotable2.Column{
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "users", "value_id.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "users", "value_username.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "users", "value_maxtime.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "users", "value_conns.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "users", "value_cooldown.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"roles": roles}, "commands", "users", "value_roles.tfx"),
					Align: gotable2.AlignLeft,
				},
			}})

			/* whenever they've been warned to a kick */
			if len(warns) >= 1 || len(kicks) >= 1 {
				tablet.Rows[len(tablet.Rows)-1].Columns[1] = &gotable2.Column{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": user.User()}, "commands", "users", "value_username_notice.tfx"),
					Align: gotable2.AlignLeft,
				}
			}

			if user.Banned {
				tablet.Rows[len(tablet.Rows)-1].Columns[1] = &gotable2.Column{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": user.User()}, "commands", "users", "value_username_banned.tfx"),
					Align: gotable2.AlignLeft,
				}
			}
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},
})
