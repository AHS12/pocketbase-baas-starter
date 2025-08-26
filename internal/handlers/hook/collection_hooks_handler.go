package hook

import (
	log "ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase/core"
)

// HandleCollectionCreate handles collection creation events
func HandleCollectionCreate(e *core.CollectionEvent) error {

	log.Info("Collection created",
		"name", e.Collection.Name,
		"id", e.Collection.Id,
		"type", e.Collection.Type,
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleCollectionUpdate handles collection update events
func HandleCollectionUpdate(e *core.CollectionEvent) error {
	// Log the collection update

	log.Info("Collection updated",
		"name", e.Collection.Name,
		"id", e.Collection.Id,
		"type", e.Collection.Type,
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}

// HandleCollectionDelete handles collection deletion events
func HandleCollectionDelete(e *core.CollectionEvent) error {
	// Log the collection deletion

	log.Info("Collection deleted",
		"name", e.Collection.Name,
		"id", e.Collection.Id,
	)

	// Add your custom logic here

	// Continue with the execution chain
	return e.Next()
}
