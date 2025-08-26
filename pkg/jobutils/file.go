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
// This is a utility function that can be used by various export handlers
func SaveExportFile(app *pocketbase.PocketBase, jobId, filename string, fileData []byte, recordCount int) (*core.Record, error) {
	return saveExportFile(app, jobId, "", filename, fileData, recordCount)
}

// SaveExportFileWithUser saves file data to the export_files collection with a specific user ID
// This variant allows specifying the user who requested the export
func SaveExportFileWithUser(app *pocketbase.PocketBase, jobId, userId, filename string, fileData []byte, recordCount int) (*core.Record, error) {
	return saveExportFile(app, jobId, userId, filename, fileData, recordCount)
}

// saveExportFile is the internal implementation for saving export files
func saveExportFile(app *pocketbase.PocketBase, jobId, userId, filename string, fileData []byte, recordCount int) (*core.Record, error) {
	// Find the export_files collection
	collection, err := app.FindCollectionByNameOrId(ExportFilesCollectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to find export_files collection for job %s: %w", jobId, err)
	}

	// Create a new record
	record := core.NewRecord(collection)

	// Get expiration days from environment variable (default: 30 days)
	expirationDays := common.GetEnvInt("EXPORT_FILE_EXPIRATION_DAYS", DefaultFileExpirationDays)
	expirationDate := time.Now().AddDate(0, 0, expirationDays)

	// Set the basic fields
	record.Set("job_id", jobId)
	record.Set("user_id", userId)
	record.Set("record_count", recordCount)
	record.Set("expires_at", expirationDate)

	// Create a filesystem.File from the file data
	file, err := filesystem.NewFileFromBytes(fileData, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create file from data for job %s: %w", jobId, err)
	}

	// Set the file field using PocketBase's file handling
	record.Set("file", file)

	// Save the record (this will automatically handle file upload)
	if err := app.Save(record); err != nil {
		return nil, fmt.Errorf("failed to save export_files record for job %s: %w", jobId, err)
	}

	return record, nil
}
