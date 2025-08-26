package hook

import (
	log "ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase/core"
)

// HandleRealtimeConnect handles realtime connection events
func HandleRealtimeConnect(e *core.RealtimeConnectRequestEvent) error {

	log.Debug("Realtime client connected",
		"client_id", e.Client.Id(),
	)

	return e.Next()
}

// HandleRealtimeSubscribe handles realtime subscription events
func HandleRealtimeSubscribe(e *core.RealtimeSubscribeRequestEvent) error {

	log.Debug("Realtime subscription created",
		"client_id", e.Client.Id(),
		"subscriptions", len(e.Subscriptions),
	)

	return e.Next()
}

// HandleRealtimeMessage handles realtime message events
func HandleRealtimeMessage(e *core.RealtimeMessageEvent) error {

	log.Debug("Realtime message sent",
		"type", e.Message.Name,
		"data_size", len(e.Message.Data),
	)

	return e.Next()
}
