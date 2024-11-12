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
			Optional:            true,
			MarkdownDescription: "The OData type of the setting instance.",
		},
		"children": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: GetChildrenAttributes(currentDepth + 1),
			},
		},
	}
}
