package sharedmodels

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SettingsCatalogProfileResourceModel holds the configuration for a Settings Catalog profile.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
type SettingsCatalogProfileResourceModel struct {
	ID                              types.String                                    `tfsdk:"id"`
	Name                            types.String                                    `tfsdk:"name"`
	Description                     types.String                                    `tfsdk:"description"`
	Platforms                       types.String                                    `tfsdk:"platforms"`
	Technologies                    []types.String                                  `tfsdk:"technologies"`
	ConfigurationPolicyTemplateType types.String                                    `tfsdk:"configuration_policy_template_type"`
	RoleScopeTagIds                 types.List                                      `tfsdk:"role_scope_tag_ids"`
	SettingsCount                   types.Int64                                     `tfsdk:"settings_count"`
	IsAssigned                      types.Bool                                      `tfsdk:"is_assigned"`
	LastModifiedDateTime            types.String                                    `tfsdk:"last_modified_date_time"`
	CreatedDateTime                 types.String                                    `tfsdk:"created_date_time"`
	Settings                        types.String                                    `tfsdk:"settings"`
	Assignments                     *SettingsCatalogSettingsAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts                        timeouts.Value                                  `tfsdk:"timeouts"`
}
