package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logLevel  LogLevel = INFO
	logOutput io.Writer = os.Stdout
)

// SetLevel sets the minimum log level that will be logged
func SetLevel(level string) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel = DEBUG
	case "INFO":
		logLevel = INFO
	case "WARN", "WARNING":
		logLevel = WARN
	case "ERROR":
		logLevel = ERROR
	case "FATAL":
		logLevel = FATAL
	default:
		logLevel = INFO
	}
}

// SetOutput sets the output destination for the logger
func SetOutput(w io.Writer) {
	logOutput = w
	log.SetOutput(w)
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	if logLevel <= INFO {
		log.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	if logLevel <= WARN {
		log.Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if logLevel <= ERROR {
		log.Printf("[ERROR] "+format, v...)
	}
}

// Fatal logs a fatal message and exits the program
func Fatal(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}

// GetLogger returns a logger with the specified prefix
func GetLogger(prefix string) *log.Logger {
	return log.New(logOutput, fmt.Sprintf("[%s] ", prefix), log.LstdFlags|log.Lshortfile)
}
