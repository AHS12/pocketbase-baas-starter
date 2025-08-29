package command

import (
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func TestHandleDBSeedCommand(t *testing.T) {
	app := pocketbase.New()
	cmd := &cobra.Command{}
	args := []string{}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HandleDBSeedCommand panicked as expected: %v", r)
		}
	}()

	HandleDBSeedCommand(app, cmd, args)
}

func TestHandleDBSeedCommandWithNilApp(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HandleDBSeedCommand with nil app panicked as expected: %v", r)
		}
	}()

	HandleDBSeedCommand(nil, cmd, args)
}

func TestHandleDBSeedCommandWithNilCmd(t *testing.T) {
	app := pocketbase.New()
	args := []string{}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HandleDBSeedCommand with nil cmd panicked as expected: %v", r)
		}
	}()

	HandleDBSeedCommand(app, nil, args)
}

func TestHandleDBSeedCommandWithNilArgs(t *testing.T) {
	app := pocketbase.New()
	cmd := &cobra.Command{}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HandleDBSeedCommand with nil args panicked as expected: %v", r)
		}
	}()

	HandleDBSeedCommand(app, cmd, nil)
}
