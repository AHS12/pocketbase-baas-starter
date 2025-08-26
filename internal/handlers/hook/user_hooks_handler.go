package hook

import (
	"encoding/json"
	"fmt"
	"time"

	"ims-pocketbase-baas-starter/pkg/cache"
	"ims-pocketbase-baas-starter/pkg/common"
	"ims-pocketbase-baas-starter/pkg/jobutils"
	log "ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase/core"
)

// HandleUserWelcomeEmail handles sending a welcome email to new users
func HandleUserWelcomeEmail(e *core.RecordEvent) error {
	appName := common.GetEnv("APP_NAME", "N/A")
	settings := e.App.Settings()
	if settings.Meta.AppName != "" {
		appName = settings.Meta.AppName
	}

	email := e.Record.GetString("email")
	name := e.Record.GetString("name")
	appUrl := common.GetEnv("APP_URL", "N/A")
	if settings.Meta.AppURL != "" {
		appUrl = settings.Meta.AppURL
	}
	if name == "" {
		name = email // Use email as name if name is not provided
	}

	payload := jobutils.EmailJobPayload{
		Type: jobutils.JobTypeEmail,
		Data: jobutils.EmailJobData{
			To:       email,
			Subject:  fmt.Sprintf("Welcome to %s!", appName),
			Template: "welcome",
			Variables: map[string]any{
				"AppName": appName,
				"Name":    name,
				"Email":   email,
				"AppURL":  appUrl,
				"Year":    time.Now().Year(),
			},
		},
		Options: jobutils.EmailJobOptions{
			RetryCount: 3,
			Timeout:    30,
		},
	}

	collection, err := e.App.FindCollectionByNameOrId("queues")
	if err != nil {
		log.Error("Failed to find queues collection", "error", err)
		return err
	}

	jobRecord := core.NewRecord(collection)
	jobRecord.Set("name", fmt.Sprintf("Welcome email for %s", email))
	jobRecord.Set("description", fmt.Sprintf("Send welcome email to new user %s", email))

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Error("Failed to marshal email payload", "error", err)
		return err
	}
	jobRecord.Set("payload", string(payloadBytes))
	jobRecord.Set("attempts", 0)

	if err := e.App.Save(jobRecord); err != nil {
		log.Error("Failed to queue welcome email job", "error", err)
		return err
	}

	log.Info("Welcome email job queued successfully",
		"user_id", e.Record.Id,
		"email", email,
		"job_id", jobRecord.Id)

	return e.Next()
}

// HandleUserCreateSettings generate default user settings
func HandleUserCreateSettings(e *core.RecordEvent) error {

	log.Info("Creating default settings for new user",
		"user_id", e.Record.Id,
		"email", e.Record.GetString("email"),
	)

	userSettingsCollection, err := e.App.FindCollectionByNameOrId("user_settings")
	if err != nil {
		log.Error("user_settings collection not found", "error", err)
		// Continue without failing if settings collection doesn't exist
		return e.Next()
	}

	// Define default user settings with their values
	defaultUserSettings := []struct {
		SettingSlug string
		Value       string
	}{
		{"theme", "light"},
		{"notifications", "true"},
	}

	for _, defaultSetting := range defaultUserSettings {
		settingRecord, err := e.App.FindFirstRecordByFilter("settings", "slug = {:slug}", map[string]any{
			"slug": defaultSetting.SettingSlug,
		})
		if err != nil {
			log.Warn("Setting not found, skipping",
				"slug", defaultSetting.SettingSlug,
				"error", err)
			continue
		}

		userSettingRecord := core.NewRecord(userSettingsCollection)
		userSettingData := map[string]any{
			"user":     e.Record.Id,
			"settings": settingRecord.Id,
			"value":    defaultSetting.Value,
		}

		userSettingRecord.Load(userSettingData)

		if err := e.App.Save(userSettingRecord); err != nil {
			log.Error("Failed to create user setting",
				"user_id", e.Record.Id,
				"setting_slug", defaultSetting.SettingSlug,
				"error", err)
			continue
		}
	}

	log.Info("Default user settings creation completed",
		"user_id", e.Record.Id)

	return e.Next()
}

// HandleUserCacheClear handles clearing user-related cache when a user is updated
func HandleUserCacheClear(e *core.RecordEvent) error {
	cacheService := cache.GetInstance()
	cacheKey := cache.CacheKey{}.UserPermissions(e.Record.Id)
	cacheService.Delete(cacheKey)

	return e.Next()
}
