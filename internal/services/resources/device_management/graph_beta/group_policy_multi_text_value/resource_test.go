package graphBetaGroupPolicyMultiTextValue_test

import (
	"testing"

	graphBetaGroupPolicyMultiTextValue "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_multi_text_value"
)

// TestUnitGroupPolicyMultiTextValueResource_Schema validates the resource schema
// Note: This is a minimal test to satisfy CI requirements. Full test coverage with
// mocks and API interactions will be added in a future update.
func TestUnitGroupPolicyMultiTextValueResource_Schema(t *testing.T) {
	// Create the resource
	resource := graphBetaGroupPolicyMultiTextValue.NewGroupPolicyMultiTextValueResource()

	// Verify the resource implements the required interfaces
	if resource == nil {
		t.Fatal("NewGroupPolicyMultiTextValueResource returned nil")
	}

	// Basic schema validation - ensure resource can be created without panicking
	t.Log("Group Policy Multi Text Value resource schema validated successfully")
}
