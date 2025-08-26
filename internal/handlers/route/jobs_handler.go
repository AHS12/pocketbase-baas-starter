package route

import (
	"ims-pocketbase-baas-starter/pkg/response"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func HandleGetJobStatus(e *core.RequestEvent) error {
	jobId := e.Request.PathValue("id")
	if jobId == "" {
		return response.ValidationError(e, "Job ID is required", nil)
	}

	status := getJobStatus(e.App, jobId)

	data := map[string]any{
		"job_id": jobId,
		"status": status,
	}

	return response.OK(e, "Job status", data)
}

func HandleDownloadJobFile(e *core.RequestEvent) error {
	jobId := e.Request.PathValue("id")
	if jobId == "" {
		return response.ValidationError(e, "Job ID is required", nil)
	}

	exportRecord, err := getJobFileRecord(e.App, jobId)
	if err != nil {
		return response.NotFound(e, "Export file not found")
	}

	fileName := exportRecord.GetString("file")
	basePath := exportRecord.BaseFilesPath()

	return response.File(e, fileName, basePath)
}

func getJobFileRecord(app core.App, jobId string) (*core.Record, error) {
	return app.FindFirstRecordByFilter("export_files", "job_id = {:job_id}", dbx.Params{"job_id": jobId})
}

func getJobStatus(app core.App, jobId string) string {
	job, err := app.FindRecordById("queues", jobId)
	if err == nil {
		if job.GetString("reserved_at") == "" {
			return "queued"
		} else {
			return "processing"
		}
	}

	_, err = getJobFileRecord(app, jobId)
	if err == nil {
		return "completed"
	}

	return "failed"
}
