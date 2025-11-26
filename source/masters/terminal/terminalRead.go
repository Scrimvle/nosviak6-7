package terminal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Message represents a terminal message
const MESSAGE      int = 10

// Read is the structure which is a child to Terminal
type Read struct {
	sync.Mutex

	Terminal *Terminal
	Prompt   string
	index    int
	read     []byte
	pos      int

	/* maximumDepth is the longest the array has been */
	maximumDepth int

	/* MaximumBufTileSize is a predefined size which represents the largest buffer size */
	MaximumBufTileSize *int

	/* MinimumBufTileSize is a predefined size which represents the smallest buffer size */
	MinimumBufTileSize *int

	/* cursorMemory is a history of all the accepted read entrys before hand */
	cursorMemory [][]byte

	/* is how we check if their is an alert being displayed */
	Alert *Alert

	/* Context allows for interaction with the signal package */
	Context context.Context

	/* KeyPressFunc is what will be triggered whenever a key is pressed and the action is performed */
	KeyPressFunc func([]byte, byte, *Read) ([]byte, bool)

	/* AutoCompleter is the object which can directly autocomplete whenever they click tab on the typist */
	AutoCompleter *AutoCompleter

	// ReaderIdle is since the last command was executed
	ReaderIdle time.Time

	// MaskChar allows us to mask the input of users in certain areas
	maskChar *[]byte
}

// AutoCompleter will maintain all the states required for the autocompletion tool
type AutoCompleter struct {
	Completer func([]byte, *Read) (bool, error)
}

// NewRead returns the new Terminal interface.
func (t *Terminal) NewRead(prompt string) *Read {
	return &Read{
		read:               make([]byte, 0),
		index:              0,
		Alert:              nil,
		Prompt:             prompt,
		Terminal:           t,
		ReaderIdle:         time.Now(),
		maximumDepth:       0,
		cursorMemory:       make([][]byte, 0),
		KeyPressFunc:       nil,
		MaximumBufTileSize: nil,
	}
}

// NewReadWithContext will implement the contextual based interface
func (t *Terminal) NewReadWithContext(prompt string, context context.Context) *Read {
	r := t.NewRead(prompt)
	r.Context = context
	return r
}

// ReadLine implements the Read interface for reading directly from the Terminal
func (r *Read) ReadLine() ([]byte, error) {
	if _, err := r.Terminal.Write([]byte(r.Prompt + string(r.read))); err != nil {
		return make([]byte, 0), err
	}

	r.pos = 0
	r.index = len(r.cursorMemory)
	r.Alert = nil
	r.Mutex.Lock()
	r.ReaderIdle = time.Now()
	defer r.Mutex.Unlock()
	r.maximumDepth = len(r.read)
	defer func() {
		r.read = make([]byte, 0)
		r.ReaderIdle = time.Now()
	}()

	for {
		if r.Terminal.Channel == nil {
			return nil, io.EOF
		}

		buf, err := r.Terminal.Signal.ReadWithContext(r.Context)
		if err != nil {
			if err == context.Canceled {
				return nil, r.Terminal.Channel.Close()
			}

			return nil, err
		}

		ok, err := r.Buf(buf, true)
		if err != nil || ok {
			return r.read, err
		}

		if r.KeyPressFunc == nil {
			continue
		}

		payload, ok := r.KeyPressFunc(r.read, buf[0], r)
		if ok {
			return payload, nil
		}
	}
}

