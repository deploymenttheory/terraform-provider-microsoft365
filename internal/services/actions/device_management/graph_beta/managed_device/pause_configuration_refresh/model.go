// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-pauseconfigurationrefresh?view=graph-rest-beta
package graphBetaPauseConfigurationRefreshManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PauseConfigurationRefreshManagedDeviceActionModel struct {
	ManagedDevices        []ManagedDevicePauseConfig   `tfsdk:"managed_devices"`
	ComanagedDevices      []ComanagedDevicePauseConfig `tfsdk:"comanaged_devices"`
	IgnorePartialFailures types.Bool                   `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool                   `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value               `tfsdk:"timeouts"`
}

type ManagedDevicePauseConfig struct {
	DeviceID                 types.String `tfsdk:"device_id"`
	PauseTimePeriodInMinutes types.Int64  `tfsdk:"pause_time_period_in_minutes"`
}

type ComanagedDevicePauseConfig struct {
	DeviceID                 types.String `tfsdk:"device_id"`
	PauseTimePeriodInMinutes types.Int64  `tfsdk:"pause_time_period_in_minutes"`
}
