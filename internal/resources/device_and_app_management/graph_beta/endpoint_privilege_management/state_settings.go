package graphBetaEndpointPrivilegeManagement

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
func MapRemoteSettingsStateToTerraform(ctx context.Context, data *EndpointPrivilegeManagementResourceModel, resp []byte) {
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

	if err := normalize.PreserveSecretSettings(configSettings, structuredContent); err != nil {
		tflog.Error(ctx, "Error stating settings catalog secret settings from HCL", map[string]interface{}{"error": err.Error()})
		return
	}

	jsonBytes, err := json.Marshal(structuredContent)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal JSON structured content during preparation for normalization", map[string]interface{}{"error": err.Error()})
		return
	}

	normalizedJSON, err := normalize.JSONAlphabetically(string(jsonBytes))
	if err != nil {
		tflog.Error(ctx, "Failed to normalize settings catalog JSON alphabetically", map[string]interface{}{"error": err.Error()})
		return
	}

	tflog.Debug(ctx, "Original settings", map[string]interface{}{"settings": string(resp)})
	tflog.Debug(ctx, "Normalized settings", map[string]interface{}{"settings": normalizedJSON})

	data.Settings = types.StringValue(normalizedJSON)
}
