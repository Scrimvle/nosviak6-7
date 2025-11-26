package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/modules/gotable2"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"
)

var UsersAttacks = reg.Users.NewCommand(&commands.Command{
	Aliases:     []string{"attacks"},
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	Description: "lookup attacks sent from a user",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		user, err := database.DB.GetUserAsParentalFigure(fmt.Sprint(context.Args[0].Values[0].Literal), session.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != session.User.ID {
			user, err := database.DB.GetUser(fmt.Sprint(context.Args[0].Values[0].Literal))
			if err != nil || user == nil {
				return session.ExecuteBranding(map[string]any{"username": fmt.Sprint(context.Args[0].Values[0].Literal)}, "commands", "users", "bad_user.tfx")
			}

			return session.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "access_denied.tfx")
		}

		attacks, err := database.DB.GetUserAttacks(user.Username)
		if err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error while trying to view %s's attacks: %v", user.Username, err)
			return session.ExecuteBranding(map[string]any{"target": user.User(), "err": err.Error()}, "commands", "users", "error_occurred.tfx")
		}

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

	Args: []*commands.Arg{{
		Name:        "user",
		Type:        commands.STRING,
		OpenEnded:   false,
		Description: "user to search under",
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
