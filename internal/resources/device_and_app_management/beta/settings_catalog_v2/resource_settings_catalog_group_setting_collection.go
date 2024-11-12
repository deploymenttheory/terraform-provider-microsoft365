package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// GroupCollectionSchemaAttributeMap defines the common type for schema attribute maps
type GroupCollectionSchemaAttributeMap map[string]schema.Attribute

// GetGroupSettingCollectionSchema returns the root schema for group collection data type
func GetGroupSettingCollectionSchema(currentDepth int) GroupCollectionSchemaAttributeMap {
	if currentDepth >= MaxDepth {
		return GroupCollectionSchemaAttributeMap{}
	}

	return GroupCollectionSchemaAttributeMap{
		"odata_type": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The OData type of the group collection setting instance.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					DeviceManagementConfigurationGroupSettingValue,
				),
			},
		},
		"children": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: GetChildrenAttributes(currentDepth + 1),
			},
			Description:         "List of child setting instances that will be included in the group collection value",
			MarkdownDescription: "Collection of child settings that will be wrapped in a single group setting value.",
		},
	}
}
