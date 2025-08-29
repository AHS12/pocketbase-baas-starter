package cache

import (
	"testing"
	"time"
)

func TestCacheService(t *testing.T) {
	cache := NewCacheService(CacheConfig{
		DefaultExpiration: 1 * time.Second,
		CleanupInterval:   2 * time.Second,
	})

	cache.Set("test_key", "test_value")

	if value, found := cache.GetString("test_key"); !found || value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s', found: %v", value, found)
	}

	testSlice := []string{"item1", "item2", "item3"}
	cache.Set("test_slice", testSlice)

	if slice, found := cache.GetStringSlice("test_slice"); !found || len(slice) != 3 {
		t.Errorf("Expected slice of length 3, got length %d, found: %v", len(slice), found)
	}

	cache.SetWithExpiration("expire_key", "expire_value", 100*time.Millisecond)

	if _, found := cache.Get("expire_key"); !found {
		t.Error("Expected key to exist immediately after setting")
	}

	time.Sleep(150 * time.Millisecond)

	if _, found := cache.Get("expire_key"); found {
		t.Error("Expected key to be expired")
	}

	cache.Set("delete_key", "delete_value")
	cache.Delete("delete_key")

	if _, found := cache.Get("delete_key"); found {
		t.Error("Expected key to be deleted")
	}

	cache.Set("pattern_test_1", "value1")
	cache.Set("pattern_test_2", "value2")
	cache.Set("other_key", "other_value")

	deleted := cache.DeletePattern("pattern_test_")
	if deleted != 2 {
		t.Errorf("Expected 2 keys deleted, got %d", deleted)
	}

	if _, found := cache.Get("pattern_test_1"); found {
		t.Error("Expected pattern_test_1 to be deleted")
	}

	if _, found := cache.Get("other_key"); !found {
		t.Error("Expected other_key to still exist")
	}
}

func TestCacheKeys(t *testing.T) {
	key := CacheKey{}

	userKey := key.UserPermissions("user123")
	expected := "user_permissions_user123"
	if userKey != expected {
		t.Errorf("Expected '%s', got '%s'", expected, userKey)
	}

	roleKey := key.RolePermissions("role456")
	expected = "role_permissions_role456"
	if roleKey != expected {
		t.Errorf("Expected '%s', got '%s'", expected, roleKey)
	}

	swaggerKey := key.SwaggerSpec()
	expected = "swagger_spec"
	if swaggerKey != expected {
		t.Errorf("Expected '%s', got '%s'", expected, swaggerKey)
	}
}

func TestCacheStats(t *testing.T) {
	cache := NewCacheService(CacheConfig{
		DefaultExpiration: 1 * time.Minute,
		CleanupInterval:   2 * time.Minute,
	})

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	stats := cache.GetStats()

	if itemCount, ok := stats["item_count"]; !ok {
		t.Error("Expected item_count in stats")
	} else if count, ok := itemCount.(int); !ok || count != 2 {
		t.Errorf("Expected item_count to be 2, got %v", itemCount)
	}

	if _, ok := stats["items"]; !ok {
		t.Error("Expected items in stats")
	}
}

func TestSingletonInstance(t *testing.T) {
	instance1 := GetInstance()
	instance2 := GetInstance()

	if instance1 != instance2 {
		t.Error("Expected GetInstance to return the same singleton instance")
	}

	instance1.Set("singleton_test", "test_value")

	if value, found := instance2.GetString("singleton_test"); !found || value != "test_value" {
		t.Error("Expected singleton instances to share the same cache")
	}
}
