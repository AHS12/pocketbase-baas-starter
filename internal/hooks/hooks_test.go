package hooks

import (
	"ims-pocketbase-baas-starter/pkg/metrics"
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestRegisterHooks(t *testing.T) {
	app := pocketbase.New()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RegisterHooks panicked: %v", r)
		}
	}()

	// Register hooks
	err := RegisterHooks(app)
	if err != nil {
		t.Fatalf("RegisterHooks failed: %v", err)
	}

	// Verify that hooks are registered by checking if the hook exists
	// Note: PocketBase doesn't provide direct access to registered hooks,
	// so we mainly test that registration doesn't fail

	// Test individual registration functions
	if err := registerRecordHooks(app); err != nil {
		t.Fatalf("registerRecordHooks failed: %v", err)
	}
	if err := registerCollectionHooks(app); err != nil {
		t.Fatalf("registerCollectionHooks failed: %v", err)
	}
	if err := registerRequestHooks(app); err != nil {
		t.Fatalf("registerRequestHooks failed: %v", err)
	}
	if err := registerMailerHooks(app); err != nil {
		t.Fatalf("registerMailerHooks failed: %v", err)
	}
	if err := registerRealtimeHooks(app); err != nil {
		t.Fatalf("registerRealtimeHooks failed: %v", err)
	}
}

func TestHookRegistrationFunctions(t *testing.T) {
	app := pocketbase.New()

	tests := []struct {
		name string
		fn   func(*pocketbase.PocketBase) error
	}{
		{"registerRecordHooks", registerRecordHooks},
		{"registerCollectionHooks", registerCollectionHooks},
		{"registerRequestHooks", registerRequestHooks},
		{"registerMailerHooks", registerMailerHooks},
		{"registerRealtimeHooks", registerRealtimeHooks},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("%s panicked: %v", tt.name, r)
				}
			}()

			err := tt.fn(app)
			if err != nil {
				t.Fatalf("%s failed: %v", tt.name, err)
			}
		})
	}
}

func TestHooksWithMetricsInstrumentation(t *testing.T) {
	app := pocketbase.New()

	metrics.InitializeProvider(metrics.Config{
		Provider: metrics.ProviderDisabled,
		Enabled:  false,
	})

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Hooks with metrics instrumentation panicked: %v", r)
		}
	}()

	err := registerRecordHooks(app) // Contains user_create_settings instrumentation
	if err != nil {
		t.Fatalf("registerRecordHooks failed: %v", err)
	}

	err = registerMailerHooks(app) // Contains email operation instrumentation
	if err != nil {
		t.Fatalf("registerMailerHooks failed: %v", err)
	}

	metrics.Reset()
}
