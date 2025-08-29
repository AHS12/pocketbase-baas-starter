package app

import (
	"testing"
)

func TestAppCreation(t *testing.T) {
	pbApp := NewApp()
	if pbApp == nil {
		t.Fatal("Expected app.NewApp() to return a non-nil app")
	}

	if pbApp.Settings() == nil {
		t.Fatal("Expected app to have settings configured")
	}

	if pbApp.OnServe() == nil {
		t.Fatal("Expected OnServe hook to be registered")
	}
}

func TestMiddlewareRegistration(t *testing.T) {
	pbApp := NewApp()

	onServeHook := pbApp.OnServe()
	if onServeHook == nil {
		t.Fatal("Expected OnServe hook to be registered")
	}
}
