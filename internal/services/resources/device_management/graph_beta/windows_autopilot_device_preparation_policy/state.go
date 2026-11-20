package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errMissingMembershipTargetStatuses = errors.New(
	"enrollment-time device membership target response is missing validation statuses",
)

// mapResourceToState maps the resource data to the state model
func mapResourceToState(
	ctx context.Context,
	stateModel *WindowsAutopilotDevicePreparationPolicyResourceModel,
	resource models.DeviceManagementConfigurationPolicyable,
) {
	if resource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(resource.GetId()),
	})

	stateModel.ID = convert.GraphToFrameworkString(resource.GetId())
	stateModel.Name = convert.GraphToFrameworkString(resource.GetName())
	stateModel.Description = convert.GraphToFrameworkString(resource.GetDescription())
	stateModel.IsAssigned = convert.GraphToFrameworkBool(resource.GetIsAssigned())
	stateModel.CreatedDateTime = convert.GraphToFrameworkTime(resource.GetCreatedDateTime())
	stateModel.LastModifiedDateTime = convert.GraphToFrameworkTime(
		resource.GetLastModifiedDateTime(),
	)
	stateModel.SettingsCount = convert.GraphToFrameworkInt32(resource.GetSettingCount())
	stateModel.RoleScopeTagIds = convert.GraphToFrameworkStringSet(
		ctx,
		resource.GetRoleScopeTagIds(),
	)

	if platforms := resource.GetPlatforms(); platforms != nil {
		stateModel.Platforms = types.StringValue(platforms.String())
	}

	if technologies := resource.GetTechnologies(); technologies != nil {
		stateModel.Technologies = types.StringValue(technologies.String())
	}

	// Determine if this is user-driven mode based on template ID
	isUserDriven := false
	if templateRef := resource.GetTemplateReference(); templateRef != nil {
		if templateId := templateRef.GetTemplateId(); templateId != nil {
			stateModel.TemplateId = convert.GraphToFrameworkString(templateRef.GetTemplateId())
			isUserDriven = *templateId == TemplateIDUserDriven
		}

		if templateFamily := templateRef.GetTemplateFamily(); templateFamily != nil {
			stateModel.TemplateFamily = types.StringValue(templateFamily.String())
		}
	}

	// Initialize deployment_settings for all modes (needed for deployment_type)
	if stateModel.DeploymentSettings == nil {
		stateModel.DeploymentSettings = &DeploymentSettingsModel{}
	}

	// Initialize OOBE settings only for user-driven policies
	if isUserDriven {
		if stateModel.OOBESettings == nil {
			stateModel.OOBESettings = &OOBESettingsModel{}
		}
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Finished mapping resource state with id %s", stateModel.ID.ValueString()),
	)
}

