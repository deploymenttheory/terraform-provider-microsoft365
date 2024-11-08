package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ChoiceSchemaAttributeMap defines the common type for schema attribute maps
type ChoiceSchemaAttributeMap map[string]schema.Attribute

// GetChoiceSchema returns the root schema for choice data type
func GetChoiceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          getChoiceValueAttributes(true),
				Description:         "Choice setting value configuration",
				MarkdownDescription: "Choice setting value (#microsoft.graph.deviceManagementConfigurationChoiceSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingValue?view=graph-rest-beta",
			},
		},
		Description:         "Choice setting instance configuration",
		MarkdownDescription: "Choice setting instance (#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingInstance?view=graph-rest-beta",
	}
}

// getChoiceValueAttributes returns choice value attributes
func getChoiceValueAttributes(includeChildren bool) ChoiceSchemaAttributeMap {
	attrs := ChoiceSchemaAttributeMap{
		"value": schema.StringAttribute{
			Optional:            true,
			Description:         "Choice setting value identifier",
			MarkdownDescription: "Choice setting value: an OptionDefinition ItemId",
		},
	}

	if includeChildren {
		attrs["children"] = getChoiceChildSettingsAttribute()
	}

	return attrs
}

// getChoiceChildSettingsAttribute returns child settings list attribute
func getChoiceChildSettingsAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getChoiceInstanceAttributes(),
		},
		Computed:            true,
		Description:         "Child setting configurations",
		MarkdownDescription: "Child settings for choice setting configuration",
	}
}

// getChoiceGroupInstance returns a group instance schema for choice settings
func getChoiceGroupInstance() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          ChoiceSchemaAttributeMap{},
		Description:         "Group setting instance configuration",
		MarkdownDescription: "Group setting instance (#microsoft.graph.deviceManagementConfigurationSettingGroupInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSettingGroupInstance?view=graph-rest-beta",
	}
}

// getChoiceGroupCollectionInstance returns a group collection instance schema
func getChoiceGroupCollectionInstance() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          ChoiceSchemaAttributeMap{},
		Description:         "Group setting collection instance configuration",
		MarkdownDescription: "Group setting collection instance (#microsoft.graph.deviceManagementConfigurationSettingGroupCollectionInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSettingGroupCollectionInstance?view=graph-rest-beta",
	}
}

// getChoiceSimpleInstance returns simple setting instance schema
func getChoiceSimpleInstance() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          getChoiceSimpleValueAttributes(),
				Description:         "Simple setting value configuration",
				MarkdownDescription: "Simple setting value (#microsoft.graph.deviceManagementConfigurationSimpleSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingValue?view=graph-rest-beta",
			},
		},
		Description:         "Simple setting instance configuration",
		MarkdownDescription: "Simple setting instance (#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingInstance?view=graph-rest-beta",
	}
}

// getChoiceSimpleCollectionInstance returns simple collection instance schema
func getChoiceSimpleCollectionInstance() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getChoiceSimpleValueAttributes(),
				},
				Description:         "Collection of simple setting values",
				MarkdownDescription: "Simple setting collection value (#microsoft.graph.deviceManagementConfigurationSimpleSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingValue?view=graph-rest-beta",
			},
		},
		Description:         "Simple setting collection instance configuration",
		MarkdownDescription: "Simple setting collection instance (#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingCollectionInstance?view=graph-rest-beta",
	}
}

// getChoiceCollectionInstance returns choice collection instance schema
func getChoiceCollectionInstance() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getChoiceValueAttributes(true),
				},
				Description:         "Collection of choice setting values",
				MarkdownDescription: "Choice setting collection value (#microsoft.graph.deviceManagementConfigurationChoiceSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingValue?view=graph-rest-beta",
			},
		},
		Description:         "Choice setting collection instance configuration",
		MarkdownDescription: "Choice setting collection instance (#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingCollectionInstance?view=graph-rest-beta",
	}
}

// Update getChoiceInstanceAttributes to include all types
func getChoiceInstanceAttributes() ChoiceSchemaAttributeMap {
	attrs := getChoiceBaseInstanceAttributes()

	attrs["choice"] = getChoiceSettingInstance(true)
	attrs["choice_collection"] = getChoiceCollectionInstance()
	attrs["group"] = getChoiceGroupSetting()
	attrs["group_collection"] = getChoiceGroupCollectionSetting()
	attrs["setting_group"] = getChoiceGroupInstance()
	attrs["setting_group_collection"] = getChoiceGroupCollectionInstance()
	attrs["simple"] = getChoiceSimpleInstance()
	attrs["simple_collection"] = getChoiceSimpleCollectionInstance()

	return attrs
}

