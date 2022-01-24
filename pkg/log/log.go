package log

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

const ReadWriteFileMode = 0666

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
