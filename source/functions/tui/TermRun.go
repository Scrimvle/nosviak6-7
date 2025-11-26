package tui

import (
	"io"
	"strings"
)

// Run is a way to execute the stuff without an error occurring
func (t *Terminal) Run() {
	t.RunWithErr()
}

// Run will attempt to execute the given ui params
func (t *Terminal) RunWithErr() error {
	if err := t.term.RequestWindowSize(); err != nil {
		return err
	}

	err := t.Draw()
	if err != nil {
		return err
	}

	defer func() {
		t.term.Write([]byte("\x1b[?1000l"))
		t.cancel()
	}()

	if _, err := t.term.Write([]byte("\033[?1000h\033[?25l")); err != nil {
		return err
	}

	for {
		select {

		case <-t.context.Done():
			return t.context.Err()

		case buf, ok := <-t.term.Signal.Queue:
			if !ok || buf == nil {
				return io.EOF
			}

			ok, err := t.handleBuf(buf)
			if err != nil || ok {
				return err
			}
		}
	}
}

// Draw will rerender the entire stack
func (t *Terminal) Draw() error {
	matrix := make([][]string, t.term.Y+1)
	for pos := range matrix {
		matrix[pos] = strings.Split(strings.Repeat(" ", int(t.term.X)), "")
	}

	for _, el := range t.entities {
		switch index := el.(type) {

		case *Button:
			index.draw(matrix, t)

		case *Input:
			index.inherit.draw(matrix, t)

		case *Text:
			index.draw(matrix, t)
		}
	}

	// clears the terminal screen
	if _, err := t.term.Write([]byte("\x1bc\033[?1000h\033[?25l")); err != nil {
		return err
	}

	// loops through the matrix to concurrently render
	for p, conc := range matrix {
		_, err := t.term.Write([]byte(strings.Join(conc, "")))
		if err != nil {
			return err
		}

		if p >= len(matrix) {
			break
		}

		if p+1 < len(matrix) {
			_, err := t.term.Write([]byte("\r\n"))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// handleBuf will directly handle any incoming bytes from the channel
func (t *Terminal) handleBuf(buf []byte) (bool, error) {
	if buf[0] == 9 {
		return t.tab(), nil
	}

	if len(buf) == 6 && buf[0] == 27 && buf[1] == 91 && buf[2] == 77 {
		switch buf[3] {

		case 35, 32:
			button, ok := t.ClickQuestion(int(buf[4:][0]-33), int(buf[4:][1]-33))
			if !ok {
				return false, nil
			}

			if _, ok := button.(*Button); ok && buf[3] == 32 {
				return false, nil
			}

			if ok := button.click(); ok {
				return ok, nil
			}
		}
	}

	return t.checkListener(buf), nil
}

// tab handles the key press for tab actions
func (t *Terminal) tab() bool {
	if err := t.Draw(); err != nil {
		return true
	}

	for t.tabPos <= len(t.entities) {
		t.tabPos++
		if t.tabPos >= len(t.entities) {
			t.tabPos = 0
		}

		// checks the type of what has been selected
		switch entity := t.entities[t.tabPos].(type) {

		/* unsupported tab case */
		default:
			continue

		case *Button:
			if len(entity.tabLabel) == 0 {
				continue
			}

			entity.tabImage()
			buf, err := t.term.Signal.ReadWithContext(t.context)
			if err != nil || buf == nil {
				return true
			}

			// checks for a key press which isn't enter
			if buf[0] != 13 {
				ok, err := t.handleBuf(buf)
				return ok || err != nil
			}

			return entity.click()

		case *Input:
			if len(entity.inherit.tabLabel) == 0 {
				continue
			}

			entity.inherit.tabImage()
			buf, err := t.term.Signal.ReadWithContext(t.context)
			if err != nil || buf == nil {
				return true
			}

			// checks for a key press which isn't enter
			if buf[0] == 9 || buf[0] == 27 {
				ok, err := t.handleBuf(buf)
				return ok || err != nil
			}

			return entity.onClickHandle()()

		}
	}

	return false
}
