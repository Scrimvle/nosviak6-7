package tui

import (
	"Nosviak4/source/swash/packages"
	"fmt"
	"strings"
)

// Button is an object that represents a clickable object
type Button struct {
	term        *Terminal
	buttonLabel []string
	X           int `swash:"x"`
	Y           int `swash:"y"`
	onClick     func() bool
	tabLabel    []string
	pos         int
}

// NewButton will create a new Button object
func (t *Terminal) NewButton(X, Y int, labels ...string) *Button {
	sections := make([]string, 0)
	for _, text := range labels {
		sections = append(sections, strings.Split(text, "\r\n")...)
	}

	el := &Button{
		term: t,
		buttonLabel: sections,
		X: X, Y: Y,
		pos: len(t.entities),
	}

	t.entities = append(t.entities, el)
	return el
}

// draw will start to draw the button to the terminal
func (button *Button) draw(matrix [][]string, t *Terminal) {
	for pos, data := range button.buttonLabel {
		if button.Y + pos > int(t.term.Y) {
			break
		}

		for charPos, token := range packages.Split(data) {
			if button.X  + charPos >= len(matrix[button.Y]) {
				continue
			}

			matrix[button.Y + pos][button.X:][charPos] = token
		}
	}
}

// properties are all the defined existing positions
func (button *Button) properties() map[int][]int {
	placeholder := make(map[int][]int)
	for pos, text := range button.buttonLabel {
		placeholder[button.Y + pos - 1] = make([]int, 0)
		for posHori := range packages.Split(text) {
			placeholder[button.Y + pos - 1] = append(placeholder[button.Y+pos - 1], button.X + posHori)
		}

	}

	return placeholder
}

// click is what is executed when the button is clicked
func (button *Button) click() bool {
	if button.onClick == nil {
		return false
	}

	return button.onClick()
}

// OnClick synchronizes with the click event
func (button *Button) OnClick(fn func() bool) {
	button.onClick = fn
}

// Tab allows you to change how the button looks when tab is pressed, when
// defining tab spaces I recommend using the same lines but adding ansi escapes
// to represent that it's selected.
func (button *Button) Tab(lines ...string) {
	sections := make([]string, 0)
	for _, text := range lines {
		sections = append(sections, strings.Split(text, "\r\n")...)
	}

	button.tabLabel = sections
}

// tabImage executes the tab query on the screen
func (button *Button) tabImage() {
	for pos, line := range button.tabLabel {
		button.term.term.Write([]byte(fmt.Sprintf("\x1b[s\033[%d;%dH%s\x1b[u\x1b[0m", pos + button.Y, button.X + 1, line)))
	}
}

// Pop will remove it's self from the array
func (b *Button) Pop() {
	for p, entity := range b.term.entities {
		if entity != b {
			continue
		}


		b.term.entities = append(b.term.entities[:p], b.term.entities[p+1:]...)
		break
	}
}