package graphBetaSettingsCatalog

import (
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const MaxDepth = 5

// ChoiceSchemaAttributeMap defines the common type for schema attribute maps
type ChoiceSchemaAttributeMap map[string]schema.Attribute

// GetChoiceSchema returns the root schema for choice data type
func GetChoiceSchema(currentDepth int) ChoiceSchemaAttributeMap {
	if currentDepth >= MaxDepth {
		return ChoiceSchemaAttributeMap{}
	}

	attrs := ChoiceSchemaAttributeMap{
		"odata_type": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The OData type of the setting instance. This is automatically set by the graph SDK during request construction.",
		},
		"string_value": schema.StringAttribute{
			Optional:    true,
			Description: "String setting configuration",
			MarkdownDescription: "Simple string setting value. With @odata.type: #microsoft.graph.deviceManagementConfigurationStringSettingValue.\n\n" +
				"For more details, see [String Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
				"api/resources/intune-deviceconfigv2-deviceManagementConfigurationStringSettingValue?view=graph-rest-beta).",
		},
		"integer_value": schema.Int32Attribute{
			Optional:            true,
			MarkdownDescription: "Simple integer setting value (#microsoft.graph.deviceManagementConfigurationIntegerSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationIntegerSettingValue?view=graph-rest-beta",
		},
	}

	if currentDepth < MaxDepth {
		// attrs["children"] = schema.SingleNestedAttribute{
		// 	Optional:            true,
		// 	Attributes:          GetChildrenAttributes(currentDepth + 1),
		// 	Description:         "Child setting configuration",
		// 	MarkdownDescription: "Child setting instance under choice setting configuration.",
		// }
		attrs["children"] = schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: GetChildrenAttributes(currentDepth + 1),
			},
			Description:         "List of child setting configurations",
			MarkdownDescription: "List of child setting instances under choice setting configuration.",
		}
	}

	return attrs
}

