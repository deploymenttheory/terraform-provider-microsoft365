package sharedValidators

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// settingsCatalogValidator validates settings catalog structure
type settingsCatalogValidator struct{}

// SettingsCatalogValidator returns a validator which ensures settings catalog is valid
func SettingsCatalogValidator() validator.List {
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

// ValidateList performs the validation on the HCL settings list.
func (v settingsCatalogValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var settings []struct {
		ID              types.String `tfsdk:"id"`
		SettingInstance struct {
			ODataType           types.String `tfsdk:"odata_type"`
			SettingDefinitionId types.String `tfsdk:"setting_definition_id"`
			SimpleSettingValue  *struct {
				ODataType  types.String `tfsdk:"odata_type"`
				Value      types.String `tfsdk:"value"`
				ValueState types.String `tfsdk:"value_state"`
			} `tfsdk:"simple_setting_value"`
		} `tfsdk:"setting_instance"`
	}

	diags := req.ConfigValue.ElementsAs(ctx, &settings, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Validate settings catalog ID sequence and initial ID value is 0
	validateSettingsIDs(req.Path, settings, resp)

	// Validate the settings hierarchy
	validateSettingsStructure(req.Path, settings, resp)

	// Validate secret settings
	validateSecretSettingsHCL(ctx, req.Path, settings, resp)
}

// validateSettingsIDs validates that settings IDs start at 0 and increment sequentially
func validateSettingsIDs(path path.Path, settings []struct {
	ID              types.String `tfsdk:"id"`
	SettingInstance struct {
		ODataType           types.String `tfsdk:"odata_type"`
		SettingDefinitionId types.String `tfsdk:"setting_definition_id"`
		SimpleSettingValue  *struct {
			ODataType  types.String `tfsdk:"odata_type"`
			Value      types.String `tfsdk:"value"`
			ValueState types.String `tfsdk:"value_state"`
		} `tfsdk:"simple_setting_value"`
	} `tfsdk:"setting_instance"`
}, resp *validator.ListResponse) {
	// First, verify that ALL settings have IDs and they are numeric
	for i, setting := range settings {
		if setting.ID.IsNull() || setting.ID.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.AtListIndex(i).AtName("id"),
				"Missing Settings ID",
				fmt.Sprintf("Setting at index %d is missing required 'id' field", i),
			)
			return
		}

		idStr := setting.ID.ValueString()
		// Validate ID is a number
		if _, err := strconv.Atoi(idStr); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.AtListIndex(i).AtName("id"),
				"Invalid Settings ID Format",
				fmt.Sprintf("Settings ID must be a number: %s", idStr),
			)
			return
		}
	}

	// Validate sequential ordering
	for i := 1; i < len(settings); i++ {
		currentID, _ := strconv.Atoi(settings[i].ID.ValueString())
		previousID, _ := strconv.Atoi(settings[i-1].ID.ValueString())

		if currentID != previousID+1 {
			resp.Diagnostics.AddAttributeError(
				path.AtListIndex(i).AtName("id"),
				"Non-sequential Settings ID",
				fmt.Sprintf("Settings IDs must increment by 1. Found ID %d after ID %d", currentID, previousID),
			)
			return
		}
	}

	// Check if first ID is 0
	if len(settings) > 0 {
		firstID, _ := strconv.Atoi(settings[0].ID.ValueString())
		if firstID != 0 {
			resp.Diagnostics.AddAttributeError(
				path.AtListIndex(0).AtName("id"),
				"Invalid First Settings ID",
				fmt.Sprintf("First setting ID must be 0, got %d", firstID),
			)
			return
		}
	}
}

// validateSettingsStructure validates that settings have the required structure
func validateSettingsStructure(path path.Path, settings []struct {
	ID              types.String `tfsdk:"id"`
	SettingInstance struct {
		ODataType           types.String `tfsdk:"odata_type"`
		SettingDefinitionId types.String `tfsdk:"setting_definition_id"`
		SimpleSettingValue  *struct {
			ODataType  types.String `tfsdk:"odata_type"`
			Value      types.String `tfsdk:"value"`
			ValueState types.String `tfsdk:"value_state"`
		} `tfsdk:"simple_setting_value"`
	} `tfsdk:"setting_instance"`
}, resp *validator.ListResponse) {
	for i, setting := range settings {
		// Check that setting_instance is present and has required fields
		if setting.SettingInstance.ODataType.IsNull() || setting.SettingInstance.ODataType.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.AtListIndex(i).AtName("setting_instance").AtName("odata_type"),
				"Missing OData Type",
				fmt.Sprintf("Setting at index %d is missing required 'odata_type' field in setting_instance", i),
			)
			return
		}

		if setting.SettingInstance.SettingDefinitionId.IsNull() || setting.SettingInstance.SettingDefinitionId.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.AtListIndex(i).AtName("setting_instance").AtName("setting_definition_id"),
				"Missing Setting Definition ID",
				fmt.Sprintf("Setting at index %d is missing required 'setting_definition_id' field in setting_instance", i),
			)
			return
		}
	}
}

// validateSecretSettingsHCL validates secret settings in the HCL structure
func validateSecretSettingsHCL(ctx context.Context, path path.Path, settings []struct {
	ID              types.String `tfsdk:"id"`
	SettingInstance struct {
		ODataType           types.String `tfsdk:"odata_type"`
		SettingDefinitionId types.String `tfsdk:"setting_definition_id"`
		SimpleSettingValue  *struct {
			ODataType  types.String `tfsdk:"odata_type"`
			Value      types.String `tfsdk:"value"`
			ValueState types.String `tfsdk:"value_state"`
		} `tfsdk:"simple_setting_value"`
	} `tfsdk:"setting_instance"`
}, resp *validator.ListResponse) {
	const (
		secretSettingODataType = "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
		expectedState          = "notEncrypted"
		invalidState           = "encryptedValueToken"
	)

	for i, setting := range settings {
		// Check if this is a simple setting value with a secret type
		if setting.SettingInstance.SimpleSettingValue != nil &&
			!setting.SettingInstance.SimpleSettingValue.ODataType.IsNull() &&
			!setting.SettingInstance.SimpleSettingValue.ODataType.IsUnknown() &&
			setting.SettingInstance.SimpleSettingValue.ODataType.ValueString() == secretSettingODataType {

			// Check if value_state is set correctly for secret settings
			if !setting.SettingInstance.SimpleSettingValue.ValueState.IsNull() &&
				!setting.SettingInstance.SimpleSettingValue.ValueState.IsUnknown() &&
				setting.SettingInstance.SimpleSettingValue.ValueState.ValueString() == invalidState {

				settingId := setting.SettingInstance.SettingDefinitionId.ValueString()
				errorMsg := fmt.Sprintf("Secret Setting Value (settingDefinitionId: %s) state must be '%s' when setting a new secret value, got '%s'",
					settingId, expectedState, invalidState)

				resp.Diagnostics.AddAttributeError(
					path.AtListIndex(i).AtName("setting_instance").AtName("simple_setting_value").AtName("value_state"),
					"Invalid Secret Setting Value State",
					errorMsg,
				)
				return
			}
		}

		// Note: For a complete validator, you would need to recursively check all nested structures
		// (choice_setting_value, choice_setting_collection_value, group_setting_collection_value, etc.)
		// This implementation focuses on simple_setting_value as an example
	}
}