// mapSettingsToState extracts settings from the response and maps them to the state model
func mapSettingsToState(
	ctx context.Context,
	stateModel *WindowsAutopilotDevicePreparationPolicyResourceModel,
	settingsResponse models.DeviceManagementConfigurationSettingCollectionResponseable,
) error {
	if settingsResponse == nil {
		tflog.Debug(ctx, "Settings response is nil")
		return nil
	}

	settings := settingsResponse.GetValue()
	if len(settings) == 0 {
		tflog.Debug(ctx, "No settings found in response")
		return nil
	}

	// Determine if this is user-driven mode based on template ID from state
	isUserDriven := false
	if !stateModel.TemplateId.IsNull() && !stateModel.TemplateId.IsUnknown() {
		isUserDriven = stateModel.TemplateId.ValueString() == TemplateIDUserDriven
	}

	// Initialize deployment_settings for all modes (needed for deployment_type)
	if stateModel.DeploymentSettings == nil {
		stateModel.DeploymentSettings = &DeploymentSettingsModel{}
	}

	// Initialize OOBE settings only for user-driven policies
	if isUserDriven {
		if stateModel.OOBESettings == nil {
			stateModel.OOBESettings = &OOBESettingsModel{}
		}
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

		tflog.Debug(
			ctx,
			fmt.Sprintf("Processing setting: %s, type: %s", *settingDefinitionId, odataType),
		)

		// Process the setting based on its definition ID
		switch *settingDefinitionId {
		// Device Security Group is configured through the enrollment-time target action.
		case "enrollment_autopilot_dpp_devicegroup":
			extractStringValue(ctx, settingInstance, &stateModel.DeviceSecurityGroup)

		// Deployment Settings
		case "enrollment_autopilot_dpp_deploymenttype":
			// deployment_type is always mapped for both modes
			extractChoiceValue(ctx, settingInstance, &stateModel.DeploymentSettings.DeploymentType)
		case "enrollment_autopilot_dpp_deploymentmode":
			// deployment_mode only for user-driven
			if isUserDriven {
				extractChoiceValue(
					ctx,
					settingInstance,
					&stateModel.DeploymentSettings.DeploymentMode,
				)
			}
		case "enrollment_autopilot_dpp_jointype":
			// join_type only for user-driven
			if isUserDriven {
				extractChoiceValue(ctx, settingInstance, &stateModel.DeploymentSettings.JoinType)
			}
		case "enrollment_autopilot_dpp_accountype":
			// account_type only for user-driven
			if isUserDriven {
				extractChoiceValue(ctx, settingInstance, &stateModel.DeploymentSettings.AccountType)
			}

		// OOBE Settings (user-driven only)
		case "enrollment_autopilot_dpp_timeout":
			if isUserDriven {
				extractIntValue(ctx, settingInstance, &stateModel.OOBESettings.TimeoutInMinutes)
			}
		case SettingDefCustomErrorMessage:
			if isUserDriven {
				extractStringValue(
					ctx,
					settingInstance,
					&stateModel.OOBESettings.CustomErrorMessage,
				)
			}
		case "enrollment_autopilot_dpp_allowskip":
			if isUserDriven {
				extractBoolValue(ctx, settingInstance, &stateModel.OOBESettings.AllowSkip)
			}
		case "enrollment_autopilot_dpp_allowdiagnostics":
			if isUserDriven {
				extractBoolValue(ctx, settingInstance, &stateModel.OOBESettings.AllowDiagnostics)
			}

		// Allowed Apps and Scripts (both policy types)
		case "enrollment_autopilot_dpp_allowedappids", "enrollment_autopilot_dpp_allowedapps":
			extractCollectionValue(ctx, settingInstance, &stateModel.AllowedApps)
		case "enrollment_autopilot_dpp_allowedscriptids", "enrollment_autopilot_dpp_allowedscripts":
			extractCollectionValue(ctx, settingInstance, &stateModel.AllowedScripts)
		default:
			tflog.Debug(ctx, fmt.Sprintf("Unknown setting definition ID: %s", *settingDefinitionId))
		}
	}

	// For automatic/self-deploying mode, ensure computed fields that don't exist are set to null
	if !isUserDriven && stateModel.DeploymentSettings != nil {
		stateModel.DeploymentSettings.JoinType = types.StringNull()
	}

	return nil
}

