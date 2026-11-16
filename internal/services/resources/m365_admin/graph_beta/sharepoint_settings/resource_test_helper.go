package graphBetaM365AdminSharepointSettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// SharepointSettingsTestResource implements the types.TestResource interface for the SharePoint settings
type SharepointSettingsTestResource struct{}

// Exists checks whether the singleton SharePoint settings are readable in Microsoft Graph
func (r SharepointSettingsTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, _ *terraform.InstanceState) error {
		_, err := client.Admin().Sharepoint().Settings().Get(ctx, nil)
		return err
	})
}
