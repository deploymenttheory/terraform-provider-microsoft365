package graphBetaWindowsCustomConfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsCustomConfigurationResourceModel describes the resource data model for
// microsoft.graph.windows10CustomConfiguration.
type WindowsCustomConfigurationResourceModel struct {
	ID              types.String   `tfsdk:"id"`
	DisplayName     types.String   `tfsdk:"display_name"`
	Description     types.String   `tfsdk:"description"`
	RoleScopeTagIds types.Set      `tfsdk:"role_scope_tag_ids"`
	OmaSettings     types.List     `tfsdk:"oma_settings"`
	Assignments     types.Set      `tfsdk:"assignments"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}

// OmaSettingResourceModel describes a single OMA-URI setting (microsoft.graph.omaSetting and subtypes).
type OmaSettingResourceModel struct {
	OdataType   types.String `tfsdk:"odata_type"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	OmaUri      types.String `tfsdk:"oma_uri"`
	Value       types.String `tfsdk:"value"`
	FileName    types.String `tfsdk:"file_name"`
}
