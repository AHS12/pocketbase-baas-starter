package logger

import (
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestLoggerSingleton(t *testing.T) {
	app1 := pocketbase.New()
	logger1 := GetLogger(app1)

	app2 := pocketbase.New()
	logger2 := GetLogger(app2)

	// Both should return the same instance
	if logger1 != logger2 {
		t.Error("Expected same logger instance, got different instances")
	}
}

func TestLoggerLevels(t *testing.T) {
	app := pocketbase.New()
	logger := GetLogger(app)

	logger.Debug("Debug message", "key", "value")
	logger.Info("Info message", "count", 42)
	logger.Warn("Warning message", "warning", "test")
	logger.Error("Error message", "error", "test")

	if !logger.IsStoringLogs() {
		t.Error("Expected logger to store logs by default")
	}

	logger.SetStoreLogs(false)
	if logger.IsStoringLogs() {
		t.Error("Expected logger to not store logs after disabling")
	}

	logger.SetStoreLogs(true)
	if !logger.IsStoringLogs() {
		t.Error("Expected logger to store logs after re-enabling")
	}
}

func TestLogLevelString(t *testing.T) {
	if DEBUG.String() != "DEBUG" {
		t.Errorf("Expected DEBUG.String() to be 'DEBUG', got %s", DEBUG.String())
	}

	if INFO.String() != "INFO" {
		t.Errorf("Expected INFO.String() to be 'INFO', got %s", INFO.String())
	}

	if WARN.String() != "WARN" {
		t.Errorf("Expected WARN.String() to be 'WARN', got %s", WARN.String())
	}

	if ERROR.String() != "ERROR" {
		t.Errorf("Expected ERROR.String() to be 'ERROR', got %s", ERROR.String())
	}

	unknownLevel := LogLevel(999)
	if unknownLevel.String() != "UNKNOWN" {
		t.Errorf("Expected unknown level to return 'UNKNOWN', got %s", unknownLevel.String())
	}
}

func TestFormatMessage(t *testing.T) {
	app := pocketbase.New()
	logger := GetLogger(app).(*pbLogger)

	result := logger.formatMessage("Simple message")
	if result != "Simple message" {
		t.Errorf("Expected 'Simple message', got %s", result)
	}

	result = logger.formatMessage("Message", "key1", "value1", "key2", "value2")
	if result != "Message key1=value1 key2=value2" {
		t.Errorf("Expected 'Message key1=value1 key2=value2', got %s", result)
	}

	result = logger.formatMessage("Message", "key1", "value1", "key2")
	if result != "Message key1=value1" {
		t.Errorf("Expected 'Message key1=value1', got %s", result)
	}
}

func TestNoOpLogger(t *testing.T) {
	noop := &noopLogger{}

	noop.Debug("debug message", "key", "value")
	noop.Info("info message", "key", "value")
	noop.Warn("warn message", "key", "value")
	noop.Error("error message", "key", "value")

	noop.SetStoreLogs(true)
	if noop.IsStoringLogs() {
		t.Error("NoOp logger should always return false for IsStoringLogs")
	}

	result := noop.FormatMessage("test message", "key", "value")
	expected := "test message key=value"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFromApp(t *testing.T) {
	app := pocketbase.New()

	logger := FromApp(app)
	if logger == nil {
		t.Error("Expected logger from valid PocketBase app, got nil")
	}

	nilLogger := FromApp(nil)
	if nilLogger != nil {
		t.Error("Expected nil logger from nil app, got non-nil")
	}
}

func TestFromAppOrDefault(t *testing.T) {
	app := pocketbase.New()

	logger := FromAppOrDefault(app)
	if logger == nil {
		t.Error("Expected logger from valid PocketBase app, got nil")
	}

	defaultLogger := FromAppOrDefault(nil)
	if defaultLogger == nil {
		t.Error("Expected default logger from nil app, got nil")
	}

	if defaultLogger.IsStoringLogs() {
		t.Error("Default logger should be noopLogger which doesn't store logs")
	}
}

func TestGlobalLoggerFunctions(t *testing.T) {
	originalGlobal := globalLogger
	globalLogger = nil
	defer func() { globalLogger = originalGlobal }()

	Debug("debug message", "key", "value")
	Info("info message", "key", "value")
	Warn("warn message", "key", "value")
	Error("error message", "key", "value")

	mockLogger := &noopLogger{}
	SetGlobalLogger(mockLogger)

	if globalLogger != mockLogger {
		t.Error("Expected global logger to be set to mock logger")
	}
}

func TestSetGlobalLogger(t *testing.T) {
	originalGlobal := globalLogger
	defer func() { globalLogger = originalGlobal }()

	mockLogger := &noopLogger{}
	SetGlobalLogger(mockLogger)

	if globalLogger != mockLogger {
		t.Error("Expected SetGlobalLogger to set the global logger")
	}

	SetGlobalLogger(nil)
	if globalLogger != nil {
		t.Error("Expected SetGlobalLogger to accept nil")
	}
}
