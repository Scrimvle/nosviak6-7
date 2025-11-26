package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/functions"
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"
	"strings"
)

/* UsersPromote supports the promotion of users to new ranks  */
var UsersDemote = reg.Users.NewCommand(&commands.Command{
	Aliases:     []string{"remove_role", "demote"},
	Description: "demote a user from a given role",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.RESELLER},
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		user, err := database.DB.GetUserAsParentalFigure(fmt.Sprint(ac.Args[0].Values[0].Literal), s.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != s.User.ID {
			user, err := database.DB.GetUser(fmt.Sprint(ac.Args[0].Values[0].Literal))
			if err != nil || user == nil {
				return s.ExecuteBranding(map[string]any{"username": fmt.Sprint(ac.Args[0].Values[0].Literal)}, "commands", "users", "bad_user.tfx")
			}

			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "access_denied.tfx")
		}

		/* when they aren't admin or mod and don't have the role they're trying to promote too. */
		if !functions.CanAccessThemPermissions(s.User, "admin", "mod") && !functions.CanAccessThemPermissions(s.User, strings.Split(ac.Header, ",")...) {
			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "alerts", "performer_access_denied_role.tfx")
		}

		/* when they're already missing the role  */
		if !functions.CanAccessThemPermissions(user, strings.Split(ac.Header, ",")...) {
			return s.ExecuteBranding(map[string]any{"target": user.User(), "role": ac.Header}, "commands", "users", "already_missing_role.tfx")
		}

		functions.RemovePermissionRights(user, ac.Header)
		if err := database.DB.EditUser(user, s.User, commands.Conn.SendWebhook); err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error while trying to demote %s from %s: %v", user.Username, ac.Header, err)
			return s.ExecuteBranding(map[string]any{"target": user.User(), "err": err.Error()}, "commands", "users", "error_occurred.tfx")
		}

		sessions.WriteToSession(sessions.IndexSessions(user.Username), func(index *sessions.Session) {
			index.User = user
			index.Reader.PostAlert(&terminal.Alert{
				AlertMessage: index.ExecuteBrandingToStringNoErr(map[string]any{"promoter": s.User.User(), "role": ac.Header}, "commands", "users", "alerts", "demoted_user.tfx"),
				AlertCode:    terminal.MESSAGE,
			})
		})

		reg.UsersLog.WriteLog("%s has demoted %s from %s", s.User.Username, user.Username, ac.Header)
		return s.ExecuteBranding(map[string]any{"target": user.User(), "role": ac.Header}, "commands", "users", "alerts", "performer_demoted_user.tfx")
	},

	Args: []*commands.Arg{{
		Name:        "user",
		Type:        commands.STRING,
		Description: "user to be promoted",
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
