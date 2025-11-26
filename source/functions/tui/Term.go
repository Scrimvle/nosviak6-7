package tui

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/masters/terminal"
	"context"
	"fmt"
	"strings"
)

// Terminal maintains every terminal element on the system
type Terminal struct {
	entities []any
	context  context.Context
	cancel   context.CancelFunc
	term     *terminal.Terminal
	update   chan []any
	tabPos   int
}

// NewTerminal will create a new terminal interface
func NewTerminal(terminal *terminal.Terminal, ctx context.Context) *Terminal {
	ctx, cancel := context.WithCancel(ctx)
	return &Terminal{
		term: terminal,
		cancel: cancel,
		context: ctx,
		entities: make([]any, 0),
		update: make(chan []any),
	}
}

// Resize will change the dimensions of the terminal to the given parameters
func (t *Terminal) Resize(x, y int) {
	t.term.Write([]byte(fmt.Sprintf("\033[8;%d;%dt", y, x)))
	t.term.X, t.term.Y = uint32(x), uint32(y)
}

// Update will push the new entities stack to render the update
func (t *Terminal) Update() {
	t.Draw()
}

// WordWrap imports the word wrapping process
func WordWrap(text, prefix string, line int) string {
	var w, s string
	w = prefix
	
	for _, x := range strings.Fields(text) {
		if gotable2.LenOf(s) + gotable2.LenOf(x) + 1 > line {
			if gotable2.LenOf(s) < line {
				s += strings.Repeat(" ", line - gotable2.LenOf(s))
			}

			w += s + "\r\n"
			s = prefix
		}

		if gotable2.LenOf(s) >= 1 {
			s += " "
		}

		s += x
	}

	if gotable2.LenOf(s) < line {
		s += strings.Repeat(" ", line - gotable2.LenOf(s))
	}

	w += s
	return w
}