// GetChildrenAttributes returns the schema attributes for nested child settings
func GetChildrenAttributes(currentDepth int) ChoiceSchemaAttributeMap {
	attrs := ChoiceSchemaAttributeMap{
		"odata_type": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The OData type of the child setting instance. This must be specified and is used to determine the specific setting instance type.",
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
		"setting_definition_id": schema.StringAttribute{
			Required:            true,
			Description:         `settingDefinitionId`,
			MarkdownDescription: "Setting Definition Id",
		},
	}

	if currentDepth >= MaxDepth {
		return attrs
	}

	attrs["choice"] = schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          GetChoiceSchema(currentDepth + 1),
				Description:         "Choice setting value configuration",
				MarkdownDescription: "Configuration of the value for choice setting.",
			},
		},
	}

	attrs["choice_collection"] = schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"string_value": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Collection of string choice values",
				MarkdownDescription: "List of string-based choice setting values.",
			},
			"int_value": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.Int32Type,
				Description:         "Collection of integer choice values",
				MarkdownDescription: "List of integer-based choice setting values.",
			},
			"children": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: GetChoiceSchema(currentDepth + 1),
				},
				Description:         "Child settings for each choice value",
				MarkdownDescription: "Nested settings configuration for choice values.",
			},
		},
	}

	attrs["group"] = schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"children": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"odata_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The OData type of the child setting instance. This must be specified and is used to determine the specific setting instance type.",
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
						"setting_definition_id": schema.StringAttribute{
							Required:            true,
							Description:         "Setting definition ID",
							MarkdownDescription: "The unique identifier for the setting definition.",
						},
						// Child type-specific attributes
						"choice_value": schema.SingleNestedAttribute{
							Optional:   true,
							Attributes: GetChoiceSchema(currentDepth + 1),
						},
						"simple_value": schema.SingleNestedAttribute{
							Optional:   true,
							Attributes: deviceManagementConfigurationSimpleSettingValueAttributes,
						},
						"group_value": schema.SingleNestedAttribute{
							Optional:   true,
							Attributes: getChoiceGroupSettingAttributes(currentDepth + 1),
						},
						"choice_collection_value": schema.SingleNestedAttribute{
							Optional:   true,
							Attributes: GetChoiceCollectionSchema(currentDepth + 1),
						},
						// "simple_collection_value": schema.SingleNestedAttribute{
						// 	Optional:   true,
						// 	Attributes: GetSimpleCollectionSchema(currentDepth + 1),
						// },
						// "group_collection_value": schema.SingleNestedAttribute{
						// 	Optional:   true,
						// 	Attributes: GetGroupCollectionSchema(currentDepth + 1),
						// },
					},
				},
				Description:         "Child settings of various types that can be nested within this group",
				MarkdownDescription: "Collection of nested device management configuration settings of different types.",
			},
		},
	}

	attrs["group_collection"] = schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"children": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"odata_type": schema.StringAttribute{
							Required:            true,
							Description:         "The OData type of the child setting instance",
							MarkdownDescription: "Specifies the type of device management configuration setting instance.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
									"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
									"#microsoft.graph.deviceManagementConfigurationGroupSettingInstance",
									"#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance",
									"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
								),
							},
						},
						"setting_definition_id": schema.StringAttribute{
							Required:            true,
							Description:         "Setting definition ID",
							MarkdownDescription: "The unique identifier for the setting definition.",
						},
						// Child type-specific attributes
						"choice_value": schema.SingleNestedAttribute{
							Optional:   true,
							Attributes: GetChoiceSchema(currentDepth + 1),
						},
						"simple_value": schema.SingleNestedAttribute{
							Optional:   true,
							Attributes: deviceManagementConfigurationSimpleSettingValueAttributes,
						},
						"group_value": schema.SingleNestedAttribute{
							Optional:   true,
							Attributes: getChoiceGroupSettingAttributes(currentDepth + 1),
						},
						"choice_collection_value": schema.SingleNestedAttribute{
							Optional:   true,
							Attributes: GetChoiceCollectionSchema(currentDepth + 1),
						},
						// "simple_collection_value": schema.SingleNestedAttribute{
						// 	Optional:   true,
						// 	Attributes: GetSimpleCollectionSchema(currentDepth + 1),
						// },
						// "group_collection_value": schema.SingleNestedAttribute{
						// 	Optional:   true,
						// 	Attributes: GetGroupCollectionSchema(currentDepth + 1),
						// },
					},
				},
				Description:         "Child settings that make up this group collection",
				MarkdownDescription: "Collection of nested device management configuration settings that will be wrapped in group setting values.",
			},
		},
	}

	attrs["setting_group"] = schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          map[string]schema.Attribute{},
		Description:         "Setting group configuration",
		MarkdownDescription: "Configuration for setting group.",
	}

	attrs["setting_group_collection"] = schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getChoiceGroupSettingAttributes(currentDepth + 1),
				},
				Description:         "Collection of setting group values",
				MarkdownDescription: "List of setting group value configurations.",
			},
		},
	}

	attrs["simple"] = schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          GetSimpleSchema().Attributes,
		Description:         "Simple setting value configuration",
		MarkdownDescription: "Configuration for simple setting value.",
	}

	attrs["simple_collection"] = schema.SingleNestedAttribute{
		Optional: true,
		Attributes: ChoiceSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: GetChoiceCollectionSchema(currentDepth + 1),
				},
				Description:         "Collection of simple setting values",
				MarkdownDescription: "List of simple setting value configurations.",
			},
		},
	}

	return attrs
}

// getChoiceGroupSettingAttributes returns group setting attributes
func getChoiceGroupSettingAttributes(currentDepth int) ChoiceSchemaAttributeMap {
	if currentDepth >= MaxDepth {
		return ChoiceSchemaAttributeMap{}
	}

	return ChoiceSchemaAttributeMap{
		"children": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: GetChildrenAttributes(currentDepth + 1),
			},
			Computed: true,
			PlanModifiers: []planmodifier.List{
				planmodifiers.DefaultListEmptyValue(),
			},
			Description:         "List of child setting instances within this GroupSetting.",
			MarkdownDescription: "Collection of child settings within a GroupSetting instance, representing grouped nested configurations.",
		},
	}
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
