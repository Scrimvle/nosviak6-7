package signal

import (
	"context"
)

// Read will continue to read until the reader is closed
func (ts *TerminalSignal) read() error {
	for {
		buffer := make([]byte, ts.readPayload)
		index, err := ts.term.Read(buffer)
		if err != nil {
			return err
		}

		ts.Queue <- buffer[:index]
	}
}

// ReadWithContext implements a contextual based io.reader reading
func (ts *TerminalSignal) ReadWithContext(ctx context.Context) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	for {
		select {

		case r, ok := <- ts.Queue:
			if !ok || r == nil {
				return nil, ErrReadQueueClosed
			}

			ctx.Done()
			return r, ctx.Err()

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}