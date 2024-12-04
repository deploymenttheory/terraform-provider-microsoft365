package normalize

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreserveSecretSettings(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		response    string
		expected    string
		expectError bool
	}{
		{
			name: "secret setting preservation",
			config: `{
                "settingValue": {
                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value": "secretValue123",
                    "valueState": "notEncrypted"
                }
            }`,
			response: `{
                "settingValue": {
                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value": "differentValue",
                    "valueState": "encryptedValueToken"
                }
            }`,
			expected: `{
                "settingValue": {
                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value": "secretValue123",
                    "valueState": "notEncrypted"
                }
            }`,
			expectError: false,
		},
		{
			name: "nested secret settings",
			config: `{
                "outer": {
                    "inner": {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                        "value": "secretValue123",
                        "valueState": "notEncrypted"
                    }
                }
            }`,
			response: `{
                "outer": {
                    "inner": {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                        "value": "differentValue",
                        "valueState": "encryptedValueToken"
                    }
                }
            }`,
			expected: `{
                "outer": {
                    "inner": {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                        "value": "secretValue123",
                        "valueState": "notEncrypted"
                    }
                }
            }`,
			expectError: false,
		},
		{
			name: "array of secret settings",
			config: `{
                "settings": [
                    {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                        "value": "secret1",
                        "valueState": "notEncrypted"
                    },
                    {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                        "value": "secret2",
                        "valueState": "notEncrypted"
                    }
                ]
            }`,
			response: `{
                "settings": [
                    {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                        "value": "different1",
                        "valueState": "encryptedValueToken"
                    },
                    {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                        "value": "different2",
                        "valueState": "encryptedValueToken"
                    }
                ]
            }`,
			expected: `{
                "settings": [
                    {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                        "value": "secret1",
                        "valueState": "notEncrypted"
                    },
                    {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                        "value": "secret2",
                        "valueState": "notEncrypted"
                    }
                ]
            }`,
			expectError: false,
		},
		{
			name: "type mismatch error",
			config: `{
                "setting": []
            }`,
			response: `{
                "setting": "not an array"
            }`,
			expectError: true,
		},
		{
			name: "non-matching structure",
			config: `{
                "setting": {
                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value": "secret"
                }
            }`,
			response: `{
                "differentKey": {
                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value": "different"
                }
            }`,
			expected: `{
                "differentKey": {
                    "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value": "different"
                }
            }`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse JSON strings to interfaces
			var configMap interface{}
			err := json.Unmarshal([]byte(tt.config), &configMap)
			require.NoError(t, err, "Failed to unmarshal config JSON")

			var responseMap interface{}
			err = json.Unmarshal([]byte(tt.response), &responseMap)
			require.NoError(t, err, "Failed to unmarshal response JSON")

			// Run the function
			err = PreserveSecretSettings(configMap, responseMap)

			// Check error expectation
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.expected != "" {
				var expectedMap interface{}
				err = json.Unmarshal([]byte(tt.expected), &expectedMap)
				require.NoError(t, err, "Failed to unmarshal expected JSON")

				// Compare the result with expected
				assert.Equal(t, expectedMap, responseMap)
			}
		})
	}
}

func TestPreserveSecretSettings_UnsupportedType(t *testing.T) {
	// Test with string type mismatch
	config := "string"
	resp := 123 // Different type to trigger error
	err := PreserveSecretSettings(config, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type mismatch")
}

func TestPreserveSecretSettings_NilValues(t *testing.T) {
	tests := []struct {
		name   string
		config interface{}
		resp   interface{}
	}{
		{
			name:   "nil config",
			config: nil,
			resp:   map[string]interface{}{},
		},
		{
			name:   "nil response",
			config: map[string]interface{}{},
			resp:   nil,
		},
		{
			name:   "both nil",
			config: nil,
			resp:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PreserveSecretSettings(tt.config, tt.resp)
			assert.Error(t, err)
		})
	}
}

func TestPreserveSecretSettings_PrimitiveTypes(t *testing.T) {
	tests := []struct {
		name        string
		config      interface{}
		resp        interface{}
		expectError bool
	}{
		{
			name:        "matching strings",
			config:      "test",
			resp:        "different",
			expectError: false,
		},
		{
			name:        "matching numbers",
			config:      123.45,
			resp:        678.90,
			expectError: false,
		},
		{
			name:        "matching booleans",
			config:      true,
			resp:        false,
			expectError: false,
		},
		{
			name:        "type mismatch string-number",
			config:      "test",
			resp:        123,
			expectError: true,
		},
		{
			name:        "type mismatch number-bool",
			config:      123,
			resp:        true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PreserveSecretSettings(tt.config, tt.resp)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
