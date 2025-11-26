package commands

import (
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
)

var Clear = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"clear", "cls"},
	Permissions: make([]string, 0),
	Description: "clears the terminal",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		return session.ExecuteBranding(make(map[string]any), "clear_splash.tfx")
	},
})
