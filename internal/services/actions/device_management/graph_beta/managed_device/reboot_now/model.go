// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rebootnow?view=graph-rest-beta
package graphBetaRebootNowManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RebootNowManagedDeviceActionModel struct {
	DeviceIDs types.List     `tfsdk:"device_ids"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
}
