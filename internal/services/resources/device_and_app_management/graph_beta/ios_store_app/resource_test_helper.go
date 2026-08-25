package graphBetaIOSStoreApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// IosStoreAppTestResource implements the types.TestResource interface for iOS store apps
type IosStoreAppTestResource struct{}

// Exists checks whether the iOS store app exists in Microsoft Graph
func (r IosStoreAppTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceAppManagement().MobileApps().ByMobileAppId(state.ID).Get(ctx, nil)
		return err
	})
}
