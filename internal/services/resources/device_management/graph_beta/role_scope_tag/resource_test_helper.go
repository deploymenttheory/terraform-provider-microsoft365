package graphBetaRoleScopeTag

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// RoleScopeTagTestResource implements the types.TestResource interface for role scope tags
type RoleScopeTagTestResource struct{}

// Exists checks whether the role scope tag exists in Microsoft Graph
func (r RoleScopeTagTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	roleScopeTagID := state.ID

	result, err := graphClient.
		DeviceManagement().
		RoleScopeTags().
		ByRoleScopeTagId(roleScopeTagID).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		// 404 means it doesn't exist
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	exists := result != nil
	return &exists, nil
}
