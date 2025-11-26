package signal

import "errors"

// ErrReadQueueClosed is returned whenever the read queue is closed
var ErrReadQueueClosed = errors.New("queue is closed")