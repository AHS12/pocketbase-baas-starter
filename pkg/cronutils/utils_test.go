package cronutils

import (
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestValidateCronExpression(t *testing.T) {
	tests := []struct {
		name        string
		cronExpr    string
		expectError bool
	}{
		// Valid expressions
		{"valid basic", "0 0 * * *", false},
		{"valid with ranges", "0-30 8-17 * * 1-5", false},
		{"valid with lists", "0,15,30,45 * * * *", false},
		{"valid with steps", "*/5 * * * *", false},
		{"valid complex", "0,15,30,45 8-17 * * 1-5", false},
		{"valid 6-field with seconds", "0 0 0 * * *", false},
		{"valid wildcard", "* * * * *", false},

		// Invalid expressions
		{"empty expression", "", true},
		{"too few fields", "0 0 *", true},
		{"too many fields", "0 0 * * * * *", true},
		{"invalid minute", "60 * * * *", true},
		{"invalid hour", "0 24 * * *", true},
		{"invalid day", "0 0 32 * *", true},
		{"invalid month", "0 0 * 13 *", true},
		{"invalid weekday", "0 0 * * 8", true},
		{"invalid range", "0 0 5-3 * *", true},
		{"invalid character", "0 0 * * X", true},
		{"negative value", "0 0 -1 * *", true},
		{"invalid range format", "0 0 1-2-3 * *", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCronExpression(tt.cronExpr)
			if tt.expectError && err == nil {
				t.Errorf("expected error for %q, but got none", tt.cronExpr)
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error for %q: %v", tt.cronExpr, err)
			}
		})
	}
}

func TestNewCronExecutionContext(t *testing.T) {
	app := pocketbase.New()
	cronID := "test-cron-123"

	ctx := NewCronExecutionContext(app, cronID)

	if ctx == nil {
		t.Fatal("NewCronExecutionContext should not return nil")
	}

	if ctx.App != app {
		t.Error("Expected context to have correct app reference")
	}

	if ctx.CronID != cronID {
		t.Errorf("Expected CronID %q, got %q", cronID, ctx.CronID)
	}

	if ctx.StartTime.IsZero() {
		t.Error("Expected StartTime to be set")
	}
}

func TestCronExecutionContext_LogMethods(t *testing.T) {
	app := pocketbase.New()
	ctx := NewCronExecutionContext(app, "test-cron")

	// Test that log methods don't panic
	ctx.LogStart("Starting test operation")
	ctx.LogEnd("Test operation completed")
	ctx.LogError(fmt.Errorf("test error"), "Test error occurred")
	ctx.LogDebug(map[string]string{"key": "value"}, "Debug information")

	// Verify context maintains state
	if ctx.CronID != "test-cron" {
		t.Errorf("Expected CronID to remain 'test-cron', got %q", ctx.CronID)
	}
}

func TestWithRecovery(t *testing.T) {
	app := pocketbase.New()
	cronID := "test-recovery"

	// Test normal execution
	executed := false
	normalFunc := func() {
		executed = true
	}

	wrappedFunc := WithRecovery(app, cronID, normalFunc)
	wrappedFunc()

	if !executed {
		t.Error("Expected wrapped function to execute normally")
	}

	// Test panic recovery
	panicFunc := func() {
		panic("test panic")
	}

	wrappedPanicFunc := WithRecovery(app, cronID, panicFunc)

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected panic to be recovered, but got panic: %v", r)
		}
	}()

	wrappedPanicFunc()
}
