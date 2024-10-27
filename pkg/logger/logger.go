package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/AyoubTahir/projects_management/config"
)

type Logger struct {
	*log.Logger
}

func New(cfg config.LoggerConfig) (*Logger, error) {
	var file *os.File
	var err error

	if cfg.File != "" {
		file, err = os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	} else {
		file = os.Stdout
	}

	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{Logger: logger}, nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Printf("[INFO] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.Printf("[DEBUG] "+format, v...)
}
