package middlewares

import (
	"testing"

	"ims-pocketbase-baas-starter/pkg/metrics"

	"github.com/pocketbase/pocketbase/core"
)

func TestNewMetricsMiddleware(t *testing.T) {
	provider := metrics.NewNoOpProvider()
	middleware := NewMetricsMiddleware(provider)

	if middleware == nil {
		t.Fatal("NewMetricsMiddleware() returned nil")
	}

	if middleware.provider != provider {
		t.Error("Expected provider to be set correctly")
	}
}

func TestNewMetricsMiddlewareWithNilProvider(t *testing.T) {
	middleware := NewMetricsMiddleware(nil)

	if middleware == nil {
		t.Fatal("NewMetricsMiddleware() with nil provider returned nil")
	}

	if middleware.provider != nil {
		t.Error("Expected nil provider to be preserved")
	}
}

func TestRequireMetricsFunc(t *testing.T) {
	provider := metrics.NewNoOpProvider()
	middleware := NewMetricsMiddleware(provider)

	metricsFunc := middleware.RequireMetricsFunc()
	if metricsFunc == nil {
		t.Fatal("RequireMetricsFunc() returned nil function")
	}
}

func TestRequireMetricsFuncWithNilProvider(t *testing.T) {
	middleware := NewMetricsMiddleware(nil)

	metricsFunc := middleware.RequireMetricsFunc()
	if metricsFunc == nil {
		t.Fatal("RequireMetricsFunc() with nil provider returned nil function")
	}
}

func TestRequireMetrics(t *testing.T) {
	provider := metrics.NewNoOpProvider()

	metricsFunc := RequireMetrics(provider)
	if metricsFunc == nil {
		t.Fatal("RequireMetrics() returned nil function")
	}
}

func TestRequireMetricsWithNilProvider(t *testing.T) {
	metricsFunc := RequireMetrics(nil)
	if metricsFunc == nil {
		t.Fatal("RequireMetrics() with nil provider returned nil function")
	}
}

func TestInstrumentHandler(t *testing.T) {
	provider := metrics.NewNoOpProvider()

	testHandler := func(e *core.RequestEvent) error {
		return nil
	}

	instrumentedHandler := InstrumentHandler(provider, "test_handler", testHandler)
	if instrumentedHandler == nil {
		t.Fatal("InstrumentHandler() returned nil function")
	}
}

func TestInstrumentHandlerWithNilProvider(t *testing.T) {
	testHandler := func(e *core.RequestEvent) error {
		return nil
	}

	instrumentedHandler := InstrumentHandler(nil, "test_handler", testHandler)
	if instrumentedHandler == nil {
		t.Fatal("InstrumentHandler() with nil provider returned nil function")
	}
}

func TestHelperFunctions(t *testing.T) {
	provider := metrics.NewNoOpProvider()
	labels := map[string]string{"test": "value"}

	RecordCustomMetric(provider, "custom_metric", 1.5, labels)
	IncrementCustomCounter(provider, "custom_counter", labels)
	SetCustomGauge(provider, "custom_gauge", 42.0, labels)

	RecordCustomMetric(nil, "custom_metric", 1.5, labels)
	IncrementCustomCounter(nil, "custom_counter", labels)
	SetCustomGauge(nil, "custom_gauge", 42.0, labels)
}
