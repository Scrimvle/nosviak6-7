package subcommands

import (
	"Nosviak4/source/masters/commands"
	reg "Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"bytes"
	"time"
)

var SessionsScreen = reg.Sessions.NewCommand(&commands.Command{
	Aliases:     []string{"screen"},
	Description: "screenshare like function",
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	CommandFunc: func(ac *commands.ArgContext, session *sessions.Session) error {
		screening := sessions.IndexSessions(ac.Args[0].Values[0].Literal.(string))
		if len(screening) == 0 || len(screening) > 1 {
			return session.ExecuteBranding(make(map[string]any), "commands", "sessions", "no_receivers.tfx")
		}

		if session.Included(screening) {
			return session.ExecuteBranding(make(map[string]any), "commands", "sessions", "selfaction_denied.tfx")
		}

		tickContext, writeContext, timeout := time.NewTicker(50*time.Millisecond), make([]byte, 0), 0

		for {
			select {

			case <-tickContext.C:
				if bytes.EqualFold(screening[0].Terminal.Screen.Bytes(), writeContext) {
					if timeout <= 400 {
						timeout++
						continue
					}

					_, err := session.Terminal.Write(append([]byte("\x1bc"), session.Terminal.Screen.Bytes()...))
					return err
				}

				/* resets the inactivity counter */
				timeout = 0

				// we employ a different method when the array size increases
				if len(screening[0].Terminal.Screen.Bytes()) > len(writeContext) && bytes.Equal(screening[0].Terminal.Screen.Bytes()[:len(writeContext)], writeContext) && len(writeContext) > 0 {
					if _, err := session.Terminal.Channel.Write(screening[0].Terminal.Screen.Bytes()[len(writeContext):]); err != nil {
						return err
					}

					writeContext = append(writeContext, screening[0].Terminal.Screen.Bytes()[len(writeContext):]...)
					continue
				}

				writeContext = screening[0].Terminal.Screen.Bytes()
				if _, err := session.Terminal.Channel.Write(writeContext); err != nil {
					return err
				}

			case buf, ok := <-session.Terminal.Signal.Queue:
				if !ok || buf == nil {
					return nil
				}

				/* escape = break from screening */
				if bytes.Equal(buf, []byte{27}) {
					_, err := session.Terminal.Write(append([]byte("\x1bc"), session.Terminal.Screen.Bytes()...))
					return err
				}

				screening[0].Terminal.Signal.Queue <- buf
			}
		}
	},

	/* params for the command */
	Args: []*commands.Arg{{
		Name:      "id",
		Type:      commands.STRING,
		OpenEnded: false,
		Callback:  func(ac *commands.ArgContext, s *sessions.Session, i int) []string { return s.Callback() },
	}},
})
