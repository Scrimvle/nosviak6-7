package terminal

import (
	"strings"
)

// Write wraps the standard ssh.Channel.Write function with the screen operation
func (t *Terminal) Write(b []byte) (int, error) {
	t.writeLock.Lock()
	defer t.writeLock.Unlock()

	if !strings.HasPrefix(string(b), "\033]0;") {
		content, err := t.Screen.Write(b)
		if err != nil {
			return content, err
		}
	}

	if t.Channel == nil {
		return 0, nil
	}

	return t.Channel.Write(b)
}

// Clear will wipe the entire buffer for screen
func (t *Terminal) Clear() (int, error) {
	t.Screen.Truncate(0)
	if t.Channel == nil {
		return 0, nil
	}

	return t.Write([]byte("\x1bc"))
}

// ClearString will wipe the entire buffer for the screen and then return the escape code
func (t *Terminal) ClearString() string {
	t.Screen.Truncate(0)
	return "\x1bc"
}
