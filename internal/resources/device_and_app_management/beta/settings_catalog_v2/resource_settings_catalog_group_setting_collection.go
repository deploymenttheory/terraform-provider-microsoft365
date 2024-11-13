package graphBetaSettingsCatalog

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// GroupCollectionSchemaAttributeMap defines the common type for schema attribute maps
type GroupCollectionSchemaAttributeMap map[string]schema.Attribute

// GetGroupSettingCollectionSchema would then use this new function
func GetGroupSettingCollectionSchema(currentDepth int) GroupCollectionSchemaAttributeMap {
	if currentDepth >= MaxDepth {
		return GroupCollectionSchemaAttributeMap{}
	}

	return GroupCollectionSchemaAttributeMap{
		"children": schema.ListNestedAttribute{
			Required:            true,
			MarkdownDescription: "List of child settings within this group collection.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"odata_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The OData type of the group setting collection setting instance. e.g #microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
					},
					"setting_definition_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The unique identifier for the setting definition. e.g diskmanagement_restrictions",
					},
					"group_setting_collection_value": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Nested group setting collection values.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"children": schema.ListNestedAttribute{
									Required:            true,
									MarkdownDescription: "List of child settings within this group collection.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"odata_type": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The OData type of the choice setting value instance. e.g #microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
											},
											"setting_definition_id": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The unique identifier for the setting definition. e.g diskmanagement_restrictions_externalstorage",
											},
											"choice_setting_value": schema.SingleNestedAttribute{
												Optional:            true,
												MarkdownDescription: "Choice setting values configuration.",
												Attributes: map[string]schema.Attribute{
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
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
