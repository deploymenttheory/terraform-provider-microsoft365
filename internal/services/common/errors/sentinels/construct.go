package sentinels

import "errors"

// Common construction-related sentinel errors used across resource construct functions.
// These errors are used when building Graph API request bodies from Terraform plan data.
var (
	// ErrSetRoleScopeTags indicates failure to set role scope tag IDs on the request body
	ErrSetRoleScopeTags = errors.New("failed to set role scope tags")

	// ErrExtractAssignments indicates failure to extract assignments from Terraform state
	ErrExtractAssignments = errors.New("failed to extract assignments")

	// ErrInvalidCatalogEntryType indicates an invalid or unsupported catalog entry type
	ErrInvalidCatalogEntryType = errors.New("invalid catalog_entry_type")

	// ErrExtractMonitoringRules indicates failure to extract monitoring rules from Terraform state
	ErrExtractMonitoringRules = errors.New("failed to extract monitoring_rules")

	// ErrCreateMonitoringRuleObject indicates failure to create monitoring rule object
	ErrCreateMonitoringRuleObject = errors.New("failed to create monitoring rule object")

	// ErrCreateMonitoringRulesSet indicates failure to create monitoring rules set
	ErrCreateMonitoringRulesSet = errors.New("failed to create monitoring rules set")

	// ErrExtractApprovalRules indicates failure to extract approval rules from Terraform state
	ErrExtractApprovalRules = errors.New("failed to extract approval_rules")

	// ErrSetClassification indicates failure to set classification field
	ErrSetClassification = errors.New("failed to set classification")

	// ErrSetCadence indicates failure to set cadence field
	ErrSetCadence = errors.New("failed to set cadence")

	// ErrCreateApprovalRuleObject indicates failure to create approval rule object
	ErrCreateApprovalRuleObject = errors.New("failed to create approval rule object")

	// ErrCreateApprovalRulesSet indicates failure to create approval rules set
	ErrCreateApprovalRulesSet = errors.New("failed to create approval rules set")
)
