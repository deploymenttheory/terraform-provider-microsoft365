package graphBetaIdentityAndAccessAuthenticationMethodsPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// AuthenticationMethodsPolicyTestResource implements the types.TestResource interface for the authentication methods policy
type AuthenticationMethodsPolicyTestResource struct{}

// Exists checks whether the singleton authentication methods policy is readable in Microsoft Graph
func (r AuthenticationMethodsPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, _ *terraform.InstanceState) error {
		_, err := client.Policies().AuthenticationMethodsPolicy().Get(ctx, nil)
		return err
	})
}
