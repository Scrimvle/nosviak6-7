package tui

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Shake will grab the terms position and move it up x and y thresholds 
// and then return to the default pos, each rep takes 50 milliseconds and
// equals to movements diagonal up and down.
func (t *Terminal) Shake(reps, size int) {
	if _, err := t.term.Write([]byte("\x1b[13t")); err != nil {
		return
	}

	time.Sleep(50 * time.Millisecond)
	payload, err := t.term.Signal.ReadWithContext(t.context)
	if err != nil || bytes.Count(payload, []byte(";")) != 2 {
		return 
	}

	queries := bytes.Split(payload, []byte{59})
	x, err := strconv.Atoi(string(queries[1]))
	if err != nil {
		return
	}

	y, err := strconv.Atoi(string(queries[2][:len(queries)]))
	if err != nil {
		return
	}

	// controls our offsets
	y -= 31
	x -= 8

	for i := 0; i < reps; i++ {
		if _, err := t.term.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", x + size, y + size))); err != nil {
			break
		}

		time.Sleep(25 * time.Millisecond)
		if _, err := t.term.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", x - size, y - size))); err != nil {
			break
		}

		time.Sleep(25 * time.Millisecond)
	}
}