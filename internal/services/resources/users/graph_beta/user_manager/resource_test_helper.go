package graphBetaUsersUserManager

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// UserManagerTestResource implements the types.TestResource interface for user manager relationships
type UserManagerTestResource struct{}

// Exists checks whether the user manager relationship exists in Microsoft Graph
func (r UserManagerTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	// The ID is the user_id for this relationship resource
	userId := state.ID

	_, err = graphClient.
		Users().
		ByUserId(userId).
		Manager().
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		// 404 means no manager is assigned
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	exists := true
	return &exists, nil
}
