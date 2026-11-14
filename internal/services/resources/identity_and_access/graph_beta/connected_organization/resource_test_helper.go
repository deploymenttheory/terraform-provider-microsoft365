package graphBetaIdentityAndAccessConnectedOrganization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// ConnectedOrganizationTestResource implements the types.TestResource interface for connected organizations
type ConnectedOrganizationTestResource struct{}

// Exists checks whether the connected organization exists in Microsoft Graph
func (r ConnectedOrganizationTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.IdentityGovernance().EntitlementManagement().ConnectedOrganizations().ByConnectedOrganizationId(state.ID).Get(ctx, nil)
		return err
	})
}
