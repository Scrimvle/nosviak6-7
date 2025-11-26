package commands

import (
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"strconv"
)

var WarnsRemove = Warns.NewCommand(&commands.Command{
	Aliases: []string{"remove"},
	Description: "removes a warn from the database",
	Permissions: []string{interactions.ADMIN},
	CommandFunc: func(ac *commands.ArgContext, session *sessions.Session) error {
		if err := database.DB.RemoveWarn(ac.Args[0].Values[0].Literal.(int)); err != nil {
			return session.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "warns", "error_occurred.tfx")
		}

		return session.ExecuteBranding(make(map[string]any), "commands", "warns", "warn_removed.tfx")
	},

	Args: []*commands.Arg{{
		Name: "warnid",
		Type: commands.NUMBER,
		Description: "warn id to be removed",
		Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
			buf := make([]string, 0)
			warns, err := database.DB.GetWarns()
			if err != nil {
				return make([]string, 0)
			}

			for _, warn := range warns {
				buf = append(buf, strconv.Itoa(warn.ID))
			}

			return buf
		},
	}},
})