package subcommands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"
)

var UsersLogins = reg.Users.NewCommand(&commands.Command{
	Aliases:     []string{"logins"},
	Description: "view all the logins from a user",
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		user, err := database.DB.GetUserAsParentalFigure(fmt.Sprint(context.Args[0].Values[0].Literal), session.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != session.User.ID {
			user, err := database.DB.GetUser(fmt.Sprint(context.Args[0].Values[0].Literal))
			if err != nil || user == nil {
				return session.ExecuteBranding(map[string]any{"username": fmt.Sprint(context.Args[0].Values[0].Literal)}, "commands", "users", "bad_user.tfx")
			}

			return session.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "access_denied.tfx")
		}

		logins, err := database.DB.GetUserLogins(user)
		if err != nil || logins == nil {
			return session.ExecuteBranding(map[string]any{"target": user.User(), "err": err.Error()}, "commands", "users", "error_occurred.tfx")
		}

		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "id.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "user.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "terminal.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "ip.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "logins", "created.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		for _, login := range logins {
			if login.User != session.User.ID {
				continue
			}

			public := login.Login()
			tablet.Append(&gotable2.Row{
				Columns: []*gotable2.Column{
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_id.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_user.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_terminal.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_ip.tfx"),
						Align: gotable2.AlignLeft,
					},
					{
						Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"row": public}, "commands", "logins", "value_created.tfx"),
						Align: gotable2.AlignLeft,
					},
				},
			})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},

	Args: []*commands.Arg{{
		Type:        commands.STRING,
		Name:        "user",
		Description: "provides the user we will view the logins for",
		Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
			child, err := database.DB.GetUsersAsParent(s.User)
			if err != nil {
				return make([]string, 0)
			}

			children := make([]string, 0)
			for _, child := range child {
				children = append(children, child.Username)
			}

			return children
		},
	}},
})
