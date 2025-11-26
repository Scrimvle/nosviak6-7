package subcommands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	command "Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"
	"time"

	"golang.org/x/exp/maps"
)

// PlansApply
var PlansApply = commands.Plans.NewCommand(&command.Command{
	Aliases: []string{"apply"},
	Description: "applies the plan to the user",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.ADMIN},
	CommandFunc: func(ac *command.ArgContext, s *sessions.Session) error {
		user, err := database.DB.GetUserAsParentalFigure(fmt.Sprint(ac.Args[1].Values[0].Literal), s.User)
		if err != nil || user == nil || database.DB.IsSuperuser(user) && user.ID != s.User.ID {
			user, err := database.DB.GetUser(fmt.Sprint(ac.Args[1].Values[0].Literal))
			if err != nil || user == nil {
				return s.ExecuteBranding(map[string]any{"username": fmt.Sprint(ac.Args[1].Values[0].Literal)}, "commands", "users", "bad_user.tfx")
			}

			return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "users", "access_denied.tfx")
		}

		plan, ok := source.Presets[ac.Args[0].Values[0].ToString()]
		if !ok || plan == nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "plans", "invalid_plan.tfx")
		}

		user.Roles = plan.Roles
		user.Theme = plan.Theme
		user.Maxtime = plan.Maxtime
		user.Conns = plan.Concurrents
		user.Cooldown = plan.Cooldown
		user.Created = time.Now().Unix()
		user.MaxAttacks = plan.DailyAttacks
		user.Sessions = source.OPTIONS.Ints("default_user", "max_sessions")
		user.Expiry = int64(plan.Days) * 86400

		/* applies the plan to the user */
		if err := database.DB.EditUser(user, s.User, command.Conn.SendWebhook); err != nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "plans", "error_occurred.tfx")
		}

		return s.ExecuteBranding(map[string]any{"target": user.User()}, "commands", "plans", "plan_applied.tfx")
	},

	Args: []*command.Arg{
		{
			Name: "plan",
			Type: command.STRING,
			Description: "the plan to apply from",
			Callback: func(ac *command.ArgContext, s *sessions.Session, i int) []string {
				return maps.Keys(source.Presets)
			},
		},

		{
			Name: "username",
			Type: command.STRING,
			Description: "the username to be used",
			Callback: func(ac *command.ArgContext, s *sessions.Session, i int) []string {
				users, err := database.DB.GetUsersAsParent(s.User)
				if err != nil || len(users) == 0 {
					return make([]string, 0)
				}

				buf := make([]string, 0)
				for _, user := range users {
					buf = append(buf, user.Username)
				}

				return buf
			},

			NotProvided: func(s *sessions.Session, _ []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "username_apply.tfx"), s.Reader.Context)
				content, err := read.ReadLine()
				if err != nil {
					return "", err
				}
	
				if len(content) == 0 {
					return "", fmt.Errorf("not allowed")
				}
	
				return string(content), nil
			},
		}, 
	},
})