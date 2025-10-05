// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsrestoredeviceenrollmentconfiguration?view=graph-rest-beta
package graphBetaWindowsBackupAndRestore

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsBackupAndRestoreResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	State    types.String   `tfsdk:"state"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
