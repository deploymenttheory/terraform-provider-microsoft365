package graphBetaGroup

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// GroupTestResource implements the types.TestResource interface for Entra ID groups
type GroupTestResource struct{}

// Exists checks whether the group exists in Microsoft Graph
// Handles both hard-delete (404) and soft-delete (deletedDateTime set) scenarios
func (r GroupTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	groupID := state.ID

	group, err := graphClient.
		Groups().
		ByGroupId(groupID).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		// 404 means it's hard-deleted (doesn't exist)
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	if group == nil {
		exists := false
		return &exists, nil
	}

	// Check if soft-deleted (for test cleanup purposes, soft-deleted is considered "destroyed")
	deletedDateTime := group.GetDeletedDateTime()
	if deletedDateTime != nil {
		exists := false // Treat soft-deleted as destroyed for test purposes
		return &exists, nil
	}

	exists := true
	return &exists, nil
}
