// settings catalog profile ref: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// settings catalog setting definition REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsettingdefinition?view=graph-rest-beta
// settings catalog setting instance REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicemanagementconfigurationsettinginstance?view=graph-rest-beta
// assignment REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicyassignment?view=graph-rest-beta
package graphBetaWindowsSettingsCatalog

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsSettingsCatalogProfileResourceModel struct to hold the configuration for a Settings Catalog profile
// https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
type WindowsSettingsCatalogProfileResourceModel struct {
	ID                   types.String                                                 `tfsdk:"id"`
	DisplayName          types.String                                                 `tfsdk:"display_name"`
	Description          types.String                                                 `tfsdk:"description"`
	Platforms            types.String                                                 `tfsdk:"platforms"`
	Technologies         types.String                                                 `tfsdk:"technologies"`
	SettingsCount        types.Int32                                                  `tfsdk:"settings_count"`
	RoleScopeTagIds      []types.String                                               `tfsdk:"role_scope_tag_ids"`
	LastModifiedDateTime types.String                                                 `tfsdk:"last_modified_date_time"`
	CreatedDateTime      types.String                                                 `tfsdk:"created_date_time"`
	Settings             []DeviceManagementConfigurationSettingResourceModel          `tfsdk:"settings"`
	Assignments          *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts             timeouts.Value                                               `tfsdk:"timeouts"`
}

// DeviceManagementConfigurationSetting struct to hold the settings catalog configuration settings.
// https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsetting?view=graph-rest-beta
type DeviceManagementConfigurationSettingResourceModel struct {
	ODataType       types.String                                  `tfsdk:"odata_type"`
	SettingInstance *DeviceManagementConfigurationSettingInstance `tfsdk:"setting_instance"`
}

// DeviceManagementConfigurationSettingInstance represents the setting instance
type DeviceManagementConfigurationSettingInstance struct {
	ODataType           types.String                                     `tfsdk:"odata_type"`
	SettingDefinitionID types.String                                     `tfsdk:"setting_definition_id"`
	ChoiceSettingValue  *DeviceManagementConfigurationChoiceSettingValue `tfsdk:"choice_setting_value"`
}

// DeviceManagementConfigurationChoiceSettingValue represents the choice setting value
type DeviceManagementConfigurationChoiceSettingValue struct {
	ODataType   types.String                                   `tfsdk:"odata_type"`
	IntValue    types.Int32                                    `tfsdk:"int_value"`
	StringValue types.String                                   `tfsdk:"string_value"`
	Children    []DeviceManagementConfigurationSettingInstance `tfsdk:"children"`
}
