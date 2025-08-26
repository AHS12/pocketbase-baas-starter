package hook

import (
	"fmt"

	log "ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase/core"
)

// HandleRecordCreate handles record creation events
func HandleRecordCreate(e *core.RecordEvent) error {
	// Log the record creation

	log.Info("Record created",
		"collection", e.Record.Collection().Name,
		"id", e.Record.Id,
		"created", e.Record.GetDateTime("created"),
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleRecordUpdate handles record update events
func HandleRecordUpdate(e *core.RecordEvent) error {
	// Log the record update
	log.Info("Record updated",
		"collection", e.Record.Collection().Name,
		"id", e.Record.Id,
		"updated", e.Record.GetDateTime("updated"),
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleRecordDelete handles record deletion events
func HandleRecordDelete(e *core.RecordEvent) error {
	// Log the record deletion

	log.Info("Record deleted",
		"collection", e.Record.Collection().Name,
		"id", e.Record.Id,
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleRecordAfterCreateSuccess handles successful record creation
func HandleRecordAfterCreateSuccess(e *core.RecordEvent) error {
	// This hook is triggered after the record is successfully persisted

	log.Debug("Record successfully persisted",
		"collection", e.Record.Collection().Name,
		"id", e.Record.Id,
	)

	// Add post-creation logic here

	return e.Next()
}

// HandleRecordAfterCreateError handles failed record creation
func HandleRecordAfterCreateError(e *core.RecordEvent) error {
	// This hook is triggered when record creation fails

	log.Error("Record creation failed",
		"collection", e.Record.Collection().Name,
		"error", fmt.Sprintf("%v", e),
	)

	// Add error handling logic here

	return e.Next()
}

// HandleUserCreate handles user-specific record creation
func HandleUserCreate(e *core.RecordEvent) error {
	// This is an example of collection-specific hook

	log.Info("New user created",
		"user_id", e.Record.Id,
		"email", e.Record.GetString("email"),
	)

	// Add user-specific logic here

	return e.Next()
}
