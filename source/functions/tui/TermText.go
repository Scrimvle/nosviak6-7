package tui

import (
	"Nosviak4/source/swash/packages"
)

type Text struct {
	Text      string
	X, Y      int
	t         *Terminal
}

// NewText will implement the Text interface to the tui package
func (t *Terminal) NewText(text string, x, y int) *Text {
	configure := &Text{
		Text: text, X: x, Y: y, t: t,
	}

	t.entities = append(t.entities, configure)
	return configure
}

// draw will render the text inside the matrix
func (t *Text) draw(matrix [][]string, term *Terminal) {
	for charPos, token := range packages.Split(t.Text) {
		if t.Y >= len(matrix) || t.X+charPos >= len(matrix[t.Y]) {
			break
		}

		matrix[t.Y][t.X:][charPos] = token
	}
}

// ChangeText will modify the content being wrote
func (t *Text) ChangeText(text string) {
	t.Text = text
}

// Pop will remove it's self from the array
func (t *Text) Pop() {
	for p, entity := range t.t.entities {
		if entity != t {
			continue
		}


		t.t.entities = append(t.t.entities[:p], t.t.entities[p+1:]...)
		break
	}
}
