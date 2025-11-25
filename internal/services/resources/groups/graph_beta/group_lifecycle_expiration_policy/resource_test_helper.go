package graphBetaGroupLifecycleExpirationPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// GroupLifecycleExpirationPolicyTestResource implements the types.TestResource interface for group lifecycle expiration policies
type GroupLifecycleExpirationPolicyTestResource struct{}

// Exists checks whether the group lifecycle expiration policy exists in Microsoft Graph
func (r GroupLifecycleExpirationPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	policyID := state.ID

	policy, err := graphClient.
		GroupLifecyclePolicies().
		ByGroupLifecyclePolicyId(policyID).
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

	if policy == nil {
		exists := false
		return &exists, nil
	}

	exists := true
	return &exists, nil
}
