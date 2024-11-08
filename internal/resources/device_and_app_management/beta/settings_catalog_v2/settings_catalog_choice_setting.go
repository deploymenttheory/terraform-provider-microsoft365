package graphBetaSettingsCatalog

import (
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const MaxDepth = 5

// ChoiceSchemaAttributeMap defines the common type for schema attribute maps
type ChoiceSchemaAttributeMap map[string]schema.Attribute

// GetChoiceSchema returns the root schema for choice data type
func GetChoiceSchema(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:    true,
				Attributes:  getChoiceValueAttributes(true, currentDepth+1),
				Description: "Choice setting value configuration",
				MarkdownDescription: "Choice setting value (#microsoft.graph.deviceManagementConfigurationChoiceSettingValue) / " +
					"For more details, see [ChoiceSettingValue Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingValue?view=graph-rest-beta).",
			},
		},
		Description: "Choice setting instance configuration",
		MarkdownDescription: "Instance configuration for choice setting (#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance) / " +
			"For more details, see [ChoiceSettingInstance Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingInstance?view=graph-rest-beta).",
	}
}

// getChoiceValueAttributes returns choice value attributes
func getChoiceValueAttributes(includeChildren bool, currentDepth int) ChoiceSchemaAttributeMap {
	attrs := ChoiceSchemaAttributeMap{
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationSettingValueTemplateReferenceAttributes,
			Description:         "Template reference for setting value",
			MarkdownDescription: "Template reference for choice setting value information in Microsoft Graph.",
		},
		"value": schema.StringAttribute{
			Optional:            true,
			Description:         "Identifier for choice setting value",
			MarkdownDescription: "Specifies the unique identifier for choice setting value.",
		},
	}

	if includeChildren && currentDepth < MaxDepth {
		attrs["children"] = getChoiceChildSettingsAttribute(currentDepth + 1)
	}

	return attrs
}

// getChoiceChildSettingsAttribute returns child settings list attribute
func getChoiceChildSettingsAttribute(currentDepth int) schema.ListNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.ListNestedAttribute{}
	}

	return schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getChoiceInstanceAttributes(currentDepth),
		},
		Computed:            true,
		Description:         "List of child setting configurations",
		MarkdownDescription: "List of child setting instances under choice setting configuration.",
	}
}

// getChoiceGroupInstance returns a group instance schema for choice settings
func getChoiceGroupInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          map[string]schema.Attribute{},
		Description:         "Group setting instance configuration within a policy.",
		MarkdownDescription: "Configuration of a GroupSetting instance, representing a group within the policy structure.",
	}
}

// getChoiceGroupCollectionInstance returns a group collection instance schema
func getChoiceGroupCollectionInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getChoiceGroupSettingAttributes(currentDepth + 1),
				},
				Description:         "Collection of GroupSetting values.",
				MarkdownDescription: "Collection of values within a GroupSettingCollection instance in Microsoft Graph.",
			},
		},
		Description:         "Configuration for a GroupSetting collection instance.",
		MarkdownDescription: "Instance configuration for a GroupSetting collection, representing multiple group settings.",
	}
}

// getChoiceSimpleInstance returns simple setting instance schema
func getChoiceSimpleInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          deviceManagementConfigurationSimpleSettingValueAttributes,
				Description:         "Configuration of simple setting value.",
				MarkdownDescription: "Value configuration for a simple setting instance in Microsoft Graph.",
			},
		},
		Description:         "Simple setting instance configuration.",
		MarkdownDescription: "Instance configuration for a simple setting, representing non-complex values.",
	}
}

// getChoiceSimpleCollectionInstance returns simple collection instance schema
func getChoiceSimpleCollectionInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: deviceManagementConfigurationSimpleSettingValueAttributes,
				},
				Description:         "Configuration of simple setting collection values.",
				MarkdownDescription: "List of values within a SimpleSettingCollection instance, each representing simple settings.",
			},
		},
		Description:         "Simple setting collection instance configuration.",
		MarkdownDescription: "Configuration for an instance of a simple setting collection, containing multiple simple settings.",
	}
}

