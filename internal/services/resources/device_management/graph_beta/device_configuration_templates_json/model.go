package graphBetaDeviceConfigurationTemplatesJson

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceConfigurationTemplatesJsonResourceModel struct {
	ID              types.String   `tfsdk:"id"`
	DisplayName     types.String   `tfsdk:"display_name"`
	Description     types.String   `tfsdk:"description"`
	RoleScopeTagIds types.Set      `tfsdk:"role_scope_tag_ids"`
	Settings        types.String   `tfsdk:"settings"`
	Assignments     types.Set      `tfsdk:"assignments"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}
