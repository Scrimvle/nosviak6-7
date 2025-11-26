package sessions

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	components "Nosviak4/source/functions"
	"Nosviak4/source/masters/terminal"
	"context"
	"sync"

	"golang.org/x/exp/maps"
)

var (
	Sessions map[int]*Session = make(map[int]*Session)
	locker   sync.Mutex
)

// Sessions are signed versions of the terminal object
type Session struct {
	Opened   int
	Theme    *source.Theme
	User     *database.User
	Reader   *terminal.Read
	Terminal *terminal.Terminal
	Cancel   context.CancelFunc
}

// NewSession will implement the logic for the session signing up.
func NewSession(terminal *terminal.Terminal, user *database.User) *Session {
	session := &Session{
		Opened:   id(),
		Terminal: terminal,
		User:     user,
	}

	index, err := source.GetTheme(user.Theme)
	if err == nil {
		session.Theme = index
	}

	locker.Lock()
	defer locker.Unlock()
	Sessions[session.Opened] = session

	// Whenever the conn is closed, we print this.
	go func(s *Session) {
		s.Terminal.Conn.Wait()
		delete(Sessions, s.Opened)
	}(session)

	return Sessions[session.Opened]
}

// ConnIP is a duplicate of the (terminal).ConnIP but introduces ip-rewrite
func (s *Session) ConnIP() string {
	return components.RewriteIP(s.Terminal.Addr.String(), s.User)
}

// ids will look for a open id within a linear search
func id() int {
	ids := maps.Keys(Sessions)
	if i := len(ids); i == 0 {
		return i
	}

	for scan := 0; scan < ids[len(ids)-1]; scan++ {
		if _, ok := Sessions[scan]; ok {
			continue
		}

		return scan
	}

	return len(ids)
}

// Broadcast is an implementation of reader shift dialect move (also known as rsdm) but for all sessions
func Broadcast(message string) {
	for _, session := range Sessions {
		session.Reader.PostAlert(&terminal.Alert{
			AlertMessage: message,
			AlertCode:    terminal.MESSAGE,
		})
	}
}
