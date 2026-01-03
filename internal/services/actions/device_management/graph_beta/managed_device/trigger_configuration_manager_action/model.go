// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-triggerconfigurationmanageraction?view=graph-rest-beta
package graphBetaTriggerConfigurationManagerActionManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TriggerConfigurationManagerActionManagedDeviceActionModel struct {
	ManagedDevices   []ManagedDeviceConfigManagerAction   `tfsdk:"managed_devices"`
	ComanagedDevices []ComanagedDeviceConfigManagerAction `tfsdk:"comanaged_devices"`
	Timeouts         timeouts.Value                       `tfsdk:"timeouts"`
}

type ManagedDeviceConfigManagerAction struct {
	DeviceID types.String `tfsdk:"device_id"`
	Action   types.String `tfsdk:"action"`
}

type ComanagedDeviceConfigManagerAction struct {
	DeviceID types.String `tfsdk:"device_id"`
	Action   types.String `tfsdk:"action"`
}
