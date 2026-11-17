package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// B2bManagementPolicyTestResource implements the types.TestResource interface for B2B management policies
type B2bManagementPolicyTestResource struct{}

// Exists checks whether the B2B management policy exists in Microsoft Graph
func (r B2bManagementPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, b2bManagementPolicyExists)
}

func b2bManagementPolicyExists(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
	reader := B2bManagementPolicyResource{client: client}
	_, err := reader.getPolicy(ctx, state.ID, true)
	return err
}