// getChoiceCollectionInstance returns choice collection instance schema
func getChoiceCollectionInstance(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getChoiceValueAttributes(true, currentDepth+1),
				},
				Description:         "Configuration for collection of choice setting values.",
				MarkdownDescription: "Instance configuration for a collection of choice setting values in Microsoft Graph.",
			},
		},
		Description:         "Choice setting collection instance configuration.",
		MarkdownDescription: "Instance configuration for a collection of choice settings, representing multiple choice values.",
	}
}

// Update getChoiceInstanceAttributes to include all types
func getChoiceInstanceAttributes(currentDepth int) ChoiceSchemaAttributeMap {
	if currentDepth >= MaxDepth {
		return getChoiceBaseInstanceAttributes()
	}

	attrs := getChoiceBaseInstanceAttributes()
	attrs["choice"] = getChoiceSettingInstance(true, currentDepth+1)
	attrs["choice_collection"] = getChoiceCollectionInstance(currentDepth + 1)
	attrs["group"] = getChoiceGroupSetting(currentDepth + 1)
	attrs["group_collection"] = getChoiceGroupCollectionSetting(currentDepth + 1)
	attrs["setting_group"] = getChoiceGroupInstance(currentDepth + 1)
	attrs["setting_group_collection"] = getChoiceGroupCollectionInstance(currentDepth + 1)
	attrs["simple"] = getChoiceSimpleInstance(currentDepth + 1)
	attrs["simple_collection"] = getChoiceSimpleCollectionInstance(currentDepth + 1)

	return attrs
}

func getChoiceBaseInstanceAttributes() ChoiceSchemaAttributeMap {
	return ChoiceSchemaAttributeMap{
		"definition_id": schema.StringAttribute{
			Required:            true,
			Description:         `settingDefinitionId`,
			MarkdownDescription: "Setting Definition Id",
		},
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationSettingInstanceTemplateReferenceAttributes,
			Description:         `settingInstanceTemplateReference`,
			MarkdownDescription: "Setting Instance Template Reference / Setting instance template reference information / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSettingInstanceTemplateReference?view=graph-rest-beta",
		},
	}
}

// getChoiceSettingInstance returns choice setting instance schema
func getChoiceSettingInstance(includeChildren bool, currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          getChoiceValueAttributes(includeChildren, currentDepth+1),
				Description:         "Configuration for choice setting value",
				MarkdownDescription: "Configuration of the value for choice setting.",
			},
		},
		Description:         "Choice setting instance configuration",
		MarkdownDescription: "Instance configuration of choice setting in Microsoft Graph.",
	}
}

// getChoiceGroupSettingAttributes returns group setting attributes
func getChoiceGroupSettingAttributes(currentDepth int) ChoiceSchemaAttributeMap {
	if currentDepth >= MaxDepth {
		return ChoiceSchemaAttributeMap{}
	}

	return ChoiceSchemaAttributeMap{
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationSettingValueTemplateReferenceAttributes,
			Description:         "Reference template for the setting value.",
			MarkdownDescription: "Template reference within the GroupSetting, providing template-based configuration options.",
		},
		"children": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: getChoiceInstanceAttributes(currentDepth + 1),
			},
			Computed:            true,
			Description:         "List of child setting instances within this GroupSetting.",
			MarkdownDescription: "Collection of child settings within a GroupSetting instance, representing grouped nested configurations.",
		},
	}
}

// getChoiceGroupSetting returns group setting schema
func getChoiceGroupSetting(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          getChoiceGroupSettingAttributes(currentDepth + 1),
				Description:         "Group setting configuration value.",
				MarkdownDescription: "Value configuration for a GroupSetting instance, containing group-specific settings.",
			},
		},
		Description:         "Configuration for a GroupSetting instance.",
		MarkdownDescription: "Configuration for a single instance of a GroupSetting, representing a set of grouped settings.",
	}
}

