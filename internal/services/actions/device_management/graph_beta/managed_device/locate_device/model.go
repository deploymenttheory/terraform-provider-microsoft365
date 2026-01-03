// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-locatedevice?view=graph-rest-beta
package graphBetaLocateManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LocateManagedDeviceActionModel struct {
	DeviceIDs types.List     `tfsdk:"device_ids"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
}
