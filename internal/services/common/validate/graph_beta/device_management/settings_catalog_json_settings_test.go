package sharedValidators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_settingsCatalogValidator(t *testing.T) {
	tests := []struct {
		name          string
		value         types.String
		expectedError string
	}{
		{
			name: "valid_hierarchy",
			value: types.StringValue(`{
				"settings": [{
					"id": "0",
					"settingInstance": {
						"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
						"settingDefinitionId": "device_vendor_msft_policy_elevationclientsettings_enableepm",
						"choiceSettingValue": {
							"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
							"value": "device_vendor_msft_policy_elevationclientsettings_enableepm_1",
							"children": []
						}
					}
				}],
				"templateReference": {
					"templateId": "e7dcaba4-959b-46ed-88f0-16ba39b14fd8_1"
				}
			}`),
			expectedError: "",
		},
		{
			name: "invalid_hierarchy_extra_field",
			value: types.StringValue(`{
					"settings": [{
							"@odata.type": "#microsoft.graph.deviceManagementConfigurationSetting",
							"id": "0",
							"settingInstance": {}
					}]
			}`),
			expectedError: `Setting at index 0 contains 3 fields (["@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting", "id" = "0", "settingInstance" = <object>]), but should only contain exactly 2 fields ('id' and 'settingInstance')`,
		},
		{
			name: "invalid_hierarchy_extra_field",
			value: types.StringValue(`{
        "settings": [{
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationSetting",
            "id": "0",
            "settingInstance": {}
        }]
    }`),
			expectedError: `Setting at index 0 contains 3 fields (["@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting", "id" = "0", "settingInstance" = <object>]), but should only contain exactly 2 fields ('id' and 'settingInstance')`,
		},
		{
			name: "valid_secret_setting",
			value: types.StringValue(`{
				"settings": [{
					"id": "0",
					"settingInstance": {
						"@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
						"settingDefinitionId": "com.apple.loginwindow_autologinpassword",
						"simpleSettingValue": {
							"@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
							"valueState": "notEncrypted",
							"value": "09047f39-e38f-4baf-a380-23ddcbb73d41"
						}
					}
				}]
			}`),
			expectedError: "",
		},
		{
			name: "invalid_secret_setting_state",
			value: types.StringValue(`{
					"settings": [{
							"id": "0",
							"settingInstance": {
									"settingDefinitionId": "com.apple.loginwindow_autologinpassword",
									"simpleSettingValue": {
											"@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
											"valueState": "encryptedValueToken",
											"value": "09047f39-e38f-4baf-a380-23ddcbb73d41",
											"settingDefinitionId": "com.apple.loginwindow_autologinpassword"
									}
							}
					}]
			}`),
			expectedError: `Secret Setting Value (settingDefinitionId: com.apple.loginwindow_autologinpassword) state must be 'notEncrypted' when setting a new secret value, got 'encryptedValueToken'`,
		},
		{
			name: "invalid_id_sequence",
			value: types.StringValue(`{
				"settings": [{
					"id": "0",
					"settingInstance": {}
				},
				{
					"id": "2",
					"settingInstance": {}
				}]
			}`),
			expectedError: "Settings IDs must increment by 1. Found ID 2 after ID 0",
		},
		{
			name: "invalid_id_format",
			value: types.StringValue(`{
					"settings": [{
							"id": "abc",
							"settingInstance": {}
					}]
			}`),
			expectedError: "Settings ID must be a number: abc",
		},
		{
			name: "missing_id",
			value: types.StringValue(`{
				"settings": [{
					"settingInstance": {}
				}]
			}`),
			expectedError: "Setting at index 0 is missing required 'id' field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := validator.StringRequest{
				ConfigValue: tt.value,
			}
			response := validator.StringResponse{}

			SettingsCatalogJSONValidator().ValidateString(context.TODO(), request, &response)

			if tt.expectedError == "" && response.Diagnostics.HasError() {
				t.Errorf("expected no error, got: %v", response.Diagnostics)
			}

			if tt.expectedError != "" {
				if !response.Diagnostics.HasError() {
					t.Errorf("expected error containing %q, got no error", tt.expectedError)
					return
				}

				diagError := response.Diagnostics.Errors()[0].Detail()
				if diagError != tt.expectedError {
					t.Errorf("expected error containing %q, got %q", tt.expectedError, diagError)
				}
			}
		})
	}
}
