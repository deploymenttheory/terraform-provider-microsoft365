package sharedValidators

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

	// Validate settings catalog ID sequence and initial ID value is 0
	validateSettingsCatalogIDSequences(req.Path, jsonData, resp)

	// validate the settings hierarchy
	validateSettingsHierarchy(req.Path, jsonData, resp)

	// validateSettingsTemplates for settings templates existance
	validateSettingsTemplates(req.Path, jsonData, resp)

}

//----------------------------------------------------------------------------------------------//

// validateSecretSettings recursively searches through the JSON structure for secret settings
func validateSecretSettings(path path.Path, data interface{}, resp *validator.StringResponse) {
	switch v := data.(type) {
	case map[string]interface{}:
		if isSecretSetting(v) {
			validateSecretSettingState(path, v, resp)
		}
		for _, value := range v {
			validateSecretSettings(path, value, resp)
		}
	case []interface{}:
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

	if id, ok := m["settingDefinitionId"].(string); ok {
		return id
	}

	if parent, ok := m["parent"].(map[string]interface{}); ok {
		if id, ok := parent["settingDefinitionId"].(string); ok {
			return id
		}
	}

	return ""
}

//----------------------------------------------------------------------------------------------//

// validateSettingsCatalogIDSequences validates that settings IDs start at 0 and increment sequentially
// It validates that all settings have an ID and that the IDs are sequential, or the correct field value type
func validateSettingsCatalogIDSequences(path path.Path, data interface{}, resp *validator.StringResponse) {

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	settingsDetails, ok := dataMap["settings"].([]interface{})
	if !ok || len(settingsDetails) == 0 {
		return
	}

	// First, verify that ALL settings have IDs and they are numeric
	for i, setting := range settingsDetails {
		settingMap, ok := setting.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if ID exists
		id, exists := settingMap["id"]
		if !exists {
			resp.Diagnostics.AddAttributeError(
				path,
				"Missing Settings ID",
				fmt.Sprintf("Setting at index %d is missing required 'id' field", i),
			)
			return
		}

		// Check if ID is a string
		idStr, ok := id.(string)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				path,
				"Invalid Settings ID Format",
				fmt.Sprintf("Settings ID must be a string: %v", id),
			)
			return
		}

		// Validate ID is a number
		if _, err := strconv.Atoi(idStr); err != nil {
			resp.Diagnostics.AddAttributeError(
				path,
				"Invalid Settings ID Format",
				fmt.Sprintf("Settings ID must be a number: %s", idStr),
			)
			return
		}
	}

	// Validate sequential ordering
	for i := 1; i < len(settingsDetails); i++ {
		setting, ok := settingsDetails[i].(map[string]interface{})
		if !ok {
			continue
		}

		currentID, ok := setting["id"].(string)
		if !ok {
			continue
		}

		prevSetting, ok := settingsDetails[i-1].(map[string]interface{})
		if !ok {
			continue
		}

		previousID, ok := prevSetting["id"].(string)
		if !ok {
			continue
		}

		// We know these are valid numbers from the first validation loop
		curr, _ := strconv.Atoi(currentID)
		prev, _ := strconv.Atoi(previousID)

		if curr != prev+1 {
			resp.Diagnostics.AddAttributeError(
				path,
				"Non-sequential Settings ID",
				fmt.Sprintf("Settings IDs must increment by 1. Found ID %d after ID %d", curr, prev),
			)
			return
		}
	}
}

//----------------------------------------------------------------------------------------------//

// validateSettingsHierarchy validates that settingsDetails entries follow the required structure and ordering
func validateSettingsHierarchy(path path.Path, data interface{}, resp *validator.StringResponse) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	settingsDetails, ok := dataMap["settings"].([]interface{})
	if !ok || len(settingsDetails) == 0 {
		return
	}

	for i, setting := range settingsDetails {
		// Get original JSON structure preserving order
		jsonBytes, err := json.Marshal(setting)
		if err != nil {
			continue
		}

		fieldOrder, keyValuePairs := extractFieldOrderAndPairs(jsonBytes)
		if len(fieldOrder) == 0 {
			continue // Skip if we couldn't parse the JSON structure
		}

		// Validate field count first
		if err := validateFieldCount(fieldOrder, keyValuePairs, i, path, resp); err != nil {
			return
		}

		// Then validate field order
		if err := validateFieldOrder(fieldOrder, keyValuePairs, i, path, resp); err != nil {
			return
		}
	}
}

// extractFieldOrderAndPairs reads the JSON and returns both the field order and formatted key-value pairs
func extractFieldOrderAndPairs(jsonBytes []byte) ([]string, []string) {
	dec := json.NewDecoder(bytes.NewReader(jsonBytes))
	var fieldOrder []string
	var keyValuePairs []string

	// Start object
	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') {
		return nil, nil
	}

	// Read field,value pairs
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			break
		}
		keyStr := key.(string)
		fieldOrder = append(fieldOrder, keyStr)

		// Read value
		var value interface{}
		if err := dec.Decode(&value); err != nil {
			break
		}

		// Format key-value pair based on type
		formattedPair := formatKeyValuePair(keyStr, value)
		keyValuePairs = append(keyValuePairs, formattedPair)
	}

	return fieldOrder, keyValuePairs
}

// formatKeyValuePair handles different value types and returns a formatted string
func formatKeyValuePair(key string, value interface{}) string {
	switch val := value.(type) {
	case string:
		return fmt.Sprintf(`"%s" = "%s"`, key, val)
	case map[string]interface{}:
		return fmt.Sprintf(`"%s" = <object>`, key)
	case []interface{}:
		return fmt.Sprintf(`"%s" = <array>`, key)
	default:
		return fmt.Sprintf(`"%s" = %v`, key, value)
	}
}

// validateFieldCount checks if there are exactly 2 fields
func validateFieldCount(fieldOrder []string, keyValuePairs []string, index int, path path.Path, resp *validator.StringResponse) error {
	if len(fieldOrder) != 2 {
		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Settings Structure",
			fmt.Sprintf("Setting at index %d contains %d fields ([%s]), but should only contain exactly 2 fields ('id' and 'settingInstance')",
				index, len(fieldOrder), strings.Join(keyValuePairs, ", ")),
		)
		return fmt.Errorf("invalid field count")
	}
	return nil
}

// validateFieldOrder ensures fields are in the correct order: id then settingInstance
func validateFieldOrder(fieldOrder []string, keyValuePairs []string, index int, path path.Path, resp *validator.StringResponse) error {
	if fieldOrder[0] != "id" || fieldOrder[1] != "settingInstance" {
		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Settings Structure",
			fmt.Sprintf("Setting at index %d has incorrect field order: found [%s] but fields must be ordered exactly as: ['id', 'settingInstance']",
				index, strings.Join(keyValuePairs, ", ")),
		)
		return fmt.Errorf("invalid field order")
	}
	return nil
}

//----------------------------------------------------------------------------------------------//

// validateSettingsTemplates validates that settingTemplates are not present in the JSON structure
func validateSettingsTemplates(path path.Path, data interface{}, resp *validator.StringResponse) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Check if settingTemplates exists in the JSON
	if templates, exists := dataMap["settingTemplates"]; exists {
		// If it exists, check if it's an array (which is the expected type if present)
		if templatesArray, ok := templates.([]interface{}); ok && len(templatesArray) > 0 {
			resp.Diagnostics.AddAttributeError(
				path,
				"Invalid Settings Configuration",
				"Settings Templates are not supported in this configuration. Please provide only settings catalog settings  with a 'settings' array and remove the 'settingTemplates' field.",
			)
			return
		}
	}
}