// getChoiceGroupCollectionSetting returns group collection setting schema
func getChoiceGroupCollectionSetting(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getChoiceGroupSettingAttributes(currentDepth + 1),
				},
				Description: "A collection of GroupSetting values.",
				MarkdownDescription: "Collection of GroupSetting values (#microsoft.graph.deviceManagementConfigurationGroupSettingValue) / " +
					"See [GroupSettingValue Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingValue?view=graph-rest-beta).",
			},
		},
		Description: "Instance of a GroupSettingCollection.",
		MarkdownDescription: "Configuration instance of a GroupSettingCollection (#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance) / " +
			"See [GroupSettingCollectionInstance Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingCollectionInstance?view=graph-rest-beta).",
	}
}

var deviceManagementConfigurationSimpleSettingValueAttributes = map[string]schema.Attribute{
	"odata_type": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The OData type of the setting instance. This is automatically set by the graph SDK during request construction.",
	},
	"integer_value": schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"value": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Value of the integer setting.",
			},
		},
		Description:         "Integer setting configuration",
		MarkdownDescription: "Simple setting value (#microsoft.graph.deviceManagementConfigurationIntegerSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationIntegerSettingValue?view=graph-rest-beta",
	},
	"reference": schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value of the string setting.",
			},
			"note": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A note that admin can use to put some contextual information",
			},
		},
		Description:         "Reference setting configuration",
		MarkdownDescription: "Schema for ReferenceSettingValue (#microsoft.graph.deviceManagementConfigurationReferenceSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationReferenceSettingValue?view=graph-rest-beta",
	},
	"secret": schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"secret_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value of the secret setting.",
			},
			"state": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("invalid", "notEncrypted", "encryptedValueToken"),
				},
				Description: `valueState`,
				MarkdownDescription: "Indicates the encryption state of the Value property, with possible values:\n\n" +
					"- **invalid**: Default invalid value.\n" +
					"- **notEncrypted**: Secret value is not encrypted.\n" +
					"- **encryptedValueToken**: A token for the encrypted value is returned by the service.\n\n" +
					"Model for SecretSettingValue with @odata.type: " +
					"#microsoft.graph.deviceManagementConfigurationSecretSettingValue.\n\n" +
					"For more details, see [Secret Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
					"api/resources/intune-deviceconfigv2-deviceManagementConfigurationSecretSettingValue?view=graph-rest-beta).",
			},
		},
		Description:         "Secret setting configuration",
		MarkdownDescription: "Graph model for a secret setting value (#microsoft.graph.deviceManagementConfigurationSecretSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSecretSettingValue?view=graph-rest-beta",
	},
	"string_value": schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value of the string setting.",
			},
		},
		Description: "String setting configuration",
		MarkdownDescription: "Value of the string setting with @odata.type: #microsoft.graph.deviceManagementConfigurationStringSettingValue.\n\n" +
			"For more details, see [String Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
			"api/resources/intune-deviceconfigv2-deviceManagementConfigurationStringSettingValue?view=graph-rest-beta).",
	},
}

var deviceManagementConfigurationSettingInstanceTemplateReferenceAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Required:            true,
		Description:         `settingInstanceTemplateId`,
		MarkdownDescription: "Setting instance template id (#microsoft.graph.deviceManagementConfigurationSettingInstanceTemplateReference)",
	},
}

var deviceManagementConfigurationSettingValueTemplateReferenceAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Required:            true,
		Description:         `settingValueTemplateId`,
		MarkdownDescription: "Setting value template id (#microsoft.graph.deviceManagementConfigurationSettingValueTemplateReference)",
	},
	"use_default": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `useTemplateDefault`,
		MarkdownDescription: "Indicates whether to update policy setting value to match template setting default value",
	},
}
