package tui

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/masters/terminal"
	"bytes"
	"fmt"
	"strings"
)

type Input struct {
	read    *terminal.Read
	inherit *Button
}

// NewInput will initialize the input with a inhertied state button
func (t *Terminal) NewInput(text string, x, y int) *Input {
	input := &Input{
		read: t.term.NewRead("").ChangeMaxLen(gotable2.LenOf(text)),
		inherit: t.NewButton(x, y, text),
	}

	input.inherit.OnClick(input.onClickHandle())
	return input
}


// Value will return the value of the text field right now
func (input *Input) Value() string {
	return string(input.read.Content())
}

// ChangeValue will attempt to change the value of the text field
func (input *Input) ChangeValue(val string) {
	input.inherit.term.term.Write([]byte(fmt.Sprintf("\033[%d;%df", input.inherit.Y, input.inherit.X + len(input.Value()) + 1)))
	input.read.ChangeInput([]byte(val))
	input.inherit.buttonLabel[0] = input.Value() + strings.Repeat(" ", *input.read.MaximumBufTileSize - gotable2.LenOf(input.Value()))
}

// Tab allows you to change how the button looks when tab is pressed, when
// defining tab spaces I recommend using the same lines but adding ansi escapes
// to represent that it's selected.
func (input *Input) Tab(lines ...string) {
	sections := make([]string, 0)
	for _, text := range lines {
		sections = append(sections, strings.Split(text, "\r\n")...)
	}

	input.inherit.tabLabel = sections
}

// onClickHandle will handle the direct feedback from an input statement
func (input *Input) onClickHandle() func() bool {
	return func() bool {
		if _, err := input.read.Terminal.Write([]byte(fmt.Sprintf("\033[%d;%df\033[?25h\033[?0c", input.inherit.Y, input.inherit.X + len(input.read.Content()) + 1))); err != nil {
			return true
		}

		defer func() {
			input.read.Terminal.Write([]byte("\033[?25l"))
			input.inherit.buttonLabel[0] = input.Value() + strings.Repeat(" ",*input.read.MaximumBufTileSize - gotable2.LenOf(input.Value()))
		}()
		
		for {
			content, err := input.read.Terminal.Signal.ReadWithContext(input.inherit.term.context)
			if err != nil {
				return true
			}

			// checks for mouse clicks
			if !bytes.HasPrefix(content, []byte{27, 91, 77}) {
				content = bytes.ReplaceAll(content, []byte{13}, []byte{130})
				ok, err := input.read.Buf(content, true)
				if err != nil {
					return true
				}

				/* on Ok event we break from the loop */
				if ok {
					break
				}

				continue
			}

			val, err := input.inherit.term.handleBuf(content)
			if err != nil {
				val = true
			}

			return val
		}

		return false
	}
}

// Pop will remove it's self from the array
func (i *Input) Pop() {
	i.inherit.Pop()
}