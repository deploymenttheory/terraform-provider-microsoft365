package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	EntraDeviceObjectIds  types.Set      `tfsdk:"entra_device_object_ids"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
