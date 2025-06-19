package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	builders "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
)

// constructResource builds the resource model for the Windows Autopilot Device Preparation Policy.
func constructResource(ctx context.Context, planModel *WindowsAutopilotDevicePreparationPolicyResourceModel) (models.DeviceManagementConfigurationPolicyable, error) {
	configurationPolicy := models.NewDeviceManagementConfigurationPolicy()

	convert.FrameworkToGraphString(planModel.Name, configurationPolicy.SetName)
	convert.FrameworkToGraphString(planModel.Description, configurationPolicy.SetDescription)

	// Set the template ID for Windows Autopilot Device Preparation Policy
	templateId := "80d33118-b7b4-40d8-b15f-81be745e053f_1"
	templateFamily := "enrollmentConfiguration"

	templateReference := models.NewDeviceManagementConfigurationPolicyTemplateReference()
	templateReference.SetTemplateId(&templateId)

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

	if err := convert.FrameworkToGraphStringSet(ctx, planModel.RoleScopeTagIds, configurationPolicy.SetRoleScopeTagIds); err != nil {
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

	// Add settings catalog template settings from the hcl plan
	if planModel.DeploymentSettings != nil {
		// Deployment Mode
		if !planModel.DeploymentSettings.DeploymentMode.IsNull() && !planModel.DeploymentSettings.DeploymentMode.IsUnknown() {
			deploymentModeSetting := builders.ConstructChoiceSettingInstance(
				"enrollment_autopilot_dpp_deploymentmode",
				planModel.DeploymentSettings.DeploymentMode.ValueString(),
				"5180aeab-886e-4589-97d4-40855c646315", // settingInstanceTemplateId
				"5874c2f6-bcf1-463b-a9eb-bee64e2f2d82", // settingValueTemplateId
			)
			settings = append(settings, deploymentModeSetting)
		}

		// Deployment Type
		if !planModel.DeploymentSettings.DeploymentType.IsNull() && !planModel.DeploymentSettings.DeploymentType.IsUnknown() {
			deploymentTypeSetting := builders.ConstructChoiceSettingInstance(
				"enrollment_autopilot_dpp_deploymenttype",
				planModel.DeploymentSettings.DeploymentType.ValueString(),
				"f4184296-fa9f-4b67-8b12-1723b3f8456b", // settingInstanceTemplateId
				"e0af022f-37f3-4a40-916d-1ab7281c88d9", // settingValueTemplateId
			)
			settings = append(settings, deploymentTypeSetting)
		}

		// Join Type
		if !planModel.DeploymentSettings.JoinType.IsNull() && !planModel.DeploymentSettings.JoinType.IsUnknown() {
			joinTypeSetting := builders.ConstructChoiceSettingInstance(
				"enrollment_autopilot_dpp_jointype",
				planModel.DeploymentSettings.JoinType.ValueString(),
				"6310e95d-6cfa-4d2f-aae0-1e7af12e2182", // settingInstanceTemplateId
				"1fa84eb3-fcfa-4ed6-9687-0f3d486402c4", // settingValueTemplateId
			)
			settings = append(settings, joinTypeSetting)
		}

		// Account Type
		if !planModel.DeploymentSettings.AccountType.IsNull() && !planModel.DeploymentSettings.AccountType.IsUnknown() {
			accountTypeSetting := builders.ConstructChoiceSettingInstance(
				"enrollment_autopilot_dpp_accountype",
				planModel.DeploymentSettings.AccountType.ValueString(),
				"d4f2a840-86d5-4162-9a08-fa8cc608b94e", // settingInstanceTemplateId
				"bf13bb47-69ef-4e06-97c1-50c2859a49c2", // settingValueTemplateId
			)
			settings = append(settings, accountTypeSetting)
		}
	}

	// Add OOBE settings
	if planModel.OOBESettings != nil {
		// Timeout in Minutes
		if !planModel.OOBESettings.TimeoutInMinutes.IsNull() && !planModel.OOBESettings.TimeoutInMinutes.IsUnknown() {
			timeoutSetting := builders.ConstructIntSimpleSettingInstance(
				"enrollment_autopilot_dpp_timeout",
				planModel.OOBESettings.TimeoutInMinutes.ValueInt64(),
				"6dec0657-dfb8-4906-a7ee-3ac6ee1edecb", // settingInstanceTemplateId
				"0bbcce5b-a55a-4e05-821a-94bf576d6cc8", // settingValueTemplateId
			)
			settings = append(settings, timeoutSetting)
		}

		// Custom Error Message
		if !planModel.OOBESettings.CustomErrorMessage.IsNull() && !planModel.OOBESettings.CustomErrorMessage.IsUnknown() {
			customErrorSetting := builders.ConstructStringSimpleSettingInstance(
				"enrollment_autopilot_dpp_customerrormessage",
				planModel.OOBESettings.CustomErrorMessage.ValueString(),
				"2ddf0619-2b7a-46de-b29b-c6191e9dda6e", // settingInstanceTemplateId
				"fe5002d5-fbe9-4920-9e2d-26bfc4b4cc97", // settingValueTemplateId
			)
			settings = append(settings, customErrorSetting)
		}

		// Allow Skip
		if !planModel.OOBESettings.AllowSkip.IsNull() && !planModel.OOBESettings.AllowSkip.IsUnknown() {
			allowSkipSetting := builders.ConstructBoolChoiceSettingInstance(
				"enrollment_autopilot_dpp_allowskip",
				planModel.OOBESettings.AllowSkip.ValueBool(),
				"2a71dc89-0f17-4ba9-bb27-af2521d34710", // settingInstanceTemplateId
				"a2323e5e-ac56-4517-8847-b0a6fdb467e7", // settingValueTemplateId
			)
			settings = append(settings, allowSkipSetting)
		}

		// Allow Diagnostics
		if !planModel.OOBESettings.AllowDiagnostics.IsNull() && !planModel.OOBESettings.AllowDiagnostics.IsUnknown() {
			allowDiagnosticsSetting := builders.ConstructBoolChoiceSettingInstance(
				"enrollment_autopilot_dpp_allowdiagnostics",
				planModel.OOBESettings.AllowDiagnostics.ValueBool(),
				"e2b7a81b-f243-4abd-bce3-c1856345f405", // settingInstanceTemplateId
				"c59d26fd-3460-4b26-b47a-f7e202e7d5a3", // settingValueTemplateId
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
			allowedAppsSetting := builders.ConstructSimpleSettingCollectionInstance(
				"enrollment_autopilot_dpp_allowedappids",
				appIds,
				"70d22a8a-a03c-4f62-b8df-dded3e327639", // settingInstanceTemplateId
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
			allowedScriptsSetting := builders.ConstructSimpleSettingCollectionInstance(
				"enrollment_autopilot_dpp_allowedscriptids",
				scriptIds,
				"1bc67702-800c-4271-8fd9-609351cc19cf", // settingInstanceTemplateId
			)
			settings = append(settings, allowedScriptsSetting)
		}
	}

	return settings, nil
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
