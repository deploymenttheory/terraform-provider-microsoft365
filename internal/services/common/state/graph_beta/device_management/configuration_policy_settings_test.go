package devicemanagement

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNormalizeSettingsCatalogJSONArray tests the normalizeSettingsCatalogJSONArray function
func TestNormalizeSettingsCatalogJSONArray(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		settingsStr    string
		apiResponse    map[string]any
		validateOutput func(t *testing.T, result string)
	}{
		{
			name:        "Simple choice setting",
			settingsStr: "",
			apiResponse: map[string]any{
				"value": []any{
					map[string]any{
						"id": "0",
						"settingInstance": map[string]any{
							"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
							"settingDefinitionId":              "test_setting_id",
							"settingInstanceTemplateReference": nil,
							"auditRuleInformation":             nil,
							"choiceSettingValue": map[string]any{
								"value":                         "test_value",
								"settingValueTemplateReference": nil,
								"children":                      []any{},
							},
						},
					},
				},
			},
			validateOutput: func(t *testing.T, result string) {
				require.NotEmpty(t, result)
				assert.Contains(t, result, "test_setting_id")
				assert.Contains(t, result, "test_value")
			},
		},
		{
			name:        "Choice setting with nested children",
			settingsStr: "",
			apiResponse: map[string]any{
				"value": []any{
					map[string]any{
						"id": "0",
						"settingInstance": map[string]any{
							"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
							"settingDefinitionId":              "parent_setting",
							"settingInstanceTemplateReference": nil,
							"auditRuleInformation":             nil,
							"choiceSettingValue": map[string]any{
								"value":                         "parent_value",
								"settingValueTemplateReference": nil,
								"children": []any{
									map[string]any{
										"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
										"settingDefinitionId":              "child_setting",
										"settingInstanceTemplateReference": nil,
										"choiceSettingValue": map[string]any{
											"value":                         "child_value",
											"settingValueTemplateReference": nil,
											"children":                      []any{},
										},
									},
								},
							},
						},
					},
				},
			},
			validateOutput: func(t *testing.T, result string) {
				require.NotEmpty(t, result)
				assert.Contains(t, result, "parent_setting")
				assert.Contains(t, result, "parent_value")
				assert.Contains(t, result, "child_setting")
				assert.Contains(t, result, "child_value")
			},
		},
		{
			name:        "Simple integer setting",
			settingsStr: "",
			apiResponse: map[string]any{
				"value": []any{
					map[string]any{
						"id": "0",
						"settingInstance": map[string]any{
							"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
							"settingDefinitionId":              "integer_setting_id",
							"settingInstanceTemplateReference": nil,
							"simpleSettingValue": map[string]any{
								"@odata.type":                   "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
								"value":                         30,
								"settingValueTemplateReference": nil,
							},
						},
					},
				},
			},
			validateOutput: func(t *testing.T, result string) {
				require.NotEmpty(t, result)
				assert.Contains(t, result, "integer_setting_id")
				assert.Contains(t, result, "30")
			},
		},
		{
			name:        "Simple string setting",
			settingsStr: "",
			apiResponse: map[string]any{
				"value": []any{
					map[string]any{
						"id": "0",
						"settingInstance": map[string]any{
							"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
							"settingDefinitionId":              "string_setting_id",
							"settingInstanceTemplateReference": nil,
							"simpleSettingValue": map[string]any{
								"@odata.type":                   "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
								"value":                         "test_string_value",
								"settingValueTemplateReference": nil,
							},
						},
					},
				},
			},
			validateOutput: func(t *testing.T, result string) {
				require.NotEmpty(t, result)
				assert.Contains(t, result, "string_setting_id")
				assert.Contains(t, result, "test_string_value")
			},
		},
		{
			name:        "Simple setting collection",
			settingsStr: "",
			apiResponse: map[string]any{
				"value": []any{
					map[string]any{
						"id": "0",
						"settingInstance": map[string]any{
							"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
							"settingDefinitionId":              "collection_setting_id",
							"settingInstanceTemplateReference": nil,
							"simpleSettingCollectionValue": []any{
								map[string]any{
									"@odata.type":                   "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
									"value":                         "item1",
									"settingValueTemplateReference": nil,
								},
								map[string]any{
									"@odata.type":                   "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
									"value":                         "item2",
									"settingValueTemplateReference": nil,
								},
							},
						},
					},
				},
			},
			validateOutput: func(t *testing.T, result string) {
				require.NotEmpty(t, result)
				assert.Contains(t, result, "collection_setting_id")
				assert.Contains(t, result, "item1")
				assert.Contains(t, result, "item2")
			},
		},
		{
			name:        "Group setting collection",
			settingsStr: "",
			apiResponse: map[string]any{
				"value": []any{
					map[string]any{
						"id": "0",
						"settingInstance": map[string]any{
							"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
							"settingDefinitionId":              "group_collection_id",
							"settingInstanceTemplateReference": nil,
							"groupSettingCollectionValue": []any{
								map[string]any{
									"settingValueTemplateReference": nil,
									"children": []any{
										map[string]any{
											"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
											"settingDefinitionId":              "child_string_id",
											"settingInstanceTemplateReference": nil,
											"simpleSettingValue": map[string]any{
												"@odata.type":                   "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
												"value":                         "group_child_value",
												"settingValueTemplateReference": nil,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			validateOutput: func(t *testing.T, result string) {
				require.NotEmpty(t, result)
				assert.Contains(t, result, "group_collection_id")
				assert.Contains(t, result, "child_string_id")
				assert.Contains(t, result, "group_child_value")
			},
		},
		{
			name:        "Multi-level nested group collection",
			settingsStr: "",
			apiResponse: map[string]any{
				"value": []any{
					map[string]any{
						"id": "0",
						"settingInstance": map[string]any{
							"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
							"settingDefinitionId":              "outer_group",
							"settingInstanceTemplateReference": nil,
							"groupSettingCollectionValue": []any{
								map[string]any{
									"settingValueTemplateReference": nil,
									"children": []any{
										map[string]any{
											"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
											"settingDefinitionId":              "inner_group",
											"settingInstanceTemplateReference": nil,
											"groupSettingCollectionValue": []any{
												map[string]any{
													"settingValueTemplateReference": nil,
													"children": []any{
														map[string]any{
															"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
															"settingDefinitionId":              "deepest_choice",
															"settingInstanceTemplateReference": nil,
															"choiceSettingValue": map[string]any{
																"value":                         "deepest_value",
																"settingValueTemplateReference": nil,
																"children":                      []any{},
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
			},
			validateOutput: func(t *testing.T, result string) {
				require.NotEmpty(t, result)
				assert.Contains(t, result, "outer_group")
				assert.Contains(t, result, "inner_group")
				assert.Contains(t, result, "deepest_choice")
				assert.Contains(t, result, "deepest_value")
			},
		},
		{
			name:        "Empty settings array",
			settingsStr: "",
			apiResponse: map[string]any{
				"value": []any{},
			},
			validateOutput: func(t *testing.T, result string) {
				assert.Contains(t, result, `"settings":[]`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the API response to bytes
			resp, err := json.Marshal(tt.apiResponse)
			require.NoError(t, err)

			result := normalizeSettingsCatalogJSONArray(ctx, tt.settingsStr, resp)

			if tt.validateOutput != nil {
				tt.validateOutput(t, result)
			}
		})
	}
}

// TestNormalizeSettingsCatalogJSONArray_SecretHandling tests secret preservation logic
func TestNormalizeSettingsCatalogJSONArray_SecretHandling(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		settingsStr string
		apiResponse map[string]any
		validateFn  func(t *testing.T, result string)
	}{
		{
			name:        "Secret setting during import (no original config)",
			settingsStr: "",
			apiResponse: map[string]any{
				"value": []any{
					map[string]any{
						"id": "0",
						"settingInstance": map[string]any{
							"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
							"settingDefinitionId":              "test_secret_setting",
							"settingInstanceTemplateReference": nil,
							"auditRuleInformation":             nil,
							"simpleSettingValue": map[string]any{
								"@odata.type":                   "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
								"value":                         "",
								"valueState":                    "notEncrypted",
								"settingValueTemplateReference": nil,
							},
						},
					},
				},
			},
			validateFn: func(t *testing.T, result string) {
				assert.Contains(t, result, "test_secret_setting")
				assert.Contains(t, result, "valueState")
				assert.Contains(t, result, "notEncrypted")
			},
		},
		{
			name: "Secret setting with original config (update scenario)",
			settingsStr: `{
				"settings": [{
					"id": "0",
					"settingInstance": {
						"@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
						"settingDefinitionId": "test_secret_setting",
						"simpleSettingValue": {
							"@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
							"value": "original_secret_value",
							"valueState": "notEncrypted"
						}
					}
				}]
			}`,
			apiResponse: map[string]any{
				"value": []any{
					map[string]any{
						"id": "0",
						"settingInstance": map[string]any{
							"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
							"settingDefinitionId":              "test_secret_setting",
							"settingInstanceTemplateReference": nil,
							"simpleSettingValue": map[string]any{
								"@odata.type":                   "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
								"value":                         "",
								"valueState":                    "encrypted",
								"settingValueTemplateReference": nil,
							},
						},
					},
				},
			},
			validateFn: func(t *testing.T, result string) {
				assert.Contains(t, result, "test_secret_setting")
				// Should preserve original value from config
				assert.Contains(t, result, "original_secret_value")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := json.Marshal(tt.apiResponse)
			require.NoError(t, err)

			result := normalizeSettingsCatalogJSONArray(ctx, tt.settingsStr, resp)
			if tt.validateFn != nil {
				tt.validateFn(t, result)
			}
		})
	}
}

// TestNormalizeSettingsCatalogJSONArray_UpdateScenario tests state preservation during updates
func TestNormalizeSettingsCatalogJSONArray_UpdateScenario(t *testing.T) {
	ctx := context.Background()

	originalConfig := `{
		"settings": [{
			"id": "0",
			"settingInstance": {
				"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
				"settingDefinitionId": "test_setting",
				"choiceSettingValue": {
					"value": "original_value",
					"children": []
				}
			}
		}]
	}`

	apiResponse := map[string]any{
		"value": []any{
			map[string]any{
				"id": "0",
				"settingInstance": map[string]any{
					"@odata.type":                      "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
					"settingDefinitionId":              "test_setting",
					"settingInstanceTemplateReference": nil,
					"choiceSettingValue": map[string]any{
						"value":                         "updated_value",
						"settingValueTemplateReference": nil,
						"children":                      []any{},
					},
				},
			},
		},
	}

	resp, err := json.Marshal(apiResponse)
	require.NoError(t, err)

	result := normalizeSettingsCatalogJSONArray(ctx, originalConfig, resp)

	// Should reflect the updated value from API
	assert.Contains(t, result, "test_setting")
	assert.Contains(t, result, "updated_value")
}
