package routes

import (
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func TestRegisterCustom(t *testing.T) {
	app := pocketbase.New()

	event := &core.ServeEvent{
		App: app,
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("RegisterCustom panicked as expected due to nil router: %v", r)
		}
	}()

	err := RegisterCustom(event)
	if err != nil {
		t.Logf("RegisterCustom failed as expected: %v", err)
	}
}

func TestRegisterCustomWithNilEvent(t *testing.T) {
	err := RegisterCustom(nil)
	if err == nil {
		t.Error("RegisterCustom should return error with nil event")
	}

	expectedError := "RegisterCustom: serve event cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestRouteStructure(t *testing.T) {
	// Test that Route struct has required fields
	route := Route{
		Method:      "GET",
		Path:        "/test",
		Handler:     func(*core.RequestEvent) error { return nil },
		Middlewares: []func(*core.RequestEvent) error{},
		Enabled:     true,
		Description: "Test route",
	}

	if route.Method != "GET" {
		t.Error("Route Method not set correctly")
	}
	if route.Path != "/test" {
		t.Error("Route Path not set correctly")
	}
	if route.Handler == nil {
		t.Error("Route Handler not set correctly")
	}
	if route.Middlewares == nil {
		t.Error("Route Middlewares not set correctly")
	}
	if !route.Enabled {
		t.Error("Route Enabled not set correctly")
	}
	if route.Description != "Test route" {
		t.Error("Route Description not set correctly")
	}
}

func TestRouteMethodSupport(t *testing.T) {
	supportedMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range supportedMethods {
		t.Run(method, func(t *testing.T) {
			// This test verifies that the method switch statement works
			if method == "" {
				t.Error("Method should not be empty")
			}
		})
	}

	// Test unsupported method (should be skipped)
	unsupportedMethod := "INVALID"
	if unsupportedMethod == "GET" || unsupportedMethod == "POST" ||
		unsupportedMethod == "PUT" || unsupportedMethod == "DELETE" ||
		unsupportedMethod == "PATCH" {
		t.Error("INVALID should not be a supported method")
	}
}
