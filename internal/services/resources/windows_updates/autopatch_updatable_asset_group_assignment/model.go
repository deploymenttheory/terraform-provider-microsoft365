package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResourceModel struct {
	ID                     types.String   `tfsdk:"id"`
	UpdatableAssetGroupId  types.String   `tfsdk:"updatable_asset_group_id"`
	EntraDeviceIds         types.Set      `tfsdk:"entra_device_ids"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
}
