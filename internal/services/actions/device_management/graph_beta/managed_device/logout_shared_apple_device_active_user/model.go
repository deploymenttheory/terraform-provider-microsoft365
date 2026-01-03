// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-logoutsharedappledeviceactiveuser?view=graph-rest-beta
package graphBetaLogoutSharedAppleDeviceActiveUser

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LogoutSharedAppleDeviceActiveUserActionModel struct {
	DeviceIDs types.List     `tfsdk:"device_ids"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
}
