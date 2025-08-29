package crons

import (
	"os"
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestRegisterCrons(t *testing.T) {
	app := pocketbase.New()

	err := RegisterCrons(app)
	if err != nil {
		t.Fatalf("RegisterCrons failed: %v", err)
	}
}

func TestRegisterCronsWithNilApp(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("RegisterCrons should panic with nil app")
		} else if r != "RegisterCrons: app cannot be nil" {
			t.Errorf("Expected panic message 'RegisterCrons: app cannot be nil', got '%v'", r)
		}
	}()

	_ = RegisterCrons(nil)
}

func TestRegisterCronsWithDisabledCrons(t *testing.T) {
	app := pocketbase.New()

	os.Setenv("ENABLE_SYSTEM_QUEUE_CRON", "false")
	os.Setenv("ENABLE_CLEAR_EXPORT_FILES_CRON", "false")
	defer func() {
		os.Unsetenv("ENABLE_SYSTEM_QUEUE_CRON")
		os.Unsetenv("ENABLE_CLEAR_EXPORT_FILES_CRON")
	}()

	err := RegisterCrons(app)
	if err != nil {
		t.Fatalf("RegisterCrons with disabled crons failed: %v", err)
	}
}

func TestCronStructure(t *testing.T) {
	cron := Cron{
		ID:          "test",
		CronExpr:    "* * * * *",
		Handler:     func() {},
		Enabled:     true,
		Description: "Test cron",
	}

	if cron.ID != "test" {
		t.Error("Cron ID not set correctly")
	}
	if cron.CronExpr != "* * * * *" {
		t.Error("Cron CronExpr not set correctly")
	}
	if cron.Handler == nil {
		t.Error("Cron Handler not set correctly")
	}
	if !cron.Enabled {
		t.Error("Cron Enabled not set correctly")
	}
	if cron.Description != "Test cron" {
		t.Error("Cron Description not set correctly")
	}
}
