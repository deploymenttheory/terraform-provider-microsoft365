package graphBetaSettingsCatalogConfigurationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const maxSchemaDepth = 15

// DeviceConfigV2Schema returns the Terraform schema for device configuration settings catalog
func DeviceConfigV2Schema() schema.Schema {
	return schema.Schema{
		Description: "Device Management Configuration Setting (Settings Catalog) resource schema",
		Attributes:  DeviceConfigV2Attributes(),
	}
}

// DeviceConfigV2Attributes returns the attributes map for device configuration settings catalog
// Use this function when you need the attributes map directly for nested schemas
func DeviceConfigV2Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"settings": schema.ListNestedAttribute{
			Description: "Array-based settings collection",
			Required:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: settingAttributes(0),
			},
		},
	}
}

// settingAttributes defines the attributes for a Setting
func settingAttributes(depth int) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Key of this setting within the policy which contains it. Automatically generated.",
			Optional:    true,
			Computed:    true,
		},
		"setting_instance": schema.SingleNestedAttribute{
			Description: "singular settings catalog instance",
			Required:    true,
			Attributes:  settingInstanceAttributes(depth),
		},
	}
}

// settingInstanceAttributes defines the attributes for SettingInstance
func settingInstanceAttributes(depth int) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Description: "OData type identifier",
			Required:    true,
		},
		"setting_definition_id": schema.StringAttribute{
			Description: "Setting definition identifier",
			Required:    true,
		},
		"setting_instance_template_reference": schema.SingleNestedAttribute{
			Description: "Setting Instance Template Reference at instance level",
			Optional:    true,
			Attributes:  settingInstanceTemplateReferenceAttributes(),
		},
		"simple_setting_value": schema.SingleNestedAttribute{
			Description: "Simple setting value",
			Optional:    true,
			Attributes:  simpleSettingAttributes(),
		},
		"simple_setting_collection_value": schema.ListNestedAttribute{
			Description: "Collection of simple settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: simpleSettingCollectionAttributes(),
			},
		},
		"choice_setting_value": schema.SingleNestedAttribute{
			Description: "Choice setting value",
			Optional:    true,
			Attributes:  choiceSettingAttributes(depth),
		},
		"choice_setting_collection_value": schema.ListNestedAttribute{
			Description: "Collection of choice settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: choiceSettingCollectionAttributes(depth),
			},
		},
		"group_setting_collection_value": schema.ListNestedAttribute{
			Description: "Collection of group settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: groupSettingCollectionAttributes(depth),
			},
		},
	}
}

// simpleSettingAttributes defines the attributes for SimpleSettingStruct
func simpleSettingAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Description: "OData type identifier",
			Required:    true,
		},
		"setting_value_template_reference": schema.SingleNestedAttribute{
			Description: "Template reference at value level",
			Optional:    true,
			Attributes:  settingValueTemplateReferenceAttributes(),
		},
		"value": schema.StringAttribute{
			Description: "Setting value (can be string, integer, or secret - all represented as strings)",
			Required:    true,
		},
		"value_state": schema.StringAttribute{
			Description: "Value state. This is only used for secret settings and for a valid request, it must always be set to 'notEncrypted'",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("notEncrypted"),
			},
		},
	}
}

// simpleSettingCollectionAttributes defines the attributes for SimpleSettingCollectionStruct
func simpleSettingCollectionAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Description: "OData type identifier",
			Required:    true,
		},
		"setting_value_template_reference": schema.SingleNestedAttribute{
			Description: "Template reference at value level",
			Optional:    true,
			Attributes:  settingValueTemplateReferenceAttributes(),
		},
		"value": schema.StringAttribute{
			Description: "Setting value",
			Required:    true,
		},
	}
}

// choiceSettingAttributes defines the attributes for ChoiceSettingStruct
func choiceSettingAttributes(depth int) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"setting_value_template_reference": schema.SingleNestedAttribute{
			Description: "Template reference at value level",
			Optional:    true,
			Attributes:  settingValueTemplateReferenceAttributes(),
		},
		"value": schema.StringAttribute{
			Description: "Choice value",
			Required:    true,
		},
	}

	// Only add children if we haven't exceeded max depth
	if depth < maxSchemaDepth {
		attrs["children"] = schema.ListNestedAttribute{
			Description: "Child elements of the choice setting",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: choiceSettingChildAttributes(depth + 1),
			},
		}
	}

	return attrs
}

// choiceSettingChildAttributes defines the attributes for ChoiceSettingChild
func choiceSettingChildAttributes(depth int) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Description: "OData type identifier",
			Required:    true,
		},
		"setting_definition_id": schema.StringAttribute{
			Description: "Setting definition identifier",
			Required:    true,
		},
		"setting_instance_template_reference": schema.SingleNestedAttribute{
			Description: "Template reference at instance level",
			Optional:    true,
			Attributes:  settingInstanceTemplateReferenceAttributes(),
		},
		"simple_setting_value": schema.SingleNestedAttribute{
			Description: "Simple setting value",
			Optional:    true,
			Attributes:  simpleSettingAttributes(),
		},
		"simple_setting_collection_value": schema.ListNestedAttribute{
			Description: "Collection of simple settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: simpleSettingCollectionAttributes(),
			},
		},
	}

	// Only add recursive choice and group settings if we haven't exceeded max depth
	if depth < maxSchemaDepth {
		attrs["choice_setting_value"] = schema.SingleNestedAttribute{
			Description: "Nested choice setting value",
			Optional:    true,
			Attributes:  choiceSettingAttributes(depth + 1),
		}
		attrs["choice_setting_collection_value"] = schema.ListNestedAttribute{
			Description: "Collection of nested choice settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: choiceSettingCollectionAttributes(depth + 1),
			},
		}
		attrs["group_setting_collection_value"] = schema.ListNestedAttribute{
			Description: "Collection of group settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: groupSettingCollectionAttributes(depth + 1),
			},
		}
	}

	return attrs
}

