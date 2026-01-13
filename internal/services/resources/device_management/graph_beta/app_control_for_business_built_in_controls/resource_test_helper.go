package graphBetaAppControlForBusinessBuiltInControls

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// AppControlForBusinessBuiltInControlsTestResource implements the types.TestResource interface
type AppControlForBusinessBuiltInControlsTestResource struct{}

// Exists checks whether the app control for business built-in controls policy exists in Microsoft Graph
func (r AppControlForBusinessBuiltInControlsTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().ConfigurationPolicies().ByDeviceManagementConfigurationPolicyId(state.ID).Get(ctx, nil)
		return err
	})
}
