package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/terminal/interactions"
	"strings"
	"time"
)

var UsersKickMake = UsersKicks.NewCommand(&commands.Command{
	Aliases:     []string{"make"},
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	Description: "kick",
	CommandFunc: func(context *commands.ArgContext, s *sessions.Session) error {
		durationOfKick := time.Unix(int64(context.Args[0].Values[0].Literal.(int)*86400), 0)
		user, err := database.DB.GetUserAsParentalFigure(context.Args[1].Values[0].ToString(), s.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != s.User.ID {
			user, err := database.DB.GetUser(context.Args[1].Values[0].ToString())
			if err != nil || user == nil {
				return s.ExecuteBranding(map[string]any{"username": context.Args[1].Values[0].ToString()}, "commands", "users", "bad_user.tfx")
			}

			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "access_denied.tfx")
		}

		if s.User == nil || user.ID == s.User.ID {
			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "cant_kick_yourself.tfx")
		}

		if err := database.DB.NewKick(&database.Kick{WeightedFor: durationOfKick.Unix(), Created: time.Now().Unix(), Reason: strings.Join(commands.ToString(context.Args[2].Values), " "), User: user.ID, Issuer: s.User.ID}); err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error when trying to kick %s for %d days: %v", user.Username, context.Args[0].Values[0].Literal, err)
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "users", "error_occurred.tfx")
		}

		sessions.WriteToSession(sessions.IndexSessions(user.Username), func(index *sessions.Session) {
			index.Reader.Reskin(index.ExecuteBrandingToStringNoErr(map[string]any{"issuer": s.User.User(), "reason": strings.Join(commands.ToString(context.Args[2].Values), " ")}, "session_user_kicked.tfx"))
			time.Sleep(4 * time.Second)
			index.Cancel()
			index.Terminal.Channel.Close()
		})

		reg.UsersLog.WriteLog("%s has kicked %s for %d days", s.User.Username, user.Username, context.Args[0].Values[0].Literal)
		return s.ExecuteBranding(map[string]any{"receiver": user.User()}, "commands", "users", "alerts", "performer_user_kicked.tfx")
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
			Description: "reason why they've been kicked",
		},
	},
})