// choiceSettingCollectionAttributes defines the attributes for ChoiceSettingCollectionStruct
func choiceSettingCollectionAttributes(depth int) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"setting_value_template_reference": schema.SingleNestedAttribute{
			Description: "Template reference at value level",
			Optional:    true,
			Attributes:  settingValueTemplateReferenceAttributes(),
		},
		"value": schema.StringAttribute{
			Description: "Choice collection value",
			Required:    true,
		},
	}

	// Only add children if we haven't exceeded max depth
	if depth < maxSchemaDepth {
		attrs["children"] = schema.ListNestedAttribute{
			Description: "Child elements of the choice setting collection",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: choiceSettingCollectionChildAttributes(depth + 1),
			},
		}
	}

	return attrs
}

// choiceSettingCollectionChildAttributes defines the attributes for ChoiceSettingCollectionChild
func choiceSettingCollectionChildAttributes(depth int) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Description: "OData type identifier",
			Required:    true,
		},
		"setting_definition_id": schema.StringAttribute{
			Description: "Setting definition identifier",
			Required:    true,
		},
		"setting_instance_template_reference": schema.SingleNestedAttribute{
			Description: "Template reference at instance level",
			Optional:    true,
			Attributes:  settingInstanceTemplateReferenceAttributes(),
		},
		"simple_setting_value": schema.SingleNestedAttribute{
			Description: "Simple setting value",
			Optional:    true,
			Attributes:  simpleSettingAttributes(),
		},
		"simple_setting_collection_value": schema.ListNestedAttribute{
			Description: "Collection of simple settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: simpleSettingCollectionAttributes(),
			},
		},
	}
}

// groupSettingCollectionAttributes defines the attributes for GroupSettingCollectionStruct
func groupSettingCollectionAttributes(depth int) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"setting_value_template_reference": schema.SingleNestedAttribute{
			Description: "Template reference at value level",
			Optional:    true,
			Attributes:  settingValueTemplateReferenceAttributes(),
		},
	}

	// Only add children if we haven't exceeded max depth
	if depth < maxSchemaDepth {
		attrs["children"] = schema.ListNestedAttribute{
			Description: "Child elements of the group setting collection",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: groupSettingCollectionChildAttributes(depth + 1),
			},
		}
	}

	return attrs
}

// groupSettingCollectionChildAttributes defines the attributes for GroupSettingCollectionChild
func groupSettingCollectionChildAttributes(depth int) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Description: "OData type identifier",
			Required:    true,
		},
		"setting_definition_id": schema.StringAttribute{
			Description: "Setting definition identifier",
			Required:    true,
		},
		"setting_instance_template_reference": schema.SingleNestedAttribute{
			Description: "Template reference at instance level",
			Optional:    true,
			Attributes:  settingInstanceTemplateReferenceAttributes(),
		},
		"simple_setting_value": schema.SingleNestedAttribute{
			Description: "Simple setting value",
			Optional:    true,
			Attributes:  simpleSettingAttributes(),
		},
		"simple_setting_collection_value": schema.ListNestedAttribute{
			Description: "Collection of simple settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: simpleSettingCollectionAttributes(),
			},
		},
	}

	// Only add recursive choice and group settings if we haven't exceeded max depth
	if depth < maxSchemaDepth {
		attrs["choice_setting_value"] = schema.SingleNestedAttribute{
			Description: "Choice setting value",
			Optional:    true,
			Attributes:  choiceSettingAttributes(depth + 1),
		}
		attrs["choice_setting_collection_value"] = schema.ListNestedAttribute{
			Description: "Collection of choice settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: choiceSettingCollectionAttributes(depth + 1),
			},
		}
		attrs["group_setting_collection_value"] = schema.ListNestedAttribute{
			Description: "Collection of nested group settings",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: groupSettingCollectionAttributes(depth + 1),
			},
		}
	}

	return attrs
}

// settingInstanceTemplateReferenceAttributes defines the attributes for SettingInstanceTemplateReference
func settingInstanceTemplateReferenceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"setting_instance_template_id": schema.StringAttribute{
			Description: "Setting instance template identifier",
			Required:    true,
		},
	}
}

// settingValueTemplateReferenceAttributes defines the attributes for SettingValueTemplateReference
func settingValueTemplateReferenceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"setting_value_template_id": schema.StringAttribute{
			Description: "Setting value template identifier",
			Required:    true,
		},
		"use_template_default": schema.BoolAttribute{
			Description: "Whether to use template default value",
			Required:    true,
		},
	}
}
