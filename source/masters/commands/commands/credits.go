package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/functions/tui"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/swash/evaluator"
	"fmt"
)

var _ = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"credits"},
	Permissions: make([]string, 0),
	Description: "Nosviak CNC credits",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		if _, err := session.Terminal.Write([]byte(fmt.Sprintf("Nosviak4 ~ %s\r\n", source.VERSION))); err != nil {
			return err
		}

		if _, err := session.Terminal.Write([]byte(tui.WordWrap("FB developed this CNC specifically for Harry's sales under the Nosviak4 brand. Programmed exclusively in Golang, it makes use of a minimal set of third-party libraries. Moreover, it features the newly introduced 'Swash' DSL, amplifying branding potential within this innovative platform. A noteworthy addition is the inclusion of a substantial 21,000 lines of Go code, underscoring its complexity.", "", int(session.Terminal.X)) + "\r\n\n")); err != nil {
			return err
		}

		if _, err := session.Terminal.Write([]byte("Other credits: Prmze & \r\n\n")); err != nil {
			return err
		}

		_, err := session.Terminal.Write([]byte(fmt.Sprintf("Swash evaluator version: %s\r\n", evaluator.VERSION)))
		return err
	},
})
