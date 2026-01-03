// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-getfilevaultkey?view=graph-rest-beta
package graphBetaGetFileVaultKeyManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GetFileVaultKeyManagedDeviceActionModel struct {
	ManagedDeviceIDs   types.List     `tfsdk:"managed_device_ids"`
	ComanagedDeviceIDs types.List     `tfsdk:"comanaged_device_ids"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}
