package route

import (
	"ims-pocketbase-baas-starter/pkg/cache"
	"ims-pocketbase-baas-starter/pkg/response"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

// HandleCacheStatus returns the current status of the global cache store
func HandleCacheStatus(e *core.RequestEvent) error {
	cacheService := cache.GetInstance()
	stats := cacheService.GetStats()

	return response.OK(e, "Cache status retrieved successfully", map[string]any{
		"status": "ok",
		"stats":  stats,
	})
}

// HandleCacheClear clears all cache entries in the system
func HandleCacheClear(e *core.RequestEvent) error {
	cacheService := cache.GetInstance()
	cacheService.Flush()

	return response.OK(e, "Cache cleared successfully", map[string]any{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
