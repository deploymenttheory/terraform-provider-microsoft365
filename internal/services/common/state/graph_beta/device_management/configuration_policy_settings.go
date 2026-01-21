package sharedStater

import (
	"context"
	"encoding/json"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/normalize"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// StateConfigurationPolicySettings states settings for the SettingsCatalogProfileResourceModel.
// It processes the response from the Graph API,
// normalizes the settings array, and preserves any secret values from the original configuration.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - data: The resource model containing the settings to be stated
//   - resp: Raw response bytes from the Graph API containing the settings array
//
// The function maps the remote state to Terraform state while maintaining the array structure
// and normalizing the JSON format for consistent state representation.
func StateConfigurationPolicySettings(ctx context.Context, data *sharedmodels.SettingsCatalogJsonResourceModel, resp []byte) {
	data.Settings = types.StringValue(normalizeSettingsCatalogJSONArray(ctx, data.Settings.ValueString(), resp))
}

// StateReusablePolicySettings maps a single Microsoft Graph reusable policy setting response
// to the Terraform state. This resource type only supports one setting instance at a time.
// It processes the response from the Graph API,wraps the single setting instance in the
// expected format, and preserves any secret values.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - data: The resource model containing the setting to be stated
//   - resp: Raw response bytes from the Graph API containing the single setting instance
//
// The function maps the remote state to Terraform state while wrapping the single
// setting instance in the expected format for consistent state representation.
func StateReusablePolicySettings(ctx context.Context, data *sharedmodels.ReuseablePolicySettingsResourceModel, resp []byte) {
	data.Settings = types.StringValue(normalizeSettingsCatalogJSON(ctx, data.Settings.ValueString(), resp))
}

// normalizeSettingsCatalogJSONArray normalizes a Microsoft Graph settings catalog array response
// into the format expected by Terraform state. This function is used typicallyy for the
// microsoft365_graph_beta_device_and_app_management_settings_catalog resource type which supports
// multiple settings configured as an array.
//
// The function performs the following steps:
// 1. Unmarshals the original settings from HCL to preserve secret values
// 2. Processes the raw API response to extract settings content
// 3. Handles various response formats (direct, value-wrapped, or settings-wrapped)
// 4. Creates a structured format with a "settings" array
// 5. Preserves any secret values from the original configuration
// 6. Normalizes the JSON alphabetically for consistent state representation
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - settingsStr: Original settings JSON string from Terraform configuration
//   - resp: Raw response bytes from the Graph API
//
// Returns:
//   - A normalized JSON string containing the settings array
//
// In case of errors during processing, returns the original settings string
// to maintain the existing state rather than potentially corrupting it.
func normalizeSettingsCatalogJSONArray(ctx context.Context, settingsStr string, resp []byte) string {
	// Parse the API response
	var rawResponse map[string]any
	if err := json.Unmarshal(resp, &rawResponse); err != nil {
		var arrayResponse []any
		if err := json.Unmarshal(resp, &arrayResponse); err != nil {
			tflog.Error(ctx, "Failed to unmarshal settings response", map[string]any{"error": err.Error()})
			return settingsStr
		}
		rawResponse = map[string]any{"value": arrayResponse}
	}

	var settingsContent any
	if value, ok := rawResponse["value"]; ok {
		settingsContent = value
	} else if settings, ok := rawResponse["settings"]; ok {
		settingsContent = settings
	} else {
		settingsContent = rawResponse
	}

	structuredContent := map[string]any{
		"settings": settingsContent,
	}

	// Preserve secrets if they exist else continue without secret preservation
	// rather than failing. User will need to manually update the settings to include the secrets.
	// they arent returned by the API.
	var configSettings map[string]any

	if settingsStr != "" {
		if err := json.Unmarshal([]byte(settingsStr), &configSettings); err != nil {
			tflog.Warn(ctx, "Failed to unmarshal config settings, skipping secret preservation", map[string]any{"error": err.Error()})
			configSettings = nil
		}
	}

	if err := normalize.PreserveSecretSettings(configSettings, structuredContent); err != nil {
		tflog.Error(ctx, "Error preserving secret settings", map[string]any{"error": err.Error()})

	}

	jsonBytes, err := json.Marshal(structuredContent)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal JSON structured content", map[string]any{"error": err.Error()})
		return settingsStr
	}

	normalizedJSON, err := normalize.JSONAlphabetically(string(jsonBytes))
	if err != nil {
		tflog.Error(ctx, "Failed to normalize settings catalog JSON", map[string]any{"error": err.Error()})
		return settingsStr
	}

	tflog.Debug(ctx, "Settings JSON normalized", map[string]any{
		"original":   string(resp),
		"normalized": normalizedJSON,
	})

	return normalizedJSON
}

// normalizeSettingsCatalogJSON handles JSON normalization for single setting instances.
// This function is designed for processing individual settings, typically used in reusable
// policy settings where only one setting instance is configured at a time.
//
// The function performs the following steps:
// 1. Unmarshals the API response to extract the setting instance
// 2. Wraps the single setting instance in the expected format with ID "0"
// 3. Creates a consistent structure matching the array format
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - settingsStr: Original settings JSON string from Terraform configuration
//   - resp: Raw response bytes from the Graph API
//
// Returns:
//   - JSON string with the setting instance wrapped in the expected format:
//     {"settings":[{"id":"0","settingInstance":{...}}]}
//
// In case of errors during processing, returns the original settings string
// to maintain the existing state rather than potentially corrupting it.
func normalizeSettingsCatalogJSON(ctx context.Context, settingsStr string, resp []byte) string {
	var responseObj map[string]any
	if err := json.Unmarshal(resp, &responseObj); err != nil {
		tflog.Error(ctx, "Failed to unmarshal response", map[string]any{"error": err.Error()})
		return settingsStr
	}

	// If we have a settingInstance, wrap it in our expected format
	if settingInstance, ok := responseObj["settingInstance"].(map[string]any); ok {
		wrappedResp := map[string]any{
			"settings": []any{
				map[string]any{
					"id":              "0",
					"settingInstance": settingInstance,
				},
			},
		}

		if newJSON, err := json.Marshal(wrappedResp); err == nil {
			return string(newJSON)
		}
	}
	return settingsStr
}
