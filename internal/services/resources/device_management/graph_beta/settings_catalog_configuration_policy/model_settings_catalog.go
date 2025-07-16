package graphBetaSettingsCatalogConfigurationPolicy

import "github.com/hashicorp/terraform-plugin-framework/types"

// DeviceConfigV2GraphServiceResourceModel is the root settings catalog model
// Updated for Terraform Framework with tfsdk tags matching schema attribute names
type DeviceConfigV2GraphServiceResourceModel struct {
	Settings []Setting `tfsdk:"settings"` // For array-based settings
}

// Setting represents a single setting detail
type Setting struct {
	ID              types.String    `tfsdk:"id"`
	SettingInstance SettingInstance `tfsdk:"setting_instance"`
}

// SettingInstance contains the core setting configuration
type SettingInstance struct {
	ODataType                        types.String                      `tfsdk:"odata_type"`
	SettingDefinitionId              types.String                      `tfsdk:"setting_definition_id"`
	SettingInstanceTemplateReference *SettingInstanceTemplateReference `tfsdk:"setting_instance_template_reference"`
	SimpleSettingValue               *SimpleSettingStruct              `tfsdk:"simple_setting_value"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct   `tfsdk:"simple_setting_collection_value"`
	ChoiceSettingValue               *ChoiceSettingStruct              `tfsdk:"choice_setting_value"`
	ChoiceSettingCollectionValue     []ChoiceSettingCollectionStruct   `tfsdk:"choice_setting_collection_value"`
	GroupSettingCollectionValue      []GroupSettingCollectionStruct    `tfsdk:"group_setting_collection_value"`
}

// SimpleSettingStruct represents a simple setting value
type SimpleSettingStruct struct {
	ODataType                     types.String                   `tfsdk:"odata_type"`
	SettingValueTemplateReference *SettingValueTemplateReference `tfsdk:"setting_value_template_reference"`
	Value                         types.String                   `tfsdk:"value"`
	ValueState                    types.String                   `tfsdk:"value_state"`
}

// SimpleSettingCollectionStruct represents a collection of simple settings
type SimpleSettingCollectionStruct struct {
	ODataType                     types.String                   `tfsdk:"odata_type"`
	SettingValueTemplateReference *SettingValueTemplateReference `tfsdk:"setting_value_template_reference"`
	Value                         types.String                   `tfsdk:"value"`
}

// ChoiceSettingChild represents a child element in a choice setting
type ChoiceSettingChild struct {
	ODataType                        types.String                      `tfsdk:"odata_type"`
	SettingDefinitionId              types.String                      `tfsdk:"setting_definition_id"`
	SettingInstanceTemplateReference *SettingInstanceTemplateReference `tfsdk:"setting_instance_template_reference"`
	ChoiceSettingValue               *ChoiceSettingStruct              `tfsdk:"choice_setting_value"`
	ChoiceSettingCollectionValue     []ChoiceSettingCollectionStruct   `tfsdk:"choice_setting_collection_value"`
	SimpleSettingValue               *SimpleSettingStruct              `tfsdk:"simple_setting_value"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct   `tfsdk:"simple_setting_collection_value"`
	GroupSettingCollectionValue      []GroupSettingCollectionStruct    `tfsdk:"group_setting_collection_value"`
}

// ChoiceSettingStruct represents a choice setting
type ChoiceSettingStruct struct {
	SettingValueTemplateReference *SettingValueTemplateReference `tfsdk:"setting_value_template_reference"`
	Value                         types.String                   `tfsdk:"value"`
	Children                      []ChoiceSettingChild           `tfsdk:"children"`
}

// ChoiceSettingCollectionChild represents a child element in a choice setting collection
type ChoiceSettingCollectionChild struct {
	ODataType                        types.String                      `tfsdk:"odata_type"`
	SettingDefinitionId              types.String                      `tfsdk:"setting_definition_id"`
	SettingInstanceTemplateReference *SettingInstanceTemplateReference `tfsdk:"setting_instance_template_reference"`
	SimpleSettingValue               *SimpleSettingStruct              `tfsdk:"simple_setting_value"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct   `tfsdk:"simple_setting_collection_value"`
}

// ChoiceSettingCollectionStruct represents a collection of choice settings
type ChoiceSettingCollectionStruct struct {
	SettingValueTemplateReference *SettingValueTemplateReference `tfsdk:"setting_value_template_reference"`
	Value                         types.String                   `tfsdk:"value"`
	Children                      []ChoiceSettingCollectionChild `tfsdk:"children"`
}

// GroupSettingCollectionChild represents a child element in a group setting collection
type GroupSettingCollectionChild struct {
	ODataType                        types.String                      `tfsdk:"odata_type"`
	SettingDefinitionId              types.String                      `tfsdk:"setting_definition_id"`
	SettingInstanceTemplateReference *SettingInstanceTemplateReference `tfsdk:"setting_instance_template_reference"`
	ChoiceSettingValue               *ChoiceSettingStruct              `tfsdk:"choice_setting_value"`
	ChoiceSettingCollectionValue     []ChoiceSettingCollectionStruct   `tfsdk:"choice_setting_collection_value"`
	GroupSettingCollectionValue      []GroupSettingCollectionStruct    `tfsdk:"group_setting_collection_value"`
	SimpleSettingValue               *SimpleSettingStruct              `tfsdk:"simple_setting_value"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct   `tfsdk:"simple_setting_collection_value"`
}

// GroupSettingCollectionStruct represents a collection of group settings
type GroupSettingCollectionStruct struct {
	SettingValueTemplateReference *SettingValueTemplateReference `tfsdk:"setting_value_template_reference"`
	Children                      []GroupSettingCollectionChild  `tfsdk:"children"`
}

// SettingInstanceTemplateReference represents the template reference at the instance level
type SettingInstanceTemplateReference struct {
	SettingInstanceTemplateId types.String `tfsdk:"setting_instance_template_id"`
}

// SettingValueTemplateReference represents the template reference at the value level
type SettingValueTemplateReference struct {
	SettingValueTemplateId types.String `tfsdk:"setting_value_template_id"`
	UseTemplateDefault     bool         `tfsdk:"use_template_default"`
}
