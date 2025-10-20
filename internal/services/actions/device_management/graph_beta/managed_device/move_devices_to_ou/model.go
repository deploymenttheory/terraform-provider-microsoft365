// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-movedevicestoou?view=graph-rest-beta
package graphBetaMoveDevicesToOUManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MoveDevicesToOUManagedDeviceActionModel struct {
	OrganizationalUnitPath types.String   `tfsdk:"organizational_unit_path"`
	ManagedDeviceIDs       types.List     `tfsdk:"managed_device_ids"`
	ComanagedDeviceIDs     types.List     `tfsdk:"comanaged_device_ids"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
}
