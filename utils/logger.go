package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the log level
type LogLevel int

const (
	// LogLevelDebug represents the debug log level
	LogLevelDebug LogLevel = iota
	// LogLevelInfo represents the info log level
	LogLevelInfo
	// LogLevelWarning represents the warning log level
	LogLevelWarning
	// LogLevelError represents the error log level
	LogLevelError
	// LogLevelFatal represents the fatal log level
	LogLevelFatal
)

// Logger represents a logger
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// NewLogger creates a new logger
func NewLogger() *Logger {
	// Create logger
	logger := log.New(os.Stdout, "", 0)
	
	// Determine log level from environment variable
	levelStr := os.Getenv("LOG_LEVEL")
	level := LogLevelInfo // Default to info
	
	switch levelStr {
	case "debug":
		level = LogLevelDebug
	case "info":
		level = LogLevelInfo
	case "warning":
		level = LogLevelWarning
	case "error":
		level = LogLevelError
	case "fatal":
		level = LogLevelFatal
	}
	
	return &Logger{
		level:  level,
		logger: logger,
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	if l.level <= LogLevelDebug {
		l.log("DEBUG", msg, keyvals...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	if l.level <= LogLevelInfo {
		l.log("INFO", msg, keyvals...)
	}
}

// Warning logs a warning message
func (l *Logger) Warning(msg string, keyvals ...interface{}) {
	if l.level <= LogLevelWarning {
		l.log("WARNING", msg, keyvals...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	if l.level <= LogLevelError {
		l.log("ERROR", msg, keyvals...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, keyvals ...interface{}) {
	if l.level <= LogLevelFatal {
		l.log("FATAL", msg, keyvals...)
		os.Exit(1)
	}
}

// log logs a message with the given level
func (l *Logger) log(level, msg string, keyvals ...interface{}) {
	// Format timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	
	// Format key-value pairs
	kvStr := ""
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			kvStr += fmt.Sprintf(" %v=%v", keyvals[i], keyvals[i+1])
		} else {
			kvStr += fmt.Sprintf(" %v=?", keyvals[i])
		}
	}
	
	// Log message
	l.logger.Printf("[%s] %s: %s%s", timestamp, level, msg, kvStr)
}
