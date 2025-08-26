package logger

import (
	"fmt"
	"log"
	"sync"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger interface defines the methods for our custom logger
type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	SetStoreLogs(store bool)
	IsStoringLogs() bool
	FormatMessage(msg string, keysAndValues ...any) string
}

// pbLogger implements the Logger interface
type pbLogger struct {
	pbApp     *pocketbase.PocketBase
	storeLogs bool
}

// noopLogger is a no-op logger that only logs to stdout
type noopLogger struct{}

// singleton instance
var (
	instance     Logger
	once         sync.Once
	globalLogger Logger
)

// GetLogger returns the singleton logger instance
func GetLogger(app *pocketbase.PocketBase) Logger {
	once.Do(func() {
		instance = &pbLogger{
			pbApp:     app,
			storeLogs: true, // Default to storing logs in DB
		}
		SetGlobalLogger(instance)
	})
	return instance
}

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

// FromApp retrieves a logger instance from a core.App interface
func FromApp(app core.App) Logger {
	if pbApp, ok := app.(*pocketbase.PocketBase); ok {
		return GetLogger(pbApp)
	}
	return nil
}

// FromAppOrDefault retrieves a logger instance from a core.App interface
func FromAppOrDefault(app core.App) Logger {
	if pbApp, ok := app.(*pocketbase.PocketBase); ok {
		return GetLogger(pbApp)
	}
	return &noopLogger{}
}

// Package-level logging functions
func Debug(msg string, keysAndValues ...any) {
	if globalLogger != nil {
		globalLogger.Debug(msg, keysAndValues...)
	} else {
		logWithLevel(DEBUG, msg, keysAndValues...)
	}
}

func Info(msg string, keysAndValues ...any) {
	if globalLogger != nil {
		globalLogger.Info(msg, keysAndValues...)
	} else {
		logWithLevel(INFO, msg, keysAndValues...)
	}
}

func Warn(msg string, keysAndValues ...any) {
	if globalLogger != nil {
		globalLogger.Warn(msg, keysAndValues...)
	} else {
		logWithLevel(WARN, msg, keysAndValues...)
	}
}

func Error(msg string, keysAndValues ...any) {
	if globalLogger != nil {
		globalLogger.Error(msg, keysAndValues...)
	} else {
		logWithLevel(ERROR, msg, keysAndValues...)
	}
}

// Implementation of Logger interface for pbLogger
func (l *pbLogger) SetStoreLogs(store bool) {
	l.storeLogs = store
}

func (l *pbLogger) IsStoringLogs() bool {
	return l.storeLogs
}

func (l *pbLogger) Debug(msg string, keysAndValues ...any) {
	l.logWithLevel(DEBUG, msg, keysAndValues...)
}

func (l *pbLogger) Info(msg string, keysAndValues ...any) {
	l.logWithLevel(INFO, msg, keysAndValues...)
}

func (l *pbLogger) Warn(msg string, keysAndValues ...any) {
	l.logWithLevel(WARN, msg, keysAndValues...)
}

func (l *pbLogger) Error(msg string, keysAndValues ...any) {
	l.logWithLevel(ERROR, msg, keysAndValues...)
}

func (l *pbLogger) FormatMessage(msg string, keysAndValues ...any) string {
	return l.formatMessage(msg, keysAndValues...)
}

// Implementation of Logger interface for noopLogger
func (n *noopLogger) Debug(msg string, keysAndValues ...any) {
	logWithLevel(DEBUG, msg, keysAndValues...)
}

func (n *noopLogger) Info(msg string, keysAndValues ...any) {
	logWithLevel(INFO, msg, keysAndValues...)
}

func (n *noopLogger) Warn(msg string, keysAndValues ...any) {
	logWithLevel(WARN, msg, keysAndValues...)
}

func (n *noopLogger) Error(msg string, keysAndValues ...any) {
	logWithLevel(ERROR, msg, keysAndValues...)
}

func (n *noopLogger) SetStoreLogs(store bool) {
	// No-op
}

func (n *noopLogger) IsStoringLogs() bool {
	return false
}

func (n *noopLogger) FormatMessage(msg string, keysAndValues ...any) string {
	return formatMessage(msg, keysAndValues...)
}

// logWithLevel is a helper method that handles logging at different levels
func (l *pbLogger) logWithLevel(level LogLevel, msg string, keysAndValues ...any) {
	if l.storeLogs && l.pbApp != nil {
		switch level {
		case DEBUG:
			l.pbApp.Logger().Debug(msg, keysAndValues...)
		case INFO:
			l.pbApp.Logger().Info(msg, keysAndValues...)
		case WARN:
			l.pbApp.Logger().Warn(msg, keysAndValues...)
		case ERROR:
			l.pbApp.Logger().Error(msg, keysAndValues...)
		}
	}
}

// logWithLevel is a helper function that logs to stdout only (for noopLogger and fallback)
func logWithLevel(level LogLevel, msg string, keysAndValues ...any) {
	formattedMsg := formatMessage(msg, keysAndValues...)
	log.Printf("[%s] %s", level.String(), formattedMsg)
}

// formatMessage formats the log message with key-value pairs for stdout
func (l *pbLogger) formatMessage(msg string, keysAndValues ...any) string {
	if len(keysAndValues) == 0 {
		return msg
	}

	result := msg
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := keysAndValues[i]
			value := keysAndValues[i+1]
			result = fmt.Sprintf("%s %v=%v", result, key, value)
		}
	}
	return result
}

// formatMessage formats the log message with key-value pairs for stdout (standalone function)
func formatMessage(msg string, keysAndValues ...any) string {
	if len(keysAndValues) == 0 {
		return msg
	}

	result := msg
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := keysAndValues[i]
			value := keysAndValues[i+1]
			result = fmt.Sprintf("%s %v=%v", result, key, value)
		}
	}
	return result
}
