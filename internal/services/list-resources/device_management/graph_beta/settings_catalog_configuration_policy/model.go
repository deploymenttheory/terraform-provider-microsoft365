package graphBetaSettingsCatalogConfigurationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SettingsCatalogListConfigModel represents the configuration for listing Settings Catalog policies
type SettingsCatalogListConfigModel struct {
	NameFilter           types.String `tfsdk:"name_filter"`
	PlatformFilter       types.List   `tfsdk:"platform_filter"`
	TemplateFamilyFilter types.String `tfsdk:"template_family_filter"`
	IsAssignedFilter     types.Bool   `tfsdk:"is_assigned_filter"`
	ODataFilter          types.String `tfsdk:"odata_filter"`
}
