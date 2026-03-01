package sentinels

import "errors"

// Common construction-related sentinel errors used across resource construct functions.
// These errors are used when building Graph API request bodies from Terraform plan data.
var (
	// ErrSetRoleScopeTags indicates failure to set role scope tag IDs on the request body
	ErrSetRoleScopeTags = errors.New("failed to set role scope tags")

	// ErrExtractAssignments indicates failure to extract assignments from Terraform state
	ErrExtractAssignments = errors.New("failed to extract assignments")
)
