package graphBetaGroupLifecycleExpirationPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// GroupLifecycleExpirationPolicyTestResource implements the types.TestResource interface for group lifecycle expiration policies
type GroupLifecycleExpirationPolicyTestResource struct{}

// Exists checks whether the group lifecycle expiration policy exists in Microsoft Graph
func (r GroupLifecycleExpirationPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.
			GroupLifecyclePolicies().
			ByGroupLifecyclePolicyId(state.ID).
			Get(ctx, nil)
		return err
	})
}
