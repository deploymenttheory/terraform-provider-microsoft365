package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"

	// Graph Beta - Device management actions
	graphBetaDeviceManagementWindows365ApplyCloudPcProvisioningPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/actions/device_management/graph_beta/windows_365/apply_cloud_pc_provisioning_policy"
	graphBetaDeviceManagementWindowsAutopilotDeviceIdentityAllowNextEnrollment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/actions/device_management/graph_beta/windows_autopilot_device_identity/allow_next_enrollment"
	graphBetaDeviceManagementWindowsAutopilotDeviceIdentityAssignUserToDevice "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/actions/device_management/graph_beta/windows_autopilot_device_identity/assign_user_to_device"
	graphBetaDeviceManagementWindowsAutopilotDeviceIdentityUnassignUserFromDevice "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/actions/device_management/graph_beta/windows_autopilot_device_identity/unassign_user_from_device"
	graphBetaDeviceManagementWindowsAutopilotDeviceIdentityUpdateDeviceProperties "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/actions/device_management/graph_beta/windows_autopilot_device_identity/update_device_properties"
)

// Actions returns a slice of functions that each return an action.Action.
// This function is a method of the M365Provider type and takes a context.Context as an argument.
// The returned slice is intended to hold the Microsoft 365 provider actions.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeout.
//
// Returns:
//
//	[]func() action.Action: A slice of functions, each returning an action.Action.
//
// Actions returns a slice of functions that each return an action.Action.
func (p *M365Provider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{
		// Graph Beta - Device management actions
		graphBetaDeviceManagementWindowsAutopilotDeviceIdentityAllowNextEnrollment.NewAllowNextEnrollmentAction,
		graphBetaDeviceManagementWindowsAutopilotDeviceIdentityAssignUserToDevice.NewAssignUserToDeviceAction,
		graphBetaDeviceManagementWindowsAutopilotDeviceIdentityUnassignUserFromDevice.NewUnassignUserFromDeviceAction,
		graphBetaDeviceManagementWindowsAutopilotDeviceIdentityUpdateDeviceProperties.NewUpdateDevicePropertiesAction,
		graphBetaDeviceManagementWindows365ApplyCloudPcProvisioningPolicy.NewApplyCloudPcProvisioningPolicyAction,

		// Add microsoft 365 provider actions here
	}
}
