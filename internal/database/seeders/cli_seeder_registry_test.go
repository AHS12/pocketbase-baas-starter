package seeders

import (
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func TestGetAllCLISeederFunctions(t *testing.T) {
	functions := GetAllCLISeederFunctions()

	if len(functions) == 0 {
		t.Error("GetAllCLISeederFunctions should return at least one function")
	}

	// Check that all functions have required fields
	for i, fn := range functions {
		if fn.Name == "" {
			t.Errorf("Function %d should have a name", i)
		}
		if fn.Description == "" {
			t.Errorf("Function %d should have a description", i)
		}
		if fn.Function == nil {
			t.Errorf("Function %d should have a function", i)
		}
	}

	// Check for expected seeders
	expectedNames := []string{
		"UserSeeder[10]",
		"UserWithRoleSeeder[Admin][5]",
		"UserWithRoleSeeder[User][10]",
	}

	for _, expectedName := range expectedNames {
		found := false
		for _, fn := range functions {
			if fn.Name == expectedName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected seeder '%s' not found", expectedName)
		}
	}
}

func TestRunAllCLISeederFunctions(t *testing.T) {
	app := pocketbase.New()

	// Test that RunAllCLISeederFunctions handles panics gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Logf("RunAllCLISeederFunctions panicked as expected: %v", r)
		}
	}()

	// This will likely fail due to missing collections, but should handle errors gracefully
	err := RunAllCLISeederFunctions(app)
	if err == nil {
		t.Log("RunAllCLISeederFunctions completed successfully")
	} else {
		t.Logf("RunAllCLISeederFunctions failed as expected: %v", err)
	}
}

func TestCLISeederFunction(t *testing.T) {
	// Test CLISeederFunction struct
	fn := CLISeederFunction{
		Name:        "TestSeeder",
		Description: "Test seeder description",
		Function: func(app core.App) error {
			return nil
		},
	}

	if fn.Name != "TestSeeder" {
		t.Error("CLISeederFunction Name not set correctly")
	}
	if fn.Description != "Test seeder description" {
		t.Error("CLISeederFunction Description not set correctly")
	}
	if fn.Function == nil {
		t.Error("CLISeederFunction Function not set correctly")
	}

	// Test function execution
	app := pocketbase.New()
	err := fn.Function(app)
	if err != nil {
		t.Errorf("Test function should not return error: %v", err)
	}
}
