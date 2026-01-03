// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-playlostmodesound?view=graph-rest-beta
package graphBetaPlayLostModeSoundManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PlayLostModeSoundManagedDeviceActionModel struct {
	ManagedDevices   []ManagedDevicePlaySound   `tfsdk:"managed_devices"`
	ComanagedDevices []ComanagedDevicePlaySound `tfsdk:"comanaged_devices"`
	Timeouts         timeouts.Value             `tfsdk:"timeouts"`
}

type ManagedDevicePlaySound struct {
	DeviceID          types.String `tfsdk:"device_id"`
	DurationInMinutes types.String `tfsdk:"duration_in_minutes"`
}

type ComanagedDevicePlaySound struct {
	DeviceID          types.String `tfsdk:"device_id"`
	DurationInMinutes types.String `tfsdk:"duration_in_minutes"`
}
