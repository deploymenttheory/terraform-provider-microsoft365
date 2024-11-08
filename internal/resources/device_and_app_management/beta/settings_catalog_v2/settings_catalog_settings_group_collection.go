package graphBetaSettingsCatalog

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

// SettingGroupCollectionSchemaAttributeMap defines the common type for schema attribute maps
type SettingGroupCollectionSchemaAttributeMap map[string]schema.Attribute

// GetSettingGroupCollectionSchema returns the root schema for setting group collection data type
func GetSettingGroupCollectionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: SettingGroupCollectionSchemaAttributeMap{
			// TODO. Define specific attributes here if needed in the future
		},
		Description: "Setting group collection instance configuration",
		MarkdownDescription: "Configuration for a setting group collection instance with @odata.type: " +
			"#microsoft.graph.deviceManagementConfigurationSettingGroupCollectionInstance. " +
			"See the [Setting Group Collection Instance Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSettingGroupCollectionInstance?view=graph-rest-beta) for more details.",
	}
}
