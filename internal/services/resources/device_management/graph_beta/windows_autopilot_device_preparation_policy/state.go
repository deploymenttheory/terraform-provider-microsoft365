package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapResourceToState maps the resource data to the state model
func mapResourceToState(ctx context.Context, stateModel *WindowsAutopilotDevicePreparationPolicyResourceModel, resource models.DeviceManagementConfigurationPolicyable) {
	if resource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(resource.GetId()),
	})

	stateModel.ID = state.StringPointerValue(resource.GetId())
	stateModel.Name = state.StringPointerValue(resource.GetName())
	stateModel.Description = state.StringPointerValue(resource.GetDescription())
	stateModel.IsAssigned = state.BoolPointerValue(resource.GetIsAssigned())
	stateModel.CreatedDateTime = state.TimeToString(resource.GetCreatedDateTime())
	stateModel.LastModifiedDateTime = state.TimeToString(resource.GetLastModifiedDateTime())
	stateModel.SettingsCount = state.Int32PtrToTypeInt64(resource.GetSettingCount())
	stateModel.RoleScopeTagIds = state.StringSliceToSet(ctx, resource.GetRoleScopeTagIds())

	// Map platform and technologies
	if platforms := resource.GetPlatforms(); platforms != nil {
		stateModel.Platforms = types.StringValue(platforms.String())
	}

	if technologies := resource.GetTechnologies(); technologies != nil {
		stateModel.Technologies = types.StringValue(technologies.String())
	}

	if templateRef := resource.GetTemplateReference(); templateRef != nil {
		if templateId := templateRef.GetTemplateId(); templateId != nil {
			stateModel.TemplateId = state.StringPointerValue(templateRef.GetTemplateId())
		}

		if templateFamily := templateRef.GetTemplateFamily(); templateFamily != nil {
			stateModel.TemplateFamily = types.StringValue(templateFamily.String())
		}
	}

	// Initialize nested objects
	if stateModel.DeploymentSettings == nil {
		stateModel.DeploymentSettings = &DeploymentSettingsModel{}
	}
	if stateModel.OOBESettings == nil {
		stateModel.OOBESettings = &OOBESettingsModel{}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource state with id %s", stateModel.ID.ValueString()))
}

// mapSettingsToState extracts settings from the response and maps them to the state model
func mapSettingsToState(ctx context.Context, stateModel *WindowsAutopilotDevicePreparationPolicyResourceModel, settingsResponse models.DeviceManagementConfigurationSettingCollectionResponseable) error {
	if settingsResponse == nil {
		tflog.Debug(ctx, "Settings response is nil")
		return nil
	}

	settings := settingsResponse.GetValue()
	if len(settings) == 0 {
		tflog.Debug(ctx, "No settings found in response")
		return nil
	}

	// Initialize the nested objects if needed
	if stateModel.DeploymentSettings == nil {
		stateModel.DeploymentSettings = &DeploymentSettingsModel{}
	}
	if stateModel.OOBESettings == nil {
		stateModel.OOBESettings = &OOBESettingsModel{}
	}

	tflog.Debug(ctx, fmt.Sprintf("Processing %d settings", len(settings)))

	// Process each setting
	for _, setting := range settings {
		if setting == nil {
			continue
		}

		settingInstance := setting.GetSettingInstance()
		if settingInstance == nil {
			continue
		}

		settingDefinitionId := settingInstance.GetSettingDefinitionId()
		if settingDefinitionId == nil {
			continue
		}

		// Get the OData type to determine how to handle the setting
		odataType := ""
		if settingInstance.GetOdataType() != nil {
			odataType = *settingInstance.GetOdataType()
		}

		tflog.Debug(ctx, fmt.Sprintf("Processing setting: %s, type: %s", *settingDefinitionId, odataType))

		// Process the setting based on its definition ID
		switch *settingDefinitionId {
		// Device Security Group
		case "enrollment_autopilot_dpp_devicegroup":
			extractStringValue(ctx, settingInstance, &stateModel.DeviceSecurityGroup)

		// Deployment Settings
		case "enrollment_autopilot_dpp_deploymentmode":
			extractChoiceValue(ctx, settingInstance, &stateModel.DeploymentSettings.DeploymentMode)
		case "enrollment_autopilot_dpp_deploymenttype":
			extractChoiceValue(ctx, settingInstance, &stateModel.DeploymentSettings.DeploymentType)
		case "enrollment_autopilot_dpp_jointype":
			extractChoiceValue(ctx, settingInstance, &stateModel.DeploymentSettings.JoinType)
		case "enrollment_autopilot_dpp_accountype":
			extractChoiceValue(ctx, settingInstance, &stateModel.DeploymentSettings.AccountType)

		// OOBE Settings
		case "enrollment_autopilot_dpp_timeout":
			extractIntValue(ctx, settingInstance, &stateModel.OOBESettings.TimeoutInMinutes)
		case "enrollment_autopilot_dpp_custonerror":
			extractStringValue(ctx, settingInstance, &stateModel.OOBESettings.CustomErrorMessage)
		case "enrollment_autopilot_dpp_allowskip":
			extractBoolValue(ctx, settingInstance, &stateModel.OOBESettings.AllowSkip)
		case "enrollment_autopilot_dpp_allowdiagnostics":
			extractBoolValue(ctx, settingInstance, &stateModel.OOBESettings.AllowDiagnostics)

		// Allowed Apps and Scripts
		case "enrollment_autopilot_dpp_allowedappids", "enrollment_autopilot_dpp_allowedapps":
			extractCollectionValue(ctx, settingInstance, &stateModel.AllowedApps)
		case "enrollment_autopilot_dpp_allowedscriptids", "enrollment_autopilot_dpp_allowedscripts":
			extractCollectionValue(ctx, settingInstance, &stateModel.AllowedScripts)
		default:
			tflog.Debug(ctx, fmt.Sprintf("Unknown setting definition ID: %s", *settingDefinitionId))
		}
	}

	return nil
}

