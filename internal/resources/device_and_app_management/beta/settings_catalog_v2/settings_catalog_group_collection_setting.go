package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// GroupCollectionSchemaAttributeMap defines the common type for schema attribute maps
type GroupCollectionSchemaAttributeMap map[string]schema.Attribute

// GetGroupCollectionSchema returns the root schema for group collection data type
func GetGroupCollectionSchema(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupCollectionSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getGroupSettingValueAttributes(currentDepth + 1),
				},
				Description:         "A collection of GroupSetting values",
				MarkdownDescription: "Group setting values within a GroupSettingCollection (#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance).",
			},
		},
		Description:         "Configuration for a GroupSettingCollection instance",
		MarkdownDescription: "A collection instance of grouped settings in Microsoft Graph (#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance).",
	}
}

// getGroupSettingValueAttributes returns attributes for group setting values
func getGroupSettingValueAttributes(currentDepth int) GroupCollectionSchemaAttributeMap {
	attrs := GroupCollectionSchemaAttributeMap{
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationSettingValueTemplateReferenceAttributes,
			Description:         "Setting value template reference",
			MarkdownDescription: "Reference for setting value template (#microsoft.graph.deviceManagementConfigurationSettingValueTemplateReference).",
		},
	}

	if currentDepth < MaxDepth {
		attrs["children"] = getGroupSettingChildrenAttribute(currentDepth + 1)
	}

	return attrs
}

// getGroupSettingChildrenAttribute returns nested attributes for child group settings
func getGroupSettingChildrenAttribute(currentDepth int) schema.ListNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.ListNestedAttribute{}
	}

	return schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getGroupSettingChildAttributes(currentDepth),
		},
		Computed:            true,
		Description:         "List of child settings within this GroupSetting",
		MarkdownDescription: "Child settings contained within the GroupSetting (#microsoft.graph.deviceManagementConfigurationGroupSettingValue).",
	}
}

// getGroupSettingChildAttributes defines attributes for nested child settings within a GroupSetting
func getGroupSettingChildAttributes(currentDepth int) GroupCollectionSchemaAttributeMap {
	attrs := GroupCollectionSchemaAttributeMap{
		"definition_id": schema.StringAttribute{
			Required:            true,
			Description:         "Setting Definition ID",
			MarkdownDescription: "Unique identifier for the setting definition (#microsoft.graph.deviceManagementConfigurationSettingInstanceTemplateReference).",
		},
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationSettingInstanceTemplateReferenceAttributes,
			Description:         "Setting instance template reference",
			MarkdownDescription: "Template reference for the setting instance (#microsoft.graph.deviceManagementConfigurationSettingInstanceTemplateReference).",
		},
	}

	if currentDepth < MaxDepth {
		attrs["choice"] = getGroupChoiceSettingInstance(true, currentDepth+1)
		attrs["group"] = getGroupSettingInstance(currentDepth + 1)
		attrs["simple"] = getGroupSimpleSettingInstance(currentDepth + 1)
	}

	return attrs
}

// getGroupChoiceSettingValueAttributes defines attributes for the value of a choice setting in Group
func getGroupChoiceSettingValueAttributes(currentDepth int) GroupCollectionSchemaAttributeMap {
	attrs := GroupCollectionSchemaAttributeMap{
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationSettingValueTemplateReferenceAttributes,
			Description:         "Template reference for choice setting value",
			MarkdownDescription: "Reference for choice setting value template (#microsoft.graph.deviceManagementConfigurationChoiceSettingValue).",
		},
	}

	if currentDepth < MaxDepth {
		attrs["children"] = getGroupSettingChildrenAttribute(currentDepth + 1)
	}

	return attrs
}

// getGroupSettingInstance returns schema for nested Group setting within GroupCollection
func getGroupSettingInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupCollectionSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          getGroupSettingValueAttributes(currentDepth + 1),
				Description:         "Configuration for group setting value",
				MarkdownDescription: "Configuration for group setting value instance (#microsoft.graph.deviceManagementConfigurationGroupSettingInstance).",
			},
		},
		Description:         "Group setting instance",
		MarkdownDescription: "Configuration of a single group setting instance within the collection (#microsoft.graph.deviceManagementConfigurationGroupSettingInstance).",
	}
}

// getGroupSimpleSettingInstance returns schema for a simple setting within GroupCollection
func getGroupSimpleSettingInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupCollectionSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          deviceManagementConfigurationSimpleSettingValueAttributes,
				Description:         "Configuration of simple setting value",
				MarkdownDescription: "Simple setting value configuration (#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance).",
			},
		},
		Description:         "Simple setting instance",
		MarkdownDescription: "Configuration for a simple setting instance in the collection (#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance).",
	}
}
