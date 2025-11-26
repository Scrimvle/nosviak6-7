package gologr

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileLogger is an individual struct for logging to files.
type FileLogger struct {
	Err             error
	file            *os.File
	origin          *GoLogr
	fileName        string
	maxFileSize     int64
	shouldAggregate bool
}

// NewFileLogger will initialize a logging instance for a file.
func (origin *GoLogr) NewFileLogger(fileName string, maxFileSize int64) *FileLogger {
	logger := &FileLogger{
		Err:             nil,
		origin:          origin,
		fileName:        fileName,
		maxFileSize:     maxFileSize,
		shouldAggregate: true,
	}

	logger.file, logger.Err = os.OpenFile(logger.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return logger
}

// WithTerminal will transition into terminal and file logging.
func (l *FileLogger) WithTerminal() *Terminal {
	return &Terminal{l}
}

// WriteLog will attempt to write the log into the file and possible the aggregated feed.
func (l *FileLogger) WriteLog(message string, args ...any) error {
	if l.shouldRotate() {
		err := l.rotateLog()
		if err != nil {
			return err
		}
	}

	format := fmt.Sprintf(time.Now().Format("[Mon Jan _2 15:04:05 2006]") + " " + message, args...) + "\r\n"
	if _, err := l.file.WriteString(format); err != nil || !l.shouldAggregate || l.origin.aggregated == nil {
		return err
	}

	return l.origin.aggregated.WriteLog(message, args...)
}

// Close will destroy the file descriptor for the FileLogger.
func (l *FileLogger) Close() error {
	if l.file == nil || l.Err != nil {
		return l.Err
	}

	defer func() {
		l.file = nil
	}()

	return l.file.Close()
}

// shouldRotate checks whether we should rotate the log file round.
func (l *FileLogger) shouldRotate() bool {

	information, err := l.file.Stat()
	if err != nil {
		return true
	}

	if l.maxFileSize <= 0 || l.maxFileSize > information.Size() {
		return false
	}

	dirLen, err := os.ReadDir(l.origin.rotateDir)
	if err != nil {
		return false
	}

	if len(dirLen) >= 20 {
		err := os.Remove(l.origin.rotateDir)
		if err != nil {
			log.Println(err)
			return false
		}

		if err := os.Mkdir(l.origin.rotateDir, 0777); err != nil {
			log.Println(err)
			return false
		}
	}

	return true
}

// rotateLog will attempt to rotate the log file around
func (l *FileLogger) rotateLog() error {
	err := l.Close()
	if err != nil {
		return err
	}

	err = os.Rename(l.fileName, filepath.Join(l.origin.rotateDir, fmt.Sprintf("%s.%s", time.Now().Format("20060102150405"), strings.Split(l.fileName, string(filepath.Separator))[strings.Count(l.fileName, string(filepath.Separator))])))
	if err != nil {
		if os.IsExist(err) {
			return err
		}

		if err := os.Mkdir(l.origin.rotateDir, 0666); err != nil {
			return err
		}

		return l.rotateLog()
	}

	l.file, err = os.OpenFile(l.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return err
}
