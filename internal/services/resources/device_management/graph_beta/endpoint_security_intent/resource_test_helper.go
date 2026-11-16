package graphBetaEndpointSecurityIntent

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// EndpointSecurityIntentTestResource implements the types.TestResource interface for endpoint security intents
type EndpointSecurityIntentTestResource struct{}

// Exists checks whether the endpoint security intent exists in Microsoft Graph
func (r EndpointSecurityIntentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().Intents().ByDeviceManagementIntentId(state.ID).Get(ctx, nil)
		return err
	})
}