// extractStringValue extracts a string value using the additional data property
func extractStringValue(ctx context.Context, settingInstance models.DeviceManagementConfigurationSettingInstanceable, target *types.String) {
	if settingInstance == nil {
		tflog.Warn(ctx, "Setting instance is nil when extracting string value")
		return
	}

	// First try to access via strongly typed methods
	if simpleInstance, ok := settingInstance.(models.DeviceManagementConfigurationSimpleSettingInstanceable); ok {
		simpleValue := simpleInstance.GetSimpleSettingValue()
		if simpleValue != nil {
			// Try strongly typed conversion
			if stringValue, ok := simpleValue.(models.DeviceManagementConfigurationStringSettingValueable); ok {
				if stringVal := stringValue.GetValue(); stringVal != nil {
					*target = types.StringValue(*stringVal)
					return
				}
			}
		}
	}

	// Fall back to using additionalData as a map
	additionalData := settingInstance.GetAdditionalData()
	if additionalData == nil {
		tflog.Warn(ctx, "No additional data when extracting string value")
		return
	}

	// Try to extract the string value from the additional data
	if simpleValue, ok := additionalData["simpleSettingValue"]; ok {
		if valueMap, ok := simpleValue.(map[string]interface{}); ok {
			if stringValue, ok := valueMap["value"].(string); ok {
				*target = types.StringValue(stringValue)
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract string value")
}

// extractIntValue extracts an integer value using the additional data property
func extractIntValue(ctx context.Context, settingInstance models.DeviceManagementConfigurationSettingInstanceable, target *types.Int64) {
	if settingInstance == nil {
		tflog.Warn(ctx, "Setting instance is nil when extracting int value")
		return
	}

	// First try to access via strongly typed methods
	if simpleInstance, ok := settingInstance.(models.DeviceManagementConfigurationSimpleSettingInstanceable); ok {
		simpleValue := simpleInstance.GetSimpleSettingValue()
		if simpleValue != nil {
			// Try strongly typed conversion
			if intValue, ok := simpleValue.(models.DeviceManagementConfigurationIntegerSettingValueable); ok {
				if intVal := intValue.GetValue(); intVal != nil {
					*target = types.Int64Value(int64(*intVal))
					return
				}
			}
		}
	}

	// Fall back to using additionalData as a map
	additionalData := settingInstance.GetAdditionalData()
	if additionalData == nil {
		tflog.Warn(ctx, "No additional data when extracting int value")
		return
	}

	// Try to extract the integer value from the additional data
	if simpleValue, ok := additionalData["simpleSettingValue"]; ok {
		if valueMap, ok := simpleValue.(map[string]interface{}); ok {
			if numValue, ok := valueMap["value"].(float64); ok {
				*target = types.Int64Value(int64(numValue))
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract int value")
}

// extractBoolValue extracts a boolean value using the additional data property
func extractBoolValue(ctx context.Context, settingInstance models.DeviceManagementConfigurationSettingInstanceable, target *types.Bool) {
	if settingInstance == nil {
		tflog.Warn(ctx, "Setting instance is nil when extracting bool value")
		return
	}

	// The SDK doesn't seem to have a specific DeviceManagementConfigurationBooleanSettingValueable type
	// So we primarily rely on the additionalData approach

	// Access the setting's backing data as a map
	additionalData := settingInstance.GetAdditionalData()
	if additionalData == nil {
		tflog.Warn(ctx, "No additional data when extracting bool value")
		return
	}

	// Try to extract the boolean value from the additional data
	if simpleValue, ok := additionalData["simpleSettingValue"]; ok {
		if valueMap, ok := simpleValue.(map[string]interface{}); ok {
			if boolValue, ok := valueMap["value"].(bool); ok {
				*target = types.BoolValue(boolValue)
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract bool value")
}

// extractChoiceValue extracts a choice value using the additional data property
func extractChoiceValue(ctx context.Context, settingInstance models.DeviceManagementConfigurationSettingInstanceable, target *types.String) {
	if settingInstance == nil {
		tflog.Warn(ctx, "Setting instance is nil when extracting choice value")
		return
	}

	// First try to access via strongly typed methods
	if choiceInstance, ok := settingInstance.(models.DeviceManagementConfigurationChoiceSettingInstanceable); ok {
		choiceValue := choiceInstance.GetChoiceSettingValue()
		if choiceValue != nil {
			if choiceVal := choiceValue.GetValue(); choiceVal != nil {
				*target = types.StringValue(*choiceVal)
				return
			}
		}
	}

	// Fall back to using additionalData as a map
	additionalData := settingInstance.GetAdditionalData()
	if additionalData == nil {
		tflog.Warn(ctx, "No additional data when extracting choice value")
		return
	}

	// Try to extract the choice value from the additional data
	if choiceValue, ok := additionalData["choiceSettingValue"]; ok {
		if valueMap, ok := choiceValue.(map[string]interface{}); ok {
			if stringValue, ok := valueMap["value"].(string); ok {
				*target = types.StringValue(stringValue)
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract choice value")
}

// extractCollectionValue extracts collection values using the additional data property
func extractCollectionValue(ctx context.Context, settingInstance models.DeviceManagementConfigurationSettingInstanceable, target interface{}) {
	if settingInstance == nil {
		tflog.Warn(ctx, "Setting instance is nil when extracting collection value")
		return
	}

	// Determine the type of the target to decide how to extract the data
	switch typedTarget := target.(type) {
	case *[]AllowedAppModel:
		// Handle app collections with ID and type
		extractAllowedAppsCollection(ctx, settingInstance, typedTarget)
	case *[]types.String:
		// Handle simple string collections
		extractSimpleStringCollection(ctx, settingInstance, typedTarget)
	default:
		tflog.Warn(ctx, "Unsupported target type for collection extraction")
	}
}

// extractAllowedAppsCollection extracts app collections with ID and type
func extractAllowedAppsCollection(ctx context.Context, settingInstance models.DeviceManagementConfigurationSettingInstanceable, target *[]AllowedAppModel) {
	if settingInstance == nil || target == nil {
		tflog.Warn(ctx, "Setting instance or target is nil when extracting app collection")
		return
	}

	// First try to access via strongly typed methods
	if collectionInstance, ok := settingInstance.(models.DeviceManagementConfigurationSimpleSettingCollectionInstanceable); ok {
		collectionValues := collectionInstance.GetSimpleSettingCollectionValue()
		if len(collectionValues) > 0 {
			var apps []AllowedAppModel
			for _, collectionValue := range collectionValues {
				if stringValue, ok := collectionValue.(models.DeviceManagementConfigurationStringSettingValueable); ok {
					if stringVal := stringValue.GetValue(); stringVal != nil {
						// Parse the JSON string to extract app ID and type
						// Format is {"id":"GUID","type":"#microsoft.graph.TYPE"}
						app := parseAppJson(ctx, *stringVal)
						if !app.AppID.IsNull() {
							apps = append(apps, app)
						}
					}
				}
			}

			if len(apps) > 0 {
				*target = apps
				return
			}
		}
	}

	// Fall back to using additionalData as a map
	additionalData := settingInstance.GetAdditionalData()
	if additionalData == nil {
		tflog.Warn(ctx, "No additional data when extracting app collection")
		return
	}

	// Try collection values directly
	if collectionValues, ok := additionalData["simpleSettingCollectionValue"]; ok {
		if collectionArray, ok := collectionValues.([]interface{}); ok {
			var apps []AllowedAppModel
			for _, item := range collectionArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if value, ok := itemMap["value"].(string); ok {
						app := parseAppJson(ctx, value)
						if !app.AppID.IsNull() {
							apps = append(apps, app)
						}
					}
				}
			}

			if len(apps) > 0 {
				*target = apps
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract app collection")
}

// parseAppJson parses a JSON string containing app ID and type
func parseAppJson(ctx context.Context, jsonStr string) AllowedAppModel {
	app := AllowedAppModel{
		AppID:   types.StringNull(),
		AppType: types.StringNull(),
	}

	// Simple parsing using string manipulation for robustness
	// Format is {"id":"GUID","type":"#microsoft.graph.TYPE"}
	idStart := strings.Index(jsonStr, "\"id\":\"")
	typeStart := strings.Index(jsonStr, "\"type\":\"")

	if idStart >= 0 && typeStart >= 0 {
		// Extract ID
		idStart += 6 // length of "\"id\":\""
		idEnd := strings.Index(jsonStr[idStart:], "\"")
		if idEnd > 0 {
			appId := jsonStr[idStart : idStart+idEnd]
			app.AppID = types.StringValue(appId)
		}

		// Extract Type
		typeStart += 8 // length of "\"type\":\""
		typeEnd := strings.Index(jsonStr[typeStart:], "\"")
		if typeEnd > 0 {
			fullType := jsonStr[typeStart : typeStart+typeEnd]
			// Remove "#microsoft.graph." prefix
			if strings.HasPrefix(fullType, "#microsoft.graph.") {
				appType := fullType[17:] // length of "#microsoft.graph."
				app.AppType = types.StringValue(appType)
			} else {
				app.AppType = types.StringValue(fullType)
			}
		}
	}

	return app
}

// extractSimpleStringCollection extracts simple string collections
func extractSimpleStringCollection(ctx context.Context, settingInstance models.DeviceManagementConfigurationSettingInstanceable, target *[]types.String) {
	if settingInstance == nil || target == nil {
		tflog.Warn(ctx, "Setting instance or target is nil when extracting string collection")
		return
	}

	// First try to access via strongly typed methods
	if collectionInstance, ok := settingInstance.(models.DeviceManagementConfigurationSimpleSettingCollectionInstanceable); ok {
		collectionValues := collectionInstance.GetSimpleSettingCollectionValue()
		if len(collectionValues) > 0 {
			var values []string
			for _, collectionValue := range collectionValues {
				if stringValue, ok := collectionValue.(models.DeviceManagementConfigurationStringSettingValueable); ok {
					if stringVal := stringValue.GetValue(); stringVal != nil {
						values = append(values, *stringVal)
					}
				}
			}

			if len(values) > 0 {
				*target = state.SliceToTypeStringSlice(values)
				return
			}
		}
	}

	// Fall back to using additionalData as a map
	additionalData := settingInstance.GetAdditionalData()
	if additionalData == nil {
		tflog.Warn(ctx, "No additional data when extracting string collection")
		return
	}

	// For collection settings, we may have stored them as comma-separated strings
	if simpleValue, ok := additionalData["simpleSettingValue"]; ok {
		if valueMap, ok := simpleValue.(map[string]interface{}); ok {
			if stringValue, ok := valueMap["value"].(string); ok {
				values := splitCommaSeparatedString(stringValue)
				*target = state.SliceToTypeStringSlice(values)
				return
			}
		}
	}

	// Also try collection values directly
	if collectionValues, ok := additionalData["simpleSettingCollectionValue"]; ok {
		if collectionArray, ok := collectionValues.([]interface{}); ok {
			var values []string
			for _, item := range collectionArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if value, ok := itemMap["value"].(string); ok {
						values = append(values, value)
					}
				}
			}

			if len(values) > 0 {
				*target = state.SliceToTypeStringSlice(values)
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract string collection")
}

// splitCommaSeparatedString splits a comma-separated string into a slice of strings
func splitCommaSeparatedString(s string) []string {
	if s == "" {
		return []string{}
	}

	// Split the string by commas and return the result
	return strings.Split(s, ",")
}
