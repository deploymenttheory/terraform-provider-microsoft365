// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-revokeapplevpplicenses?view=graph-rest-beta
package graphBetaRevokeAppleVppLicenses

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RevokeAppleVppLicensesActionModel struct {
	ManagedDeviceIDs   types.List     `tfsdk:"managed_device_ids"`
	ComanagedDeviceIDs types.List     `tfsdk:"comanaged_device_ids"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}
