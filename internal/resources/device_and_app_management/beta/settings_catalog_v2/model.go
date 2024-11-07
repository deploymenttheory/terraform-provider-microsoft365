package graphBetaSettingsCatalog

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsSettingsCatalogProfileResourceModel holds the configuration for a Settings Catalog profile.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
type WindowsSettingsCatalogProfileResourceModel struct {
	ID                   types.String                                                 `tfsdk:"id"`
	Name                 types.String                                                 `tfsdk:"name"`
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

// DeviceManagementConfigurationSettingResourceModel holds individual settings for the catalog configuration policy.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsetting?view=graph-rest-beta
type DeviceManagementConfigurationSettingResourceModel struct {
	ODataType       types.String                                  `tfsdk:"odata_type"`
	SettingInstance *DeviceManagementConfigurationSettingInstance `tfsdk:"setting_instance"`
}

// DeviceManagementConfigurationSettingInstance represents a setting instance within the catalog.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicemanagementconfigurationsettinginstance?view=graph-rest-beta
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

// DeviceManagementConfigurationChoiceSettingValueResourceModel represents the choice setting value.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingValue?view=graph-rest-beta
type DeviceManagementConfigurationChoiceSettingValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	IntValue                      types.Int32                                                  `tfsdk:"int_value"`
	StringValue                   types.String                                                 `tfsdk:"string_value"`
	SecretValue                   types.String                                                 `tfsdk:"secret_value"`
	State                         types.String                                                 `tfsdk:"state"`    // Encryption state for secrets
	Children                      []*DeviceManagementConfigurationSettingInstance              `tfsdk:"children"` // Recursive handling
}

// DeviceManagementConfigurationChoiceCollectionValueResourceModel represents the choice collection setting value.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingCollectionValue?view=graph-rest-beta
type DeviceManagementConfigurationChoiceCollectionValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	IntValue                      []types.Int32                                                `tfsdk:"int_value"`
	StringValue                   []types.String                                               `tfsdk:"string_value"`
	SecretValue                   types.String                                                 `tfsdk:"secret_value"`
	State                         types.String                                                 `tfsdk:"state"`    // Encryption state for secrets
	Children                      []*DeviceManagementConfigurationSettingInstance              `tfsdk:"children"` // Recursive handling
}

// DeviceManagementConfigurationSimpleSettingValueResourceModel represents the simple setting value.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingValue?view=graph-rest-beta
type DeviceManagementConfigurationSimpleSettingValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	IntValue                      types.Int32                                                  `tfsdk:"int_value"`
	StringValue                   types.String                                                 `tfsdk:"string_value"`
	SecretValue                   types.String                                                 `tfsdk:"secret_value"`
	State                         types.String                                                 `tfsdk:"state"`    // Encryption state for secrets
	Children                      []*DeviceManagementConfigurationSettingInstance              `tfsdk:"children"` // Recursive handling
}

// DeviceManagementConfigurationSimpleCollectionValueResourceModel represents the simple collection setting value.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingCollectionValue?view=graph-rest-beta
type DeviceManagementConfigurationSimpleCollectionValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	IntValue                      []types.Int32                                                `tfsdk:"int_value"`
	StringValue                   []types.String                                               `tfsdk:"string_value"`
	SecretValue                   types.String                                                 `tfsdk:"secret_value"`
	State                         types.String                                                 `tfsdk:"state"`    // Encryption state for secrets
	Children                      []*DeviceManagementConfigurationSettingInstance              `tfsdk:"children"` // Recursive handling
}

// DeviceManagementConfigurationGroupSettingValueResourceModel represents the group setting value.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingValue?view=graph-rest-beta
type DeviceManagementConfigurationGroupSettingValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	Children                      []*DeviceManagementConfigurationSettingInstance              `tfsdk:"children"` // Recursive handling
}

// DeviceManagementConfigurationGroupCollectionValueResourceModel represents the group collection setting value.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingCollectionValue?view=graph-rest-beta
type DeviceManagementConfigurationGroupCollectionValueResourceModel struct {
	ODataType                     types.String                                                 `tfsdk:"odata_type"`
	SettingValueTemplateReference *DeviceManagementConfigurationTemplateReferenceResourceModel `tfsdk:"template_reference"`
	Children                      []*DeviceManagementConfigurationSettingInstance              `tfsdk:"children"` // Recursive handling
}

// DeviceManagementConfigurationTemplateReferenceResourceModel represents the setting instance template reference.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSettingInstanceTemplateReference?view=graph-rest-beta
type DeviceManagementConfigurationTemplateReferenceResourceModel struct {
	SettingInstanceTemplateId types.String `tfsdk:"setting_instance_template_id"`
	UseTemplateDefault        types.Bool   `tfsdk:"use_template_default"`
}
