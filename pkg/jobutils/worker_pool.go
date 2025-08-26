package jobutils

import (
	"context"
	"fmt"
	"ims-pocketbase-baas-starter/pkg/cronutils"
	log "ims-pocketbase-baas-starter/pkg/logger"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// WorkerPool manages a pool of persistent workers for job processing
type WorkerPool struct {
	workers     []*Worker
	jobQueue    chan *core.Record
	resultQueue chan WorkerJobResult
	quit        chan bool
	wg          sync.WaitGroup
	maxWorkers  int
	app         *pocketbase.PocketBase
	registry    *JobRegistry
	isShutdown  bool
	mu          sync.RWMutex
}

// Worker represents a single worker in the pool
type Worker struct {
	id          int
	jobQueue    chan *core.Record
	resultQueue chan WorkerJobResult
	quit        chan bool
	app         *pocketbase.PocketBase
	registry    *JobRegistry
}

// WorkerJobResult represents the result of job processing
type WorkerJobResult struct {
	JobID string
	Error error
}

// NewWorkerPool creates a new persistent worker pool
func NewWorkerPool(app *pocketbase.PocketBase, registry *JobRegistry, maxWorkers int) *WorkerPool {
	if maxWorkers <= 0 {
		maxWorkers = 5
	}

	jobQueueSize := maxWorkers * 10
	resultQueueSize := maxWorkers * 10

	pool := &WorkerPool{
		workers:     make([]*Worker, 0, maxWorkers),
		jobQueue:    make(chan *core.Record, jobQueueSize),
		resultQueue: make(chan WorkerJobResult, resultQueueSize),
		quit:        make(chan bool),
		maxWorkers:  maxWorkers,
		app:         app,
		registry:    registry,
		isShutdown:  false,
	}

	for i := 0; i < maxWorkers; i++ {
		worker := &Worker{
			id:          i,
			jobQueue:    pool.jobQueue,
			resultQueue: pool.resultQueue,
			quit:        make(chan bool),
			app:         app,
			registry:    registry,
		}
		pool.workers = append(pool.workers, worker)
		pool.wg.Add(1)
		go worker.start(&pool.wg)
	}

	log.Info("Worker pool started", "workers", maxWorkers, "job_queue_size", jobQueueSize)
	return pool
}

// IsShutdown returns whether the worker pool has been shut down
func (wp *WorkerPool) IsShutdown() bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.isShutdown
}

// ProcessJobs processes a batch of jobs using the worker pool
func (wp *WorkerPool) ProcessJobs(jobs []*core.Record) []error {
	// Check if pool is shutdown
	if wp.IsShutdown() {
		err := fmt.Errorf("worker pool is shutdown")
		results := make([]error, len(jobs))
		for i := range results {
			results[i] = err
		}
		return results
	}

	if len(jobs) == 0 {
		return nil
	}

	jobIndexMap := make(map[string]int, len(jobs))
	for i, job := range jobs {
		jobIndexMap[job.Id] = i
	}

	// Send jobs to workers
	sendErrors := make([]error, len(jobs))
	jobsSent := 0
	for i, job := range jobs {
		select {
		case wp.jobQueue <- job:
			jobsSent++
		case <-time.After(30 * time.Second):
			err := fmt.Errorf("job queue timeout for job %s", job.Id)
			log.Error("Job queue timeout", "job_id", job.Id)
			sendErrors[i] = err
		}
	}

	// If we couldn't send any jobs, return the send errors
	if jobsSent == 0 {
		log.Warn("No jobs were sent to worker pool", "total_jobs", len(jobs))
		return sendErrors
	}

	// Collect results
	results := make([]error, len(jobs))
	copy(results, sendErrors)

	for i := 0; i < jobsSent; i++ {
		select {
		case result := <-wp.resultQueue:
			if jobIndex, exists := jobIndexMap[result.JobID]; exists {
				results[jobIndex] = result.Error
			} else {
				log.Warn("Received result for unknown job", "job_id", result.JobID)
			}
		case <-time.After(5 * time.Minute):
			err := fmt.Errorf("job processing timeout")
			log.Error("Job processing timeout")
			for j := range results {
				if results[j] == nil {
					results[j] = err
					break
				}
			}
		}
	}

	// Log processing summary
	successCount := 0
	failureCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		} else {
			failureCount++
		}
	}

	log.Info("Worker pool job processing completed",
		"total_jobs", len(jobs),
		"jobs_sent", jobsSent,
		"successful", successCount,
		"failed", failureCount)

	return results
}

