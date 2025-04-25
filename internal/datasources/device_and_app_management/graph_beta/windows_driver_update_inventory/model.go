package graphBetaWindowsDriverUpdateInventory

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsDriverUpdateInventoryDataSourceModel defines the data source model
type WindowsDriverUpdateInventoryDataSourceModel struct {
	ID                           types.String   `tfsdk:"id"`
	Name                         types.String   `tfsdk:"name"`
	WindowsDriverUpdateProfileID types.String   `tfsdk:"windows_driver_update_profile_id"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
}
