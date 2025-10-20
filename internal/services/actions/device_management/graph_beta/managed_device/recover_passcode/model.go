// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-recoverpasscode?view=graph-rest-beta
package graphBetaRecoverManagedDevicePasscode

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RecoverManagedDevicePasscodeActionModel struct {
	DeviceIDs types.List     `tfsdk:"device_ids"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
}
