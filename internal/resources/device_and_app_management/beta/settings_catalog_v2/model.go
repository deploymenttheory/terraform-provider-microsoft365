// settings catalog profile ref: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// settings catalog setting ref: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsetting?view=graph-rest-beta
// settings catalog setting definition REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsettingdefinition?view=graph-rest-beta
// settings catalog setting instance REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicemanagementconfigurationsettinginstance?view=graph-rest-beta
// assignment REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicyassignment?view=graph-rest-beta
package graphBetaSettingsCatalog

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
	ODataType                        types.String                                                 `tfsdk:"odata_type"`
	SettingDefinitionID              types.String                                                 `tfsdk:"setting_definition_id"`
	SettingInstanceTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`

	// Different setting type values
	ChoiceSettingValue    *DeviceManagementConfigurationChoiceSettingValueResourceModel    `tfsdk:"choice"`
	ChoiceCollectionValue *DeviceManagementConfigurationChoiceCollectionValueResourceModel `tfsdk:"choice_collection"`
	SimpleSettingValue    *DeviceManagementConfigurationSimpleSettingValueResourceModel    `tfsdk:"simple"`
	SimpleCollectionValue *DeviceManagementConfigurationSimpleCollectionValueResourceModel `tfsdk:"simple_collection"`
	GroupSettingValue     *DeviceManagementConfigurationGroupSettingValueResourceModel     `tfsdk:"group"`
	GroupCollectionValue  *DeviceManagementConfigurationGroupCollectionValueResourceModel  `tfsdk:"group_collection"`
}

// DeviceManagementConfigurationChoiceSettingValueResourceModel represents the choice setting value
type DeviceManagementConfigurationChoiceSettingValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	IntValue                      types.Int32                                                  `tfsdk:"int_value"`
	StringValue                   types.String                                                 `tfsdk:"string_value"`
	Children                      []DeviceManagementConfigurationSettingInstance               `tfsdk:"children"`
}

// DeviceManagementConfigurationChoiceCollectionValueResourceModel represents the choice collection setting value
type DeviceManagementConfigurationChoiceCollectionValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	IntValue                      []types.Int32                                                `tfsdk:"int_value"`
	StringValue                   []types.String                                               `tfsdk:"string_value"`
	Children                      []DeviceManagementConfigurationSettingInstance               `tfsdk:"children"`
}

// DeviceManagementConfigurationSimpleSettingValueResourceModel represents the simple setting value
type DeviceManagementConfigurationSimpleSettingValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	IntValue                      types.Int32                                                  `tfsdk:"int_value"`
	StringValue                   types.String                                                 `tfsdk:"string_value"`
	Children                      []DeviceManagementConfigurationSettingInstance               `tfsdk:"children"`
}

// DeviceManagementConfigurationSimpleCollectionValueResourceModel represents the simple collection setting value
type DeviceManagementConfigurationSimpleCollectionValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	IntValue                      []types.Int32                                                `tfsdk:"int_value"`
	StringValue                   []types.String                                               `tfsdk:"string_value"`
	Children                      []DeviceManagementConfigurationSettingInstance               `tfsdk:"children"`
}

// DeviceManagementConfigurationGroupSettingValueResourceModel represents the group setting value
type DeviceManagementConfigurationGroupSettingValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	Children                      []DeviceManagementConfigurationSettingInstance               `tfsdk:"children"`
}

// DeviceManagementConfigurationGroupCollectionValueResourceModel represents the group collection setting value
type DeviceManagementConfigurationGroupCollectionValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	Children                      []DeviceManagementConfigurationSettingInstance               `tfsdk:"children"`
}

// DeviceManagementConfigurationTemplateReferenceResourceModel represents the setting instance template reference
type DeviceManagementConfigurationTemplateReferenceResourceModel struct {
	SettingInstanceTemplateId types.String `tfsdk:"setting_instance_template_id"`
	UseTemplateDefault        types.Bool   `tfsdk:"use_template_default"`
}
