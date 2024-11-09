package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SimpleCollectionSchemaAttributeMap defines the common type for schema attribute maps
type SimpleCollectionSchemaAttributeMap map[string]schema.Attribute

// GetSimpleCollectionSchema returns the root schema for simple setting collection data type
func GetSimpleCollectionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: SimpleCollectionSchemaAttributeMap{
			"odata_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The OData type of the setting instance. This is automatically set by the graph SDK during request construction.",
			},
			"string_value": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "String setting configuration",
				MarkdownDescription: "List of simple string setting values.With @odata.type: #microsoft.graph.deviceManagementConfigurationStringSettingValue.\n\n" +
					"For more details, see [String Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
					"api/resources/intune-deviceconfigv2-deviceManagementConfigurationStringSettingValue?view=graph-rest-beta).",
			},
			"integer_value": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.Int32Type,
				MarkdownDescription: "Simple integer setting value (#microsoft.graph.deviceManagementConfigurationIntegerSettingValue) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationIntegerSettingValue?view=graph-rest-beta",
			},
			"secret_value": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of secret string setting values.",
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
		Description: "Simple setting collection instance configuration",
		MarkdownDescription: "Configuration for a simple setting collection instance with @odata.type: " +
			"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance. " +
			"See the [Simple Setting Collection Instance Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingCollectionInstance?view=graph-rest-beta) for more details.",
	}
}
