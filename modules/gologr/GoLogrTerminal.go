package gologr

import (
	"bytes"
	"fmt"
	"time"
)

// Terminal directly links with the FileLogger object
type Terminal struct {
	logger *FileLogger
}

// WriteLog will write to the file, and then the aggregated file
func (t *Terminal) WriteLog(mode *TerminalMode, message string, args ...any) {
	if mode.ModeFunc != nil && !mode.ModeFunc() {
		return
	}
	
	if err := t.logger.WriteLog(message, args...); err != nil {
		return
	}

	if mode.Modify != nil {
		message = mode.Modify(message)
	}

	/* prepares formatting for the message with the mode defined params */
	payload := bytes.NewBufferString(fmt.Sprintf(mode.TimestampLeft+"%s"+mode.TimestampRight, time.Now().Format(mode.TimestampFormat)))
	payload.WriteString(fmt.Sprintf(mode.BodyLeft+"%s"+mode.BodyRight, fmt.Sprintf(message, args...)))
	if _, err := t.logger.origin.writer.Write(payload.Bytes()); err != nil {
		return
	}
}
