package devicemanagement

import (
	"context"
	"encoding/json"
	"testing"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// TestNormalizeSettingsCatalogJSONArray_ErrorPaths tests error handling paths
func TestNormalizeSettingsCatalogJSONArray_ErrorPaths(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		settingsStr string
		resp        []byte
		validate    func(t *testing.T, result string)
	}{
		{
			name:        "Invalid JSON response returns array format",
			settingsStr: `{"settings":[]}`,
			resp:        []byte(`invalid json`),
			validate: func(t *testing.T, result string) {
				// Should attempt to parse as array and return original on failure
				assert.Equal(t, `{"settings":[]}`, result)
			},
		},
		{
			name:        "Response as direct array (no wrapper)",
			settingsStr: "",
			resp: []byte(`[
				{
					"id": "1",
					"settingDefinitionId": "test_setting"
				}
			]`),
			validate: func(t *testing.T, result string) {
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err)
				
				// Should wrap in "settings" key
				settings, ok := parsed["settings"]
				assert.True(t, ok, "Should have settings key")
				assert.NotNil(t, settings)
			},
		},
		{
			name:        "Response with 'settings' key instead of 'value'",
			settingsStr: "",
			resp: []byte(`{
				"settings": [
					{
						"id": "1",
						"settingDefinitionId": "test_setting"
					}
				]
			}`),
			validate: func(t *testing.T, result string) {
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err)
				
				settings := parsed["settings"]
				assert.NotNil(t, settings)
			},
		},
		{
			name:        "Response with neither 'value' nor 'settings' key",
			settingsStr: "",
			resp: []byte(`{
				"id": "1",
				"settingDefinitionId": "direct_setting"
			}`),
			validate: func(t *testing.T, result string) {
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err)
				
				// Should use the whole response as settings content
				settings := parsed["settings"]
				assert.NotNil(t, settings)
			},
		},
		{
			name:        "Invalid settings string (malformed JSON)",
			settingsStr: `{invalid json}`,
			resp: []byte(`{
				"value": [
					{
						"id": "1",
						"settingDefinitionId": "test_setting"
					}
				]
			}`),
			validate: func(t *testing.T, result string) {
				// Should skip secret preservation and continue
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err, "Should still produce valid JSON")
				assert.NotNil(t, parsed["settings"])
			},
		},
		{
			name: "Complex secret preservation with invalid config",
			settingsStr: `{
				"settings": [
					{
						"id": "1",
						"simpleSettingValue": {
							"@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
							"value": "original_secret"
						}
					}
				]
			}`,
			resp: []byte(`{
				"value": [
					{
						"id": "1",
						"simpleSettingValue": {
							"@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
							"valueState": "notEncrypted"
						}
					}
				]
			}`),
			validate: func(t *testing.T, result string) {
				// Should preserve the secret value
				assert.Contains(t, result, "original_secret")
				
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err)
			},
		},
		{
			name:        "Empty response object",
			settingsStr: `{"settings":[]}`,
			resp:        []byte(`{}`),
			validate: func(t *testing.T, result string) {
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err)
				
				// Should create settings key with empty object
				settings := parsed["settings"]
				assert.NotNil(t, settings)
			},
		},
		{
			name:        "Null settings in original config",
			settingsStr: "",
			resp: []byte(`{
				"value": [
					{
						"id": "1",
						"settingDefinitionId": "test"
					}
				]
			}`),
			validate: func(t *testing.T, result string) {
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err)
				assert.NotNil(t, parsed["settings"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSettingsCatalogJSONArray(ctx, tt.settingsStr, tt.resp)
			tt.validate(t, result)
		})
	}
}

