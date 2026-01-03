// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-setdevicename?view=graph-rest-beta
package graphBetaSetDeviceNameManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SetDeviceNameManagedDeviceActionModel struct {
	ManagedDevices   []ManagedDeviceSetName   `tfsdk:"managed_devices"`
	ComanagedDevices []ComanagedDeviceSetName `tfsdk:"comanaged_devices"`
	Timeouts         timeouts.Value           `tfsdk:"timeouts"`
}

type ManagedDeviceSetName struct {
	DeviceID   types.String `tfsdk:"device_id"`
	DeviceName types.String `tfsdk:"device_name"`
}

type ComanagedDeviceSetName struct {
	DeviceID   types.String `tfsdk:"device_id"`
	DeviceName types.String `tfsdk:"device_name"`
}
