package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// ConditionalAccessPolicyTestResource implements the types.TestResource interface for conditional access policies
type ConditionalAccessPolicyTestResource struct{}

// Exists checks whether the conditional access policy exists in Microsoft Graph
func (r ConditionalAccessPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.Identity().ConditionalAccess().Policies().ByConditionalAccessPolicyId(state.ID).Get(ctx, nil)
		return err
	})
}
