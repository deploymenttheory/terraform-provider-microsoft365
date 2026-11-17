package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// SecurityDefaultsEnforcementPolicyTestResource implements the types.TestResource interface for the security defaults enforcement policy
type SecurityDefaultsEnforcementPolicyTestResource struct{}

// Exists checks whether the singleton security defaults enforcement policy is readable in Microsoft Graph
func (r SecurityDefaultsEnforcementPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, _ *terraform.InstanceState) error {
		_, err := client.Policies().IdentitySecurityDefaultsEnforcementPolicy().Get(ctx, nil)
		return err
	})
}
