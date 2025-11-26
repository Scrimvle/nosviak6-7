package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"fmt"
	"net/url"
)

var Web = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"web"},
	Permissions: []string{interactions.ADMIN},
	Description: "view the conns details for the web panel",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		url := url.URL{
			Scheme: "http",
			Path: "/",
			Host: fmt.Sprintf("%s:%d", source.OPTIONS.String("web", "address"), source.OPTIONS.Ints("web", "port")),
		}

		// if https is enabled, we change the url to suit it
		if source.OPTIONS.Bool("web", "ssl", "enabled") {
			url.Scheme = "https"
		}

		status := session.ExecuteBrandingToStringNoErr(make(map[string]any), "true.tfx")
		if !source.OPTIONS.Bool("web", "enabled") {
			status = session.ExecuteBrandingToStringNoErr(make(map[string]any), "false.tfx")
		}

		return session.ExecuteBranding(map[string]any{"status": status, "url": url.String()}, "web.tfx")
	},
})
