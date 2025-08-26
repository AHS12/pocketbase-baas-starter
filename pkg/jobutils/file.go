package jobutils

import (
	"fmt"
	"time"

	"ims-pocketbase-baas-starter/pkg/common"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Constants for export files
const (
	ExportFilesCollectionName = "export_files"
	DefaultFileExpirationDays = 30
)

// SaveExportFile saves file data to the export_files collection
func SaveExportFile(app *pocketbase.PocketBase, jobId, filename string, fileData []byte, recordCount int) (*core.Record, error) {
	return saveExportFile(app, jobId, "", filename, fileData, recordCount)
}

// SaveExportFileWithUser saves file data to the export_files collection with a specific user ID
func SaveExportFileWithUser(app *pocketbase.PocketBase, jobId, userId, filename string, fileData []byte, recordCount int) (*core.Record, error) {
	return saveExportFile(app, jobId, userId, filename, fileData, recordCount)
}

// saveExportFile is the internal implementation for saving export files
func saveExportFile(app *pocketbase.PocketBase, jobId, userId, filename string, fileData []byte, recordCount int) (*core.Record, error) {
	collection, err := app.FindCollectionByNameOrId(ExportFilesCollectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to find export_files collection for job %s: %w", jobId, err)
	}

	record := core.NewRecord(collection)

	expirationDays := common.GetEnvInt("EXPORT_FILE_EXPIRATION_DAYS", DefaultFileExpirationDays)
	expirationDate := time.Now().AddDate(0, 0, expirationDays)

	record.Set("job_id", jobId)
	record.Set("user_id", userId)
	record.Set("record_count", recordCount)
	record.Set("expires_at", expirationDate)

	file, err := filesystem.NewFileFromBytes(fileData, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create file from data for job %s: %w", jobId, err)
	}

	record.Set("file", file)

	if err := app.Save(record); err != nil {
		return nil, fmt.Errorf("failed to save export_files record for job %s: %w", jobId, err)
	}

	return record, nil
}
