package gotable2

import (
	"fmt"
	"strings"
)

// Align corresponds to the alignment of the text
type Align int

const (
	/* different modes */
	AlignCenter 		Align = iota
	AlignLeft
	AlignRight
)

// Pad will position the text in the position provided for the alignment.
func (align Align) Pad(size int, text string) string {
	switch align {

	case AlignLeft:
		sum := size - LenOf(text)
		if sum <= 0 {
			return text
		}


		return text + strings.Repeat(" ", sum)

	case AlignRight:
		sum := size - LenOf(text)
		if sum <= 0 {
			return text
		}

		return strings.Repeat(" ", sum) + text

	case AlignCenter:
		padSize := (size - LenOf(text))
		lft := padSize / 2
		rgt := padSize - lft
		left := strings.Repeat(" ", lft)
		right := strings.Repeat(" ", rgt)
		return fmt.Sprintf("%s%s%s", left, text, right)
	}

	return text
}

// PadCustom does the same as Pad but introduces a custom padding element
func (align Align) PadCustom(size int, text, custom string) string {
	switch align {

	case AlignLeft:
		sum := size - LenOf(text)
		if sum <= 0 {
			return text
		}


		return text + strings.Repeat(custom, sum)

	case AlignRight:
		sum := size - LenOf(text)
		if sum <= 0 {
			return text
		}

		return strings.Repeat(custom, sum) + text

	case AlignCenter:
		padSize := (size - LenOf(text))
		lft := padSize / 2
		rgt := padSize - lft
		left := strings.Repeat(custom, lft)
		right := strings.Repeat(custom, rgt)
		return fmt.Sprintf("%s%s%s", left, text, right)
	}

	return text
}