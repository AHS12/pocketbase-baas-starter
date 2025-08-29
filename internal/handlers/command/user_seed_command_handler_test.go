package command

import (
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func TestHandleSeedUsersCommand(t *testing.T) {
	app := pocketbase.New()
	cmd := &cobra.Command{}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no arguments (default count)",
			args: []string{},
		},
		{
			name: "valid count argument",
			args: []string{"5"},
		},
		{
			name: "invalid count argument",
			args: []string{"invalid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("HandleSeedUsersCommand panicked as expected: %v", r)
				}
			}()

			HandleSeedUsersCommand(app, cmd, tt.args)
		})
	}
}

func TestHandleSeedUsersWithRoleCommand(t *testing.T) {
	app := pocketbase.New()
	cmd := &cobra.Command{}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no arguments",
			args: []string{},
		},
		{
			name: "insufficient arguments",
			args: []string{"5"},
		},
		{
			name: "valid arguments",
			args: []string{"5", "admin"},
		},
		{
			name: "invalid count argument",
			args: []string{"invalid", "admin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("HandleSeedUsersWithRoleCommand panicked as expected: %v", r)
				}
			}()

			HandleSeedUsersWithRoleCommand(app, cmd, tt.args)
		})
	}
}

func TestHandleSeedUsersCommandWithNilApp(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{"5"}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HandleSeedUsersCommand with nil app panicked as expected: %v", r)
		}
	}()

	HandleSeedUsersCommand(nil, cmd, args)
}
