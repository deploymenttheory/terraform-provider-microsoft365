package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SimpleSchemaAttributeMap defines the common type for schema attribute maps
type SimpleSchemaAttributeMap map[string]schema.Attribute

// GetSimpleSchema returns the root schema for simple settings data type
func GetSimpleSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: SimpleSchemaAttributeMap{
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
			"secret_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value of the secret string setting.",
			},
			"state": schema.StringAttribute{
				Optional: true,
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
		Description:         "Simple setting instance configuration",
		MarkdownDescription: "Simple setting instance / ReferenceSettingValue (#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingInstance?view=graph-rest-beta",
	}
}

var deviceManagementConfigurationSimpleSettingValueAttributes = map[string]schema.Attribute{
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
	"reference": schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"string_value": schema.StringAttribute{
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
			"string_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value of the secret string setting.",
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
}
