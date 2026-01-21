package graphBetaSettingsCatalogConfigurationPolicyJson

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// SettingsCatalogJsonTestResource implements the TestResource interface for settings catalog policies
type SettingsCatalogJsonTestResource struct{}

// Exists checks if a settings catalog policy exists in Microsoft Graph
func (r SettingsCatalogJsonTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().ConfigurationPolicies().ByDeviceManagementConfigurationPolicyId(state.ID).Get(ctx, nil)
		return err
	})
}
