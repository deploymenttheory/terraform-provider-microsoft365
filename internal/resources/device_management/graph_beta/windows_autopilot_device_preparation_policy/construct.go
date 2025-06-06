package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource builds the resource model for the Windows Autopilot Device Preparation Policy.
func constructResource(ctx context.Context, planModel *WindowsAutopilotDevicePreparationPolicyResourceModel) (models.DeviceManagementConfigurationPolicyable, error) {
	configurationPolicy := models.NewDeviceManagementConfigurationPolicy()

	// Set the basic properties using constructors
	constructors.SetStringProperty(planModel.Name, configurationPolicy.SetName)
	constructors.SetStringProperty(planModel.Description, configurationPolicy.SetDescription)

	// Set the template ID for Windows Autopilot Device Preparation Policy
	templateId := "80d33118-b7b4-40d8-b15f-81be745e053f_1"
	templateFamily := "enrollmentConfiguration"

	templateReference := models.NewDeviceManagementConfigurationPolicyTemplateReference()
	templateReference.SetTemplateId(&templateId)

	// Parse the template family string into the proper enum type
	parsedTemplateFamily, _ := models.ParseDeviceManagementConfigurationTemplateFamily(templateFamily)
	if parsedFamily, ok := parsedTemplateFamily.(*models.DeviceManagementConfigurationTemplateFamily); ok && parsedFamily != nil {
		templateReference.SetTemplateFamily(parsedFamily)
	}

	configurationPolicy.SetTemplateReference(templateReference)

	platformStr := "windows10"
	parsedPlatform, _ := models.ParseDeviceManagementConfigurationPlatforms(platformStr)
	if platform, ok := parsedPlatform.(*models.DeviceManagementConfigurationPlatforms); ok && platform != nil {
		configurationPolicy.SetPlatforms(platform)
	}

	techStr := "enrollment"
	parsedTech, _ := models.ParseDeviceManagementConfigurationTechnologies(techStr)
	if tech, ok := parsedTech.(*models.DeviceManagementConfigurationTechnologies); ok && tech != nil {
		configurationPolicy.SetTechnologies(tech)
	}

	if err := constructors.SetStringSet(ctx, planModel.RoleScopeTagIds, configurationPolicy.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	settings, err := constructAutopilotDevicePreparationPolicySettings(ctx, planModel)
	if err != nil {
		return nil, err
	}
	configurationPolicy.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), configurationPolicy); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return configurationPolicy, nil
}

