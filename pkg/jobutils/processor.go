package jobutils

import (
	"encoding/json"
	"fmt"
	"ims-pocketbase-baas-starter/pkg/common"
	"ims-pocketbase-baas-starter/pkg/cronutils"
	log "ims-pocketbase-baas-starter/pkg/logger"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const QueuesCollection = "queues"

// NewJobRegistry creates a new job registry
func NewJobRegistry() *JobRegistry {
	return &JobRegistry{
		handlers: make(map[string]JobHandler),
	}
}

// Register adds a job handler to the registry
func (r *JobRegistry) Register(handler JobHandler) error {
	if handler == nil {
		return fmt.Errorf("job handler cannot be nil")
	}

	jobType := handler.GetJobType()
	if jobType == "" {
		return fmt.Errorf("job handler must return a non-empty job type")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[jobType]; exists {
		return fmt.Errorf("job handler for type '%s' is already registered", jobType)
	}

	r.handlers[jobType] = handler
	return nil
}

// GetHandler retrieves a job handler by job type
func (r *JobRegistry) GetHandler(jobType string) (JobHandler, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.handlers[jobType]
	if !exists {
		return nil, fmt.Errorf("no handler registered for job type '%s'", jobType)
	}

	return handler, nil
}

// ListHandlers returns a list of all registered job types
func (r *JobRegistry) ListHandlers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.handlers))
	for jobType := range r.handlers {
		types = append(types, jobType)
	}
	return types
}

// ParseJobDataFromRecord extracts JobData from a PocketBase record
func ParseJobDataFromRecord(record *core.Record) (*JobData, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	// Parse the JSON payload
	var payload map[string]any
	payloadStr := record.GetString("payload")
	if payloadStr != "" {
		if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
			return nil, fmt.Errorf("failed to parse job payload: %w", err)
		}
	}

	// Extract job type from payload
	jobType := ""
	if payload != nil {
		if typeVal, ok := payload["type"]; ok {
			if typeStr, ok := typeVal.(string); ok {
				jobType = typeStr
			}
		}
	}

	if jobType == "" {
		return nil, fmt.Errorf("job payload must contain a 'type' field")
	}

	// Parse reserved_at timestamp
	var reservedAt *time.Time
	if reservedAtStr := record.GetString("reserved_at"); reservedAtStr != "" {
		if parsed, err := time.Parse(time.RFC3339, reservedAtStr); err == nil {
			reservedAt = &parsed
		}
	}

	return &JobData{
		ID:          record.Id,
		Name:        record.GetString("name"),
		Description: record.GetString("description"),
		Type:        jobType,
		Payload:     payload,
		Attempts:    int(record.GetFloat("attempts")),
		ReservedAt:  reservedAt,
		CreatedAt:   record.GetDateTime("created").Time(),
		UpdatedAt:   record.GetDateTime("updated").Time(),
	}, nil
}

// ValidateJobPayload validates that a job payload has the required structure
func ValidateJobPayload(payload map[string]any) error {
	if payload == nil {
		return fmt.Errorf("job payload cannot be nil")
	}

	typeVal, exists := payload["type"]
	if !exists {
		return fmt.Errorf("job payload must contain a 'type' field")
	}

	typeStr, ok := typeVal.(string)
	if !ok || typeStr == "" {
		return fmt.Errorf("job payload 'type' field must be a non-empty string")
	}

	if dataVal, exists := payload["data"]; exists {
		if _, ok := dataVal.(map[string]any); !ok {
			return fmt.Errorf("job payload 'data' field must be an object")
		}
	}

	if optionsVal, exists := payload["options"]; exists {
		if _, ok := optionsVal.(map[string]any); !ok {
			return fmt.Errorf("job payload 'options' field must be an object")
		}
	}

	return nil
}

// NewJobProcessor creates a new job processor with initialized registry and worker pool
func NewJobProcessor(app *pocketbase.PocketBase) *JobProcessor {
	if app == nil {
		panic("NewJobProcessor: app cannot be nil")
	}

	registry := NewJobRegistry()

	return &JobProcessor{
		app:        app,
		registry:   registry,
		workerPool: NewWorkerPool(app, registry, common.GetEnvInt("JOB_MAX_WORKERS", 5)),
	}
}

// GetRegistry returns the job registry for handler registration
func (p *JobProcessor) GetRegistry() *JobRegistry {
	return p.registry
}

