package route

import (
	"encoding/json"
	"ims-pocketbase-baas-starter/pkg/jobutils"
	"ims-pocketbase-baas-starter/pkg/response"

	"github.com/pocketbase/pocketbase/core"
)

func HandleUserExport(e *core.RequestEvent) error {
	payload := jobutils.DataProcessingJobPayload{
		Type: jobutils.JobTypeDataProcessing,
		Data: jobutils.DataProcessingJobData{
			Operation: jobutils.DataProcessingOperationExport,
			Source:    jobutils.DataProcessingCollectionUsers,
			Target:    jobutils.DataProcessingFileCSV,
		},
		Options: jobutils.DataProcessingJobOptions{
			Timeout: 900, // 15 minutes
		},
	}

	queuesCollection, err := e.App.FindCollectionByNameOrId("queues")
	if err != nil {
		return response.InternalServerError(e, "Queue system unavailable", nil)
	}

	job := core.NewRecord(queuesCollection)
	job.Set("name", "User Export")
	job.Set("description", "Export users to CSV")

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return response.InternalServerError(e, "Failed to create job", nil)
	}
	job.Set("payload", string(payloadJSON))

	if err := e.App.Save(job); err != nil {
		return response.InternalServerError(e, "Failed to queue export job", nil)
	}

	data := map[string]any{
		"job_id": job.Id,
		"status": "queued",
	}
	return response.OK(e, "User export job queued successfully", data)
}
