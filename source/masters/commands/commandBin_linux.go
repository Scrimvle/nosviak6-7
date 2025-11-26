//go:build linux

package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/masters/sessions"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

// executeBinCommand will execute the bin command
func executeBinCommand(bin *source.Bin, session *sessions.Session) error {

	env := make([]string, 0)
	for _, line := range bin.Env {
		line, err := session.Terminal.ExecuteStringToString(session.AppendDefaultSession(make(map[string]any)), line)
		if err != nil {
			return err
		}

		env = append(env, line)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, strings.Split(bin.Runtime, " ")[0], strings.Split(bin.Runtime, " ")[1:]...)
	cmd.Env = append(os.Environ(), env...)

	tty, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	defer tty.Close()
	go func() {
		io.Copy(session.Terminal, tty)
	}()

	// handles our writing
	go func() {
		for {
			select {

			// reading from the terminal
			case buf, ok := <-session.Terminal.Signal.Queue:
				if !ok || buf == nil {
					cancel()
				}

				// writes to the file descriptor
				if _, err := tty.Write(buf); err != nil {
					session.Reader.Buf(buf, false)
					cancel()
					return
				}

			// broken from the context
			case <-ctx.Done():
				return
			}
		}
	}()

	return cmd.Wait()
}