// mapEnrollmentTimeDeviceMembershipTargetToState maps the enrollment-time target to Terraform state.
func mapEnrollmentTimeDeviceMembershipTargetToState(
	ctx context.Context,
	stateModel *WindowsAutopilotDevicePreparationPolicyResourceModel,
	body []byte,
) error {
	var result struct {
		Statuses []struct {
			TargetID string `json:"targetId"`
		} `json:"enrollmentTimeDeviceMembershipTargetValidationStatuses"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to decode enrollment-time device membership target: %w", err)
	}

	if result.Statuses == nil {
		return errMissingMembershipTargetStatuses
	}

	stateModel.DeviceSecurityGroup = types.StringNull()
	if len(result.Statuses) > 0 && result.Statuses[0].TargetID != "" {
		stateModel.DeviceSecurityGroup = types.StringValue(result.Statuses[0].TargetID)
	}
	tflog.Debug(ctx, "Mapped enrollment-time device membership target")
	return nil
}

// extractStringValue extracts a string value using the additional data property
func extractStringValue(
	ctx context.Context,
	settingInstance models.DeviceManagementConfigurationSettingInstanceable,
	target *types.String,
) {
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
					*target = convert.GraphToFrameworkString(stringVal)
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
		if valueMap, ok := simpleValue.(map[string]any); ok {
			if stringValue, ok := valueMap["value"].(string); ok {
				*target = types.StringValue(stringValue)
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract string value")
}

// extractIntValue extracts an integer value using the additional data property
func extractIntValue(
	ctx context.Context,
	settingInstance models.DeviceManagementConfigurationSettingInstanceable,
	target *types.Int64,
) {
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
		if valueMap, ok := simpleValue.(map[string]any); ok {
			if numValue, ok := valueMap["value"].(float64); ok {
				*target = types.Int64Value(int64(numValue))
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract int value")
}

// extractBoolValue extracts a boolean from a typed choice or legacy additional data.
func extractBoolValue(
	ctx context.Context,
	settingInstance models.DeviceManagementConfigurationSettingInstanceable,
	target *types.Bool,
) {
	if settingInstance == nil {
		tflog.Warn(ctx, "Setting instance is nil when extracting bool value")
		return
	}

	if definitionID := settingInstance.GetSettingDefinitionId(); definitionID != nil {
		choice := types.StringNull()
		extractChoiceValue(ctx, settingInstance, &choice)
		switch choice.ValueString() {
		case *definitionID + "_0":
			*target = types.BoolValue(false)
			return
		case *definitionID + "_1":
			*target = types.BoolValue(true)
			return
		}
	}

	// Preserve support for legacy boolean values stored in additional data.

	// Access the setting's backing data as a map
	additionalData := settingInstance.GetAdditionalData()
	if additionalData == nil {
		tflog.Warn(ctx, "No additional data when extracting bool value")
		return
	}

	// Try to extract the boolean value from the additional data
	if simpleValue, ok := additionalData["simpleSettingValue"]; ok {
		if valueMap, ok := simpleValue.(map[string]any); ok {
			if boolValue, ok := valueMap["value"].(bool); ok {
				*target = types.BoolValue(boolValue)
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract bool value")
}

// extractChoiceValue extracts a choice value using the additional data property
func extractChoiceValue(
	ctx context.Context,
	settingInstance models.DeviceManagementConfigurationSettingInstanceable,
	target *types.String,
) {
	if settingInstance == nil {
		tflog.Warn(ctx, "Setting instance is nil when extracting choice value")
		return
	}

	// First try to access via strongly typed methods
	if choiceInstance, ok := settingInstance.(models.DeviceManagementConfigurationChoiceSettingInstanceable); ok {
		choiceValue := choiceInstance.GetChoiceSettingValue()
		if choiceValue != nil {
			if choiceVal := choiceValue.GetValue(); choiceVal != nil {
				*target = convert.GraphToFrameworkString(choiceVal)
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
		if valueMap, ok := choiceValue.(map[string]any); ok {
			if stringValue, ok := valueMap["value"].(string); ok {
				*target = types.StringValue(stringValue)
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract choice value")
}

// extractCollectionValue extracts collection values using the additional data property
func extractCollectionValue(
	ctx context.Context,
	settingInstance models.DeviceManagementConfigurationSettingInstanceable,
	target any,
) {
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
func extractAllowedAppsCollection(
	ctx context.Context,
	settingInstance models.DeviceManagementConfigurationSettingInstanceable,
	target *[]AllowedAppModel,
) {
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
		if collectionArray, ok := collectionValues.([]any); ok {
			var apps []AllowedAppModel
			for _, item := range collectionArray {
				if itemMap, ok := item.(map[string]any); ok {
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
func extractSimpleStringCollection(
	ctx context.Context,
	settingInstance models.DeviceManagementConfigurationSettingInstanceable,
	target *[]types.String,
) {
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
				*target = convert.GraphToFrameworkStringSlice(values)
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
		if valueMap, ok := simpleValue.(map[string]any); ok {
			if stringValue, ok := valueMap["value"].(string); ok {
				values := helpers.SplitCommaSeparatedString(stringValue)
				*target = convert.GraphToFrameworkStringSlice(values)
				return
			}
		}
	}

	// Also try collection values directly
	if collectionValues, ok := additionalData["simpleSettingCollectionValue"]; ok {
		if collectionArray, ok := collectionValues.([]any); ok {
			var values []string
			for _, item := range collectionArray {
				if itemMap, ok := item.(map[string]any); ok {
					if value, ok := itemMap["value"].(string); ok {
						values = append(values, value)
					}
				}
			}

			if len(values) > 0 {
				*target = convert.GraphToFrameworkStringSlice(values)
				return
			}
		}
	}

	tflog.Warn(ctx, "Failed to extract string collection")
}
