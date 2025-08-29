package response

import (
	"testing"
)

func TestPackageLevelFunctions(t *testing.T) {
	_ = OK
	_ = Created
	_ = BadRequest
	_ = Unauthorized
	_ = Forbidden
	_ = NotFound
	_ = InternalServerError
	_ = ValidationError
	_ = Success
	_ = Error
	_ = File
}

func TestPocketBaseResponseStructure(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		message string
		data    map[string]any
	}{
		{
			name:    "success response",
			status:  200,
			message: "Operation successful",
			data:    map[string]any{"result": "success"},
		},
		{
			name:    "error response",
			status:  400,
			message: "Validation failed",
			data:    map[string]any{"errors": []string{"field required"}},
		},
		{
			name:    "empty data",
			status:  204,
			message: "No content",
			data:    map[string]any{},
		},
		{
			name:    "nil data",
			status:  200,
			message: "Success",
			data:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := PocketBaseResponse{
				Status:  tt.status,
				Message: tt.message,
				Data:    tt.data,
			}

			if response.Status != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, response.Status)
			}

			if response.Message != tt.message {
				t.Errorf("Expected message %q, got %q", tt.message, response.Message)
			}

			if tt.data == nil && response.Data != nil {
				t.Errorf("Expected nil data, got %v", response.Data)
			}

			if tt.data != nil {
				if response.Data == nil {
					t.Error("Expected non-nil data, got nil")
				} else {
					for key, expectedValue := range tt.data {
						if actualValue, exists := response.Data[key]; !exists {
							t.Errorf("Expected key %q to exist in data", key)
						} else {
							switch v := expectedValue.(type) {
							case []string:
								if actualSlice, ok := actualValue.([]string); ok {
									if len(actualSlice) != len(v) {
										t.Errorf("Expected data[%q] slice length %d, got %d", key, len(v), len(actualSlice))
									} else {
										for i, item := range v {
											if i < len(actualSlice) && actualSlice[i] != item {
												t.Errorf("Expected data[%q][%d] = %v, got %v", key, i, item, actualSlice[i])
											}
										}
									}
								} else {
									t.Errorf("Expected data[%q] to be []string, got %T", key, actualValue)
								}
							default:
								if actualValue != expectedValue {
									t.Errorf("Expected data[%q] = %v, got %v", key, expectedValue, actualValue)
								}
							}
						}
					}
				}
			}
		})
	}
}
