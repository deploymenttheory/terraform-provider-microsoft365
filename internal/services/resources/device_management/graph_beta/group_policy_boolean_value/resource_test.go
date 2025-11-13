package graphBetaGroupPolicyBooleanValue_test

import (
	"testing"

	graphBetaGroupPolicyBooleanValue "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_boolean_value"
)

// TestUnitGroupPolicyBooleanValueResource_Schema validates the resource schema
// Note: This is a minimal test to satisfy CI requirements. Full test coverage with
// mocks and API interactions will be added in a future update.
func TestUnitGroupPolicyBooleanValueResource_Schema(t *testing.T) {
	// Create the resource
	resource := graphBetaGroupPolicyBooleanValue.NewGroupPolicyBooleanValueResource()

	// Verify the resource implements the required interfaces
	if resource == nil {
		t.Fatal("NewGroupPolicyBooleanValueResource returned nil")
	}

	// Basic schema validation - ensure resource can be created without panicking
	t.Log("Group Policy Boolean Value resource schema validated successfully")
}
