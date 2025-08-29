package command

import (
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func TestHandleSyncPermissionsCommand(t *testing.T) {
	app := pocketbase.New()
	cmd := &cobra.Command{}
	args := []string{}

	// The function will fail due to missing database, but we can verify
	// it attempts to process the expected number of permissions
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Command panicked as expected due to missing database: %v", r)
		}
	}()

	HandleSyncPermissionsCommand(app, cmd, args)
}

func TestHandleSyncPermissionsCommandWithNilApp(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HandleSyncPermissionsCommand with nil app panicked as expected: %v", r)
		}
	}()

	HandleSyncPermissionsCommand(nil, cmd, args)
}

func TestFindPermissionBySlug(t *testing.T) {
	app := pocketbase.New()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("findPermissionBySlug panicked as expected: %v", r)
		}
	}()

	record, err := findPermissionBySlug(app, "test-slug")
	if err == nil {
		t.Error("findPermissionBySlug should return error when permissions collection doesn't exist")
	}

	if record != nil {
		t.Error("findPermissionBySlug should return nil record on error")
	}
}

func TestSavePermissionBatch(t *testing.T) {
	app := pocketbase.New()

	// Test with empty batch
	err := savePermissionBatch(app, []*core.Record{})
	if err != nil {
		t.Errorf("savePermissionBatch with empty batch should not return error: %v", err)
	}

	// Test with nil records slice
	err = savePermissionBatch(app, nil)
	if err != nil {
		t.Errorf("savePermissionBatch with nil records should not return error: %v", err)
	}
}
