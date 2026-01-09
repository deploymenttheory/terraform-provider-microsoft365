package graphBetaWindowsRemediationScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// WindowsRemediationScriptTestResource implements the types.TestResource interface for Windows remediation scripts
type WindowsRemediationScriptTestResource struct{}

// Exists checks whether the Windows remediation script exists in Microsoft Graph
func (r WindowsRemediationScriptTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId(state.ID).Get(ctx, nil)
		return err
	})
}
