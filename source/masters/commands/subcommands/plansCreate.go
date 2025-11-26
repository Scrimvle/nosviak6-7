package subcommands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	command "Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"

	"golang.org/x/exp/maps"
)

var PlansCreate = commands.Plans.NewCommand(&command.Command{
	Aliases: []string{"create", "new"},
	Description: "create a new user",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.RESELLER},
	CommandFunc: func(ac *command.ArgContext, s *sessions.Session) error {
		plan, ok := source.Presets[ac.Args[0].Values[0].ToString()]
		if !ok || plan == nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "plans", "invalid_plan.tfx")
		}

		if err := database.DB.NewUser(&database.User{Username: fmt.Sprint(ac.Args[1].Values[0].Literal), Password: []byte(fmt.Sprint(ac.Args[2].Values[0].Literal)), API: false, Roles: plan.Roles, Theme: plan.Theme, NewUser: true, Maxtime: plan.Maxtime, Conns: plan.Concurrents, Cooldown: plan.Cooldown, Expiry: int64(plan.Days) * 86400, Sessions: source.OPTIONS.Ints("default_user", "max_sessions"), MaxAttacks: plan.DailyAttacks}, s.User, command.Conn.SendWebhook); err != nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "plans", "error_occurred.tfx")
		}

		return s.ExecuteBranding(make(map[string]any), "commands", "plans", "success.tfx")
	},

	Args: []*command.Arg{
		{
			Name: "plan",
			Type: command.STRING,
			Description: "the plan to create from",
			Callback: func(ac *command.ArgContext, s *sessions.Session, i int) []string {
				return maps.Keys(source.Presets)
			},
		},
		{
			Type: command.STRING,
			Name: "username",
			NotProvided: func(s *sessions.Session, _ []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "username.tfx"), s.Reader.Context)
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
		{
			Type: command.STRING,
			Name: "password",
			NotProvided: func(s *sessions.Session, _ []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "password.tfx"), s.Reader.Context).ChangeMaxLen(source.OPTIONS.Ints("maximum_password_length")).ChangeMinLen(source.OPTIONS.Ints("minimum_password_length"))
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