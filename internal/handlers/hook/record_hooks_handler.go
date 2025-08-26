package hook

import (
	"fmt"
	log "ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase/core"
)

// HandleRecordCreate handles record creation events
func HandleRecordCreate(e *core.RecordEvent) error {
	log.Info("Record created",
		"collection", e.Record.Collection().Name,
		"id", e.Record.Id,
		"created", e.Record.GetDateTime("created"),
	)

	return e.Next()
}

// HandleRecordUpdate handles record update events
func HandleRecordUpdate(e *core.RecordEvent) error {
	log.Info("Record updated",
		"collection", e.Record.Collection().Name,
		"id", e.Record.Id,
		"updated", e.Record.GetDateTime("updated"),
	)

	return e.Next()
}

// HandleRecordDelete handles record deletion events
func HandleRecordDelete(e *core.RecordEvent) error {
	log.Info("Record deleted",
		"collection", e.Record.Collection().Name,
		"id", e.Record.Id,
	)

	return e.Next()
}

// HandleRecordAfterCreateSuccess handles successful record creation
func HandleRecordAfterCreateSuccess(e *core.RecordEvent) error {
	log.Info("Record successfully persisted",
		"collection", e.Record.Collection().Name,
		"id", e.Record.Id,
	)

	return e.Next()
}

// HandleRecordAfterCreateError handles failed record creation
func HandleRecordAfterCreateError(e *core.RecordEvent) error {
	log.Error("Record creation failed",
		"collection", e.Record.Collection().Name,
		"error", fmt.Sprintf("%v", e),
	)

	return e.Next()
}

// HandleUserCreate handles user-specific record creation
func HandleUserCreate(e *core.RecordEvent) error {
	log.Info("New user created",
		"user_id", e.Record.Id,
		"email", e.Record.GetString("email"),
	)

	return e.Next()
}
