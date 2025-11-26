package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/masters/terminal/interactions"
	"strings"
	"time"
)

var WarnsMake = Warns.NewCommand(&commands.Command{
	Aliases:     []string{"make"},
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	Description: "warn a user about a matter of importance",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		durationOfWarn := time.Unix(int64(context.Args[0].Values[0].Literal.(int)*86400), 0)
		user, err := database.DB.GetUserAsParentalFigure(context.Args[1].Values[0].ToString(), session.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != session.User.ID {
			user, err := database.DB.GetUser(context.Args[1].Values[0].ToString())
			if err != nil || user == nil {
				return session.ExecuteBranding(map[string]any{"username": context.Args[1].Values[0].ToString()}, "commands", "warns", "bad_user.tfx")
			}

			return session.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "warns", "access_denied.tfx")
		}

		if session.User == nil || user.ID == session.User.ID {
			return session.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "warns", "cant_warn_yourself.tfx")
		}

		if err := database.DB.NewWarn(&database.Warn{WeightedFor: durationOfWarn.Unix(), Created: time.Now().Unix(), Reason: strings.Join(commands.ToString(context.Args[2].Values), " "), User: user.ID, Issuer: session.User.ID}); err != nil {
			return session.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "warns", "error_occurred.tfx")
		}

		ongoingWarns, err := database.DB.GetOngoingWarnings(user)
		if err != nil {
			return session.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "warns", "error_occurred.tfx")
		}

		/* whenever they hit this marker, we kick their account */
		if source.OPTIONS.Ints("default", "warns_before_kick") < len(ongoingWarns) {
			sessions.WriteToSession(sessions.IndexSessions(user.Username), func(index *sessions.Session) {
				index.Reader.Reskin(index.ExecuteBrandingToStringNoErr(map[string]any{"issuer": session.User.User(), "reason": strings.Join(commands.ToString(context.Args[2].Values), " ")}, "session_warned_threshold.tfx"))
				time.Sleep(4 * time.Second)
				index.Cancel()
				index.Terminal.Channel.Close()
			})

			return session.ExecuteBranding(map[string]any{"receiver": user.User()}, "commands", "warns", "success_temp_kick.tfx")
		}

		sessions.WriteToSession(sessions.IndexSessions(user.Username), func(index *sessions.Session) {
			index.Reader.PostAlert(&terminal.Alert{
				AlertMessage: index.ExecuteBrandingToStringNoErr(map[string]any{"issuer": session.User.User(), "reason": strings.Join(commands.ToString(context.Args[2].Values), " ")}, "session_warned.tfx"),
				AlertCode:    terminal.MESSAGE,
			})
		})

		return session.ExecuteBranding(map[string]any{"receiver": user.User()}, "commands", "warns", "success.tfx")
	},

	Args: []*commands.Arg{
		{
			Type:        commands.NUMBER,
			Name:        "length",
			Description: "length in days",
		},
		{
			Type:        commands.STRING,
			Name:        "user",
			Description: "length in days",
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
		},
		{
			Type:        commands.STRING,
			Name:        "reason",
			Description: "reason why they're warned",
		},
	},
})
