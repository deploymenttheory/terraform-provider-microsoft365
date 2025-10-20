// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-shutdown?view=graph-rest-beta
package graphBetaShutdownManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ShutdownManagedDeviceActionModel struct {
	DeviceIDs types.List     `tfsdk:"device_ids"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
}
