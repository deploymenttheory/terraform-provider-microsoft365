package graphBetaSettingsCatalog

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

// SettingGroupSchemaAttributeMap defines the common type for schema attribute maps
type SettingGroupSchemaAttributeMap map[string]schema.Attribute

// GetSettingGroupSchema returns the root schema for setting group instance data type
func GetSettingGroupSchema(currentDepth int) schema.SingleNestedAttribute {
	if currentDepth >= MaxDepth {
		return schema.SingleNestedAttribute{}
	}

	return schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: SettingGroupSchemaAttributeMap{
			// TODO. Define specific attributes here if needed in the future
		},
		Description: "Setting group instance configuration within a policy",
		MarkdownDescription: "Configuration for a setting group instance with @odata.type: " +
			"#microsoft.graph.deviceManagementConfigurationSettingGroupInstance. " +
			"See the [Setting Group Instance Documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationSettingGroupInstance?view=graph-rest-beta) for more details.",
	}
}
