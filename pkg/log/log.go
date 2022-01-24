package log

import (
	"fmt"
	"log"
	"os"
)

// Logger allows log info by level.
type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

// ReadWriteFileMode presents a code for ReadWrite file mode.
const ReadWriteFileMode = 0666

// NewLogger creates a new instance of Logger.
func NewLogger() (*Logger, error) {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, ReadWriteFileMode)
	if err != nil {
		return &Logger{}, fmt.Errorf("cannot create logger: %w", err)
	}

	return &Logger{
		InfoLogger:  log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
		ErrorLogger: log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
	}, nil
}
