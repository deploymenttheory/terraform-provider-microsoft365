// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-enablelostmode?view=graph-rest-beta
package graphBetaEnableLostModeManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnableLostModeManagedDeviceActionModel struct {
	ManagedDevices   []ManagedDeviceLostMode   `tfsdk:"managed_devices"`
	ComanagedDevices []ComanagedDeviceLostMode `tfsdk:"comanaged_devices"`
	Timeouts         timeouts.Value            `tfsdk:"timeouts"`
}

type ManagedDeviceLostMode struct {
	DeviceID    types.String `tfsdk:"device_id"`
	Message     types.String `tfsdk:"message"`
	PhoneNumber types.String `tfsdk:"phone_number"`
	Footer      types.String `tfsdk:"footer"`
}

type ComanagedDeviceLostMode struct {
	DeviceID    types.String `tfsdk:"device_id"`
	Message     types.String `tfsdk:"message"`
	PhoneNumber types.String `tfsdk:"phone_number"`
	Footer      types.String `tfsdk:"footer"`
}
