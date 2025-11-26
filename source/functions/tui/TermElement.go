package tui

import (

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Element represents a field which can support clicks
type Element interface {
	properties() map[int][]int
	click() bool
}

// ClickQuestion will decide whether the click occurred on a button or not
func (t *Terminal) ClickQuestion(x, y int) (Element, bool) {
	for _, entity := range t.entities {
		switch index := entity.(type) {

		case Element:
			coordinates := index.properties()
			if !slices.Contains(maps.Keys(coordinates), y) || !slices.Contains(coordinates[y], x) {
				continue
			}

			return index, true
		}
	}

	return nil, false
}