package commands

import (
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestRegisterCommands(t *testing.T) {
	app := pocketbase.New()

	err := RegisterCommands(app)
	if err != nil {
		t.Fatalf("RegisterCommands failed: %v", err)
	}

	rootCmd := app.RootCmd
	if rootCmd == nil {
		t.Fatal("App should have root command")
	}

	expectedCommands := []string{"health", "sync-permissions", "db-seed", "seed-users"}
	commands := rootCmd.Commands()

	for _, expectedCmd := range expectedCommands {
		found := false
		for _, cmd := range commands {
			if cmd.Use == expectedCmd || cmd.Use == expectedCmd+" [count]" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command '%s' not found", expectedCmd)
		}
	}
}

func TestRegisterCommandsWithNilApp(t *testing.T) {
	err := RegisterCommands(nil)
	if err == nil {
		t.Error("RegisterCommands should return error with nil app")
	}

	expectedError := "RegisterCommands: app cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestCommandStructure(t *testing.T) {
	// Test that Command struct has required fields
	cmd := Command{
		ID:      "test",
		Use:     "test",
		Short:   "Test command",
		Long:    "Test command description",
		Handler: nil,
		Enabled: true,
	}

	if cmd.ID != "test" {
		t.Error("Command ID not set correctly")
	}
	if cmd.Use != "test" {
		t.Error("Command Use not set correctly")
	}
	if cmd.Short != "Test command" {
		t.Error("Command Short not set correctly")
	}
	if cmd.Long != "Test command description" {
		t.Error("Command Long not set correctly")
	}
	if !cmd.Enabled {
		t.Error("Command Enabled not set correctly")
	}
}
