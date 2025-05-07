// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsdriverupdateinventory?view=graph-rest-beta
package graphBetaWindowsDriverUpdateInventory

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsDriverUpdateInventoryResourceModel struct {
	ID                           types.String   `tfsdk:"id"`
	Name                         types.String   `tfsdk:"name"`
	Version                      types.String   `tfsdk:"version"`
	Manufacturer                 types.String   `tfsdk:"manufacturer"`
	ReleaseDateTime              types.String   `tfsdk:"release_date_time"`
	DriverClass                  types.String   `tfsdk:"driver_class"`
	ApplicableDeviceCount        types.Int32    `tfsdk:"applicable_device_count"`
	ApprovalStatus               types.String   `tfsdk:"approval_status"`
	Category                     types.String   `tfsdk:"category"`
	DeployDateTime               types.String   `tfsdk:"deploy_date_time"`
	WindowsDriverUpdateProfileID types.String   `tfsdk:"windows_driver_update_profile_id"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
}
