// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-activatedeviceesim?view=graph-rest-beta
package graphBetaActivateDeviceEsimManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ActivateDeviceEsimManagedDeviceActionModelV2 describes the action data model.
// Note: Attributes are at root level - Terraform Core handles the config{} wrapper in HCL
type ActivateDeviceEsimManagedDeviceActionModelV2 struct {
	ManagedDevices        []ManagedDeviceActivateEsim   `tfsdk:"managed_devices"`
	ComanagedDevices      []ComanagedDeviceActivateEsim `tfsdk:"comanaged_devices"`
	Timeouts              timeouts.Value                `tfsdk:"timeouts"`
	IgnorePartialFailures types.Bool                    `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool                    `tfsdk:"validate_device_exists"`
}

type ManagedDeviceActivateEsim struct {
	DeviceID   types.String `tfsdk:"device_id"`
	CarrierURL types.String `tfsdk:"carrier_url"`
}

type ComanagedDeviceActivateEsim struct {
	DeviceID   types.String `tfsdk:"device_id"`
	CarrierURL types.String `tfsdk:"carrier_url"`
}
