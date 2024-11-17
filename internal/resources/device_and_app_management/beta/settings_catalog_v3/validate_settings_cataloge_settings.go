package graphBetaSettingsCatalog

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
// )

// // settingsCatalogValidator validates settings catalog json structure
// type settingsCatalogValidator struct{}

// // SettingsCatalogValidator returns a validator which ensures settings catalog json is valid
// func SettingsCatalogValidator() validator.String {
// 	return &settingsCatalogValidator{}
// }

// // Description describes the validation in plain text formatting.
// func (v settingsCatalogValidator) Description(_ context.Context) string {
// 	return "validates settings catalog configuration"
// }

// // MarkdownDescription describes the validation in Markdown formatting.
// func (v settingsCatalogValidator) MarkdownDescription(ctx context.Context) string {
// 	return v.Description(ctx)
// }

// // Validate performs the validation.
// // Validate performs the validation.
// func (v settingsCatalogValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
// 	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
// 		return
// 	}

// 	// Parse JSON directly into the existing settingsCatalogModel var
// 	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &SettingsCatalogModel); err != nil {
// 		resp.Diagnostics.AddAttributeError(
// 			req.Path,
// 			"Invalid Settings Catalog JSON",
// 			fmt.Sprintf("Error parsing settings catalog JSON: %s", err),
// 		)
// 		return
// 	}

// 	// Helper function to check secret setting value state
// 	checkSecretValueState := func(path string, i int, j *int, valueState string) {
// 		expectedState := "notEncrypted"
// 		if valueState != expectedState {
// 			errorMsg := fmt.Sprintf("Setting %d: For Secret Setting Value got '%s', expected '%s'. You must use 'notEncrypted' when trying to set a new secret value.", i, valueState, expectedState)
// 			if j != nil {
// 				errorMsg = fmt.Sprintf("Setting %d, Child %d: For Secret Setting Value got '%s', expected '%s'. You must use 'notEncrypted' when trying to set a new secret value.", i, *j, valueState, expectedState)
// 			}
// 			resp.Diagnostics.AddAttributeError(
// 				path,
// 				"Invalid Secret Setting Value State",
// 				errorMsg,
// 			)
// 		}
// 	}

// 	// Validate each setting
// 	for i, setting := range SettingsCatalogModel.SettingsDetails {
// 		// Check root level SimpleSettingValue
// 		if setting.SettingInstance.SimpleSettingValue != nil &&
// 			setting.SettingInstance.SimpleSettingValue.ODataType == "#microsoft.graph.deviceManagementConfigurationSecretSettingValue" {
// 			checkSecretValueState(req.Path, i, nil, setting.SettingInstance.SimpleSettingValue.ValueState)
// 		}

// 		// Check ChoiceSettingValue children
// 		if setting.SettingInstance.ChoiceSettingValue != nil {
// 			for j, child := range setting.SettingInstance.ChoiceSettingValue.Children {
// 				if child.SimpleSettingValue != nil &&
// 					child.SimpleSettingValue.ODataType == "#microsoft.graph.deviceManagementConfigurationSecretSettingValue" {
// 					jCopy := j // Create copy for pointer
// 					checkSecretValueState(req.Path, i, &jCopy, child.SimpleSettingValue.ValueState)
// 				}
// 			}
// 		}

// 		// Check GroupSettingCollectionValue
// 		if setting.SettingInstance.GroupSettingCollectionValue != nil {
// 			for j, group := range setting.SettingInstance.GroupSettingCollectionValue {
// 				for k, child := range group.Children {
// 					if child.SimpleSettingValue != nil &&
// 						child.SimpleSettingValue.ODataType == "#microsoft.graph.deviceManagementConfigurationSecretSettingValue" {
// 						kCopy := k
// 						checkSecretValueState(req.Path, i, &kCopy, child.SimpleSettingValue.ValueState)
// 					}
// 				}
// 			}
// 		}
// 	}
// }
