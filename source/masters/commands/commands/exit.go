package commands

import (
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
)

var Exit = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"exit", "close", "leave"},
	Permissions: make([]string, 0),
	Description: "closes your session",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		defer func() {
			session.Terminal.Channel.Close()
			session.Cancel()
			delete(sessions.Sessions, session.Opened)
		}()

		return session.ExecuteBranding(make(map[string]any), "exit.tfx")
	},
})
