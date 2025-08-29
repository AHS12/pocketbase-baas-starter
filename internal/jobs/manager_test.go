package jobs

import (
	"sync"
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestGetJobManager(t *testing.T) {
	once = sync.Once{}
	globalJobManager = nil

	manager1 := GetJobManager()
	manager2 := GetJobManager()

	if manager1 == manager2 {
		t.Log("GetJobManager returns singleton instance")
	} else {
		t.Error("GetJobManager should return the same instance")
	}

	if manager1 == nil {
		t.Fatal("GetJobManager should not return nil")
	}
}

func TestJobManager_Initialize(t *testing.T) {
	once = sync.Once{}
	globalJobManager = nil

	manager := GetJobManager()
	app := pocketbase.New()

	// Test initialization
	err := manager.Initialize(app)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if !manager.IsInitialized() {
		t.Error("Manager should be initialized after Initialize call")
	}

	// Test double initialization (should not error)
	err = manager.Initialize(app)
	if err != nil {
		t.Errorf("Double initialization should not error: %v", err)
	}
}

func TestJobManager_GetProcessor(t *testing.T) {
	once = sync.Once{}
	globalJobManager = nil

	manager := GetJobManager()
	app := pocketbase.New()

	// Before initialization
	processor := manager.GetProcessor()
	if processor != nil {
		t.Error("GetProcessor should return nil before initialization")
	}

	// After initialization
	err := manager.Initialize(app)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	processor = manager.GetProcessor()
	if processor == nil {
		t.Error("GetProcessor should return processor after initialization")
	}
}

func TestJobManager_IsInitialized(t *testing.T) {
	once = sync.Once{}
	globalJobManager = nil

	manager := GetJobManager()
	app := pocketbase.New()

	// Before initialization
	if manager.IsInitialized() {
		t.Error("IsInitialized should return false before initialization")
	}

	// After initialization
	err := manager.Initialize(app)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if !manager.IsInitialized() {
		t.Error("IsInitialized should return true after initialization")
	}
}
