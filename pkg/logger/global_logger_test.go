package logger

import (
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestGlobalLogger(t *testing.T) {
	app := pocketbase.New()

	logger := GetLogger(app)

	Info("Test info message")
	Error("Test error message")
	Warn("Test warn message")
	Debug("Test debug message")

	if logger == nil {
		t.Error("Expected logger instance, got nil")
	}

	if !logger.IsStoringLogs() {
		t.Error("Expected logger to store logs by default")
	}
}
