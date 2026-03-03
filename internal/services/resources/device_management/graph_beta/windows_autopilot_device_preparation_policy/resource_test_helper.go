package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// WindowsAutopilotDevicePreparationPolicyTestResource implements the types.TestResource interface
// for Windows Autopilot Device Preparation Policies.
type WindowsAutopilotDevicePreparationPolicyTestResource struct{}

// Exists checks whether the Windows Autopilot Device Preparation Policy exists in Microsoft Graph.
func (r WindowsAutopilotDevicePreparationPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().ConfigurationPolicies().ByDeviceManagementConfigurationPolicyId(state.ID).Get(ctx, nil)
		return err
	})
}
