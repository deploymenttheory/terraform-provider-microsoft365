package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