// RegisterHandler is a convenience method to register a job handler
func (p *JobProcessor) RegisterHandler(handler JobHandler) error {
	return p.registry.Register(handler)
}

// ProcessJob processes a single job with complete lifecycle management
func (p *JobProcessor) ProcessJob(record *core.Record) error {
	if record == nil || record.Id == "" || record.Collection().Name != QueuesCollection {
		return fmt.Errorf("invalid job record")
	}

	originalReservedAt := record.GetString("reserved_at")
	if originalReservedAt != "" && p.isJobReserved(record) {
		return fmt.Errorf("job %s is already reserved", record.Id)
	}

	now := time.Now()
	record.Set("reserved_at", now.Format(time.RFC3339))

	if err := p.app.Save(record); err != nil {
		record.Set("reserved_at", originalReservedAt)
		return fmt.Errorf("failed to reserve job %s: %w", record.Id, err)
	}

	jobData, err := ParseJobDataFromRecord(record)
	if err != nil {
		failErr := p.failJob(record, fmt.Errorf("failed to parse job data: %w", err))
		if failErr != nil {
			log.Error("Failed to mark job as failed", "job_id", record.Id, "error", failErr)
		}
		return err
	}

	if err := ValidateJobPayload(jobData.Payload); err != nil {
		failErr := p.failJob(record, fmt.Errorf("invalid job payload: %w", err))
		if failErr != nil {
			log.Error("Failed to mark job as failed", "job_id", record.Id, "error", failErr)
		}
		return err
	}

	handler, err := p.registry.GetHandler(jobData.Type)
	if err != nil {
		failErr := p.failJob(record, fmt.Errorf("no handler found for job type '%s': %w", jobData.Type, err))
		if failErr != nil {
			log.Error("Failed to mark job as failed", "job_id", record.Id, "error", failErr)
		}
		return err
	}

	ctx := cronutils.NewCronExecutionContext(p.app, record.Id)
	var jobErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				jobErr = fmt.Errorf("job handler panicked: %v", r)
				ctx.LogError(jobErr, "Job handler panic recovered")
			}
		}()

		ctx.LogStart(fmt.Sprintf("Processing %s job: %s", jobData.Type, jobData.Name))
		jobErr = handler.Handle(ctx, jobData)
	}()

	if jobErr != nil {
		ctx.LogError(jobErr, "Job processing failed")
		failErr := p.failJob(record, jobErr)
		if failErr != nil {
			ctx.LogError(failErr, "Failed to mark job as failed")
		}
		return jobErr
	}

	ctx.LogEnd("Job processed successfully")

	if err := p.app.Delete(record); err != nil {
		return fmt.Errorf("failed to delete completed job %s: %w", record.Id, err)
	}

	log.Info("Job completed and removed from queue", "job_id", record.Id, "job_name", record.GetString("name"))
	return nil
}

// ProcessJobsConcurrently processes multiple jobs concurrently using the persistent worker pool
func (p *JobProcessor) ProcessJobsConcurrently(records []*core.Record, maxWorkers int) []error {
	if len(records) == 0 {
		return nil
	}

	// Use the persistent worker pool for better performance
	if p.workerPool != nil {
		return p.workerPool.ProcessJobs(records)
	}

	// Fallback to sequential processing
	return p.ProcessJobs(records)
}

// ProcessJobs processes multiple jobs sequentially
func (p *JobProcessor) ProcessJobs(records []*core.Record) []error {
	errors := make([]error, len(records))
	for i, record := range records {
		errors[i] = p.ProcessJob(record)
	}
	return errors
}

func (p *JobProcessor) isJobReserved(record *core.Record) bool {
	reservedAtStr := record.GetString("reserved_at")
	if reservedAtStr == "" {
		return false
	}

	reservedAt, err := time.Parse(time.RFC3339, reservedAtStr)
	if err != nil {
		return false
	}

	return time.Since(reservedAt) < 5*time.Minute
}

func (p *JobProcessor) failJob(record *core.Record, jobErr error) error {
	currentAttempts := int(record.GetFloat("attempts"))
	record.Set("attempts", currentAttempts+1)
	record.Set("reserved_at", "")

	if err := p.app.Save(record); err != nil {
		log.Error("Failed to update failed job record", "job_id", record.Id, "error", err)
		return fmt.Errorf("failed to update failed job %s: %w", record.Id, err)
	}

	log.Error("Job failed", "job_id", record.Id, "job_name", record.GetString("name"), "attempts", currentAttempts+1, "error", jobErr)
	return jobErr
}
