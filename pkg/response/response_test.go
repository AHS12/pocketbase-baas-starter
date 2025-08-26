package response

import (
	"testing"
)

// Test that all package-level functions are available
func TestPackageLevelFunctions(t *testing.T) {
	// Just testing that functions exist and are accessible
	// Actual functionality would require mocking core.RequestEvent which is complex

	// These tests just ensure the functions are defined
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
