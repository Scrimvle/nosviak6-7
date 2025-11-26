package sessions

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// IndexSessions will attempt to retrieve the user from the sessions
func IndexSessions(index string) []*Session {
	destination := make([]*Session, 0)
	for _, s := range Sessions {
		if strings.Split(index, ":")[0] != s.User.Username {
			continue
		}

		if !strings.Contains(index, ":") || strings.Split(index, ":")[1] == fmt.Sprint(s.Opened) {
			destination = append(destination, s)
		}
	}

	return destination
}

// WriteToSession offers a way to concurrently write to all the sessions
func WriteToSession(sessions []*Session, messager func(*Session)) error {
	for _, session := range sessions {
		if session == nil {
			continue
		}

		go messager(session)
	}

	return nil
}

// Included will check if our sessions are contained within the array
func (session *Session) Included(sessions []*Session) bool {
	for _, s := range sessions {
		if s.User.Username == session.User.Username {
			return true
		}

		continue
	}
	return false
}

// Callback will return a list of all the sessions
func (session *Session) Callback() []string {
	dest := make([]string, 0)
	for _, indexSession := range Sessions {
		if indexSession.User.Username == session.User.Username {
			continue
		}

		dest = append(dest, fmt.Sprintf("%s:%d", indexSession.User.Username, indexSession.Opened))
		if slices.Contains(dest, indexSession.User.Username) {
			dest = append(dest, indexSession.User.Username)
		}
	}

	return dest
}