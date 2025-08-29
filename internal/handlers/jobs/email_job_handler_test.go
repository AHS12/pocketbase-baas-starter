package jobs

import (
	"ims-pocketbase-baas-starter/pkg/cronutils"
	"ims-pocketbase-baas-starter/pkg/jobutils"
	"ims-pocketbase-baas-starter/pkg/metrics"
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestNewEmailJobHandler(t *testing.T) {
	app := pocketbase.New()
	handler := NewEmailJobHandler(app)

	if handler == nil {
		t.Fatal("NewEmailJobHandler should not return nil")
	}

	if handler.app != app {
		t.Error("Handler should store app reference")
	}
}

func TestEmailJobHandler_GetJobType(t *testing.T) {
	app := pocketbase.New()
	handler := NewEmailJobHandler(app)

	jobType := handler.GetJobType()
	if jobType != jobutils.JobTypeEmail {
		t.Errorf("Expected job type '%s', got '%s'", jobutils.JobTypeEmail, jobType)
	}
}

func TestEmailJobHandler_validateEmailPayload(t *testing.T) {
	app := pocketbase.New()
	handler := NewEmailJobHandler(app)

	validPayload := &jobutils.EmailJobPayload{
		Type: jobutils.JobTypeEmail,
		Data: jobutils.EmailJobData{
			To:      "test@example.com",
			Subject: "Test Subject",
		},
	}

	err := handler.validateEmailPayload(validPayload)
	if err != nil {
		t.Errorf("validateEmailPayload() error = %v, expected nil", err)
	}

	invalidPayload := &jobutils.EmailJobPayload{
		Type: "invalid_type",
		Data: jobutils.EmailJobData{
			To:      "test@example.com",
			Subject: "Test Subject",
		},
	}

	err = handler.validateEmailPayload(invalidPayload)
	if err == nil {
		t.Error("validateEmailPayload() expected error for invalid type")
	}
}

func TestEmailJobHandler_Handle_InvalidPayload(t *testing.T) {
	metrics.InitializeProvider(metrics.Config{
		Provider: metrics.ProviderDisabled,
		Enabled:  false,
	})
	defer metrics.Reset()

	app := pocketbase.New()
	handler := NewEmailJobHandler(app)
	ctx := cronutils.NewCronExecutionContext(app, "test-job")

	invalidJob := &jobutils.JobData{
		ID:      "test-job",
		Type:    jobutils.JobTypeEmail,
		Payload: map[string]any{"invalid": "payload"},
	}

	err := handler.Handle(ctx, invalidJob)
	if err == nil {
		t.Error("Handle should return error for invalid payload")
	}
}

func TestEmailJobHandler_processSingleTemplate_FileNotFound(t *testing.T) {
	app := pocketbase.New()
	handler := NewEmailJobHandler(app)

	payload := &jobutils.EmailJobPayload{
		Data: jobutils.EmailJobData{
			Template: "nonexistent_template",
		},
	}

	content, err := handler.processSingleTemplate(payload, ".html")
	if err == nil {
		t.Error("processSingleTemplate should return error for nonexistent template")
	}

	if content != "" {
		t.Error("processSingleTemplate should return empty content on error")
	}
}
