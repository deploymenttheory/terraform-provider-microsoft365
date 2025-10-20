// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-windowsdefenderscan?view=graph-rest-beta
package graphBetaWindowsDefenderScan

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsDefenderScanActionModel struct {
	ManagedDevices   []ManagedDeviceScan   `tfsdk:"managed_devices"`
	ComanagedDevices []ComanagedDeviceScan `tfsdk:"comanaged_devices"`
	Timeouts         timeouts.Value        `tfsdk:"timeouts"`
}

type ManagedDeviceScan struct {
	DeviceID  types.String `tfsdk:"device_id"`
	QuickScan types.Bool   `tfsdk:"quick_scan"`
}

type ComanagedDeviceScan struct {
	DeviceID  types.String `tfsdk:"device_id"`
	QuickScan types.Bool   `tfsdk:"quick_scan"`
}
