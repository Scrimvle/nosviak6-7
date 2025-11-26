package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/functions"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"

	"golang.org/x/exp/maps"
)

// ThemesChange is one of a very short list of subcommands kept inside the main commands directory.
var ThemesChange = Themes.NewCommand(&commands.Command{
	Aliases:     []string{"change", "apply"},
	Description: "change your theme",
	Permissions: make([]string, 0),
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		if !functions.CanAccessThemPermissions(s.User, source.OPTIONS.Strings("theme_changer")...) {
			return s.ExecuteBranding(make(map[string]any), "commands", "themes", "access_denied.tfx")
		}

		context, ok := source.Themes[ac.Args[0].Values[0].ToString()]
		if !ok || context == nil {
			return s.ExecuteBranding(map[string]any{"theme": s.User.Theme}, "commands", "themes", "unknown_theme.tfx")
		}

		s.User.Theme = ac.Args[0].Values[0].ToString()
		if err := database.DB.EditUser(s.User, s.User, commands.Conn.SendWebhook); err != nil {
			return s.ExecuteBranding(map[string]any{"theme": s.User.Theme}, "commands", "themes", "error_occurred.tfx")
		}

		sessions.WriteToSession(sessions.IndexSessions(s.User.Username), func(s *sessions.Session) {
			s.User.Theme, s.Theme = ac.Args[0].Values[0].ToString(), context
			prompt, err := s.ExecuteBrandingToString(make(map[string]any), "prompt.tfx")
			if err != nil {
				return
			}

			/* changes the prompt design and then enables what tab engine to use. */
			s.Reader.Reskin(prompt)
		})

		s.User.Theme, s.Theme = ac.Args[0].Values[0].ToString(), context
		return s.ExecuteBranding(map[string]any{"theme": s.User.Theme}, "home_splash.tfx")
	},

	Args: []*commands.Arg{{
		Name: "theme",
		Type: commands.STRING,
		Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
			return maps.Keys(source.Themes)
		},
	}},
})
