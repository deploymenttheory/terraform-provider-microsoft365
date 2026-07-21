package graphBetaDeviceAndAppManagementAndroidManagedAppProtection

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// AndroidManagedAppProtectionTestResource is used by acceptance tests to verify
// the resource exists or has been destroyed in the remote API.
type AndroidManagedAppProtectionTestResource struct{}

func (r AndroidManagedAppProtectionTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.
			DeviceAppManagement().
			AndroidManagedAppProtections().
			ByAndroidManagedAppProtectionId(state.ID).
			Get(ctx, nil)
		return err
	})
}
