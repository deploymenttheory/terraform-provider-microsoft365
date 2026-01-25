package graphBetaAppControlForBusinessPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// AppControlForBusinessPolicyTestResource implements the types.TestResource interface for app control policies
type AppControlForBusinessPolicyTestResource struct{}

// Exists checks whether the app control for business policy exists in Microsoft Graph
func (r AppControlForBusinessPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(state.ID).
			Get(ctx, nil)
		return err
	})
}
