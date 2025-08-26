package hook

import (
	log "ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase/core"
)

// HandleRecordListRequest handles record list request events
func HandleRecordListRequest(e *core.RecordsListRequestEvent) error {

	log.Debug("Record list requested",
		"collection", e.Collection.Name,
		"user_ip", e.Request.RemoteAddr,
		"user_agent", e.Request.UserAgent(),
	)

	return e.Next()
}

// HandleRecordViewRequest handles record view request events
func HandleRecordViewRequest(e *core.RecordRequestEvent) error {

	log.Debug("Record view requested",
		"collection", e.Collection.Name,
		"record_id", e.Record.Id,
		"user_ip", e.Request.RemoteAddr,
	)

	return e.Next()
}

// HandleRecordCreateRequest handles record create request events
func HandleRecordCreateRequest(e *core.RecordRequestEvent) error {

	log.Debug("Record create requested",
		"collection", e.Collection.Name,
		"user_ip", e.Request.RemoteAddr,
	)

	return e.Next()
}

// HandleRecordUpdateRequest handles record update request events
func HandleRecordUpdateRequest(e *core.RecordRequestEvent) error {

	log.Debug("Record update requested",
		"collection", e.Collection.Name,
		"record_id", e.Record.Id,
		"user_ip", e.Request.RemoteAddr,
	)

	return e.Next()
}

// HandleRecordDeleteRequest handles record delete request events
func HandleRecordDeleteRequest(e *core.RecordRequestEvent) error {

	log.Debug("Record delete requested",
		"collection", e.Collection.Name,
		"record_id", e.Record.Id,
		"user_ip", e.Request.RemoteAddr,
	)

	return e.Next()
}

// HandleUserListRequest handles user-specific list requests
func HandleUserListRequest(e *core.RecordsListRequestEvent) error {
	log.Debug("User list requested",
		"user_ip", e.Request.RemoteAddr,
		"query_params", e.Request.URL.RawQuery,
	)

	return e.Next()
}
