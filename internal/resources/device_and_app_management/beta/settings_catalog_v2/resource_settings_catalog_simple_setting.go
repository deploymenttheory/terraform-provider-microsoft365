package graphBetaSettingsCatalog

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

// SimpleSchemaAttributeMap defines the common type for schema attribute maps
type SimpleSchemaAttributeMap map[string]schema.Attribute

// GetSimpleSchema returns the root schema for simple settings data type
func GetSimpleSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: SimpleSchemaAttributeMap{
			"value": schema.SingleNestedAttribute{
				Required:            true,
				Attributes:          deviceManagementConfigurationSimpleSettingValueAttributes,
				Description:         "simpleSettingValue",
				MarkdownDescription: "Simple setting instance value / Simple setting value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingValue?view=graph-rest-beta",
			},
		},
		Description:         "Simple setting instance configuration",
		MarkdownDescription: "Simple setting instance / ReferenceSettingValue (#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingInstance?view=graph-rest-beta",
	}
}