// buf handles the case on every key press.
func (r *Read) Buf(payload []byte, init bool) (bool, error) {
	switch payload[0] {

	case 9:
		if r.AutoCompleter == nil || !init {
			return false, nil
		}

		// removes any previous rendered alerts
		r.DisgardAlert()

		return r.AutoCompleter.Completer(r.read, r)

	case 127:
		if len(r.read) <= 0 || r.maximumDepth == 0 || len(r.read[:len(r.read)-r.pos]) == 0 {
			return false, nil
		}

		r.read = append(r.read[:len(r.read)-r.pos-1], r.read[len(r.read)-r.pos:]...)
		payload := append([]byte{8}, append([]byte(strings.Repeat(" ", r.maximumDepth-(len(r.read)-r.pos)+1)+fmt.Sprintf("\x1b[%dD", r.maximumDepth-(len(r.read)-r.pos)+1)), r.read[len(r.read)-r.pos:]...)...)
		if r.pos > 0 {
			payload = append(payload, []byte(fmt.Sprintf("\033[%dD", r.pos))...)
		}

		_, err := r.Terminal.Write(payload)
		if err != nil {
			return false, err
		}

	case 13, 130:
		if r.MinimumBufTileSize != nil && len(r.read) < *r.MinimumBufTileSize {
			return false, nil
		}

		// checks for the specific use case where the new line isn't printed.
		if payload[0] == 13 {
			_, err := r.Terminal.Write([]byte("\r\n"))
			if err != nil {
				return false, err
			}
		}

		r.cursorMemory = append(r.cursorMemory, r.read)
		return true, nil

	case 27:
		if len(payload[1:]) == 0 {
			return false, nil
		}

		switch payload[1:][0] {

		case 91: // arrow keys
			if len(payload[1:]) < 2 {
				return false, nil
			}

			// 66 down , 65 up
			switch payload[1:][1] {

			case 66, 65:
				if len(r.cursorMemory) <= 0 || r.pos >= 1 {
					return false, nil
				}

				switch payload[1:][1] {

				case 66: // comes back into time
					if r.index + 1 >= len(r.cursorMemory) {
						r.ChangeInput(make([]byte, 0))
						return false, nil
					}

					r.index++
					if r.index == len(r.cursorMemory) {
						if _, err := r.ChangeInput(r.cursorMemory[r.index]); err != nil {
							return true, err
						}

						return false, nil
					}

					if _, err := r.ChangeInput(r.cursorMemory[r.index]); err != nil {
						return true, err
					}

				case 65: // sees into the future
					if r.index <= 0 {
						return false, nil
					}

					r.index--
					if _, err := r.ChangeInput(r.cursorMemory[r.index]); err != nil {
						return true, err
					}
				}

			case 68: // left arrow
				if r.pos >= len(r.read) {
					return false, nil
				}

				r.pos++
				if _, err := r.Terminal.Write([]byte("\x1b[1D")); err != nil {
					return false, err
				}

			case 67: // right arrow
				if r.pos <= 0 {
					return false, nil
				}

				r.pos--
				if _, err := r.Terminal.Write([]byte("\x1b[1C")); err != nil {
					return false, err
				}

			}
		}

	case 32, 96, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 45, 61, 194, 172, 33, 34, 163, 36, 37, 94, 38, 42, 40, 41, 95, 43, 113, 119, 101, 114, 116, 121, 117, 105, 111, 112, 91, 93, 81, 87, 69, 82, 84, 89, 85, 73, 79, 80, 123, 125, 97, 115, 100, 102, 103, 104, 106, 107, 108, 59, 39, 35, 65, 83, 68, 70, 71, 72, 74, 75, 76, 58, 64, 126, 92, 122, 120, 99, 118, 98, 110, 109, 44, 46, 47, 124, 90, 88, 67, 86, 66, 78, 77, 60, 62, 63:
		if r.MaximumBufTileSize != nil && *r.MaximumBufTileSize <= len(r.read) || len(payload) > 1 {
			for pos := 0; len(payload) > 1 && pos < len(payload); pos++ {
				ok, err := r.Buf([]byte{payload[pos]}, init)
				if err != nil || ok {
					return true, nil
				}
			}

			return false, nil
		}

		if r.pos <= 0 {
			r.read = append(r.read, payload...)
		} else {
			r.read = append(r.read[:len(r.read)-r.pos], append(payload, r.read[len(r.read)-r.pos:]...)...)
		}

		// uses our masking char instead
		if r.maskChar != nil && len(*r.maskChar) >= 1 {
			payload = *r.maskChar
		}

		payload := append(payload, r.read[len(r.read)-r.pos:]...)
		if r.pos > 0 {
			payload = append(payload, []byte(fmt.Sprintf("\033[%dD", r.pos))...)
		}

		/* writes the payload with the ansi escapes implemented */
		if _, err := r.Terminal.Write(payload); err != nil {
			return false, nil
		}

		if len(r.read) > r.maximumDepth {
			r.maximumDepth = len(r.read)
		}
	}

	return false, nil
}

// Alert is how we can maintain the alerts being displayed
type Alert struct {
	AlertCode    int
	AlertMessage string
}

// PostAlert will post the alert onto the Terminals channel
func (r *Read) PostAlert(a *Alert) error {
	payload := append(append(r.deskinPrompt(), "\r\x1b[K"+a.AlertMessage+"\r\n"...), []byte(r.Prompt+string(r.read))...)
	if r.pos > 0 {
		payload = append(payload, []byte(fmt.Sprintf("\033[%dD", r.pos))...)
	}

	if _, err := r.Terminal.Write(payload); err != nil {
		return err
	}

	r.Alert = a
	return nil
}

// DisgardAlert will remove any alert which is being displayed
func (r *Read) DisgardAlert() error {
	if r.Alert == nil {
		return errors.New("no alerts being displayed")
	}

	payload := append(append(r.deskinPrompt(), "\x1b[1A\x1b[K"...), []byte(r.Prompt+string(r.read))...)
	if r.pos > 0 {
		payload = append(payload, []byte(fmt.Sprintf("\033[%dD", r.pos))...)
	}

	if _, err := r.Terminal.Write(payload); err != nil {
		return err
	}

	r.Alert = nil
	return nil
}

