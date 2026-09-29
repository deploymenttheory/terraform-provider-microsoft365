package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// B2bManagementPolicyTestResource implements the types.TestResource interface for B2B management policies
type B2bManagementPolicyTestResource struct{}

// Exists checks whether the B2B management policy exists in Microsoft Graph
func (r B2bManagementPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.Policies().B2bManagementPolicies().ByB2bManagementPolicyId(state.ID).Get(ctx, nil)
		return err
	})
}
