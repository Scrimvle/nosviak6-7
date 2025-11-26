package commands

import (
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
)

var Home = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"home", "reset"},
	Permissions: make([]string, 0),
	Description: "go back to home",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		return session.ExecuteBranding(make(map[string]any), "home_splash.tfx")
	},
})
