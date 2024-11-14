package validators

import (
	"context"
	"encoding/json"
	"fmt"

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

	// Parse JSON
	var settingsData struct {
		SettingsDetails []struct {
			SettingInstance struct {
				ODataType          string `json:"@odata.type"`
				SimpleSettingValue *struct {
					ODataType  string `json:"@odata.type"`
					ValueState string `json:"valueState,omitempty"`
					Value      string `json:"value"`
				} `json:"simpleSettingValue,omitempty"`
			} `json:"settingInstance"`
		} `json:"settingsDetails"`
	}

	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &settingsData); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Settings Catalog JSON",
			fmt.Sprintf("Error parsing settings catalog JSON: %s", err),
		)
		return
	}

	// Validate each setting
	for i, setting := range settingsData.SettingsDetails {
		// Check for secret settings
		if setting.SettingInstance.SimpleSettingValue != nil &&
			setting.SettingInstance.SimpleSettingValue.ODataType == "#microsoft.graph.deviceManagementConfigurationSecretSettingValue" {

			valueState := setting.SettingInstance.SimpleSettingValue.ValueState
			expectedState := "notEncrypted"

			if valueState != expectedState {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid Secret Setting Value State",
					fmt.Sprintf("Setting %d: For Secret Setting Value got '%s', expected '%s'. You must use 'notEncrypted' when trying to set a new secret value.",
						i, valueState, expectedState),
				)
			}
		}
	}
}
