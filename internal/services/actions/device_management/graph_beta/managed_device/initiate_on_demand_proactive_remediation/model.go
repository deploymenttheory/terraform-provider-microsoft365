// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiateondemandproactiveremediation?view=graph-rest-beta
package graphBetaInitiateOnDemandProactiveRemediationManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InitiateOnDemandProactiveRemediationManagedDeviceActionModel struct {
	ManagedDevices        []ManagedDeviceProactiveRemediation   `tfsdk:"managed_devices"`
	ComanagedDevices      []ComanagedDeviceProactiveRemediation `tfsdk:"comanaged_devices"`
	IgnorePartialFailures types.Bool                            `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool                            `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value                        `tfsdk:"timeouts"`
}

type ManagedDeviceProactiveRemediation struct {
	DeviceID       types.String `tfsdk:"device_id"`
	ScriptPolicyID types.String `tfsdk:"script_policy_id"`
}

type ComanagedDeviceProactiveRemediation struct {
	DeviceID       types.String `tfsdk:"device_id"`
	ScriptPolicyID types.String `tfsdk:"script_policy_id"`
}
