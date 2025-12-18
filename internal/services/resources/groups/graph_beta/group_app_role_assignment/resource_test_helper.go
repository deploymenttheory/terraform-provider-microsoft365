package graphBetaGroupAppRoleAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// GroupAppRoleAssignmentTestResource implements the types.TestResource interface for group app role assignments
type GroupAppRoleAssignmentTestResource struct{}

// Exists checks whether the group app role assignment exists in Microsoft Graph
func (r GroupAppRoleAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	assignmentID := state.ID
	groupID := state.Attributes["target_group_id"]

	assignment, err := graphClient.
		Groups().
		ByGroupId(groupID).
		AppRoleAssignments().
		ByAppRoleAssignmentId(assignmentID).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		// 404 means it doesn't exist
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	if assignment == nil {
		exists := false
		return &exists, nil
	}

	exists := true
	return &exists, nil
}
