package normalize

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeJSONAlphabetically(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple flat object",
			input:    `{"z": 1, "a": 2, "b": 3}`,
			expected: `{"a":2,"b":3,"z":1}`,
			wantErr:  false,
		},
		{
			name:     "Nested object",
			input:    `{"z": {"y": 1, "a": 2}, "a": 3}`,
			expected: `{"a":3,"z":{"a":2,"y":1}}`,
			wantErr:  false,
		},
		{
			name:     "Array values",
			input:    `{"b": [{"z": 1, "a": 2}, {"y": 3, "x": 4}], "a": 5}`,
			expected: `{"a":5,"b":[{"a":2,"z":1},{"x":4,"y":3}]}`,
			wantErr:  false,
		},
		{
			name: "Complex nested structure",
			input: `{
				"settingsDetails": [
					{
						"id": "1",
						"settingInstance": {
							"simpleSettingValue": {
								"value": "test",
								"@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
								"settingValueTemplateReference": null
							},
							"@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
							"settingDefinitionId": "test_id",
							"settingInstanceTemplateReference": null
						}
					}
				]
			}`,
			expected: `{"settingsDetails":[{"id":"1","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"test_id","settingInstanceTemplateReference":null,"simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","settingValueTemplateReference":null,"value":"test"}}}]}`,
			wantErr:  false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"invalid": json}`,
			wantErr: true,
		},
		{
			name:     "Empty object",
			input:    `{}`,
			expected: `{}`,
			wantErr:  false,
		},
		{
			name:     "Null values",
			input:    `{"b": null, "a": null}`,
			expected: `{"a":null,"b":null}`,
			wantErr:  false,
		},
		{
			name:     "Mixed types",
			input:    `{"string": "value", "number": 123, "bool": true, "null": null}`,
			expected: `{"bool":true,"null":null,"number":123,"string":"value"}`,
			wantErr:  false,
		},
		{
			name: "Choice setting collection with nested settings",
			input: `{
					"settingsDetails": [{
							"id": "1",
							"settingInstance": {
									"choiceSettingCollectionValue": [{
											"children": [{
													"simpleSettingValue": {
															"value": "test",
															"@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
															"settingValueTemplateReference": null
													},
													"@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
													"settingDefinitionId": "test_id",
													"settingInstanceTemplateReference": null
											}],
											"settingValueTemplateReference": null,
											"value": "test_value"
									}]
							}
					}]
			}`,
			expected: `{"settingsDetails":[{"id":"1","settingInstance":{"choiceSettingCollectionValue":[{"children":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"test_id","settingInstanceTemplateReference":null,"simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","settingValueTemplateReference":null,"value":"test"}}],"settingValueTemplateReference":null,"value":"test_value"}]}}]}`,
			wantErr:  false,
		},
		{
			name: "Deep nested group settings with all types",
			input: `{
					"settingsDetails": [{
							"id": "1",
							"settingInstance": {
									"groupSettingCollectionValue": [{
											"children": [{
													"groupSettingCollectionValue": [{
															"children": [{
																	"choiceSettingValue": {
																			"children": [{
																					"simpleSettingValue": {
																							"value": 123,
																							"@odata.type": "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
																							"settingValueTemplateReference": null,
																							"valueState": "valid"
																					}
																			}],
																			"value": "choice_value"
																	}
															}]
													}]
											}],
											"settingValueTemplateReference": null
									}],
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
							}
					}]
			}`,
			expected: `{"settingsDetails":[{"id":"1","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance","groupSettingCollectionValue":[{"children":[{"groupSettingCollectionValue":[{"children":[{"choiceSettingValue":{"children":[{"simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationIntegerSettingValue","settingValueTemplateReference":null,"value":123,"valueState":"valid"}}],"value":"choice_value"}}]}]}],"settingValueTemplateReference":null}]}}]}`,
			wantErr:  false,
		},
		{
			name: "Multiple setting types in one structure",
			input: `{
					"settingsDetails": [{
							"id": "1",
							"settingInstance": {
									"simpleSettingValue": {
											"value": "simple",
											"@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
									},
									"choiceSettingValue": {
											"children": [],
											"value": "choice"
									},
									"simpleSettingCollectionValue": [{
											"value": "collection",
											"@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
									}],
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationInstance"
							}
					}]
			}`,
			expected: `{"settingsDetails":[{"id":"1","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationInstance","choiceSettingValue":{"children":[],"value":"choice"},"simpleSettingCollectionValue":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","value":"collection"}],"simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","value":"simple"}}}]}`,
			wantErr:  false,
		},
		{
			name: "MacOS Complex Group Setting Collection",
			input: `{
					"settingsDetails": [{
							"id": "0",
							"settingInstance": {
									"groupSettingCollectionValue": [{
											"children": [{
													"choiceSettingValue": {
															"value": "com.apple.example_setting_true",
															"children": [],
															"settingValueTemplateReference": null
													},
													"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
													"settingDefinitionId": "com.apple.example_setting",
													"settingInstanceTemplateReference": null
											}],
											"settingValueTemplateReference": null
									}],
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
									"settingDefinitionId": "com.apple.example",
									"settingInstanceTemplateReference": null
							}
					}]
			}`,
			expected: `{"settingsDetails":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance","groupSettingCollectionValue":[{"children":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","choiceSettingValue":{"children":[],"settingValueTemplateReference":null,"value":"com.apple.example_setting_true"},"settingDefinitionId":"com.apple.example_setting","settingInstanceTemplateReference":null}],"settingValueTemplateReference":null}],"settingDefinitionId":"com.apple.example","settingInstanceTemplateReference":null}}]}`,
			wantErr:  false,
		},
		{
			name: "Mixed Setting Types In Collection",
			input: `{
					"settingsDetails": [{
							"id": "1",
							"settingInstance": {
									"groupSettingCollectionValue": [{
											"children": [{
													"simpleSettingValue": {
															"value": 1,
															"@odata.type": "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
															"settingValueTemplateReference": null
													},
													"@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
													"settingDefinitionId": "com.apple.numeric_setting",
													"settingInstanceTemplateReference": null
											},
											{
													"simpleSettingCollectionValue": [{
															"value": "thing",
															"@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
															"settingValueTemplateReference": null
													}],
													"@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
													"settingDefinitionId": "com.apple.array_setting",
													"settingInstanceTemplateReference": null
											}],
											"settingValueTemplateReference": null
									}]
							}
					}]
			}`,
			expected: `{"settingsDetails":[{"id":"1","settingInstance":{"groupSettingCollectionValue":[{"children":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"com.apple.numeric_setting","settingInstanceTemplateReference":null,"simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationIntegerSettingValue","settingValueTemplateReference":null,"value":1}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.apple.array_setting","settingInstanceTemplateReference":null,"simpleSettingCollectionValue":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","settingValueTemplateReference":null,"value":"thing"}]}],"settingValueTemplateReference":null}]}}]}`,
			wantErr:  false,
		},
		{
			name: "Deeply Nested Group Collections",
			input: `{
					"settingsDetails": [{
							"id": "6",
							"settingInstance": {
									"groupSettingCollectionValue": [{
											"children": [{
													"groupSettingCollectionValue": [{
															"children": [{
																	"groupSettingCollectionValue": [{
																			"children": [{
																					"simpleSettingValue": {
																							"value": "test",
																							"@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
																							"settingValueTemplateReference": null
																					}
																			}]
																	}]
															}]
													}]
											}]
									}],
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
							}
					}]
			}`,
			expected: `{"settingsDetails":[{"id":"6","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance","groupSettingCollectionValue":[{"children":[{"groupSettingCollectionValue":[{"children":[{"groupSettingCollectionValue":[{"children":[{"simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","settingValueTemplateReference":null,"value":"test"}}]}]}]}]}]}]}}]}`,
			wantErr:  false,
		},
		{
			name: "Choice Settings With Empty Children",
			input: `{
					"settingsDetails": [{
							"id": "5",
							"settingInstance": {
									"choiceSettingValue": {
											"value": "com.apple.setting_true",
											"children": [],
											"settingValueTemplateReference": null
									},
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
							}
					}]
			}`,
			expected: `{"settingsDetails":[{"id":"5","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","choiceSettingValue":{"children":[],"settingValueTemplateReference":null,"value":"com.apple.setting_true"}}}]}`,
			wantErr:  false,
		},
		{
			name: "Multiple Choice Settings With Empty Children",
			input: `{
                "settingsDetails": [{
                    "id": "0",
                    "settingInstance": {
                        "choiceSettingValue": {
                            "children": [],
                            "value": "device_vendor_msft_policy_config_something_1",
                            "settingValueTemplateReference": null
                        },
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                        "settingDefinitionId": "device_vendor_msft_policy_config_something",
                        "settingInstanceTemplateReference": null
                    }
                },
                {
                    "id": "1",
                    "settingInstance": {
                        "choiceSettingValue": {
                            "children": [],
                            "value": "user_vendor_msft_policy_config_something_1",
                            "settingValueTemplateReference": null
                        },
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                        "settingDefinitionId": "user_vendor_msft_policy_config_something",
                        "settingInstanceTemplateReference": null
                    }
                }]
            }`,
			expected: `{"settingsDetails":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","choiceSettingValue":{"children":[],"settingValueTemplateReference":null,"value":"device_vendor_msft_policy_config_something_1"},"settingDefinitionId":"device_vendor_msft_policy_config_something","settingInstanceTemplateReference":null}},{"id":"1","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","choiceSettingValue":{"children":[],"settingValueTemplateReference":null,"value":"user_vendor_msft_policy_config_something_1"},"settingDefinitionId":"user_vendor_msft_policy_config_something","settingInstanceTemplateReference":null}}]}`,
			wantErr:  false,
		},
		{
			name: "Choice Setting With Nested Choice Settings Array",
			input: `{
                "settingsDetails": [{
                    "id": "0",
                    "settingInstance": {
                        "choiceSettingValue": {
                            "children": [
                                {
                                    "choiceSettingValue": {
                                        "children": [],
                                        "value": "user_vendor_msft_policy_config_nested_1",
                                        "settingValueTemplateReference": null
                                    },
                                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                                    "settingDefinitionId": "user_vendor_msft_policy_config_nested",
                                    "settingInstanceTemplateReference": null
                                },
                                {
                                    "choiceSettingValue": {
                                        "children": [],
                                        "value": "user_vendor_msft_policy_config_nested2_1",
                                        "settingValueTemplateReference": null
                                    },
                                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                                    "settingDefinitionId": "user_vendor_msft_policy_config_nested2",
                                    "settingInstanceTemplateReference": null
                                }
                            ],
                            "value": "parent_setting_value",
                            "settingValueTemplateReference": null
                        },
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                        "settingDefinitionId": "parent_setting",
                        "settingInstanceTemplateReference": null
                    }
                }]
            }`,
			expected: `{"settingsDetails":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","choiceSettingValue":{"children":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","choiceSettingValue":{"children":[],"settingValueTemplateReference":null,"value":"user_vendor_msft_policy_config_nested_1"},"settingDefinitionId":"user_vendor_msft_policy_config_nested","settingInstanceTemplateReference":null},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","choiceSettingValue":{"children":[],"settingValueTemplateReference":null,"value":"user_vendor_msft_policy_config_nested2_1"},"settingDefinitionId":"user_vendor_msft_policy_config_nested2","settingInstanceTemplateReference":null}],"settingValueTemplateReference":null,"value":"parent_setting_value"},"settingDefinitionId":"parent_setting","settingInstanceTemplateReference":null}}]}`,
			wantErr:  false,
		},
		{
			name: "Choice Setting With Mixed GroupSettingCollection and SimpleSettingValue",
			input: `{
                "settingsDetails": [{
                    "id": "0",
                    "settingInstance": {
                        "choiceSettingValue": {
                            "children": [
                                {
                                    "groupSettingCollectionValue": [
                                        {
                                            "settingValueTemplateReference": null,
                                            "children": [
                                                {
                                                    "simpleSettingValue": {
                                                        "value": "key_value",
                                                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                                                        "settingValueTemplateReference": null
                                                    },
                                                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                                                    "settingDefinitionId": "key_setting",
                                                    "settingInstanceTemplateReference": null
                                                },
                                                {
                                                    "simpleSettingValue": {
                                                        "value": "value_string",
                                                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                                                        "settingValueTemplateReference": null
                                                    },
                                                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                                                    "settingDefinitionId": "value_setting",
                                                    "settingInstanceTemplateReference": null
                                                }
                                            ]
                                        }
                                    ],
                                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
                                    "settingDefinitionId": "collection_setting",
                                    "settingInstanceTemplateReference": null
                                }
                            ],
                            "value": "parent_value",
                            "settingValueTemplateReference": null
                        },
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                        "settingDefinitionId": "parent_setting",
                        "settingInstanceTemplateReference": null
                    }
                }]
            }`,
			expected: `{"settingsDetails":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","choiceSettingValue":{"children":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance","groupSettingCollectionValue":[{"children":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"key_setting","settingInstanceTemplateReference":null,"simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","settingValueTemplateReference":null,"value":"key_value"}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"value_setting","settingInstanceTemplateReference":null,"simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","settingValueTemplateReference":null,"value":"value_string"}}],"settingValueTemplateReference":null}],"settingDefinitionId":"collection_setting","settingInstanceTemplateReference":null}],"settingValueTemplateReference":null,"value":"parent_value"},"settingDefinitionId":"parent_setting","settingInstanceTemplateReference":null}}]}`,
			wantErr:  false,
		},
		{
			name: "SimpleSettingCollectionValue With Multiple String Values",
			input: `{
                "settingsDetails": [{
                    "settingInstance": {
                        "simpleSettingCollectionValue": [
                            {
                                "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                                "settingValueTemplateReference": null,
                                "value": "thing"
                            },
                            {
                                "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                                "settingValueTemplateReference": null,
                                "value": "thing2"
                            }
                        ],
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
                        "settingDefinitionId": "some_setting_alloweddomains"
                    }
                }]
            }`,
			expected: `{"settingsDetails":[{"settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"some_setting_alloweddomains","simpleSettingCollectionValue":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","settingValueTemplateReference":null,"value":"thing"},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","settingValueTemplateReference":null,"value":"thing2"}]}}]}`,
			wantErr:  false,
		},
		{
			name: "Mixed Integer and String Values in Group Collection",
			input: `{
                "settingsDetails": [{
                    "settingInstance": {
                        "groupSettingCollectionValue": [{
                            "children": [
                                {
                                    "simpleSettingValue": {
                                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                                        "value": 1,
                                        "settingValueTemplateReference": null
                                    },
                                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                                    "settingDefinitionId": "numeric_setting"
                                },
                                {
                                    "simpleSettingValue": {
                                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                                        "value": "2024-10-31T00:00:00",
                                        "settingValueTemplateReference": null
                                    },
                                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                                    "settingDefinitionId": "date_setting"
                                }
                            ],
                            "settingValueTemplateReference": null
                        }],
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
                        "settingDefinitionId": "mixed_settings"
                    }
                }]
            }`,
			expected: `{"settingsDetails":[{"settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance","groupSettingCollectionValue":[{"children":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"numeric_setting","simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationIntegerSettingValue","settingValueTemplateReference":null,"value":1}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"date_setting","simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","settingValueTemplateReference":null,"value":"2024-10-31T00:00:00"}}],"settingValueTemplateReference":null}],"settingDefinitionId":"mixed_settings"}}]}`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONAlphabetically(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			t.Logf("Expected length: %d", len(tt.expected))
			t.Logf("Got length: %d", len(got))

			// Print characters if lengths don't match
			if len(tt.expected) != len(got) {
				t.Logf("Expected: %v", []byte(tt.expected))
				t.Logf("Got: %v", []byte(got))
			}

			// Try unmarshaling both to compare structure
			var expectedMap, gotMap map[string]interface{}
			err = json.Unmarshal([]byte(tt.expected), &expectedMap)
			assert.NoError(t, err, "Failed to unmarshal expected JSON")

			err = json.Unmarshal([]byte(got), &gotMap)
			assert.NoError(t, err, "Failed to unmarshal actual JSON")

			assert.Equal(t, expectedMap, gotMap)
		})
	}
}
