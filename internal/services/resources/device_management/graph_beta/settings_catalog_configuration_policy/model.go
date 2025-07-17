package graphBetaSettingsCatalogConfigurationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SettingsCatalogProfileResourceModel holds the configuration for a Settings Catalog profile.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
type SettingsCatalogProfileResourceModel struct {
	ID                   types.String                                    `tfsdk:"id"`
	Name                 types.String                                    `tfsdk:"name"`
	Description          types.String                                    `tfsdk:"description"`
	Platforms            types.String                                    `tfsdk:"platforms"`
	Technologies         types.List                                      `tfsdk:"technologies"`
	RoleScopeTagIds      types.Set                                       `tfsdk:"role_scope_tag_ids"`
	SettingsCount        types.Int32                                     `tfsdk:"settings_count"`
	IsAssigned           types.Bool                                      `tfsdk:"is_assigned"`
	LastModifiedDateTime types.String                                    `tfsdk:"last_modified_date_time"`
	CreatedDateTime      types.String                                    `tfsdk:"created_date_time"`
	ConfigurationPolicy  *DeviceConfigV2GraphServiceResourceModel        `tfsdk:"configuration_policy"`
	Assignments          *SettingsCatalogSettingsAssignmentResourceModel `tfsdk:"assignments"`
	TemplateReference    *TemplateReferenceResourceModel                 `tfsdk:"template_reference"`
	Timeouts             timeouts.Value                                  `tfsdk:"timeouts"`
}

// SettingsCatalogSettingsAssignmentResourceModel struct to hold device configuation assignment configuration
type SettingsCatalogSettingsAssignmentResourceModel struct {
	AllDevices           types.Bool     `tfsdk:"all_devices"`
	AllDevicesFilterType types.String   `tfsdk:"all_devices_filter_type"`
	AllDevicesFilterId   types.String   `tfsdk:"all_devices_filter_id"`
	AllUsers             types.Bool     `tfsdk:"all_users"`
	AllUsersFilterType   types.String   `tfsdk:"all_users_filter_type"`
	AllUsersFilterId     types.String   `tfsdk:"all_users_filter_id"`
	IncludeGroups        []IncludeGroup `tfsdk:"include_groups"`
	ExcludeGroupIds      []types.String `tfsdk:"exclude_group_ids"`
}

// IncludeGroup represents a group with its corresponding filter type and filter group ID
type IncludeGroup struct {
	GroupId                 types.String `tfsdk:"group_id"`
	IncludeGroupsFilterType types.String `tfsdk:"include_groups_filter_type"`
	IncludeGroupsFilterId   types.String `tfsdk:"include_groups_filter_id"`
}

// DeviceManagementScriptAssignmentResourceModel struct to hold platform script assignment configuration
type DeviceManagementScriptAssignmentResourceModel struct {
	AllDevices      types.Bool     `tfsdk:"all_devices"`
	AllUsers        types.Bool     `tfsdk:"all_users"`
	IncludeGroupIds []types.String `tfsdk:"include_group_ids"`
	ExcludeGroupIds []types.String `tfsdk:"exclude_group_ids"`
}

// TemplateReferenceResourceModel struct to hold template reference configuration
type TemplateReferenceResourceModel struct {
	TemplateId             types.String `tfsdk:"template_id"`
	TemplateFamily         types.String `tfsdk:"template_family"`
	TemplateDisplayName    types.String `tfsdk:"template_display_name"`
	TemplateDisplayVersion types.String `tfsdk:"template_display_version"`
}
