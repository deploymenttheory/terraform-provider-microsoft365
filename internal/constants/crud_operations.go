package constants

// Terraform provider operation contexts for error handling and state management
const (
	// TfOperationCreate - Creating a new resource
	// 404: Error (resource should exist after POST)
	// State: Never remove on error
	TfOperationCreate = "Create"

	// TfOperationRead - Reading existing resource state
	// 404: Remove from state (resource deleted externally)
	// State: Remove on 404/400
	TfOperationRead = "Read"

	// TfTfOperationUpdate - Updating existing resource
	// 404: Remove from state (resource deleted externally)
	// State: Remove on 404
	TfOperationUpdate = "Update"

	// TfOperationDelete - Deleting a resource
	// 404: Success (idempotent, already deleted)
	// State: Always remove
	TfOperationDelete = "Delete"
)
