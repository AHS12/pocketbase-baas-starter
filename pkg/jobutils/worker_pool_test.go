package jobutils

import (
	"context"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func TestNewWorkerPool(t *testing.T) {
	app := pocketbase.New()
	registry := NewJobRegistry()

	tests := []struct {
		name       string
		maxWorkers int
		expected   int
	}{
		{
			name:       "default workers",
			maxWorkers: 3,
			expected:   3,
		},
		{
			name:       "zero workers defaults to 5",
			maxWorkers: 0,
			expected:   5,
		},
		{
			name:       "negative workers defaults to 5",
			maxWorkers: -1,
			expected:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewWorkerPool(app, registry, tt.maxWorkers)

			if pool == nil {
				t.Fatal("NewWorkerPool should not return nil")
			}

			if pool.maxWorkers != tt.expected {
				t.Errorf("expected %d workers, got %d", tt.expected, pool.maxWorkers)
			}

			if len(pool.workers) != tt.expected {
				t.Errorf("expected %d worker instances, got %d", tt.expected, len(pool.workers))
			}

			// Clean shutdown
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			_ = pool.Shutdown(ctx)
		})
	}
}

func TestWorkerPool_IsShutdown(t *testing.T) {
	app := pocketbase.New()
	registry := NewJobRegistry()
	pool := NewWorkerPool(app, registry, 2)

	// Initially not shutdown
	if pool.IsShutdown() {
		t.Error("expected pool to not be shutdown initially")
	}

	// After shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := pool.Shutdown(ctx)
	if err != nil {
		t.Errorf("unexpected shutdown error: %v", err)
	}

	if !pool.IsShutdown() {
		t.Error("expected pool to be shutdown after Shutdown() call")
	}
}

func TestWorkerPool_ProcessJobsEmpty(t *testing.T) {
	app := pocketbase.New()
	registry := NewJobRegistry()
	pool := NewWorkerPool(app, registry, 2)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_ = pool.Shutdown(ctx)
	}()

	errors := pool.ProcessJobs([]*core.Record{})
	if errors != nil {
		t.Errorf("expected nil errors for empty jobs, got %v", errors)
	}
}

func TestWorkerPool_ProcessJobsShutdown(t *testing.T) {
	app := pocketbase.New()
	registry := NewJobRegistry()
	pool := NewWorkerPool(app, registry, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_ = pool.Shutdown(ctx)

	if !pool.IsShutdown() {
		t.Error("expected pool to be shutdown")
	}
}

func TestWorkerPool_ProcessJobsConcurrently(t *testing.T) {
	app := pocketbase.New()
	registry := NewJobRegistry()
	pool := NewWorkerPool(app, registry, 3)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_ = pool.Shutdown(ctx)
	}()

	errors := pool.ProcessJobsConcurrently([]*core.Record{}, 5)
	if errors != nil {
		t.Errorf("expected nil errors for empty jobs, got %v", errors)
	}

	if pool.maxWorkers != 3 {
		t.Errorf("expected pool to maintain its configured worker count of 3, got %d", pool.maxWorkers)
	}
}

func TestWorkerPool_Shutdown(t *testing.T) {
	app := pocketbase.New()
	registry := NewJobRegistry()

	tests := []struct {
		name          string
		timeout       time.Duration
		expectError   bool
		shutdownTwice bool
	}{
		{
			name:        "normal shutdown",
			timeout:     5 * time.Second,
			expectError: false,
		},
		{
			name:          "shutdown twice",
			timeout:       5 * time.Second,
			expectError:   false,
			shutdownTwice: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh pool for each test
			testPool := NewWorkerPool(app, registry, 2)

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			err := testPool.Shutdown(ctx)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Test shutdown twice
			if tt.shutdownTwice {
				ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel2()

				err2 := testPool.Shutdown(ctx2)
				if err2 != nil {
					t.Errorf("second shutdown should not error, got: %v", err2)
				}
			}
		})
	}
}

func TestWorkerJobResult(t *testing.T) {
	result := WorkerJobResult{
		JobID: "test-job-123",
		Error: nil,
	}

	if result.JobID != "test-job-123" {
		t.Errorf("expected JobID 'test-job-123', got %q", result.JobID)
	}

	if result.Error != nil {
		t.Errorf("expected nil error, got %v", result.Error)
	}
}

func TestWorker(t *testing.T) {
	app := pocketbase.New()
	registry := NewJobRegistry()
	jobQueue := make(chan *core.Record, 10)
	resultQueue := make(chan WorkerJobResult, 10)
	quit := make(chan bool)

	worker := &Worker{
		id:          1,
		jobQueue:    jobQueue,
		resultQueue: resultQueue,
		quit:        quit,
		app:         app,
		registry:    registry,
	}

	if worker.id != 1 {
		t.Errorf("expected worker ID 1, got %d", worker.id)
	}

	if worker.app != app {
		t.Error("expected worker to have correct app reference")
	}

	if worker.registry != registry {
		t.Error("expected worker to have correct registry reference")
	}

	close(jobQueue)
	close(resultQueue)
	close(quit)
}
