package hook

import (
	log "ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase/core"
)

// HandleRecordListRequest handles record list request events
func HandleRecordListRequest(e *core.RecordsListRequestEvent) error {
	// Log the record list request

	log.Debug("Record list requested",
		"collection", e.Collection.Name,
		"user_ip", e.Request.RemoteAddr,
		"user_agent", e.Request.UserAgent(),
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleRecordViewRequest handles record view request events
func HandleRecordViewRequest(e *core.RecordRequestEvent) error {
	// Log the record view request

	log.Debug("Record view requested",
		"collection", e.Collection.Name,
		"record_id", e.Record.Id,
		"user_ip", e.Request.RemoteAddr,
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleRecordCreateRequest handles record create request events
func HandleRecordCreateRequest(e *core.RecordRequestEvent) error {
	// Log the record create request

	log.Debug("Record create requested",
		"collection", e.Collection.Name,
		"user_ip", e.Request.RemoteAddr,
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleRecordUpdateRequest handles record update request events
func HandleRecordUpdateRequest(e *core.RecordRequestEvent) error {
	// Log the record update request

	log.Debug("Record update requested",
		"collection", e.Collection.Name,
		"record_id", e.Record.Id,
		"user_ip", e.Request.RemoteAddr,
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleRecordDeleteRequest handles record delete request events
func HandleRecordDeleteRequest(e *core.RecordRequestEvent) error {
	// Log the record delete request

	log.Debug("Record delete requested",
		"collection", e.Collection.Name,
		"record_id", e.Record.Id,
		"user_ip", e.Request.RemoteAddr,
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleUserListRequest handles user-specific list requests
func HandleUserListRequest(e *core.RecordsListRequestEvent) error {
	// This is an example of collection-specific request hook

	log.Debug("User list requested",
		"user_ip", e.Request.RemoteAddr,
		"query_params", e.Request.URL.RawQuery,
	)

	// Add user-specific logic here

	return e.Next()
}