// ProcessJobsConcurrently processes jobs concurrently using the worker pool
func (wp *WorkerPool) ProcessJobsConcurrently(jobs []*core.Record, maxWorkers int) []error {
	if len(jobs) == 0 {
		return nil
	}

	log.Debug("ProcessJobsConcurrently called with maxWorkers parameter (ignored)",
		"requested_workers", maxWorkers,
		"configured_workers", wp.maxWorkers)

	return wp.ProcessJobs(jobs)
}

// Shutdown gracefully shuts down the worker pool
func (wp *WorkerPool) Shutdown(ctx context.Context) error {
	wp.mu.Lock()
	if wp.isShutdown {
		wp.mu.Unlock()
		log.Info("Worker pool already shutdown")
		return nil
	}
	wp.isShutdown = true
	wp.mu.Unlock()

	log.Info("Shutting down worker pool")
	close(wp.jobQueue)

	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info("Worker pool shutdown completed")
		return nil
	case <-ctx.Done():
		close(wp.quit)
		log.Warn("Worker pool force shutdown due to timeout")
		return ctx.Err()
	}
}

func (w *Worker) start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-w.jobQueue:
			if !ok {
				return
			}

			err := w.processJob(job)

			select {
			case w.resultQueue <- WorkerJobResult{JobID: job.Id, Error: err}:
			case <-time.After(10 * time.Second):
				log.Error("Timeout sending result to result queue", "job_id", job.Id, "worker_id", w.id)
			}

		case <-w.quit:
			return
		}
	}
}

func (w *Worker) processJob(record *core.Record) error {
	if record == nil || record.Id == "" || record.Collection().Name != QueuesCollection {
		return fmt.Errorf("invalid job record")
	}

	originalReservedAt := record.GetString("reserved_at")
	if originalReservedAt != "" && w.isJobReserved(record) {
		return fmt.Errorf("job %s is already reserved", record.Id)
	}

	now := time.Now()
	record.Set("reserved_at", now.Format(time.RFC3339))

	if err := w.app.Save(record); err != nil {
		record.Set("reserved_at", originalReservedAt)
		return fmt.Errorf("failed to reserve job %s: %w", record.Id, err)
	}

	jobData, err := ParseJobDataFromRecord(record)
	if err != nil {
		failErr := w.failJob(record, fmt.Errorf("failed to parse job data: %w", err))
		if failErr != nil {
			log.Error("Failed to mark job as failed", "job_id", record.Id, "worker_id", w.id, "error", failErr)
		}
		return err
	}

	handler, err := w.registry.GetHandler(jobData.Type)
	if err != nil {
		failErr := w.failJob(record, fmt.Errorf("no handler for job type '%s': %w", jobData.Type, err))
		if failErr != nil {
			log.Error("Failed to mark job as failed", "job_id", record.Id, "worker_id", w.id, "error", failErr)
		}
		return err
	}

	// Execute job with panic recovery
	var jobErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				jobErr = fmt.Errorf("job handler panicked: %v", r)
				log.Error("Job handler panic", "job_id", record.Id, "worker_id", w.id, "panic", r)
			}
		}()

		ctx := cronutils.NewCronExecutionContext(w.app, record.Id)
		ctx.LogStart(fmt.Sprintf("Processing %s job: %s", jobData.Type, jobData.Name))
		jobErr = handler.Handle(ctx, jobData)

		if jobErr == nil {
			ctx.LogEnd("Job processed successfully")
		}
	}()

	if jobErr != nil {
		log.Error("Job failed", "job_id", record.Id, "worker_id", w.id, "job_type", jobData.Type, "error", jobErr)
		failErr := w.failJob(record, jobErr)
		if failErr != nil {
			log.Error("Failed to mark job as failed", "job_id", record.Id, "worker_id", w.id, "error", failErr)
		}
		return jobErr
	}

	if err := w.app.Delete(record); err != nil {
		log.Error("Failed to complete job", "job_id", record.Id, "worker_id", w.id, "error", err)
		return err
	}

	log.Info("Job completed successfully", "job_id", record.Id, "worker_id", w.id, "job_type", jobData.Type)
	return nil
}

func (w *Worker) isJobReserved(record *core.Record) bool {
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

func (w *Worker) failJob(record *core.Record, jobErr error) error {
	currentAttempts := int(record.GetFloat("attempts"))
	record.Set("attempts", currentAttempts+1)
	record.Set("reserved_at", "")

	if err := w.app.Save(record); err != nil {
		log.Error("Failed to update failed job", "job_id", record.Id, "error", err)
		return fmt.Errorf("failed to update failed job: %w", err)
	}

	return jobErr
}
