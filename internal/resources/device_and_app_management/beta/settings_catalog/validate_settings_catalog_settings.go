package graphBetaSettingsCatalog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// settingsCatalogValidator validates settings catalog json structure
type settingsCatalogValidator struct{}

// SettingsCatalogValidator returns a validator which ensures settings catalog json is valid
func SettingsCatalogValidator() validator.String {
	return &settingsCatalogValidator{}
}

// Description describes the validation in plain text formatting.
func (v settingsCatalogValidator) Description(_ context.Context) string {
	return "validates settings catalog configuration"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v settingsCatalogValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v settingsCatalogValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Parse the JSON into a generic map structure
	var jsonData interface{}
	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &jsonData); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Settings Catalog JSON",
			fmt.Sprintf("Error parsing settings catalog JSON: %s", err),
		)
		return
	}

	// Validate all secret settings in the JSON structure
	validateSecretSettings(req.Path, jsonData, resp)
}

// validateSecretSettings recursively searches through the JSON structure for secret settings
func validateSecretSettings(path path.Path, data interface{}, resp *validator.StringResponse) {
	switch v := data.(type) {
	case map[string]interface{}:
		// Check if this is a secret setting value
		if isSecretSetting(v) {
			validateSecretSettingState(path, v, resp)
		}
		// Recursively check all map values
		for _, value := range v {
			validateSecretSettings(path, value, resp)
		}
	case []interface{}:
		// Recursively check all array elements
		for _, elem := range v {
			validateSecretSettings(path, elem, resp)
		}
	}
}

// isSecretSetting checks if the current map represents a secret setting
func isSecretSetting(m map[string]interface{}) bool {
	odataType, ok := m["@odata.type"].(string)
	return ok && odataType == "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
}

// validateSecretSettingState checks the valueState of a secret setting
func validateSecretSettingState(path path.Path, secretSetting map[string]interface{}, resp *validator.StringResponse) {
	const (
		expectedState = "notEncrypted"
		invalidState  = "encryptedValueToken"
	)

	valueState, ok := secretSetting["valueState"].(string)
	if !ok {
		return // valueState is not present or not a string
	}

	if valueState == invalidState {
		// Get the settingDefinitionId if available (walking up the structure)
		settingId := findSettingDefinitionId(secretSetting)

		errorMsg := fmt.Sprintf("Secret Setting Value state must be '%s' when setting a new secret value, got '%s'",
			expectedState, invalidState)
		if settingId != "" {
			errorMsg = fmt.Sprintf("Secret Setting Value (settingDefinitionId: %s) state must be '%s' when setting a new secret value, got '%s'",
				settingId, expectedState, invalidState)
		}

		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Secret Setting Value State",
			errorMsg,
		)
	}
}

// findSettingDefinitionId attempts to find the settingDefinitionId associated with a secret setting
func findSettingDefinitionId(m map[string]interface{}) string {
	// Check if settingDefinitionId is directly in the parent object
	if id, ok := m["settingDefinitionId"].(string); ok {
		return id
	}

	// If not found and there's a parent property, check the parent
	if parent, ok := m["parent"].(map[string]interface{}); ok {
		if id, ok := parent["settingDefinitionId"].(string); ok {
			return id
		}
	}

	return ""
}