// constructAutopilotDevicePreparationPolicySettings builds the settings for the Windows Autopilot Device Preparation Policy.
//
// With the specific:
// - settingDefinitionId string extracted from graph x-ray testing.
// - value string from the HCL, e.g enrollment_autopilot_dpp_deploymentmode_0 or enrollment_autopilot_dpp_accountype_0 which is the settings catalog value. Extracted from graph x-ray testing.
// - settingInstanceTemplateId string extracted from graph x-ray testing.
// - settingValueTemplateId string extracted from graph x-ray testing.
//
// It then calls the generic helpers to construct the setting instance and value
// This func is unique per policy type and therefore per terraform resource type that uses templates.
func constructAutopilotDevicePreparationPolicySettings(ctx context.Context, planModel *WindowsAutopilotDevicePreparationPolicyResourceModel) ([]models.DeviceManagementConfigurationSettingable, error) {
	var settings []models.DeviceManagementConfigurationSettingable

	// Add deployment settings
	if planModel.DeploymentSettings != nil {
		// Deployment Mode
		if !planModel.DeploymentSettings.DeploymentMode.IsNull() && !planModel.DeploymentSettings.DeploymentMode.IsUnknown() {
			deploymentModeSetting := constructChoiceSetting(
				"enrollment_autopilot_dpp_deploymentmode",
				planModel.DeploymentSettings.DeploymentMode.ValueString(),
				"5180aeab-886e-4589-97d4-40855c646315",
				"5874c2f6-bcf1-463b-a9eb-bee64e2f2d82",
			)
			settings = append(settings, deploymentModeSetting)
		}

		// Deployment Type
		if !planModel.DeploymentSettings.DeploymentType.IsNull() && !planModel.DeploymentSettings.DeploymentType.IsUnknown() {
			deploymentTypeSetting := constructChoiceSetting(
				"enrollment_autopilot_dpp_deploymenttype",
				planModel.DeploymentSettings.DeploymentType.ValueString(),
				"f4184296-fa9f-4b67-8b12-1723b3f8456b",
				"e0af022f-37f3-4a40-916d-1ab7281c88d9",
			)
			settings = append(settings, deploymentTypeSetting)
		}

		// Join Type
		if !planModel.DeploymentSettings.JoinType.IsNull() && !planModel.DeploymentSettings.JoinType.IsUnknown() {
			joinTypeSetting := constructChoiceSetting(
				"enrollment_autopilot_dpp_jointype",
				planModel.DeploymentSettings.JoinType.ValueString(),
				"6310e95d-6cfa-4d2f-aae0-1e7af12e2182",
				"1fa84eb3-fcfa-4ed6-9687-0f3d486402c4",
			)
			settings = append(settings, joinTypeSetting)
		}

		// Account Type
		if !planModel.DeploymentSettings.AccountType.IsNull() && !planModel.DeploymentSettings.AccountType.IsUnknown() {
			accountTypeSetting := constructChoiceSetting(
				"enrollment_autopilot_dpp_accountype",
				planModel.DeploymentSettings.AccountType.ValueString(),
				"d4f2a840-86d5-4162-9a08-fa8cc608b94e",
				"bf13bb47-69ef-4e06-97c1-50c2859a49c2",
			)
			settings = append(settings, accountTypeSetting)
		}
	}

	// Add OOBE settings
	if planModel.OOBESettings != nil {
		// Timeout in Minutes
		if !planModel.OOBESettings.TimeoutInMinutes.IsNull() && !planModel.OOBESettings.TimeoutInMinutes.IsUnknown() {
			timeoutSetting := constructIntSetting(
				"enrollment_autopilot_dpp_timeout",
				planModel.OOBESettings.TimeoutInMinutes.ValueInt64(),
				"6dec0657-dfb8-4906-a7ee-3ac6ee1edecb",
				"0bbcce5b-a55a-4e05-821a-94bf576d6cc8",
			)
			settings = append(settings, timeoutSetting)
		}

		// Custom Error Message
		if !planModel.OOBESettings.CustomErrorMessage.IsNull() && !planModel.OOBESettings.CustomErrorMessage.IsUnknown() {
			customErrorSetting := constructSimpleSetting(
				"enrollment_autopilot_dpp_customerrormessage",
				planModel.OOBESettings.CustomErrorMessage.ValueString(),
				"2ddf0619-2b7a-46de-b29b-c6191e9dda6e",
				"fe5002d5-fbe9-4920-9e2d-26bfc4b4cc97",
			)
			settings = append(settings, customErrorSetting)
		}

		// Allow Skip
		if !planModel.OOBESettings.AllowSkip.IsNull() && !planModel.OOBESettings.AllowSkip.IsUnknown() {
			allowSkipSetting := constructBoolSetting(
				"enrollment_autopilot_dpp_allowskip",
				planModel.OOBESettings.AllowSkip.ValueBool(),
				"2a71dc89-0f17-4ba9-bb27-af2521d34710",
				"a2323e5e-ac56-4517-8847-b0a6fdb467e7",
			)
			settings = append(settings, allowSkipSetting)
		}

		// Allow Diagnostics
		if !planModel.OOBESettings.AllowDiagnostics.IsNull() && !planModel.OOBESettings.AllowDiagnostics.IsUnknown() {
			allowDiagnosticsSetting := constructBoolSetting(
				"enrollment_autopilot_dpp_allowdiagnostics",
				planModel.OOBESettings.AllowDiagnostics.ValueBool(),
				"e2b7a81b-f243-4abd-bce3-c1856345f405",
				"c59d26fd-3460-4b26-b47a-f7e202e7d5a3",
			)
			settings = append(settings, allowDiagnosticsSetting)
		}
	}

	// Add allowed apps
	if len(planModel.AllowedApps) > 0 {
		var appIds []string
		for _, app := range planModel.AllowedApps {
			if !app.AppID.IsNull() && !app.AppID.IsUnknown() && !app.AppType.IsNull() && !app.AppType.IsUnknown() {

				appId := app.AppID.ValueString()

				appTypeStr := app.AppType.ValueString()
				graphAppType := fmt.Sprintf("#microsoft.graph.%s", appTypeStr)

				appJson := fmt.Sprintf("{\"id\":\"%s\",\"type\":\"%s\"}", appId, graphAppType)
				appIds = append(appIds, appJson)
			}
		}

		if len(appIds) > 0 {
			allowedAppsSetting := constructCollectionSetting(
				"enrollment_autopilot_dpp_allowedappids",
				appIds,
				"70d22a8a-a03c-4f62-b8df-dded3e327639",
			)
			settings = append(settings, allowedAppsSetting)
		}
	}

	// Add allowed scripts
	if len(planModel.AllowedScripts) > 0 {
		var scriptIds []string
		for _, script := range planModel.AllowedScripts {
			if !script.IsNull() && !script.IsUnknown() {
				scriptIds = append(scriptIds, script.ValueString())
			}
		}

		if len(scriptIds) > 0 {
			allowedScriptsSetting := constructCollectionSetting(
				"enrollment_autopilot_dpp_allowedscriptids",
				scriptIds,
				"1bc67702-800c-4271-8fd9-609351cc19cf",
			)
			settings = append(settings, allowedScriptsSetting)
		}
	}

	return settings, nil
}

