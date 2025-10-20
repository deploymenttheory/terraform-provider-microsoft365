// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-cleanwindowsdevice?view=graph-rest-beta
package graphBetaCleanWindowsManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CleanWindowsManagedDeviceActionModel struct {
	ManagedDevices        []ManagedDeviceCleanWindows   `tfsdk:"managed_devices"`
	ComanagedDevices      []ComanagedDeviceCleanWindows `tfsdk:"comanaged_devices"`
	IgnorePartialFailures types.Bool                    `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool                    `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value                `tfsdk:"timeouts"`
}

type ManagedDeviceCleanWindows struct {
	DeviceID     types.String `tfsdk:"device_id"`
	KeepUserData types.Bool   `tfsdk:"keep_user_data"`
}

type ComanagedDeviceCleanWindows struct {
	DeviceID     types.String `tfsdk:"device_id"`
	KeepUserData types.Bool   `tfsdk:"keep_user_data"`
}
