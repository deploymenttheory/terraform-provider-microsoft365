package graphBetaSettingsCatalog

import (
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// GroupCollectionSchemaAttributeMap defines the common type for schema attribute maps
type GroupCollectionSchemaAttributeMap map[string]schema.Attribute

// GetGroupCollectionSchema returns the root schema for group collection data type
func GetGroupCollectionSchema(currentDepth int) GroupCollectionSchemaAttributeMap {
	if currentDepth >= MaxDepth {
		return GroupCollectionSchemaAttributeMap{}
	}

	return GroupCollectionSchemaAttributeMap{
		"odata_type": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The OData type of the setting instance. This is automatically set by the Graph SDK during request construction.",
			PlanModifiers: []planmodifier.String{
				planmodifiers.UseStateForUnknownString(),
			},
		},
		"children": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: GetChildrenAttributes(currentDepth + 1),
			},
			Description:         "List of child setting configurations",
			MarkdownDescription: "List of child setting instances under group setting configuration.",
		},
	}
}
