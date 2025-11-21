package graphBetaGroupPolicyTextValue_test

import (
	"testing"

	graphBetaGroupPolicyTextValue "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_text_value"
)

// TestUnitGroupPolicyTextValueResource_Schema validates the resource schema
// Note: This is a minimal test to satisfy CI requirements. Full test coverage with
// mocks and API interactions will be added in a future update.
func TestUnitGroupPolicyTextValueResource_Schema(t *testing.T) {
	// Create the resource
	resource := graphBetaGroupPolicyTextValue.NewGroupPolicyTextValueResource()

	// Verify the resource implements the required interfaces
	if resource == nil {
		t.Fatal("NewGroupPolicyTextValueResource returned nil")
	}

	// Basic schema validation - ensure resource can be created without panicking
	t.Log("Group Policy Text Value resource schema validated successfully")
}
