package command

import (
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"

	log "ims-pocketbase-baas-starter/pkg/logger"
)

// HandleHealthCheckCommand performs a basic health check
func HandleHealthCheckCommand(app *pocketbase.PocketBase, cmd *cobra.Command, args []string) {
	log.Info("Running health check...")

	if app == nil {
		log.Error("Application is not properly initialized")
		return
	}

	if app.Settings() == nil {
		log.Error("Application settings are not loaded")
		return
	}

	log.Info("OK: Application is properly initialized")
	log.Info("OK: Application settings are loaded")
	log.Info("OK: Application is healthy")
}
