package graphBetaSettingsCatalogConfigurationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SettingsCatalogProfileResourceModel holds the configuration for a Settings Catalog profile.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
type SettingsCatalogProfileResourceModel struct {
	ID                   types.String                             `tfsdk:"id"`
	Name                 types.String                             `tfsdk:"name"`
	Description          types.String                             `tfsdk:"description"`
	Platforms            types.String                             `tfsdk:"platforms"`
	Technologies         types.List                               `tfsdk:"technologies"`
	RoleScopeTagIds      types.Set                                `tfsdk:"role_scope_tag_ids"`
	SettingsCount        types.Int32                              `tfsdk:"settings_count"`
	IsAssigned           types.Bool                               `tfsdk:"is_assigned"`
	LastModifiedDateTime types.String                             `tfsdk:"last_modified_date_time"`
	CreatedDateTime      types.String                             `tfsdk:"created_date_time"`
	ConfigurationPolicy  *DeviceConfigV2GraphServiceResourceModel `tfsdk:"configuration_policy"`
	Assignments          types.Set                                `tfsdk:"assignments"`
	TemplateReference    *TemplateReferenceResourceModel          `tfsdk:"template_reference"`
	Timeouts             timeouts.Value                           `tfsdk:"timeouts"`
}

// TemplateReferenceResourceModel struct to hold template reference configuration
type TemplateReferenceResourceModel struct {
	TemplateId             types.String `tfsdk:"template_id"`
	TemplateFamily         types.String `tfsdk:"template_family"`
	TemplateDisplayName    types.String `tfsdk:"template_display_name"`
	TemplateDisplayVersion types.String `tfsdk:"template_display_version"`
}
