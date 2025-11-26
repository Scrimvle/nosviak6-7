package signal

import (
	"golang.org/x/crypto/ssh"
)

// TerminalSignal will represent a signaller for the terminal
type TerminalSignal struct {
	term        ssh.Channel
	Queue       chan []byte
	readPayload int
}

// NewSignaller creates a new TerminalSignal which also initializes the required routine
func NewSignaller(channel ssh.Channel) *TerminalSignal {
	signal := &TerminalSignal{
		term:        channel,
		Queue:       make(chan []byte),
		readPayload: 128,
	}

	go signal.read()
	return signal
}