// TestStateConfigurationPolicySettings tests the StateConfigurationPolicySettings wrapper function
func TestStateConfigurationPolicySettings(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		initialSettings  string
		apiResponse      []byte
		expectedContains []string
	}{
		{
			name:            "Basic settings array",
			initialSettings: `{"settings":[]}`,
			apiResponse: []byte(`{
				"value": [
					{
						"id": "1",
						"settingDefinitionId": "device_vendor_msft_policy_config_test",
						"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
					}
				]
			}`),
			expectedContains: []string{`"settings"`, `"settingDefinitionId"`, `"device_vendor_msft_policy_config_test"`},
		},
		{
			name:            "Empty initial settings",
			initialSettings: "",
			apiResponse: []byte(`{
				"value": [
					{
						"id": "1",
						"settingDefinitionId": "test_setting"
					}
				]
			}`),
			expectedContains: []string{`"settings"`, `"test_setting"`},
		},
		{
			name:            "Settings with secrets preserved",
			initialSettings: `{"settings":[{"id":"1","settingDefinitionId":"test","simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSecretSettingValue","valueState":"encryptedValueToken","value":"secret123"}}]}`,
			apiResponse: []byte(`{
				"value": [
					{
						"id": "1",
						"settingDefinitionId": "test",
						"simpleSettingValue": {
							"@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
							"valueState": "notEncrypted"
						}
					}
				]
			}`),
			expectedContains: []string{`"encryptedValueToken"`, `"secret123"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &sharedmodels.SettingsCatalogJsonResourceModel{
				Settings: types.StringValue(tt.initialSettings),
			}

			StateConfigurationPolicySettings(ctx, data, tt.apiResponse)

			result := data.Settings.ValueString()
			for _, expected := range tt.expectedContains {
				assert.Contains(t, result, expected)
			}

			// Verify it's valid JSON
			var resultMap map[string]any
			err := json.Unmarshal([]byte(result), &resultMap)
			require.NoError(t, err, "Result should be valid JSON")
		})
	}
}

// TestStateReusablePolicySettings tests the StateReusablePolicySettings wrapper function
func TestStateReusablePolicySettings(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		initialSettings  string
		apiResponse      []byte
		expectedContains []string
	}{
		{
			name:            "Single setting instance",
			initialSettings: `{"settings":[{"id":"0","settingInstance":{}}]}`,
			apiResponse: []byte(`{
				"settingInstance": {
					"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
					"settingDefinitionId": "device_vendor_msft_policy_config_test",
					"choiceSettingValue": {
						"value": "device_vendor_msft_policy_config_test_enabled"
					}
				}
			}`),
			expectedContains: []string{`"settings"`, `"id":"0"`, `"settingInstance"`, `"settingDefinitionId"`},
		},
		{
			name:            "Empty initial settings",
			initialSettings: "",
			apiResponse: []byte(`{
				"settingInstance": {
					"settingDefinitionId": "test_setting",
					"@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
				}
			}`),
			expectedContains: []string{`"settings"`, `"test_setting"`, `"id":"0"`},
		},
		{
			name:            "Complex setting instance",
			initialSettings: `{}`,
			apiResponse: []byte(`{
				"settingInstance": {
					"@odata.type": "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
					"settingDefinitionId": "device_vendor_msft_policy_group_test",
					"groupSettingCollectionValue": [
						{
							"children": [
								{
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
									"settingDefinitionId": "nested_setting"
								}
							]
						}
					]
				}
			}`),
			expectedContains: []string{`"settingInstance"`, `"groupSettingCollectionValue"`, `"children"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &sharedmodels.ReuseablePolicySettingsResourceModel{
				Settings: types.StringValue(tt.initialSettings),
			}

			StateReusablePolicySettings(ctx, data, tt.apiResponse)

			result := data.Settings.ValueString()
			for _, expected := range tt.expectedContains {
				assert.Contains(t, result, expected)
			}

			// Verify it's valid JSON
			var resultMap map[string]any
			err := json.Unmarshal([]byte(result), &resultMap)
			require.NoError(t, err, "Result should be valid JSON")
		})
	}
}

// TestNormalizeSettingsCatalogJSON tests the single-setting normalizer
func TestNormalizeSettingsCatalogJSON(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		settingsStr string
		resp        []byte
		validate    func(t *testing.T, result string)
	}{
		{
			name:        "Valid setting instance wrapped correctly",
			settingsStr: `{}`,
			resp: []byte(`{
				"settingInstance": {
					"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
					"settingDefinitionId": "test_setting",
					"choiceSettingValue": {
						"value": "test_value"
					}
				}
			}`),
			validate: func(t *testing.T, result string) {
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err)

				settings, ok := parsed["settings"].([]any)
				require.True(t, ok, "Should have settings array")
				require.Len(t, settings, 1, "Should have one setting")

				setting := settings[0].(map[string]any)
				assert.Equal(t, "0", setting["id"])
				assert.NotNil(t, setting["settingInstance"])
			},
		},
		{
			name:        "Invalid JSON returns original",
			settingsStr: `{"original":"settings"}`,
			resp:        []byte(`invalid json`),
			validate: func(t *testing.T, result string) {
				assert.Equal(t, `{"original":"settings"}`, result)
			},
		},
		{
			name:        "Response without settingInstance returns original",
			settingsStr: `{"original":"settings"}`,
			resp: []byte(`{
				"someOtherField": "value"
			}`),
			validate: func(t *testing.T, result string) {
				assert.Equal(t, `{"original":"settings"}`, result)
			},
		},
		{
			name:        "Empty settings string with valid response",
			settingsStr: "",
			resp: []byte(`{
				"settingInstance": {
					"settingDefinitionId": "test"
				}
			}`),
			validate: func(t *testing.T, result string) {
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err)
				assert.NotNil(t, parsed["settings"])
			},
		},
		{
			name:        "Complex nested setting instance",
			settingsStr: `{}`,
			resp: []byte(`{
				"settingInstance": {
					"@odata.type": "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
					"settingDefinitionId": "group_setting",
					"groupSettingCollectionValue": [
						{
							"children": [
								{
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
									"settingDefinitionId": "nested_choice",
									"choiceSettingValue": {
										"value": "nested_value",
										"children": []
									}
								}
							]
						}
					]
				}
			}`),
			validate: func(t *testing.T, result string) {
				var parsed map[string]any
				err := json.Unmarshal([]byte(result), &parsed)
				require.NoError(t, err)

				settings := parsed["settings"].([]any)
				setting := settings[0].(map[string]any)
				instance := setting["settingInstance"].(map[string]any)

				assert.Equal(t, "group_setting", instance["settingDefinitionId"])
				assert.NotNil(t, instance["groupSettingCollectionValue"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSettingsCatalogJSON(ctx, tt.settingsStr, tt.resp)
			tt.validate(t, result)
		})
	}
}
