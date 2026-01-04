// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatefilevaultkey?view=graph-rest-beta
package graphBetaRotateFileVaultKeyManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RotateFileVaultKeyManagedDeviceActionModel struct {
	ManagedDeviceIDs      types.List     `tfsdk:"managed_device_ids"`
	ComanagedDeviceIDs    types.List     `tfsdk:"comanaged_device_ids"`
	IgnorePartialFailures types.Bool     `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool     `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