// ChangeMaxLen will modify the maximum length of the prompt
func (r *Read) ChangeMaxLen(i int) *Read {
	r.MaximumBufTileSize = &i
	return r
}

// ChangeMinLen will modify the minimum length of the prompt
func (r *Read) ChangeMinLen(i int) *Read {
	r.MinimumBufTileSize = &i
	return r
}

// Mask will toggle the masking function on the reader
func (r *Read) Mask(mask []byte) *Read {
	r.maskChar = &mask
	return r
}

// ChangeInput will modify the current prompt value on the Terminal
func (r *Read) ChangeInput(callback []byte) (int, error) {
	if r.maximumDepth <= 0 {
		r.read = callback
		r.maximumDepth = len(callback)
		return r.Terminal.Write(callback)
	}

	payload := make([]byte, 0)
	if len(r.read) >= 1 {
		payload = append(payload, fmt.Sprintf("\x1b[%dD", len(r.read))...)
	}

	payload = append(payload, append([]byte(fmt.Sprintf(strings.Repeat(" ", r.maximumDepth)+"\x1b[%dD", r.maximumDepth)), callback...)...)
	if r.pos > 0 {
		payload = append([]byte(fmt.Sprintf("\x1b[%dD", r.pos)), payload...)
	}

	/* rebuilds the entire Terminal env */
	r.pos = 0
	r.read = callback
	r.maximumDepth = len(r.read) + 1
	return r.Terminal.Write(payload)
}

// deskinPrompt will return the bytes in ansi sequences to remove the prompt from the Terminal entirely.
func (r *Read) deskinPrompt() []byte {
	promptLines, payload := strings.Split(r.Prompt, "\n"), make([]byte, 0)
	sort.Sort(sort.Reverse(sort.StringSlice(promptLines)))
	for linePos := range promptLines {
		if linePos == 0 {
			payload = append(payload, "\r\x1b[K"...)
			continue
		}

		payload = append(payload, "\x1b[1A\x1b[K"...)
	}

	return payload
}

// Reskin will remove the entire prompt
func (r *Read) Reskin(prompt string) (int, error) {
	payload := append(r.deskinPrompt(), append([]byte(prompt), r.read...)...)
	if r.pos > 0 {
		payload = append(payload, []byte(fmt.Sprintf("\x1b[%dD", r.pos))...)
	}

	r.Prompt = prompt
	return r.Terminal.Write(payload)
}

// Content will return the current bytes of the reader
func (r *Read) Content() []byte {
	return r.read
}

// requestWindowSize will return the current size of the window
func (t *Terminal) RequestWindowSize() error {
	if _, err := t.Channel.Write([]byte("\x1b[18t")); err != nil {
		return err
	}

	/* implements a timeout so if the Terminal doesn't return the correct thing we break */
	ctx, ok := context.WithTimeout(context.TODO(), 50 * time.Millisecond)
	defer ok()

	data, err := t.Signal.ReadWithContext(ctx)
	if err != nil && err != context.DeadlineExceeded {
		return io.EOF
	}

	interested := strings.Split(strings.ReplaceAll(string(data), "t", ""), ";")[1:]
	if interested == nil || len(interested) != 2 {
		return nil
	}

	height, err := strconv.Atoi(interested[0])
	if err != nil {
		return err
	}

	width, err := strconv.Atoi(interested[1])
	if err != nil {
		return err
	}

	t.XTerm = true
	t.X, t.Y = uint32(width), uint32(height)
	return nil
}

// RequestCursorSize will return the current position of the cursor
func (t *Terminal) RequestCursorSize() (int, int, error) {
	t.XTerm = false
	if _, err := t.Channel.Write([]byte("\x1b[6n")); err != nil {
		return 0, 0, err
	}

	/* implements a timeout so if the Terminal doesn't return the correct thing we break */
	ctx, ok := context.WithTimeout(context.TODO(), 1 * time.Second)
	defer ok()

	data, err := t.Signal.ReadWithContext(ctx)
	if err != nil && err != context.DeadlineExceeded {
		return 0, 0, io.EOF
	}

	if len(strings.Split(strings.ReplaceAll(string(data), "t", ""), "[")) <= 1 {
		return 0, 0, nil
	}

	interested := strings.Split(strings.Split(strings.ReplaceAll(string(data), "t", ""), "[")[1], ";")
	if interested == nil {
		return 0, 0, nil
	}

	height, err := strconv.Atoi(interested[0])
	if err != nil {
		return 0, 0, err
	}

	width, err := strconv.Atoi(interested[1][:len(interested[1]) - 1])
	if err != nil {
		return 0, 0, err
	}

	t.XTerm = true
	return width, height, nil
}
