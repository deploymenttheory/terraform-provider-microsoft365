package graphBetaSettingsCatalog

import (
	"context"
	"encoding/json"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/normalize"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteSettingsStateToTerraform maps the remote settings catalog settings state to the Terraform state
// using the shared DeviceConfigV2GraphServiceMode used also for requests.
// Results are normalized and stored alphabetically in the Terraform state.
// func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, resp []byte) {
// 	if err := json.Unmarshal(resp, &DeviceConfigV2GraphServiceModel); err != nil {
// 		tflog.Error(ctx, "Failed to unmarshal settings", map[string]interface{}{"error": err.Error()})
// 		return
// 	}

// 	normalizedJSON, err := normalize.JSONAlphabetically(string(resp))
// 	if err != nil {
// 		tflog.Error(ctx, "Failed to normalize JSON", map[string]interface{}{"error": err.Error()})
// 		return
// 	}

// 	data.Settings = types.StringValue(normalizedJSON)
// }

func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, resp []byte) {
	// Parse the raw response
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(resp, &rawResponse); err != nil {
		// Try parsing as array if map fails
		var arrayResponse []interface{}
		if err := json.Unmarshal(resp, &arrayResponse); err != nil {
			tflog.Error(ctx, "Failed to unmarshal settings response", map[string]interface{}{"error": err.Error()})
			return
		}
		rawResponse = map[string]interface{}{"value": arrayResponse}
	}

	// Extract the actual settings content
	var settingsContent interface{}
	if value, ok := rawResponse["value"]; ok {
		settingsContent = value
	} else if details, ok := rawResponse["settingsDetails"]; ok {
		settingsContent = details
	} else {
		// If neither exists, assume the content is the settings
		settingsContent = rawResponse
	}

	// Create properly structured content
	structuredContent := map[string]interface{}{
		"settingsDetails": settingsContent,
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(structuredContent)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal structured content", map[string]interface{}{"error": err.Error()})
		return
	}

	// Apply alphabetical normalization
	normalizedJSON, err := normalize.JSONAlphabetically(string(jsonBytes))
	if err != nil {
		tflog.Error(ctx, "Failed to normalize JSON alphabetically", map[string]interface{}{"error": err.Error()})
		return
	}

	// Log the before/after for debugging if needed
	tflog.Debug(ctx, "Original settings", map[string]interface{}{"settings": string(resp)})
	tflog.Debug(ctx, "Normalized settings", map[string]interface{}{"settings": normalizedJSON})

	data.Settings = types.StringValue(normalizedJSON)
}
