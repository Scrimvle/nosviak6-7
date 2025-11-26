package gotable2

import (
	"Nosviak4/source/swash"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

// ansi is the regexp string for stripping away them codes
var ansi = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

// lenOf will return the length of v
func LenOf(v string) int {
	return runewidth.StringWidth(strings.ReplaceAll(ansi.ReplaceAllString(v, ""), "<escape>", ""))
}

// String will build the GoTable into a table representation
func (gotable *GoTable) String(dest []string) []string {
	dest = append(dest, gotable.headerString()...)
	dest = append(dest, gotable.valueString(make([]string, 0))...)
	for _, line := range dest {
		length := LenOf(line)
		if length > gotable.LongestLine {
			gotable.LongestLine = length
		}
	}

	for pos := range dest {
		dest[pos] = swash.Strip(dest[pos])
	}

	return dest
}

// stringPointer converts the dest into the string representation pointer
func stringPointer(dest string) *string {
	return &dest
}