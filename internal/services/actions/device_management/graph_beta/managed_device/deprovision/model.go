// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deprovision?view=graph-rest-beta
package graphBetaDeprovisionManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeprovisionManagedDeviceActionModel struct {
	ManagedDevices        []ManagedDeviceDeprovision   `tfsdk:"managed_devices"`
	ComanagedDevices      []ComanagedDeviceDeprovision `tfsdk:"comanaged_devices"`
	IgnorePartialFailures types.Bool                   `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool                   `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value               `tfsdk:"timeouts"`
}

type ManagedDeviceDeprovision struct {
	DeviceID          types.String `tfsdk:"device_id"`
	DeprovisionReason types.String `tfsdk:"deprovision_reason"`
}

type ComanagedDeviceDeprovision struct {
	DeviceID          types.String `tfsdk:"device_id"`
	DeprovisionReason types.String `tfsdk:"deprovision_reason"`
}
