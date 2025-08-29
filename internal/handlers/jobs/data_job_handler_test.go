package jobs

import (
	"ims-pocketbase-baas-starter/pkg/cronutils"
	"ims-pocketbase-baas-starter/pkg/jobutils"
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestNewDataProcessingJobHandler(t *testing.T) {
	app := pocketbase.New()
	handler := NewDataProcessingJobHandler(app)

	if handler == nil {
		t.Fatal("NewDataProcessingJobHandler should not return nil")
	}

	if handler.app != app {
		t.Error("Handler should store app reference")
	}
}

func TestDataProcessingJobHandler_GetJobType(t *testing.T) {
	app := pocketbase.New()
	handler := NewDataProcessingJobHandler(app)

	jobType := handler.GetJobType()
	if jobType != jobutils.JobTypeDataProcessing {
		t.Errorf("Expected job type '%s', got '%s'", jobutils.JobTypeDataProcessing, jobType)
	}
}

func TestDataProcessingJobHandler_validateDataProcessingPayload(t *testing.T) {
	app := pocketbase.New()
	handler := NewDataProcessingJobHandler(app)

	tests := []struct {
		name    string
		payload *jobutils.DataProcessingJobPayload
		wantErr bool
	}{
		{
			name: "valid export payload",
			payload: &jobutils.DataProcessingJobPayload{
				Type: jobutils.JobTypeDataProcessing,
				Data: jobutils.DataProcessingJobData{
					Operation: jobutils.DataProcessingOperationExport,
					Source:    "users",
					Target:    "csv",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid job type",
			payload: &jobutils.DataProcessingJobPayload{
				Type: "invalid_type",
				Data: jobutils.DataProcessingJobData{
					Operation: jobutils.DataProcessingOperationExport,
					Source:    "users",
					Target:    "csv",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid operation",
			payload: &jobutils.DataProcessingJobPayload{
				Type: jobutils.JobTypeDataProcessing,
				Data: jobutils.DataProcessingJobData{
					Operation: "invalid_operation",
					Source:    "users",
					Target:    "csv",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateDataProcessingPayload(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDataProcessingPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDataProcessingJobHandler_Handle_InvalidPayload(t *testing.T) {
	app := pocketbase.New()
	handler := NewDataProcessingJobHandler(app)
	ctx := cronutils.NewCronExecutionContext(app, "test-job")

	// Test with invalid job data
	invalidJob := &jobutils.JobData{
		ID:      "test-job",
		Type:    jobutils.JobTypeDataProcessing,
		Payload: map[string]any{"invalid": "payload"},
	}

	err := handler.Handle(ctx, invalidJob)
	if err == nil {
		t.Error("Handle should return error for invalid payload")
	}
}

func TestDataProcessingJobHandler_handleTransformOperation(t *testing.T) {
	app := pocketbase.New()
	handler := NewDataProcessingJobHandler(app)
	ctx := cronutils.NewCronExecutionContext(app, "test-job")

	payload := &jobutils.DataProcessingJobPayload{
		Data: jobutils.DataProcessingJobData{
			Operation: jobutils.DataProcessingOperationTransform,
			Source:    "source_table",
			Target:    "target_table",
		},
	}

	err := handler.handleTransformOperation(ctx, payload)
	if err != nil {
		t.Errorf("handleTransformOperation should not return error: %v", err)
	}
}

func TestDataProcessingJobHandler_handleAggregateOperation(t *testing.T) {
	app := pocketbase.New()
	handler := NewDataProcessingJobHandler(app)
	ctx := cronutils.NewCronExecutionContext(app, "test-job")

	payload := &jobutils.DataProcessingJobPayload{
		Data: jobutils.DataProcessingJobData{
			Operation: jobutils.DataProcessingOperationAggregate,
			Source:    "source_table",
			Target:    "target_table",
		},
	}

	err := handler.handleAggregateOperation(ctx, payload)
	if err != nil {
		t.Errorf("handleAggregateOperation should not return error: %v", err)
	}
}

func TestDataProcessingJobHandler_handleImportOperation(t *testing.T) {
	app := pocketbase.New()
	handler := NewDataProcessingJobHandler(app)
	ctx := cronutils.NewCronExecutionContext(app, "test-job")

	payload := &jobutils.DataProcessingJobPayload{
		Data: jobutils.DataProcessingJobData{
			Operation: jobutils.DataProcessingOperationImport,
			Source:    "source_file",
			Target:    "target_table",
		},
	}

	err := handler.handleImportOperation(ctx, payload)
	if err != nil {
		t.Errorf("handleImportOperation should not return error: %v", err)
	}
}
