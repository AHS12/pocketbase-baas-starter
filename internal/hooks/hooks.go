package hooks

import (
	"fmt"
	"ims-pocketbase-baas-starter/internal/handlers/hook"
	log "ims-pocketbase-baas-starter/pkg/logger"
	"ims-pocketbase-baas-starter/pkg/metrics"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterHooks registers all custom event hooks
func RegisterHooks(app *pocketbase.PocketBase) error {
	if app == nil {
		return fmt.Errorf("RegisterHooks: app cannot be nil")
	}

	log.Info("Registering custom event hooks")

	// Register Record hooks
	if err := registerRecordHooks(app); err != nil {
		return fmt.Errorf("failed to register record hooks: %w", err)
	}

	// Register Collection hooks
	if err := registerCollectionHooks(app); err != nil {
		return fmt.Errorf("failed to register collection hooks: %w", err)
	}

	// Register Request hooks
	if err := registerRequestHooks(app); err != nil {
		return fmt.Errorf("failed to register request hooks: %w", err)
	}

	// Register Mailer hooks
	if err := registerMailerHooks(app); err != nil {
		return fmt.Errorf("failed to register mailer hooks: %w", err)
	}

	// Register Realtime hooks
	if err := registerRealtimeHooks(app); err != nil {
		return fmt.Errorf("failed to register realtime hooks: %w", err)
	}

	log.Info("Custom event hooks registration completed")
	return nil
}

// registerRecordHooks registers all record-related event hooks
func registerRecordHooks(app *pocketbase.PocketBase) error {
	// Example: Log all record creations
	app.OnRecordCreate().BindFunc(func(e *core.RecordEvent) error {
		return hook.HandleRecordCreate(e)
	})

	// Example: Log all record updates
	app.OnRecordUpdate().BindFunc(func(e *core.RecordEvent) error {
		return hook.HandleRecordUpdate(e)
	})

	// Example: Handle record deletions
	app.OnRecordDelete().BindFunc(func(e *core.RecordEvent) error {
		return hook.HandleRecordDelete(e)
	})

	// Example: Collection-specific hooks
	// app.OnRecordCreate("users").BindFunc(func(e *core.RecordEvent) error {
	//     return hook.HandleUserCreate(e)
	// })

	// Example: Additional hook registrations (uncomment to enable)
	// app.OnRecordAfterCreateSuccess().BindFunc(func(e *core.RecordEvent) error {
	//     return hook.HandleAuditLog(e)
	// })

	// app.OnRecordCreateRequest().BindFunc(func(e *core.RecordCreateRequestEvent) error {
	//     return hook.HandleDataValidation(&core.RecordEvent{
	//         App:    e.App,
	//         Record: e.Record,
	//     })
	// })

	// app.OnRecordUpdate().BindFunc(func(e *core.RecordEvent) error {
	//     return hook.HandleCacheInvalidation(e)
	// })

	// Send welcome email to new users
	app.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		return hook.HandleUserWelcomeEmail(e)
	})

	// Create user default settings (with metrics instrumentation)
	app.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		metricsProvider := metrics.GetInstance()

		return metrics.InstrumentHook(metricsProvider, "user_create_settings", func() error {
			return hook.HandleUserCreateSettings(e)
		})
	})

	// Invalidate user permission cache when user is updated
	app.OnRecordUpdate("users").BindFunc(func(e *core.RecordEvent) error {
		return hook.HandleUserCacheClear(e)
	})

	log.Debug("Record hooks registered")
	return nil
}

// registerCollectionHooks registers all collection-related event hooks
func registerCollectionHooks(app *pocketbase.PocketBase) error {
	// Example: Log collection creations
	app.OnCollectionCreate().BindFunc(func(e *core.CollectionEvent) error {
		return hook.HandleCollectionCreate(e)
	})

	// Example: Log collection updates
	app.OnCollectionUpdate().BindFunc(func(e *core.CollectionEvent) error {
		return hook.HandleCollectionUpdate(e)
	})

	log.Debug("Collection hooks registered")
	return nil
}

// registerRequestHooks registers all request-related event hooks
func registerRequestHooks(app *pocketbase.PocketBase) error {
	// Example: Log all record list requests
	app.OnRecordsListRequest().BindFunc(func(e *core.RecordsListRequestEvent) error {
		return hook.HandleRecordListRequest(e)
	})

	// Example: Log all record view requests
	app.OnRecordViewRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return hook.HandleRecordViewRequest(e)
	})

	// Example: Collection-specific request hooks
	// app.OnRecordListRequest("users").BindFunc(func(e *core.RecordListRequestEvent) error {
	//     return hook.HandleUserListRequest(e)
	// })

	log.Debug("Request hooks registered")
	return nil
}

// registerMailerHooks registers all mailer-related event hooks
func registerMailerHooks(app *pocketbase.PocketBase) error {
	// Example: Log all email sends (with metrics instrumentation)
	app.OnMailerSend().BindFunc(func(e *core.MailerEvent) error {
		// Get the metrics provider instance
		metricsProvider := metrics.GetInstance()

		// Instrument the email operation with metrics collection
		return metrics.InstrumentEmailOperation(metricsProvider, func() error {
			return hook.HandleMailerSend(e)
		})
	})

	log.Debug("Mailer hooks registered")
	return nil
}

// registerRealtimeHooks registers all realtime-related event hooks
func registerRealtimeHooks(app *pocketbase.PocketBase) error {
	// Example: Log realtime connections
	app.OnRealtimeConnectRequest().BindFunc(func(e *core.RealtimeConnectRequestEvent) error {
		return hook.HandleRealtimeConnect(e)
	})

	// Example: Log realtime disconnections
	app.OnRealtimeSubscribeRequest().BindFunc(func(e *core.RealtimeSubscribeRequestEvent) error {
		return hook.HandleRealtimeSubscribe(e)
	})

	log.Debug("Realtime hooks registered")
	return nil
}