// getChoiceSettingInstance returns choice setting instance schema
func getChoiceSettingInstance(includeChildren bool) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          getChoiceValueAttributes(includeChildren),
				Description:         "Choice setting value configuration",
				MarkdownDescription: "Choice setting value (#microsoft.graph.deviceManagementConfigurationChoiceSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingValue?view=graph-rest-beta",
			},
		},
		Description:         "Choice setting instance configuration",
		MarkdownDescription: "Choice setting instance (#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingInstance?view=graph-rest-beta",
	}
}

// getChoiceGroupSettingAttributes returns group setting attributes
func getChoiceGroupSettingAttributes() ChoiceSchemaAttributeMap {
	return ChoiceSchemaAttributeMap{
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          getChoiceTemplateReferenceAttributes(),
			Description:         "Template reference for group setting",
			MarkdownDescription: "Setting value template reference (#microsoft.graph.deviceManagementConfigurationSettingValueTemplateReference) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSettingValueTemplateReference?view=graph-rest-beta",
		},
		"children": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: getChoiceInstanceAttributes(),
			},
			Computed:            true,
			Description:         "Collection of child setting instances contained within this GroupSetting",
			MarkdownDescription: "Collection of child setting instances contained within this GroupSetting",
		},
	}
}

// getChoiceGroupSetting returns group setting schema
func getChoiceGroupSetting() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          getChoiceGroupSettingAttributes(),
				Description:         "Group setting configuration",
				MarkdownDescription: "GroupSetting value (#microsoft.graph.deviceManagementConfigurationGroupSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingValue?view=graph-rest-beta",
			},
		},
		Description:         "Instance of a GroupSetting",
		MarkdownDescription: "Instance of a GroupSetting (#microsoft.graph.deviceManagementConfigurationGroupSettingInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingInstance?view=graph-rest-beta",
	}
}

// getChoiceGroupCollectionSetting returns group collection setting schema
func getChoiceGroupCollectionSetting() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getChoiceGroupSettingAttributes(),
				},
				Description:         "A collection of GroupSetting values",
				MarkdownDescription: "A collection of GroupSetting values (#microsoft.graph.deviceManagementConfigurationGroupSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingValue?view=graph-rest-beta",
			},
		},
		Description:         "Instance of a GroupSettingCollection",
		MarkdownDescription: "Instance of a GroupSettingCollection (#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingCollectionInstance?view=graph-rest-beta",
	}
}

// getChoiceSimpleValueAttributes returns attributes for simple value settings within choice context
func getChoiceSimpleValueAttributes() ChoiceSchemaAttributeMap {
	return ChoiceSchemaAttributeMap{
		"integer": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: ChoiceSchemaAttributeMap{
				"value": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "Value of the integer setting.",
				},
			},
			Description:         "Integer setting value configuration",
			MarkdownDescription: "Simple setting value (#microsoft.graph.deviceManagementConfigurationIntegerSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationIntegerSettingValue?view=graph-rest-beta",
		},
		"reference": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: ChoiceSchemaAttributeMap{
				"value": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Value of the string setting.",
				},
				"note": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "A note that admin can use to put some contextual information",
				},
			},
			Description:         "Reference setting value configuration",
			MarkdownDescription: "Model for ReferenceSettingValue (#microsoft.graph.deviceManagementConfigurationReferenceSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationReferenceSettingValue?view=graph-rest-beta",
		},
		"secret": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: ChoiceSchemaAttributeMap{
				"value": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Value of the secret setting.",
				},
				"state": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("invalid", "notEncrypted", "encryptedValueToken"),
					},
					Description:         `valueState`,
					MarkdownDescription: "Gets or sets a value indicating the encryption state of the Value property. / type tracking the encryption state of a secret setting value; possible values are: `invalid` (default invalid value), `notEncrypted` (secret value is not encrypted), `encryptedValueToken` (a token for the encrypted value is returned by the service)",
				},
			},
			Description:         "Secret setting value configuration",
			MarkdownDescription: "Graph model for a secret setting value (#microsoft.graph.deviceManagementConfigurationSecretSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSecretSettingValue?view=graph-rest-beta",
		},
		"string": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: ChoiceSchemaAttributeMap{
				"value": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Value of the string setting.",
				},
			},
			Description:         "String setting value configuration",
			MarkdownDescription: "Simple setting value (#microsoft.graph.deviceManagementConfigurationStringSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationStringSettingValue?view=graph-rest-beta",
		},
	}
}
