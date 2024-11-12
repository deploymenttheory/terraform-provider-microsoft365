package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// GroupSchemaAttributeMap defines the common type for schema attribute maps
type GroupSchemaAttributeMap map[string]schema.Attribute

// GetGroupSchema returns the root schema for group data type
func GetGroupSchema(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:    true,
				Attributes:  getGroupValueAttributes(true, currentDepth+1),
				Description: "Group setting value configuration",
				MarkdownDescription: "Group setting value (#microsoft.graph.deviceManagementConfigurationGroupSettingValue) / " +
					"For more details, see [GroupSettingValue Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingValue?view=graph-rest-beta).",
			},
		},
		Description: "Group setting instance configuration",
		MarkdownDescription: "Instance configuration for group setting (#microsoft.graph.deviceManagementConfigurationGroupSettingInstance) / " +
			"For more details, see [GroupSettingInstance Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingInstance?view=graph-rest-beta).",
	}
}

// getGroupValueAttributes returns group value attributes
func getGroupValueAttributes(includeChildren bool, currentDepth int) GroupSchemaAttributeMap {
	attrs := GroupSchemaAttributeMap{
		"value": schema.StringAttribute{
			Optional:            true,
			Description:         "Identifier for group setting value",
			MarkdownDescription: "Specifies the unique identifier for group setting value.",
		},
	}

	if includeChildren && currentDepth < MaxDepth {
		attrs["children"] = getGroupChildSettingsAttribute(currentDepth + 1)
	}

	return attrs
}

// getGroupChildSettingsAttribute returns child settings list attribute
func getGroupChildSettingsAttribute(currentDepth int) schema.ListNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.ListNestedAttribute{}
	}

	return schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getGroupInstanceAttributes(currentDepth),
		},
		Computed:            true,
		Description:         "List of child setting configurations",
		MarkdownDescription: "List of child setting instances under group setting configuration.",
	}
}

// getGroupInstanceAttributes to include all nested types within the group
func getGroupInstanceAttributes(currentDepth int) GroupSchemaAttributeMap {
	if currentDepth >= MaxDepth {
		return getGroupBaseInstanceAttributes()
	}

	attrs := getGroupBaseInstanceAttributes()
	attrs["choice"] = getGroupChoiceSettingInstance(true, currentDepth+1)
	attrs["choice_collection"] = getGroupChoiceCollectionInstance(currentDepth + 1)
	attrs["group"] = getGroupInstance(currentDepth + 1)
	attrs["group_collection"] = getGroupCollectionInstance(currentDepth + 1)
	attrs["simple"] = getGroupSimpleInstance(currentDepth + 1)
	attrs["simple_collection"] = getGroupSimpleCollectionInstance(currentDepth + 1)

	return attrs
}

func getGroupBaseInstanceAttributes() GroupSchemaAttributeMap {
	return GroupSchemaAttributeMap{
		"setting_definition_id": schema.StringAttribute{
			Required:            true,
			Description:         `settingDefinitionId`,
			MarkdownDescription: "Setting Definition Id (#microsoft.graph.deviceManagementConfigurationSettingInstanceTemplateReference)",
		},
	}
}

// getGroupChoiceSettingInstance returns choice setting instance schema for groups
func getGroupChoiceSettingInstance(includeChildren bool, currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          getGroupValueAttributes(includeChildren, currentDepth+1),
				Description:         "Configuration for choice setting value",
				MarkdownDescription: "Configuration of the value for group choice setting (#microsoft.graph.deviceManagementConfigurationChoiceSettingValue).",
			},
		},
		Description:         "Group choice setting instance configuration",
		MarkdownDescription: "Instance configuration of choice setting in group (#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance).",
	}
}

// getGroupChoiceCollectionInstance returns choice collection instance schema for groups
func getGroupChoiceCollectionInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getGroupValueAttributes(true, currentDepth+1),
				},
				Description:         "Configuration for collection of choice setting values.",
				MarkdownDescription: "Instance configuration for a collection of choice setting values (#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance).",
			},
		},
		Description:         "Choice setting collection instance configuration.",
		MarkdownDescription: "Instance configuration for a collection of choice settings in group (#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance).",
	}
}

// getGroupInstance returns group instance schema
func getGroupInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          getGroupValueAttributes(true, currentDepth+1),
				Description:         "Configuration for group setting value.",
				MarkdownDescription: "Configuration of a group setting value (#microsoft.graph.deviceManagementConfigurationGroupSettingInstance).",
			},
		},
		Description:         "Group setting instance configuration.",
		MarkdownDescription: "Configuration for a single instance of group setting (#microsoft.graph.deviceManagementConfigurationGroupSettingInstance).",
	}
}

// getGroupCollectionInstance returns group collection instance schema
func getGroupCollectionInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getGroupValueAttributes(true, currentDepth+1),
				},
				Description:         "Configuration for collection of group setting values.",
				MarkdownDescription: "Instance configuration for a collection of group setting values (#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance).",
			},
		},
		Description:         "Group setting collection instance configuration.",
		MarkdownDescription: "Instance configuration for a collection of group settings in group (#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance).",
	}
}

// getGroupSimpleInstance returns simple setting instance schema for groups
func getGroupSimpleInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          deviceManagementConfigurationSimpleSettingValueAttributes,
				Description:         "Configuration of simple setting value.",
				MarkdownDescription: "Simple setting instance value in group (#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance).",
			},
		},
		Description:         "Simple setting instance configuration.",
		MarkdownDescription: "Configuration for a simple setting instance in group (#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance).",
	}
}

// getGroupSimpleCollectionInstance returns simple collection instance schema for groups
func getGroupSimpleCollectionInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: GroupSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: deviceManagementConfigurationSimpleSettingValueAttributes,
				},
				Description:         "Configuration of simple setting collection values.",
				MarkdownDescription: "List of values within a SimpleSettingCollection instance in group (#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance).",
			},
		},
		Description:         "Simple setting collection instance configuration.",
		MarkdownDescription: "Configuration for an instance of a simple setting collection in group (#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance).",
	}
}
