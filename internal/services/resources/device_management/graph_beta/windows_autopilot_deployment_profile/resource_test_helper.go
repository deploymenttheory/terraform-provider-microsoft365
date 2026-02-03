package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// WindowsAutopilotDeploymentProfileTestResource implements the types.TestResource interface for Windows Autopilot deployment profiles
type WindowsAutopilotDeploymentProfileTestResource struct{}

// Exists checks whether the Windows Autopilot deployment profile exists in Microsoft Graph
func (r WindowsAutopilotDeploymentProfileTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().WindowsAutopilotDeploymentProfiles().ByWindowsAutopilotDeploymentProfileId(state.ID).Get(ctx, nil)
		return err
	})
}
