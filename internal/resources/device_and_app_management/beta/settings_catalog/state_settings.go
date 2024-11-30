package graphBetaSettingsCatalog

import (
	"context"
	"encoding/json"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/normalize"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteSettingsStateToTerraform maps the remote settings catalog settings state to the Terraform state
// taking the raw json from a custom GET request, normalizing the content and then stating. The stating logic:
// 1. Parses the original HCL settings to preserve secret values and states
// 2. Parses the raw response and extracts the settings content
// 3. Structures the content under settingsDetails like the PUT request
// 4. Recursively preserves secret setting values and states from the original HCL config
// 5. Converts the structured content to JSON and normalizes it alphabetically
// 6. States the normalized JSON in the Terraform state
func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, resp []byte) {
	var configSettings map[string]interface{}
	if err := json.Unmarshal([]byte(data.Settings.ValueString()), &configSettings); err != nil {
		tflog.Error(ctx, "Failed to unmarshal config settings", map[string]interface{}{"error": err.Error()})
		return
	}

	var rawResponse map[string]interface{}
	if err := json.Unmarshal(resp, &rawResponse); err != nil {
		var arrayResponse []interface{}
		if err := json.Unmarshal(resp, &arrayResponse); err != nil {
			tflog.Error(ctx, "Failed to unmarshal settings response", map[string]interface{}{"error": err.Error()})
			return
		}
		rawResponse = map[string]interface{}{"value": arrayResponse}
	}

	var settingsContent interface{}
	if value, ok := rawResponse["value"]; ok {
		settingsContent = value
	} else if details, ok := rawResponse["settingsDetails"]; ok {
		settingsContent = details
	} else {
		settingsContent = rawResponse
	}

	structuredContent := map[string]interface{}{
		"settingsDetails": settingsContent,
	}

	preserveSecretSettings(configSettings, structuredContent)

	jsonBytes, err := json.Marshal(structuredContent)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal structured content", map[string]interface{}{"error": err.Error()})
		return
	}

	normalizedJSON, err := normalize.JSONAlphabetically(string(jsonBytes))
	if err != nil {
		tflog.Error(ctx, "Failed to normalize JSON alphabetically", map[string]interface{}{"error": err.Error()})
		return
	}

	tflog.Debug(ctx, "Original settings", map[string]interface{}{"settings": string(resp)})
	tflog.Debug(ctx, "Normalized settings", map[string]interface{}{"settings": normalizedJSON})

	data.Settings = types.StringValue(normalizedJSON)
}

// preserveSecretSettings recursively searches through settings catalog HCL JSON structure for secret settings
// and preserves the value and valueState from the config settings. This is used to ensure that secret values
// within the state match the original config settings and do not cause unnecessary updates.
func preserveSecretSettings(config, resp interface{}) {
	switch configV := config.(type) {
	case map[string]interface{}:
		respV, ok := resp.(map[string]interface{})
		if !ok {
			return
		}

		if odataType, ok := configV["@odata.type"].(string); ok &&
			odataType == "#microsoft.graph.deviceManagementConfigurationSecretSettingValue" {
			if value, ok := configV["value"]; ok {
				respV["value"] = value
			}
			if valueState, ok := configV["valueState"]; ok {
				respV["valueState"] = valueState
			}
			return
		}

		for k, v := range configV {
			if respChild, ok := respV[k]; ok {
				preserveSecretSettings(v, respChild)
			}
		}

	case []interface{}:
		respV, ok := resp.([]interface{})
		if !ok {
			return
		}
		for i := range configV {
			if i < len(respV) {
				preserveSecretSettings(configV[i], respV[i])
			}
		}
	}
}
