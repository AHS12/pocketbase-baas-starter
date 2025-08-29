package jobutils

import (
	"ims-pocketbase-baas-starter/pkg/cronutils"
	"strings"
	"testing"
)

func TestValidateJobPayload(t *testing.T) {
	tests := []struct {
		name        string
		payload     map[string]any
		expectError bool
	}{
		{
			name:        "nil payload",
			payload:     nil,
			expectError: true,
		},
		{
			name:        "missing type field",
			payload:     map[string]any{},
			expectError: true,
		},
		{
			name: "valid payload with type",
			payload: map[string]any{
				"type": "test_job",
			},
			expectError: false,
		},
		{
			name: "valid payload with type and data",
			payload: map[string]any{
				"type": "test_job",
				"data": map[string]any{"key": "value"},
			},
			expectError: false,
		},
		{
			name: "invalid type field - not string",
			payload: map[string]any{
				"type": 123,
			},
			expectError: true,
		},
		{
			name: "empty type field",
			payload: map[string]any{
				"type": "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJobPayload(tt.payload)
			if tt.expectError && err == nil {
				t.Errorf("expected error for %q, but got none", tt.name)
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error for %q: %v", tt.name, err)
			}
		})
	}
}

func TestNewJobRegistry(t *testing.T) {
	registry := NewJobRegistry()
	if registry == nil {
		t.Error("NewJobRegistry should not return nil")
	}

	handlers := registry.ListHandlers()
	if len(handlers) != 0 {
		t.Errorf("new registry should have 0 handlers, got %d", len(handlers))
	}
}

type MockJobHandler struct {
	jobType string
	err     error
}

func (m *MockJobHandler) Handle(ctx *cronutils.CronExecutionContext, job *JobData) error {
	return m.err
}

func (m *MockJobHandler) GetJobType() string {
	return m.jobType
}

func TestJobRegistry_Register(t *testing.T) {
	tests := []struct {
		name        string
		handler     JobHandler
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil handler",
			handler:     nil,
			expectError: true,
			errorMsg:    "job handler cannot be nil",
		},
		{
			name:        "empty job type",
			handler:     &MockJobHandler{jobType: ""},
			expectError: true,
			errorMsg:    "job handler must return a non-empty job type",
		},
		{
			name:        "valid handler",
			handler:     &MockJobHandler{jobType: "test_job"},
			expectError: false,
		},
		{
			name:        "duplicate handler",
			handler:     &MockJobHandler{jobType: "test_job"}, // Same type as above
			expectError: true,
			errorMsg:    "job handler for type 'test_job' is already registered",
		},
	}

	registry := NewJobRegistry()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.Register(tt.handler)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestJobRegistry_GetHandler(t *testing.T) {
	registry := NewJobRegistry()
	handler := &MockJobHandler{jobType: "test_job"}

	// Register handler
	err := registry.Register(handler)
	if err != nil {
		t.Fatalf("failed to register handler: %v", err)
	}

	tests := []struct {
		name        string
		jobType     string
		expectError bool
	}{
		{
			name:        "existing handler",
			jobType:     "test_job",
			expectError: false,
		},
		{
			name:        "non-existing handler",
			jobType:     "unknown_job",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := registry.GetHandler(tt.jobType)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if result != nil {
					t.Errorf("expected nil handler but got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected handler but got nil")
				}
			}
		})
	}
}

func TestJobRegistry_ListHandlers(t *testing.T) {
	registry := NewJobRegistry()

	handlers := registry.ListHandlers()
	if len(handlers) != 0 {
		t.Errorf("expected 0 handlers, got %d", len(handlers))
	}

	handler1 := &MockJobHandler{jobType: "job1"}
	handler2 := &MockJobHandler{jobType: "job2"}

	_ = registry.Register(handler1)
	_ = registry.Register(handler2)

	handlers = registry.ListHandlers()
	if len(handlers) != 2 {
		t.Errorf("expected 2 handlers, got %d", len(handlers))
	}

	typeMap := make(map[string]bool)
	for _, jobType := range handlers {
		typeMap[jobType] = true
	}

	if !typeMap["job1"] || !typeMap["job2"] {
		t.Errorf("expected both job1 and job2 in handlers list, got %v", handlers)
	}
}

func TestParseJobDataFromRecord(t *testing.T) {
	_, err := ParseJobDataFromRecord(nil)
	if err == nil {
		t.Error("expected error for nil record")
	}
}

func TestValidateJobPayloadExtended(t *testing.T) {
	tests := []struct {
		name        string
		payload     map[string]any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil payload",
			payload:     nil,
			expectError: true,
			errorMsg:    "job payload cannot be nil",
		},
		{
			name:        "missing type field",
			payload:     map[string]any{},
			expectError: true,
			errorMsg:    "job payload must contain a 'type' field",
		},
		{
			name: "invalid data field type",
			payload: map[string]any{
				"type": "test_job",
				"data": "invalid_data_type", // Should be object
			},
			expectError: true,
			errorMsg:    "job payload 'data' field must be an object",
		},
		{
			name: "invalid options field type",
			payload: map[string]any{
				"type":    "test_job",
				"options": "invalid_options_type", // Should be object
			},
			expectError: true,
			errorMsg:    "job payload 'options' field must be an object",
		},
		{
			name: "valid complete payload",
			payload: map[string]any{
				"type": "test_job",
				"data": map[string]any{
					"key": "value",
				},
				"options": map[string]any{
					"timeout": 300,
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJobPayload(tt.payload)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
