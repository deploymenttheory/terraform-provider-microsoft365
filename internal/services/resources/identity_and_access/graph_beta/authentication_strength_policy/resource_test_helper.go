package graphBetaAuthenticationStrengthPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// AuthenticationStrengthTestResource implements the types.TestResource interface for authentication strength policies
type AuthenticationStrengthTestResource struct{}

// Exists checks whether the authentication strength policy exists in Microsoft Graph
func (r AuthenticationStrengthTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.Identity().ConditionalAccess().AuthenticationStrength().Policies().ByAuthenticationStrengthPolicyId(state.ID).Get(ctx, nil)
		return err
	})
}
