package sessions

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/functions/tui"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/web/propagator"
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// ExecuteBranding wraps a terminal function
func (s *Session) ExecuteBranding(vars map[string]any, basePath ...string) error {
	if s.Theme == nil {
		defaultTheme, err := source.GetTheme("default")
		if err != nil {
			return err
		}

		s.Theme = defaultTheme
	}

	content := append(strings.Split(filepath.Join(s.Theme.Branding), string(filepath.Separator)), basePath...)
	if _, ok := source.OPTIONS.Config.Renders[filepath.Join(content...)]; !ok {
		theme, err := source.GetTheme(source.OPTIONS.String("default_theme"))
		if err != nil {
			return err
		}

		content = append(strings.Split(filepath.Join(theme.Branding), string(filepath.Separator)), basePath...)
	}

	return s.Terminal.ExecuteBranding(s.AppendDefaultSession(vars), content...)
}

// ExecuteBrandingToString wraps a terminal function
func (s *Session) ExecuteBrandingToString(vars map[string]any, basePath ...string) (string, error) {
	if s.Theme == nil {
		defaultTheme, err := source.GetTheme("default")
		if err != nil {
			return "", err
		}

		s.Theme = defaultTheme
	}

	content := append(strings.Split(filepath.Join(s.Theme.Branding), string(filepath.Separator)), basePath...)
	if _, ok := source.OPTIONS.Config.Renders[filepath.Join(content...)]; !ok {
		theme, err := source.GetTheme(source.OPTIONS.String("default_theme"))
		if err != nil {
			return "", err
		}

		content = append(strings.Split(filepath.Join(theme.Branding), string(filepath.Separator)), basePath...)
	}

	return s.Terminal.ExecuteBrandingToString(s.AppendDefaultSession(vars), content...)
}

// ExecuteBrandingToStringNoErr wraps the ExecuteBrandingToString method to execute
func (s *Session) ExecuteBrandingToStringNoErr(vars map[string]any, content ...string) string {
	context, err := s.ExecuteBrandingToString(vars, content...)
	if err != nil {
		return ""
	}

	return context
}

// appendDefaultSession will append all the default app vars based on a session
func (s *Session) AppendDefaultSession(vars map[string]any) map[string]any {
	vars["ip"] = s.ConnIP()
	vars["db"] = database.DBFields
	vars["user"] = s.User.User()
	vars["session"] = sessionObject{
		Total: func() int { return len(Sessions) },
		X:     int(s.Terminal.X), Y: int(s.Terminal.Y),
	}

	vars["bool"] = s.Bool

	/* checks if the user is online or not */
	vars["is_online"] = func(user string) bool {
		if user == s.User.Username {
			return true
		}

		return len(IndexSessions(user)) >= 1
	}

	if s.Reader == nil || s.Reader.Context == nil {
		s.Reader = &terminal.Read{
			Terminal: s.Terminal,
			Context:  context.TODO(),
		}
	}

	vars["term"] = tui.NewTerminal(s.Terminal, s.Reader.Context)
	vars["exec"] = func(cmd string) {
		go func() {
			s.Reader.Prompt = ""
			s.Terminal.Signal.Queue <- append([]byte(cmd), 130)
			s.Reader.Prompt = s.ExecuteBrandingToStringNoErr(make(map[string]any), "prompt.tfx")
		}()
	}

	vars["prop"] = func(prop string) string {
		val, ok := propagator.PropagatedFields[prop]
		if !ok {
			return "0"
		}

		return fmt.Sprintf("%v", val)
	}

	return vars
}

// sessionObject is the swash recognized object
type sessionObject struct {
	Total func() int `swash:"total"`
	X     int        `swash:"width"`
	Y     int        `swash:"height"`
}

// Bool is a fancy print for boolean values
func (s *Session) Bool(x bool) string {
	switch x {

	case true:
		return s.ExecuteBrandingToStringNoErr(make(map[string]any), "true.tfx")
	default:
		return s.ExecuteBrandingToStringNoErr(make(map[string]any), "false.tfx")
	}
}
