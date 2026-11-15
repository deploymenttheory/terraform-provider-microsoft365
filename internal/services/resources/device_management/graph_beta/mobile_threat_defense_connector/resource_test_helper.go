package graphBetaMobileThreatDefenseConnector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// MobileThreatDefenseConnectorTestResource implements the types.TestResource interface for mobile threat defense connectors
type MobileThreatDefenseConnectorTestResource struct{}

// Exists checks whether the mobile threat defense connector exists in Microsoft Graph
func (r MobileThreatDefenseConnectorTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().MobileThreatDefenseConnectors().ByMobileThreatDefenseConnectorId(state.ID).Get(ctx, nil)
		return err
	})
}
