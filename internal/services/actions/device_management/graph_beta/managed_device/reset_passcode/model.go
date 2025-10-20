// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-resetpasscode?view=graph-rest-beta
package graphBetaResetManagedDevicePasscode

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResetManagedDevicePasscodeActionModel struct {
	DeviceIDs types.List     `tfsdk:"device_ids"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
}
