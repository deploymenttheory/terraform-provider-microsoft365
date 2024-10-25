// settings catalog profile ref: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// settings catalog setting definition REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsettingdefinition?view=graph-rest-beta
// settings catalog setting instance REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicemanagementconfigurationsettinginstance?view=graph-rest-beta
// assignment REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicyassignment?view=graph-rest-beta
package graphBetaWindowsSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsSettingsCatalogProfileResourceModel struct to hold the configuration for a Settings Catalog profile
// https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
type WindowsSettingsCatalogProfileResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	DisplayName          types.String   `tfsdk:"display_name"`
	Description          types.String   `tfsdk:"description"`
	Platforms            types.String   `tfsdk:"platforms"`
	Technologies         types.String   `tfsdk:"technologies"`
	SettingsCount        types.Int32    `tfsdk:"settings_count"`
	CreationSource       types.String   `tfsdk:"creation_source"`
	RoleScopeTagIds      []types.String `tfsdk:"role_scope_tag_ids"`
	IsAssigned           types.Bool     `tfsdk:"is_assigned"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	//TemplateReference    TemplateReference                                   `tfsdk:"template_reference"`
	Settings    []DeviceManagementConfigurationSettingResourceModel `tfsdk:"settings"`
	Assignments *SettingsCatalogSettingsAssignmentResourceModel     `tfsdk:"assignments"`
	Timeouts    timeouts.Value                                      `tfsdk:"timeouts"`
}

// TemplateReference struct to hold template reference information
// type TemplateReference struct {
// 	TemplateID             types.String `tfsdk:"template_id"`
// 	TemplateFamily         types.String `tfsdk:"template_family"`
// 	TemplateDisplayName    types.String `tfsdk:"template_display_name"`
// 	TemplateDisplayVersion types.String `tfsdk:"template_display_version"`
// }

// DeviceManagementConfigurationSetting struct to hold the settings catalog configuration settings.
// https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsetting?view=graph-rest-beta
type DeviceManagementConfigurationSettingResourceModel struct {
	ODataType       types.String                                  `tfsdk:"odata_type"`
	ID              types.String                                  `tfsdk:"id"`
	SettingInstance *DeviceManagementConfigurationSettingInstance `tfsdk:"setting_instance"`
}

// DeviceManagementConfigurationSettingInstance represents the setting instance
type DeviceManagementConfigurationSettingInstance struct {
	ODataType                        types.String                                                   `tfsdk:"odata_type"`
	SettingDefinitionID              types.String                                                   `tfsdk:"setting_definition_id"`
	SettingInstanceTemplateReference *DeviceManagementConfigurationSettingInstanceTemplateReference `tfsdk:"setting_instance_template_reference"`
	ChoiceSettingValue               *DeviceManagementConfigurationChoiceSettingValue               `tfsdk:"choice_setting_value"`
}

// DeviceManagementConfigurationSettingInstanceTemplateReference represents the template reference
type DeviceManagementConfigurationSettingInstanceTemplateReference struct {
	ODataType                 types.String `tfsdk:"odata_type"`
	SettingInstanceTemplateID types.String `tfsdk:"setting_instance_template_id"`
}

// DeviceManagementConfigurationChoiceSettingValue represents the choice setting value
type DeviceManagementConfigurationChoiceSettingValue struct {
	ODataType                     types.String                                                `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationSettingValueTemplateReference `tfsdk:"setting_value_template_reference"`
	IntValue                      types.Int32                                                 `tfsdk:"int_value"`
	StringValue                   types.String                                                `tfsdk:"string_value"`
	Children                      []DeviceManagementConfigurationSettingInstance              `tfsdk:"children"`
}

// DeviceManagementConfigurationSettingValueTemplateReference represents the value template reference
type DeviceManagementConfigurationSettingValueTemplateReference struct {
	ODataType              types.String `tfsdk:"odata_type"`
	SettingValueTemplateID types.String `tfsdk:"setting_value_template_id"`
	UseTemplateDefault     types.Bool   `tfsdk:"use_template_default"`
}

// SettingsCatalogSettingsAssignmentResourceModel struct to hold assignment configuration
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
