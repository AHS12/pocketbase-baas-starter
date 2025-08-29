package factories

import (
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestNewUserFactory(t *testing.T) {
	app := pocketbase.New()
	factory := NewUserFactory(app)

	if factory == nil {
		t.Fatal("NewUserFactory should not return nil")
	}

	if factory.app != app {
		t.Error("Factory should store app reference")
	}
}

func TestUserFactory_GenerateMany_ZeroCount(t *testing.T) {
	app := pocketbase.New()
	factory := NewUserFactory(app)

	records, err := factory.GenerateMany(0)
	if err != nil {
		t.Errorf("GenerateMany(0) should not return error: %v", err)
	}

	if len(records) != 0 {
		t.Errorf("GenerateMany(0) should return empty slice, got %d records", len(records))
	}
}
