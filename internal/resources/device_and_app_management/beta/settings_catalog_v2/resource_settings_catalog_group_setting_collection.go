package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			MarkdownDescription: "The OData type of the setting instance.",
		},
		"group_setting_collection_value": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"children": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: GetChildrenAttributes(currentDepth + 1),
						},
						Description:         "List of child setting instances that will be included in the group collection value",
						MarkdownDescription: "Collection of child settings that will be wrapped in a single group setting value.",
					},
				},
			},
			Description: "Array of group setting collection values",
		},
	}
}
