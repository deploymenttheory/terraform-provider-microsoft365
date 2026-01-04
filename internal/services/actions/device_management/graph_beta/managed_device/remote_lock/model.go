// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-remotelock?view=graph-rest-beta
package graphBetaRemoteLockManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RemoteLockManagedDeviceActionModel struct {
	DeviceIDs             types.List     `tfsdk:"device_ids"`
	IgnorePartialFailures types.Bool     `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool     `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
