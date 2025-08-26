package apidoc

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"ims-pocketbase-baas-starter/pkg/cache"
	log "ims-pocketbase-baas-starter/pkg/logger"
	"sync"
	"time"
)

// Cache keys for swagger
const (
	SwaggerSpecKey         = "swagger_spec"
	SwaggerCollectionsHash = "swagger_collections_hash"
)

var (
	cacheTTL   = 5 * time.Minute
	generating sync.Mutex
)

// GenerateSpecWithCache generates OpenAPI spec with centralized caching and automatic invalidation
func GenerateSpecWithCache(generator *Generator) (*CombinedOpenAPISpec, error) {
	if spec := getCachedSpec(); spec != nil {
		if !hasCollectionsChanged(generator) {
			return spec, nil
		}
		InvalidateCache()
	}

	generating.Lock()
	defer generating.Unlock()

	if spec := getCachedSpec(); spec != nil {
		if !hasCollectionsChanged(generator) {
			return spec, nil
		}
	}

	spec, err := generator.GenerateSpec()
	if err != nil {
		return nil, err
	}

	cache.GetInstance().SetWithExpiration(SwaggerSpecKey, spec, cacheTTL)
	updateCollectionsHash(generator)

	return spec, nil
}

// getCachedSpec retrieves the cached OpenAPI spec if it exists
func getCachedSpec() *CombinedOpenAPISpec {
	if cachedSpec, found := cache.GetInstance().Get(SwaggerSpecKey); found {
		if spec, ok := cachedSpec.(*CombinedOpenAPISpec); ok {
			return spec
		}
	}
	return nil
}

// InvalidateCache clears the cached spec and collections hash
func InvalidateCache() {
	cache.GetInstance().Delete(SwaggerSpecKey)
	cache.GetInstance().Delete(SwaggerCollectionsHash)
}

// GetCacheStatus returns cache information including collection change detection
func GetCacheStatus(generator *Generator) map[string]any {
	_, specCached := cache.GetInstance().Get(SwaggerSpecKey)
	collectionsHash, _ := cache.GetInstance().GetString(SwaggerCollectionsHash)

	status := map[string]any{
		"cached":              specCached,
		"cache_ttl":           cacheTTL.String(),
		"collections_hash":    collectionsHash,
		"collections_changed": hasCollectionsChanged(generator),
		"cache_stats":         cache.GetInstance().GetStats(),
	}

	return status
}

// hasCollectionsChanged checks if collections have changed since last cache
func hasCollectionsChanged(generator *Generator) bool {
	currentHash, err := generateCollectionsHash(generator)
	if err != nil {
		log.Warn("Failed to generate collections hash", "error", err)
		return true
	}

	cachedHash, found := cache.GetInstance().GetString(SwaggerCollectionsHash)

	if !found || cachedHash == "" {
		return true
	}

	return currentHash != cachedHash
}

// generateCollectionsHash creates a hash of collection metadata for change detection
func generateCollectionsHash(generator *Generator) (string, error) {
	collections, err := generator.discovery.DiscoverCollections()
	if err != nil {
		return "", fmt.Errorf("failed to discover collections for hashing: %w", err)
	}

	jsonData, err := json.Marshal(collections)
	if err != nil {
		return "", fmt.Errorf("failed to marshal collections for hashing: %w", err)
	}

	hash := md5.Sum(jsonData)
	return fmt.Sprintf("%x", hash), nil
}

// updateCollectionsHash updates the stored collections hash
func updateCollectionsHash(generator *Generator) {
	hash, err := generateCollectionsHash(generator)
	if err != nil {
		log.Warn("Failed to generate collections hash", "error", err)
		return
	}

	cache.GetInstance().SetWithExpiration(SwaggerCollectionsHash, hash, cacheTTL)
}

// CheckAndInvalidateIfChanged checks for collection changes and invalidates cache if needed
func CheckAndInvalidateIfChanged(generator *Generator) bool {
	if hasCollectionsChanged(generator) {
		InvalidateCache()
		return true
	}
	return false
}
