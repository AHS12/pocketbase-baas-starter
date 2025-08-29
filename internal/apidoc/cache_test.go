package apidoc

import (
	"ims-pocketbase-baas-starter/pkg/cache"
	"testing"
	"time"
)

func TestCacheBasicFunctionality(t *testing.T) {
	InvalidateCache()

	if _, found := cache.GetInstance().Get(SwaggerSpecKey); found {
		t.Error("Expected swagger spec cache to be cleared after invalidation")
	}

	if _, found := cache.GetInstance().Get(SwaggerCollectionsHash); found {
		t.Error("Expected collections hash cache to be cleared after invalidation")
	}
}

func TestCacheStatus(t *testing.T) {
	generator := NewGenerator(nil, DefaultConfig())

	status := GetCacheStatus(generator)

	if status["cached"].(bool) {
		t.Error("Expected cached to be false initially")
	}

	// Check that cache_stats is included
	if _, exists := status["cache_stats"]; !exists {
		t.Error("Expected cache_stats to be included in status")
	}
}

func TestCacheTTL(t *testing.T) {
	originalTTL := cacheTTL
	defer func() { cacheTTL = originalTTL }()

	cacheTTL = 1 * time.Millisecond

	cache.GetInstance().SetWithExpiration(SwaggerSpecKey, &CombinedOpenAPISpec{}, 1*time.Millisecond)

	time.Sleep(2 * time.Millisecond)

	// Check that cache is expired
	if _, found := cache.GetInstance().Get(SwaggerSpecKey); found {
		t.Error("Expected cache to be expired after TTL")
	}
}

func TestCollectionChanges(t *testing.T) {
	generator := NewGenerator(nil, DefaultConfig())

	changed := hasCollectionsChanged(generator)

	// Should be true initially (no previous hash)
	if !changed {
		t.Error("Expected collections to be considered changed initially")
	}

	invalidated := CheckAndInvalidateIfChanged(generator)

	// Should be true since collections changed
	if !invalidated {
		t.Error("Expected cache to be invalidated due to collection changes")
	}
}

func TestCacheService(t *testing.T) {
	cacheService := cache.GetInstance()
	if cacheService == nil {
		t.Error("Expected to get the centralized cache service")
	}

	expectedSpecKey := "swagger_spec"
	if SwaggerSpecKey != expectedSpecKey {
		t.Errorf("Expected spec key to be '%s', got '%s'", expectedSpecKey, SwaggerSpecKey)
	}

	expectedHashKey := "swagger_collections_hash"
	if SwaggerCollectionsHash != expectedHashKey {
		t.Errorf("Expected hash key to be '%s', got '%s'", expectedHashKey, SwaggerCollectionsHash)
	}
}

func TestGetCachedSpec(t *testing.T) {
	InvalidateCache()

	// Should return nil when no cache exists
	spec := getCachedSpec()
	if spec != nil {
		t.Error("Expected nil when no cache exists")
	}

	testSpec := &CombinedOpenAPISpec{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	}
	cache.GetInstance().SetWithExpiration(SwaggerSpecKey, testSpec, 5*time.Minute)

	cachedSpec := getCachedSpec()
	if cachedSpec == nil {
		t.Error("Expected cached spec to be returned")
		return // Add early return to prevent nil pointer dereference
	}

	if cachedSpec.OpenAPI != "3.0.0" {
		t.Error("Expected cached spec to have correct OpenAPI version")
	}
}
