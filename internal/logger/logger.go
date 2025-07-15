package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logger  *log.Logger
	logFile *os.File

	// DebugEnabled controls whether debug messages are printed to stdout.
	DebugEnabled bool

	// mu protects the logger from concurrent writes.
	mu sync.Mutex
)

// Init sets up the logger to write to a daily log file.
func Init(logDir string) error {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("creating log directory %q: %w", logDir, err)
	}

	today := time.Now().Format("2006-01-02")
	path := filepath.Join(logDir, "jobctl-"+today+".log")

	var err error
	logFile, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	multi := io.MultiWriter(logFile)

	logger = log.New(multi, "[LOG] ", log.LstdFlags|log.Lmsgprefix)

	return nil
}

// Log writes a message to the log file and, if DebugEnabled is true, to stdout.
func Log(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)

	mu.Lock()
	defer mu.Unlock()

	logger.Println(msg)

	// If debug mode is on, also print the log to the console.
	if DebugEnabled {
		fmt.Fprint(os.Stdout, "\r\n")
		fmt.Fprintf(os.Stdout, "[LOG] %s\n", msg)
		fmt.Fprint(os.Stdout, "jobctl> ")
	}
}

// Close closes the log file.
func Close() error {
	if logFile != nil {
		return logFile.Close()
	}

	return nil
}
