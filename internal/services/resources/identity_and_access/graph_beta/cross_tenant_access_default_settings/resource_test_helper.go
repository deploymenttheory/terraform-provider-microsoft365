package graphBetaCrossTenantAccessDefaultSettings

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// CrossTenantAccessDefaultSettingsTestResource implements the types.TestResource interface.
// Because this is a singleton resource, Exists simply verifies the GET endpoint is reachable;
// the resource always exists in the tenant and cannot be deleted.
type CrossTenantAccessDefaultSettingsTestResource struct{}

// Exists checks whether the cross-tenant access policy default settings are readable from
// Microsoft Graph.
func (r CrossTenantAccessDefaultSettingsTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, _ *terraform.InstanceState) error {
		_, err := client.Policies().CrossTenantAccessPolicy().DefaultEscaped().Get(ctx, nil)
		return err
	})
}
