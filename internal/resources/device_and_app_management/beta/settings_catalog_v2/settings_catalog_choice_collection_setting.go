package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// ChoiceCollectionSchemaAttributeMap defines the schema attribute map for choice collection settings
type ChoiceCollectionSchemaAttributeMap map[string]schema.Attribute

// GetChoiceCollectionSchema returns the root schema for choice collection data type
func GetChoiceCollectionSchema(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceCollectionSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getChoiceCollectionValueAttributes(true, currentDepth+1),
				},
				Description:         "Choice setting collection values",
				MarkdownDescription: "Collection of values within a ChoiceSettingCollection instance in Microsoft Graph.",
			},
		},
		Description: "Choice setting collection instance configuration",
		MarkdownDescription: "Configuration for a ChoiceSettingCollection instance (#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance) / " +
			"For more details, see [ChoiceSettingCollectionInstance Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingCollectionInstance?view=graph-rest-beta).",
	}
}

// getChoiceCollectionValueAttributes returns choice collection value attributes
func getChoiceCollectionValueAttributes(includeChildren bool, currentDepth int) ChoiceCollectionSchemaAttributeMap {
	attrs := ChoiceCollectionSchemaAttributeMap{
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationSettingValueTemplateReferenceAttributes,
			Description:         "Template reference for choice collection setting value",
			MarkdownDescription: "Template reference within ChoiceSettingCollection, providing template-based configuration options.",
		},
		"value": schema.StringAttribute{
			Optional:            true,
			Description:         "Identifier for choice setting collection value",
			MarkdownDescription: "Specifies the unique identifier for choice setting collection value.",
		},
	}

	if includeChildren && currentDepth < MaxDepth {
		attrs["children"] = getChoiceCollectionChildSettingsAttribute(currentDepth + 1)
	}

	return attrs
}

// getChoiceCollectionChildSettingsAttribute returns the list of child settings within the choice collection
func getChoiceCollectionChildSettingsAttribute(currentDepth int) schema.ListNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.ListNestedAttribute{}
	}

	return schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getChoiceInstanceAttributes(currentDepth),
		},
		Computed:            true,
		Description:         "List of child setting configurations within choice collection",
		MarkdownDescription: "List of child settings within a ChoiceSettingCollection instance.",
	}
}