// constructChoiceSetting creates a choice setting.
func constructChoiceSetting(
	settingDefinitionId string,
	value string,
	settingInstanceTemplateId string,
	settingValueTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	// Create setting instance
	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	// Add OData type for the instance
	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Create choice setting value
	choiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	valuePtr := value
	choiceValue.SetValue(&valuePtr)

	// Add OData type for the value
	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceValue.SetOdataType(&odataTypeValue)

	// Set empty children array
	var children []models.DeviceManagementConfigurationSettingInstanceable
	choiceValue.SetChildren(children)

	// Set template references if provided
	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		choiceValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	settingInstance.SetChoiceSettingValue(choiceValue)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// constructSimpleSetting creates a simple string setting.
func constructSimpleSetting(
	settingDefinitionId string,
	value string,
	settingInstanceTemplateId string,
	settingValueTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	simpleSettingValue := models.NewDeviceManagementConfigurationStringSettingValue()
	valuePtr := value
	simpleSettingValue.SetValue(&valuePtr)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
	simpleSettingValue.SetOdataType(&odataTypeValue)

	// Set template references if provided
	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		simpleSettingValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	settingInstance.SetSimpleSettingValue(simpleSettingValue)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// constructIntSetting creates an integer setting.
func constructIntSetting(
	settingDefinitionId string,
	value int64,
	settingInstanceTemplateId string,
	settingValueTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	intSettingValue := models.NewDeviceManagementConfigurationIntegerSettingValue()
	intValue := int32(value)
	intSettingValue.SetValue(&intValue)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
	intSettingValue.SetOdataType(&odataTypeValue)

	// Set template references if provided
	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		intSettingValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	settingInstance.SetSimpleSettingValue(intSettingValue)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// constructCollectionSetting creates a collection setting for string values.
func constructCollectionSetting(
	settingDefinitionId string,
	values []string,
	settingInstanceTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set template reference if provided
	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	// Create simple setting collection values
	var simpleSettingCollectionValues []models.DeviceManagementConfigurationSimpleSettingValueable
	for _, val := range values {
		simpleSettingValue := models.NewDeviceManagementConfigurationStringSettingValue()
		valuePtr := val
		simpleSettingValue.SetValue(&valuePtr)

		odataTypeValue := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
		simpleSettingValue.SetOdataType(&odataTypeValue)

		simpleSettingCollectionValues = append(simpleSettingCollectionValues, simpleSettingValue)
	}

	settingInstance.SetSimpleSettingCollectionValue(simpleSettingCollectionValues)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// constructBoolSetting creates a boolean setting using the choice setting format.
func constructBoolSetting(settingDefinitionId string, value bool, settingInstanceTemplateId string, settingValueTemplateId string) models.DeviceManagementConfigurationSettingable {
	// Convert bool to appropriate format for the API
	// The Graph API uses strings with _0 or _1 suffixes to represent boolean values
	strValue := fmt.Sprintf("%s_0", settingDefinitionId) // Default to false value
	if value {
		strValue = fmt.Sprintf("%s_1", settingDefinitionId) // True value
	}

	// Use constructChoiceSetting since boolean settings are represented as choice settings in the API
	return constructChoiceSetting(settingDefinitionId, strValue, settingInstanceTemplateId, settingValueTemplateId)
}

// constructAssignment creates an assignment for the policy.
func constructAssignment(ctx context.Context, assignment *WindowsAutopilotDevicePreparationAssignment) (graphdevicemanagement.ConfigurationPoliciesItemAssignPostRequestBodyable, error) {
	if assignment == nil {
		return nil, fmt.Errorf("assignment cannot be nil")
	}

	requestBody := graphdevicemanagement.NewConfigurationPoliciesItemAssignPostRequestBody()
	assignments := make([]models.DeviceManagementConfigurationPolicyAssignmentable, 0)

	if len(assignment.IncludeGroupIds) > 0 {
		for _, groupId := range assignment.IncludeGroupIds {
			if !groupId.IsNull() && !groupId.IsUnknown() && groupId.ValueString() != "" {
				policyAssignment := models.NewDeviceManagementConfigurationPolicyAssignment()
				target := models.NewGroupAssignmentTarget()

				// Set the group ID
				groupIDStr := groupId.ValueString()
				target.SetGroupId(&groupIDStr)

				policyAssignment.SetTarget(target)
				assignments = append(assignments, policyAssignment)
			}
		}
	}

	requestBody.SetAssignments(assignments)
	return requestBody, nil
}
