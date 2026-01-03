// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-wipe?view=graph-rest-beta
package graphBetaWipeManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WipeManagedDeviceActionModel struct {
	DeviceIDs            types.List     `tfsdk:"device_ids"`
	KeepEnrollmentData   types.Bool     `tfsdk:"keep_enrollment_data"`
	KeepUserData         types.Bool     `tfsdk:"keep_user_data"`
	MacOsUnlockCode      types.String   `tfsdk:"macos_unlock_code"`
	ObliterationBehavior types.String   `tfsdk:"obliteration_behavior"`
	PersistEsimDataPlan  types.Bool     `tfsdk:"persist_esim_data_plan"`
	UseProtectedWipe     types.Bool     `tfsdk:"use_protected_wipe"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
