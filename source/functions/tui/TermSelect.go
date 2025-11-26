package tui

import (
	"Nosviak4/modules/gotable2"
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// DropDown is a drop down helper function which on clicks brings a bunch of references down into display
type DropDown struct {
	inherited *Input
	options   []string
	onClick   func() bool

	/* relates to the current cursor pos */
	cursor    int
}

// NewDropDown will deploy a brand new DropDown feature set
// spawnx relates to the point which spawns the dropdown box from.
func (t *Terminal) NewDropdown(label string, x, y int) *DropDown {
	dropDown := &DropDown{
		inherited: t.NewInput(label, x, y),
		options: make([]string, 0),
	}

	dropDown.inherited.inherit.OnClick(dropDown.onClickHandle(t))
	return dropDown
}

// Append will drop the item into the list of items
func (dropDown *DropDown) Append(item string) {
	if gotable2.LenOf(item) + 2 > *dropDown.inherited.read.MaximumBufTileSize {
		dropDown.inherited.read.ChangeMaxLen(gotable2.LenOf(item))
	}
	
	/* if it's our first value, we remove */
	if len(dropDown.options) == 0 {
		dropDown.inherited.ChangeValue(item)
	}

	dropDown.options = append(dropDown.options, item)
}

// Value will return the current selected method
func (dropDown *DropDown) Value() string {
	return dropDown.options[dropDown.cursor]
}

// Clear will remove all the items from the list
func (DropDown *DropDown) Clear() {
	DropDown.options = make([]string, 0)
}

// OnSelect triggers on the selection
func (dropDown *DropDown) OnSelect(fn func() bool) {
	dropDown.onClick = fn
}

// Pop will remove it's self from the array
func (dropDown *DropDown) Pop() {
	dropDown.inherited.inherit.Pop()
}

// onClickHandle will directly interact within the button and help us show the menu
func (dropDown *DropDown) onClickHandle(t *Terminal) func() bool {
	return func() bool {
		if err := dropDown.render(0); err != nil {
			return err != nil
		}

		// x stores all the required valid x coordinates
		x := dropDown.inherited.inherit.properties()[dropDown.inherited.inherit.Y - 1]
		y := make([]int, 0)
		for i := dropDown.inherited.inherit.Y; i < dropDown.inherited.inherit.Y + len(dropDown.options); i++ {
			y = append(y, i)
		}

		for {
			content, err := t.term.Signal.ReadWithContext(t.context)
			if err != nil {
				return err != nil
			}

			// only handles mouse clicks from this point down
			if !bytes.HasPrefix(content, []byte{27, 91, 77, 32}) {
				continue
			}

			// click within the boundaries
			if slices.Contains(y, int(content[4:][1] - 33)) && slices.Contains(x, int(content[4:][0] - 33)) {
				if dropDown.cursor == (int(content[4:][1] - 33) - dropDown.inherited.inherit.Y) {
					if err := t.Draw(); err != nil {
						return true
					}

					if dropDown.onClick == nil {
						return false
					}

					return dropDown.onClick()
				}

				dropDown.render(int(content[4:][1] - 33) - dropDown.inherited.inherit.Y)
				continue
			}
			
			/* any error happens we return false */
			if err := t.Draw(); err != nil {
				return err != nil
			}

			ok, _ := t.handleBuf(content)
			return ok
		}
	}
}

// render will print all the options available to the terminal, often resulting in different outputs depending on the cursor pos
func (dropDown *DropDown) render(cursorPos int) error {
	dropDown.cursor = cursorPos
	for pos, item := range dropDown.options {
		y := pos + dropDown.inherited.inherit.Y + 1
		filler := fmt.Sprintf(" %s", item)
		filler += strings.Repeat(" ", *dropDown.inherited.read.MaximumBufTileSize - gotable2.LenOf(filler))

		if pos == cursorPos {
			dropDown.inherited.read.Terminal.Write([]byte(fmt.Sprintf("\033[%d;%df%s" + filler + "\x1b[0m", y, dropDown.inherited.inherit.X + 1, "\x1b[48;5;235;38;5;15m")))
			continue
		}	

		dropDown.inherited.read.Terminal.Write([]byte(fmt.Sprintf("\033[%d;%df%s" + filler + "\x1b[0m", y, dropDown.inherited.inherit.X + 1, "\x1b[48;5;255;38;5;16m")))
	}

	dropDown.inherited.ChangeValue(dropDown.options[cursorPos])
	return nil
}