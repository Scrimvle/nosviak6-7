package sessions

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source"
	"Nosviak4/source/swash/packages"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

// Page will page the results to the session
func (s *Session) Page(content []string) error {
	if len(content) < int(s.Terminal.Y) {
		return nil
	}

	content = append([]string{s.ExecuteBrandingToStringNoErr(make(map[string]any), "pager_top.tfx")}, content...)
	content = append(content, s.ExecuteBrandingToStringNoErr(make(map[string]any), "pager_bottom.tfx"))

	// defines information about rendering
	axis, payload := make([][]string, 0), bytes.NewBuffer(make([]byte, 0))

	// defines information about current render scope
	y, x, mx := 0, 0, 0

	payload.Write([]byte("\x1bc"))
	for rpos, line := range content {
		axis = append(axis, packages.Split(strings.ReplaceAll(line, packages.EscapeCode, "")))
		if i := len(axis[len(axis)-1]); i > mx {
			mx = i
		}

		// renders it at the same time
		if rpos < int(s.Terminal.Y) {
			line := axis[len(axis)-1]
			if len(axis[len(axis)-1]) > int(s.Terminal.X) {
				line = line[:int(s.Terminal.X)]
			}

			payload.Write([]byte(strings.Join(line, "") + "\x1b[0m"))
			if rpos+1 < int(s.Terminal.Y) {
				payload.Write([]byte("\r\n"))
			}
		}
	}

	if _, err := s.Terminal.Channel.Write(payload.Bytes()); err != nil {
		return err
	}

	for {
		data, err := s.Terminal.Signal.ReadWithContext(s.Reader.Context)
		if err != nil || len(data) == 0 {
			break
		}

		if s.Terminal.Y > uint32(len(axis)) || data[0] == 113 || data[0] == 81 {
			_, err := s.Terminal.Write(s.Terminal.Screen.Bytes())
			return err
		}

		switch data[0] {
			
		case 87, 119: // up
			if y-1 < 0 || y-1 == 0 && x != 0 {
				continue
			}

			y--

		case 83, 115: // down
			if y+1 > len(axis)-int(s.Terminal.Y) || y+1 == len(axis)-int(s.Terminal.Y) && x != 0 {
				continue
			}

			y++

		case 65, 97:
			if x-1 < 0 {
				continue
			}

			x--

		case 68, 100:
			if x+1 > mx-int(s.Terminal.X) || y <= 0 || y == len(axis)-int(s.Terminal.Y) {
				continue
			}

			x++

		default:
			continue
		}

		buf := bytes.NewBuffer(make([]byte, 0))
		buf.Write([]byte("\033[0;0f\033[J"))

		for p, rbuf := range axis[y:][:s.Terminal.Y] {
			if p >= int(s.Terminal.Y) {
				continue
			}

			// we add some blank space chars
			if len(rbuf) < mx {
				rbuf = append(rbuf, strings.Split(strings.Repeat(" ", mx-len(rbuf)-1), " ")...)
			}

			buf.Write([]byte("\x1b[38;5;15m" + strings.Join(rbuf[x:][:s.Terminal.X], "") + "\x1b[0m"))
			if p+1 < int(s.Terminal.Y) {
				buf.Write([]byte("\r\n"))
			}
		}

		if _, err := s.Terminal.Channel.Write(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

// Table will use the gotable system to introduce coloured tables etc
func (s *Session) Table(table *gotable2.GoTable, name string) error {
	content, ok := source.OPTIONS.Config.Renders[filepath.Join(source.ASSETS, source.COMMANDS, "configs", name+".toml")]
	if !ok || content == nil {
		return fmt.Errorf("missing table config: %s", name+".toml")
	}

	dest := make(map[string]any)
	err := toml.Unmarshal(content, &dest)
	if err != nil {
		return err
	}

	style, ok := source.TABLEVIEWS[fmt.Sprint(dest["name"])]
	if !ok || style == nil {
		return fmt.Errorf("missing table style: %v", dest["name"])
	}

	table.SetStyle(style)
	gradient := packages.NewDerivative()
	for _, i := range s.Theme.Glamour.Colours {
		gradient.AppendRgbToGradient(i[0], i[1], i[2])
	}

	if !s.Theme.Glamour.Enabled {
		gradient.ResetRGB()
		gradient.AppendRgbToGradient(252, 252, 252)
		gradient.AppendRgbToGradient(252, 252, 252)
	}

	// gets the gradient content
	render := table.String(make([]string, 0))
	if len(render) > int(s.Terminal.Y) {
		return s.Page(render)
	}

	// performs the curve steps
	rgb := gradient.Curve(table.LongestLine)

	// apply's in an async manor
	var wg sync.WaitGroup
	for p := range render {
		wg.Add(1)

		go func(pos int) {
			render[pos] = gradient.Apply(render[pos], rgb, packages.Foreground, int(s.Terminal.X)-1)
			wg.Done()
		}(p)
	}

	wg.Wait()

	if err := s.ExecuteBranding(map[string]any{"table": name}, "before_table.tfx"); err != nil {
		return err
	}

	for _, table := range render {
		_, err := s.Terminal.Write([]byte(table + "\x1b[0m\r\n"))
		if err != nil {
			return err
		}
	}

	return s.ExecuteBranding(map[string]any{"table": name}, "after_table.tfx")
}
