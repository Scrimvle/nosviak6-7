package tui

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// AddMarquee will concurrently introduce a marquee scroll into the button on the y axis
// AddMarquee($y: number, $duration: number, $label: string)
// this will also introduce an additional goroutine which may affect performance
func (button *Button) AddMarquee(y, duration int, text string) {
	coordinates := button.properties()
	if !slices.Contains(maps.Keys(coordinates), y) || coordinates[y][len(coordinates[y]) - 1] - coordinates[y][0] - 3 <= utf8.RuneCountInString(text) {
		return
	}

	/* sleeps for each frame */
	go func() {
		tick := time.NewTicker(time.Duration(duration) * time.Millisecond)

		/* default start pos */
		pos := button.X - utf8.RuneCountInString(text)

		/* concurrency handler, means it keeps shifting position */
		for range tick.C {
			if button.term.context.Err() != nil {
				tick.Stop()
				return
			}

			/* begins to remove the first character of the marquee of the display */
			if pos + utf8.RuneCountInString(text) + 1 > coordinates[y][len(coordinates[y]) - 1] - 1{
				if pos + utf8.RuneCountInString(text) + 1 < coordinates[y][len(coordinates[y]) - 1] + utf8.RuneCountInString(text) {
					pos++
					button.term.term.Write([]byte(fmt.Sprintf("\033[s\033[%d;%df" + strings.Repeat(" ", pos - button.X) + strings.Join(strings.Split(text, "")[:coordinates[y][len(coordinates[y]) - 1] - pos - 1], "") + "\x1b[u",  y, button.X + 2)))
					continue
				}

				pos = button.X - utf8.RuneCountInString(text)
				button.term.term.Write([]byte(fmt.Sprintf("\033[s\033[%d;%df" + strings.Repeat(" ", coordinates[y][len(coordinates[y]) - 1] - coordinates[y][0] - 1) + "\x1b[u",  y, button.X + 2)))
			}

			/* begin to bring the marquee into site */
			if pos < button.X {
				pos++
				button.term.term.Write([]byte(fmt.Sprintf("\033[s\033[%d;%df" + strings.Join(strings.Split(text, "")[button.X - pos:], "") + "\x1b[u",  y, button.X + 2)))
				if pos != button.X {
					continue
				}
			}

			pos++
			button.term.term.Write([]byte(fmt.Sprintf("\033[s\033[%d;%df" + strings.Repeat(" ", pos - button.X) + text + "\x1b[u",  y, button.X + 2)))
		}
	}()
}