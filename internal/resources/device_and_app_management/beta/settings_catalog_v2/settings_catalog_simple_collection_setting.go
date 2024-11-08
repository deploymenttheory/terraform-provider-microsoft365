package graphBetaSettingsCatalog

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

// SimpleCollectionSchemaAttributeMap defines the common type for schema attribute maps
type SimpleCollectionSchemaAttributeMap map[string]schema.Attribute

// GetSimpleCollectionSchema returns the root schema for simple setting collection data type
func GetSimpleCollectionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: SimpleCollectionSchemaAttributeMap{
			"values": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: deviceManagementConfigurationSimpleSettingValueAttributes,
				},
				Description: "simpleSettingCollectionValue",
				MarkdownDescription: "Simple setting collection instance value with @odata.type: " +
					"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance. " +
					"See the [Simple Setting Value Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingValue?view=graph-rest-beta) for more details.",
			},
		},
		Description: "Simple setting collection instance configuration",
		MarkdownDescription: "Configuration for a simple setting collection instance with @odata.type: " +
			"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance. " +
			"See the [Simple Setting Collection Instance Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingCollectionInstance?view=graph-rest-beta) for more details.",
	}
}
