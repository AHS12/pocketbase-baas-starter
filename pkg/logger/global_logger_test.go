package logger

import (
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestGlobalLogger(t *testing.T) {
	// Create a new PocketBase app
	app := pocketbase.New()

	// Get the logger (this should automatically set the global logger)
	logger := GetLogger(app)

	// Verify that the global logger is set by calling package-level functions
	// These should not panic
	Info("Test info message")
	Error("Test error message")
	Warn("Test warn message")
	Debug("Test debug message")

	// Verify the logger instance
	if logger == nil {
		t.Error("Expected logger instance, got nil")
	}

	// Verify that the global logger is storing logs by default
	if !logger.IsStoringLogs() {
		t.Error("Expected logger to store logs by default")
	}
}
