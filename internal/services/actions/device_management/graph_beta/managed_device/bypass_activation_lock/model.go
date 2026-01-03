// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-bypassactivationlock?view=graph-rest-beta
package graphBetaBypassActivationLockManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BypassActivationLockManagedDeviceActionModel struct {
	DeviceIDs             types.List     `tfsdk:"device_ids"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
	IgnorePartialFailures types.Bool     `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool     `tfsdk:"validate_device_exists"`
}
