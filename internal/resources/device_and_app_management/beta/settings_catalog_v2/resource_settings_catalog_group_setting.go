package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// GroupSchemaAttributeMap defines the common type for schema attribute maps
type GroupSchemaAttributeMap map[string]schema.Attribute

// GetGroupSettingSchema returns the root schema for group data type
func GetGroupSettingSchema(currentDepth int) GroupSchemaAttributeMap {
	if currentDepth >= MaxDepth {
		return GroupSchemaAttributeMap{}
	}

	return GroupSchemaAttributeMap{
		"odata_type": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The OData type of the group setting value.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					DeviceManagementConfigurationChoiceSettingInstance,
					DeviceManagementConfigurationChoiceSettingCollectionInstance,
					DeviceManagementConfigurationSimpleSettingInstance,
					DeviceManagementConfigurationSimpleSettingCollectionInstance,
					DeviceManagementConfigurationSettingGroupInstance,
					DeviceManagementConfigurationGroupSettingInstance,
					DeviceManagementConfigurationSettingGroupCollectionInstance,
					DeviceManagementConfigurationGroupSettingCollectionInstance,
				),
			},
		},
		"children": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: GetChildrenAttributes(currentDepth + 1),
			},
			Description:         "List of child setting instances within this group",
			MarkdownDescription: "Collection of child settings that will be included in the group value.",
		},
	}
}
