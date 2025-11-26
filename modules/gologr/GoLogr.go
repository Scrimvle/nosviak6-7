package gologr

import (
	"io"
)


// GoLogr is the header struct for the package.
type GoLogr struct {
	aggregated *FileLogger
	rotateDir  string
	writer     io.Writer
}

// NewGoLogr initializes a new struct for logging
func NewGoLogr(rotateDir string, wr io.Writer) *GoLogr {
	return &GoLogr{
		rotateDir: rotateDir,
		writer:    wr,
	}
}

// SetAggregatedLogger will define all the params for the logger
func (origin *GoLogr) SetAggregatedLogger(fileName string, maxFileSize int64) error {
	logger := origin.NewFileLogger(fileName, maxFileSize)
	if logger.Err != nil {
		return logger.Err
	}

	state, err := logger.file.Stat()
	if err != nil {
		return err
	}

	if state.Size() > 0 {
		err := logger.rotateLog()
		if err != nil {
			return err
		}
	}

	/* defines the params inside origin and logger */
	logger.shouldAggregate = false
	origin.aggregated = logger
	return nil
}

func (origin *GoLogr) AggregateTerminal() *Terminal {
	return origin.aggregated.WithTerminal()
}
