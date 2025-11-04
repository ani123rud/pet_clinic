package logger

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"petclinic/data"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

// ContextKey is the type for context keys used by the logger
type ContextKey string

// CtxUserIDKey is the exported key to put/get user id from context
var CtxUserIDKey ContextKey = "user_id"
// CtxUserEmailKey is the exported key to put/get user email from context
var CtxUserEmailKey ContextKey = "user_email"

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logLevel  LogLevel  = DEBUG
	logOutput io.Writer = os.Stdout
)

// getCallerInfo returns the file and function name of the caller
func getCallerInfo() (string, string) {
	// Skip 3 levels to get the actual caller (1 for getCallerInfo, 1 for the log function, 1 for the actual caller)
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown", "unknown"
	}

	// Get just the filename from the full path
	filename := filepath.Base(file)

	// Get the function name
	funcName := runtime.FuncForPC(pc).Name()
	// Get just the package.function name
	parts := strings.Split(funcName, ".")
	if len(parts) > 1 {
		funcName = parts[len(parts)-1]
	}

	return fmt.Sprintf("%s:%d", filename, line), funcName
}

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

// SetDB sets the database connection for logging
func SetDB(database *sql.DB) {
	dbOnce.Do(func() {
		db = database
		// Optionally initialize logs table if environment allows
		if os.Getenv("LOG_INIT_TABLE") == "true" {
			if err := data.InitLogsTable(db); err != nil {
				log.Printf("Failed to initialize logs table: %v", err)
			}
		}
	})
}

// logInternal is the internal logging function that adds caller info
func logInternal(level, format string, v ...interface{}) {
	fileInfo, funcName := getCallerInfo()
	message := fmt.Sprintf(format, v...)
	prefix := fmt.Sprintf("[%s] [%s] [%s] %s", level, fileInfo, funcName, message)
	
	// Log to console
	log.Print(prefix)
	
	// Log to database if DB is set
	if db != nil {
		go func() {
			// Use a goroutine to prevent blocking
			if err := data.SaveLog(db, level, message, fileInfo, funcName); err != nil {
				log.Printf("Failed to save log to database: %v", err)
			}
		}()
	}
}

// logInternalWithUser logs and persists with optional user id and email
func logInternalWithUser(userID *int, userEmail *string, level, format string, v ...interface{}) {
	fileInfo, funcName := getCallerInfo()
	message := fmt.Sprintf(format, v...)
	// Add [user:<id>] or [user:<email>] suffix only for console clarity when present
	userSuffix := ""
	if userEmail != nil && *userEmail != "" {
		userSuffix = fmt.Sprintf(" [user:%s]", *userEmail)
	} else if userID != nil {
		userSuffix = fmt.Sprintf(" [user:%d]", *userID)
	}
	prefix := fmt.Sprintf("[%s] [%s] [%s] %s%s", level, fileInfo, funcName, message, userSuffix)
	log.Print(prefix)
	if db != nil {
		go func() {
			var uidNull sql.NullInt64
			if userID != nil {
				uidNull = sql.NullInt64{Int64: int64(*userID), Valid: true}
			}
			var emailNull sql.NullString
			if userEmail != nil && *userEmail != "" {
				emailNull = sql.NullString{String: *userEmail, Valid: true}
			}
			if err := data.SaveLogWithUser(db, level, message, fileInfo, funcName, uidNull, emailNull); err != nil {
				log.Printf("Failed to save log to database: %v", err)
			}
		}()
	}
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		logInternal("DEBUG", format, v...)
	}
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	if logLevel <= INFO {
		logInternal("INFO", format, v...)
	}
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	if logLevel <= WARN {
		logInternal("WARN", format, v...)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if logLevel <= ERROR {
		logInternal("ERROR", format, v...)
	}
}

// Fatal logs a fatal message and exits the program
func Fatal(format string, v ...interface{}) {
	logInternal("FATAL", format, v...)
	os.Exit(1)
}

// GetLogger returns a logger with the specified prefix
func GetLogger(prefix string) *log.Logger {
	return log.New(logOutput, fmt.Sprintf("[%s] ", prefix), log.LstdFlags|log.Lshortfile)
}

// Context-aware variants
func DebugCtx(ctx context.Context, format string, v ...interface{}) {
	if logLevel <= DEBUG {
		var uid *int
		var uemail *string
		if ctx != nil {
			if v, ok := ctx.Value(CtxUserIDKey).(int); ok {
				uid = &v
			}
			if e, ok := ctx.Value(CtxUserEmailKey).(string); ok && e != "" {
				uemail = &e
			}
		}
		logInternalWithUser(uid, uemail, "DEBUG", format, v...)
	}
}

func InfoCtx(ctx context.Context, format string, v ...interface{}) {
	if logLevel <= INFO {
		var uid *int
		var uemail *string
		if ctx != nil {
			if v, ok := ctx.Value(CtxUserIDKey).(int); ok {
				uid = &v
			}
			if e, ok := ctx.Value(CtxUserEmailKey).(string); ok && e != "" {
				uemail = &e
			}
		}
		logInternalWithUser(uid, uemail, "INFO", format, v...)
	}
}

func WarnCtx(ctx context.Context, format string, v ...interface{}) {
	if logLevel <= WARN {
		var uid *int
		var uemail *string
		if ctx != nil {
			if v, ok := ctx.Value(CtxUserIDKey).(int); ok {
				uid = &v
			}
			if e, ok := ctx.Value(CtxUserEmailKey).(string); ok && e != "" {
				uemail = &e
			}
		}
		logInternalWithUser(uid, uemail, "WARN", format, v...)
	}
}

func ErrorCtx(ctx context.Context, format string, v ...interface{}) {
	if logLevel <= ERROR {
		var uid *int
		var uemail *string
		if ctx != nil {
			if v, ok := ctx.Value(CtxUserIDKey).(int); ok {
				uid = &v
			}
			if e, ok := ctx.Value(CtxUserEmailKey).(string); ok && e != "" {
				uemail = &e
			}
		}
		logInternalWithUser(uid, uemail, "ERROR", format, v...)
	}
}

func FatalCtx(ctx context.Context, format string, v ...interface{}) {
	var uid *int
	var uemail *string
	if ctx != nil {
		if v, ok := ctx.Value(CtxUserIDKey).(int); ok {
			uid = &v
		}
		if e, ok := ctx.Value(CtxUserEmailKey).(string); ok && e != "" {
			uemail = &e
		}
	}
	logInternalWithUser(uid, uemail, "FATAL", format, v...)
	os.Exit(1)
}
