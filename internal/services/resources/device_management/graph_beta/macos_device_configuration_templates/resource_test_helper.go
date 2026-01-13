package graphBetaMacosDeviceConfigurationTemplates

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type MacosDeviceConfigurationTemplatesTestResource struct{}

// Exists checks whether the macOS device configuration template exists in Microsoft Graph
func (r MacosDeviceConfigurationTemplatesTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().DeviceConfigurations().ByDeviceConfigurationId(state.ID).Get(ctx, nil)
		return err
	})
}
