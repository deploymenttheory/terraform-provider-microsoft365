package graphBetaMacOSPlatformScript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// MacOSPlatformScriptTestResource implements the types.TestResource interface for macOS platform scripts
type MacOSPlatformScriptTestResource struct{}

// Exists checks whether the macOS platform script exists in Microsoft Graph
func (r MacOSPlatformScriptTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().DeviceShellScripts().ByDeviceShellScriptId(state.ID).Get(ctx, nil)
		return err
	})
}
