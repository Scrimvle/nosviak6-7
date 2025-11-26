package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"
	"time"
)

var UsersBan = reg.Users.NewCommand(&commands.Command{
	Aliases:     []string{"ban"},
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	Description: "ban a user from the cnc",
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		user, err := database.DB.GetUserAsParentalFigure(fmt.Sprint(ac.Args[0].Values[0].Literal), s.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != s.User.ID {
			user, err := database.DB.GetUser(fmt.Sprint(ac.Args[0].Values[0].Literal))
			if err != nil || user == nil {
				return s.ExecuteBranding(map[string]any{"username": fmt.Sprint(ac.Args[0].Values[0].Literal)}, "commands", "users", "bad_user.tfx")
			}

			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "access_denied.tfx")
		}

		if user.Banned {
			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "already_banned.tfx")
		}

		if err := database.DB.BanUser(user); err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error while trying to ban %s: %v", user.Username, err)
			return s.ExecuteBranding(map[string]any{"target": user.User(), "err": err.Error()}, "commands", "users", "error_occurred.tfx")
		}

		sessions.WriteToSession(sessions.IndexSessions(user.Username), func(session *sessions.Session) {
			session.Reader.Reskin(session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "alerts", "user_banned.tfx"))
			time.Sleep(4 * time.Second)
			session.Terminal.Channel.Close()
		})

		reg.UsersLog.WriteLog("%s has banned %s", s.User.Username, user.Username)
		commands.Conn.SendWebhook("user_banned", map[string]any{"user": user.User(), "actor": s.User.User()})
		return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "alerts", "performer_user_ban.tfx")
	},

	Args: []*commands.Arg{{
		Name:        "user",
		Type:        commands.STRING,
		OpenEnded:   false,
		Description: "user to be banned",
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
