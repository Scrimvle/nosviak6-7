package gologr

import (
	"fmt"
	"runtime"
	"time"
)

var DEBUGENABLED = false

// TerminalMode allows for different formatting depending on the mode
type TerminalMode struct {
	*TerminalModeTimestamp

	BodyLeft  string
	BodyRight string
	ModeFunc  TerminalModeFunc
	Modify    TerminalModeFuncMessageModifier
}

// TerminalModeFunc is a decider on whether
type TerminalModeFunc func() bool

// TerminalModeFuncMessageModifier allows you to modify the string depending on certain cases
type TerminalModeFuncMessageModifier func(string) string

// TerminalModeTimestamp allows for different formatting depending on the timestamp
type TerminalModeTimestamp struct {
	TimestampLeft   string
	TimestampRight  string
	TimestampFormat string
}

var (
	DEFAULT *TerminalMode = &TerminalMode{
		TerminalModeTimestamp: &TerminalModeTimestamp{
			TimestampLeft:   "\x1b[38;5;16;48;5;10m ",
			TimestampRight:  " \x1b[0m",
			TimestampFormat: time.RFC3339,
		},

		BodyLeft:  " ",
		BodyRight: "\r\n",
	}

	ERROR *TerminalMode = &TerminalMode{
		TerminalModeTimestamp: &TerminalModeTimestamp{
			TimestampLeft:   "\x1b[38;5;16;48;5;9m ",
			TimestampRight:  " \x1b[0m",
			TimestampFormat: time.RFC3339,
		},

		BodyLeft:  " ",
		BodyRight: "\r\n",
		Modify: func(s string) string {
			pc, _, _, ok := runtime.Caller(2)
			details := runtime.FuncForPC(pc)
			if ok && details != nil && DEBUGENABLED {
				return fmt.Sprintf("%s: %s", details.Name(), s)
			}

			return s
		},
	}

	DEBUG *TerminalMode = &TerminalMode{
		TerminalModeTimestamp: &TerminalModeTimestamp{
			TimestampLeft:   "\x1b[38;5;16;48;5;11m ",
			TimestampRight:  " \x1b[0m",
			TimestampFormat: time.RFC3339,
		},

		BodyLeft:  " ",
		BodyRight: "\r\n",
		Modify: func(s string) string {
			pc, _, _, ok := runtime.Caller(2)
			details := runtime.FuncForPC(pc)
			if ok && details != nil {
				return fmt.Sprintf("%s: %s", details.Name(), s)
			}

			return s
		},

		ModeFunc: func() bool {
			return DEBUGENABLED
		},
	}
	
	ALERT *TerminalMode = &TerminalMode{
		TerminalModeTimestamp: &TerminalModeTimestamp{
			TimestampLeft:   "\x1b[38;5;16;48;2;33;150;243m ",
			TimestampRight:  " \x1b[0m",
			TimestampFormat: time.RFC3339,
		},

		BodyLeft:  " ",
		BodyRight: "\r\n",
	}
)
