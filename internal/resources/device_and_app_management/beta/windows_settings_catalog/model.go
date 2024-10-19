// setting catalog settings ref: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsetting?view=graph-rest-beta
// settings catalog setting definition REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsettingdefinition?view=graph-rest-beta
package graphBetaWindowsSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsSettingsCatalogProfileResourceModel struct to hold the configuration for a Settings Catalog profile
type WindowsSettingsCatalogProfileResourceModel struct {
	ID                   types.String      `tfsdk:"id"`
	DisplayName          types.String      `tfsdk:"display_name"`
	Description          types.String      `tfsdk:"description"`
	Platforms            types.String      `tfsdk:"platforms"`
	Technologies         types.String      `tfsdk:"technologies"`
	SettingsCount        types.Int64       `tfsdk:"settings_count"`
	Name                 types.String      `tfsdk:"name"`
	CreationSource       types.String      `tfsdk:"creation_source"`
	RoleScopeTagIds      []types.String    `tfsdk:"role_scope_tag_ids"`
	IsAssigned           types.Bool        `tfsdk:"is_assigned"`
	LastModifiedDateTime types.String      `tfsdk:"last_modified_date_time"`
	CreatedDateTime      types.String      `tfsdk:"created_date_time"`
	TemplateReference    TemplateReference `tfsdk:"template_reference"`
	Timeouts             timeouts.Value    `tfsdk:"timeouts"`
}

// WindowsSettingsCatalogSettingsResourceModel represents the top-level structure for settings
type WindowsSettingsCatalogSettingsResourceModel struct {
	Settings []DeviceManagementConfigurationSetting `tfsdk:"settings"`
}

// DeviceManagementConfigurationSetting represents a single setting
type DeviceManagementConfigurationSetting struct {
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
	Value                         types.String                                                `tfsdk:"value"`
	Children                      []DeviceManagementConfigurationSettingInstance              `tfsdk:"children"`
}

// DeviceManagementConfigurationSettingValueTemplateReference represents the value template reference
type DeviceManagementConfigurationSettingValueTemplateReference struct {
	ODataType              types.String `tfsdk:"odata_type"`
	SettingValueTemplateID types.String `tfsdk:"setting_value_template_id"`
	UseTemplateDefault     types.Bool   `tfsdk:"use_template_default"`
}

// SettingDefinition struct to hold setting definition information
type SettingDefinition struct {
	SettingInstanceTemplateID types.String              `tfsdk:"setting_instance_template_id"`
	SettingDefinitionID       types.String              `tfsdk:"setting_definition_id"`
	DisplayName               types.String              `tfsdk:"display_name"`
	Description               types.String              `tfsdk:"description"`
	VersionInfo               types.String              `tfsdk:"version_info"`
	CategoryID                types.String              `tfsdk:"category_id"`
	Metadata                  SettingDefinitionMetadata `tfsdk:"metadata"`
	RootDefinitionID          types.String              `tfsdk:"root_definition_id"`
}

// SettingDefinitionMetadata struct to hold metadata for setting definitions
type SettingDefinitionMetadata struct {
	DefaultValue        types.String   `tfsdk:"default_value"`
	DependsOn           []types.String `tfsdk:"depends_on"`
	DependentOn         []types.String `tfsdk:"dependent_on"`
	EnforcementType     types.String   `tfsdk:"enforcement_type"`
	HelpText            types.String   `tfsdk:"help_text"`
	Keywords            []types.String `tfsdk:"keywords"`
	Placeholder         types.String   `tfsdk:"placeholder"`
	SupportedDefinition types.String   `tfsdk:"supported_definition"`
	UrlInfo             []UrlInfo      `tfsdk:"url_info"`
}

// UrlInfo struct to hold URL information for settings
type UrlInfo struct {
	Url         types.String `tfsdk:"url"`
	DisplayName types.String `tfsdk:"display_name"`
}

// TemplateReference struct to hold template reference information
type TemplateReference struct {
	TemplateID             types.String `tfsdk:"template_id"`
	TemplateDisplayName    types.String `tfsdk:"template_display_name"`
	TemplateDisplayVersion types.String `tfsdk:"template_display_version"`
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
