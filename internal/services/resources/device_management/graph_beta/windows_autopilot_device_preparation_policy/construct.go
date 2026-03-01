package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphdevicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	builders "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
)

// constructResource builds the resource model for the Windows Autopilot Device Preparation Policy.
func constructResource(
	ctx context.Context,
	planModel *WindowsAutopilotDevicePreparationPolicyResourceModel,
) (models.DeviceManagementConfigurationPolicyable, error) {
	configurationPolicy := models.NewDeviceManagementConfigurationPolicy()

	convert.FrameworkToGraphString(planModel.Name, configurationPolicy.SetName)
	convert.FrameworkToGraphString(planModel.Description, configurationPolicy.SetDescription)

	// Set the template ID for Windows Autopilot Device Preparation Policy
	templateId := "80d33118-b7b4-40d8-b15f-81be745e053f_1"
	templateFamily := "enrollmentConfiguration"

	templateReference := models.NewDeviceManagementConfigurationPolicyTemplateReference()
	templateReference.SetTemplateId(&templateId)

	parsedTemplateFamily, _ := models.ParseDeviceManagementConfigurationTemplateFamily(
		templateFamily,
	)
	if parsedFamily, ok := parsedTemplateFamily.(*models.DeviceManagementConfigurationTemplateFamily); ok &&
		parsedFamily != nil {
		templateReference.SetTemplateFamily(parsedFamily)
	}

	configurationPolicy.SetTemplateReference(templateReference)

	platformStr := "windows10"
	parsedPlatform, _ := models.ParseDeviceManagementConfigurationPlatforms(platformStr)
	if platform, ok := parsedPlatform.(*models.DeviceManagementConfigurationPlatforms); ok &&
		platform != nil {
		configurationPolicy.SetPlatforms(platform)
	}

	techStr := "enrollment"
	parsedTech, _ := models.ParseDeviceManagementConfigurationTechnologies(techStr)
	if tech, ok := parsedTech.(*models.DeviceManagementConfigurationTechnologies); ok &&
		tech != nil {
		configurationPolicy.SetTechnologies(tech)
	}

	if err := convert.FrameworkToGraphStringSet(
		ctx,
		planModel.RoleScopeTagIds,
		configurationPolicy.SetRoleScopeTagIds,
	); err != nil {
		return nil, fmt.Errorf("%w: %w", sentinels.ErrSetRoleScopeTags, err)
	}

	settings, err := constructAutopilotDevicePreparationPolicySettings(ctx, planModel)
	if err != nil {
		return nil, err
	}
	configurationPolicy.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		configurationPolicy,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
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
func constructAutopilotDevicePreparationPolicySettings(
	ctx context.Context,
	planModel *WindowsAutopilotDevicePreparationPolicyResourceModel,
) ([]models.DeviceManagementConfigurationSettingable, error) {
	var settings []models.DeviceManagementConfigurationSettingable

	if planModel.DeploymentSettings != nil {
		settings = appendDeploymentSettings(settings, planModel.DeploymentSettings)
	}

	if planModel.OOBESettings != nil {
		settings = appendOOBESettings(settings, planModel.OOBESettings)
	}

	settings = appendAllowedApps(settings, planModel.AllowedApps)
	settings = appendAllowedScripts(settings, planModel.AllowedScripts)

	return settings, nil
}

// appendDeploymentSettings adds deployment-related settings catalog template settings.
func appendDeploymentSettings(
	settings []models.DeviceManagementConfigurationSettingable,
	ds *DeploymentSettingsModel,
) []models.DeviceManagementConfigurationSettingable {
	if !ds.DeploymentMode.IsNull() && !ds.DeploymentMode.IsUnknown() {
		settings = append(settings, builders.ConstructChoiceSettingInstance(
			"enrollment_autopilot_dpp_deploymentmode",
			ds.DeploymentMode.ValueString(),
			"5180aeab-886e-4589-97d4-40855c646315", // settingInstanceTemplateId
			"5874c2f6-bcf1-463b-a9eb-bee64e2f2d82", // settingValueTemplateId
		))
	}

	if !ds.DeploymentType.IsNull() && !ds.DeploymentType.IsUnknown() {
		settings = append(settings, builders.ConstructChoiceSettingInstance(
			"enrollment_autopilot_dpp_deploymenttype",
			ds.DeploymentType.ValueString(),
			"f4184296-fa9f-4b67-8b12-1723b3f8456b", // settingInstanceTemplateId
			"e0af022f-37f3-4a40-916d-1ab7281c88d9", // settingValueTemplateId
		))
	}

	if !ds.JoinType.IsNull() && !ds.JoinType.IsUnknown() {
		settings = append(settings, builders.ConstructChoiceSettingInstance(
			"enrollment_autopilot_dpp_jointype",
			ds.JoinType.ValueString(),
			"6310e95d-6cfa-4d2f-aae0-1e7af12e2182", // settingInstanceTemplateId
			"1fa84eb3-fcfa-4ed6-9687-0f3d486402c4", // settingValueTemplateId
		))
	}

	if !ds.AccountType.IsNull() && !ds.AccountType.IsUnknown() {
		settings = append(settings, builders.ConstructChoiceSettingInstance(
			"enrollment_autopilot_dpp_accountype",
			ds.AccountType.ValueString(),
			"d4f2a840-86d5-4162-9a08-fa8cc608b94e", // settingInstanceTemplateId
			"bf13bb47-69ef-4e06-97c1-50c2859a49c2", // settingValueTemplateId
		))
	}

	return settings
}

// appendOOBESettings adds OOBE-related settings.
func appendOOBESettings(
	settings []models.DeviceManagementConfigurationSettingable,
	oobe *OOBESettingsModel,
) []models.DeviceManagementConfigurationSettingable {
	if !oobe.TimeoutInMinutes.IsNull() && !oobe.TimeoutInMinutes.IsUnknown() {
		settings = append(settings, builders.ConstructIntSimpleSettingInstance(
			"enrollment_autopilot_dpp_timeout",
			oobe.TimeoutInMinutes.ValueInt64(),
			"6dec0657-dfb8-4906-a7ee-3ac6ee1edecb", // settingInstanceTemplateId
			"0bbcce5b-a55a-4e05-821a-94bf576d6cc8", // settingValueTemplateId
		))
	}

	if !oobe.CustomErrorMessage.IsNull() && !oobe.CustomErrorMessage.IsUnknown() {
		settings = append(settings, builders.ConstructStringSimpleSettingInstance(
			"enrollment_autopilot_dpp_customerrormessage",
			oobe.CustomErrorMessage.ValueString(),
			"2ddf0619-2b7a-46de-b29b-c6191e9dda6e", // settingInstanceTemplateId
			"fe5002d5-fbe9-4920-9e2d-26bfc4b4cc97", // settingValueTemplateId
		))
	}

	if !oobe.AllowSkip.IsNull() && !oobe.AllowSkip.IsUnknown() {
		settings = append(settings, builders.ConstructBoolChoiceSettingInstance(
			"enrollment_autopilot_dpp_allowskip",
			oobe.AllowSkip.ValueBool(),
			"2a71dc89-0f17-4ba9-bb27-af2521d34710", // settingInstanceTemplateId
			"a2323e5e-ac56-4517-8847-b0a6fdb467e7", // settingValueTemplateId
		))
	}

	if !oobe.AllowDiagnostics.IsNull() && !oobe.AllowDiagnostics.IsUnknown() {
		settings = append(settings, builders.ConstructBoolChoiceSettingInstance(
			"enrollment_autopilot_dpp_allowdiagnostics",
			oobe.AllowDiagnostics.ValueBool(),
			"e2b7a81b-f243-4abd-bce3-c1856345f405", // settingInstanceTemplateId
			"c59d26fd-3460-4b26-b47a-f7e202e7d5a3", // settingValueTemplateId
		))
	}

	return settings
}

// appendAllowedApps adds allowed app settings.
func appendAllowedApps(
	settings []models.DeviceManagementConfigurationSettingable,
	allowedApps []AllowedAppModel,
) []models.DeviceManagementConfigurationSettingable {
	if len(allowedApps) == 0 {
		return settings
	}

	var appIds []string
	for _, app := range allowedApps {
		if !app.AppID.IsNull() && !app.AppID.IsUnknown() && !app.AppType.IsNull() &&
			!app.AppType.IsUnknown() {
			appId := app.AppID.ValueString()
			graphAppType := fmt.Sprintf("#microsoft.graph.%s", app.AppType.ValueString())
			appJson := fmt.Sprintf("{\"id\":\"%s\",\"type\":\"%s\"}", appId, graphAppType)
			appIds = append(appIds, appJson)
		}
	}

	if len(appIds) > 0 {
		settings = append(settings, builders.ConstructSimpleSettingCollectionInstance(
			"enrollment_autopilot_dpp_allowedappids",
			appIds,
			"70d22a8a-a03c-4f62-b8df-dded3e327639", // settingInstanceTemplateId
		))
	}

	return settings
}

// appendAllowedScripts adds allowed script settings.
func appendAllowedScripts(
	settings []models.DeviceManagementConfigurationSettingable,
	allowedScripts []types.String,
) []models.DeviceManagementConfigurationSettingable {
	if len(allowedScripts) == 0 {
		return settings
	}

	var scriptIds []string
	for _, script := range allowedScripts {
		if !script.IsNull() && !script.IsUnknown() {
			scriptIds = append(scriptIds, script.ValueString())
		}
	}

	if len(scriptIds) > 0 {
		settings = append(settings, builders.ConstructSimpleSettingCollectionInstance(
			"enrollment_autopilot_dpp_allowedscriptids",
			scriptIds,
			"1bc67702-800c-4271-8fd9-609351cc19cf", // settingInstanceTemplateId
		))
	}

	return settings
}

// constructAssignment constructs and returns a ConfigurationPoliciesItemAssignPostRequestBody
func constructAssignment(
	ctx context.Context,
	data *WindowsAutopilotDevicePreparationPolicyResourceModel,
) (graphdevicemanagement.ConfigurationPoliciesItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Starting assignment construction")

	requestBody := graphdevicemanagement.NewConfigurationPoliciesItemAssignPostRequestBody()
	policyAssignments := make([]models.DeviceManagementConfigurationPolicyAssignmentable, 0)

	if data.Assignments.IsNull() || data.Assignments.IsUnknown() {
		tflog.Debug(ctx, "Assignments is null or unknown, creating empty assignments array")
		requestBody.SetAssignments(policyAssignments)
		return requestBody, nil
	}

	var terraformAssignments []sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel
	diags := data.Assignments.ElementsAs(ctx, &terraformAssignments, false)
	if diags.HasError() {
		return nil, fmt.Errorf("%w: %v", sentinels.ErrExtractAssignments, diags.Errors())
	}

	for idx, assignment := range terraformAssignments {
		tflog.Debug(ctx, "Processing assignment", map[string]any{
			"index": idx,
		})

		graphAssignment := models.NewDeviceManagementConfigurationPolicyAssignment()

		if assignment.Type.IsNull() || assignment.Type.IsUnknown() {
			tflog.Error(ctx, "Assignment target type is missing or invalid", map[string]any{
				"index": idx,
			})
			continue
		}

		targetType := assignment.Type.ValueString()

		target := constructAssignmentTarget(ctx, targetType, assignment)
		if target == nil {
			tflog.Error(ctx, "Failed to create target", map[string]any{
				"index":      idx,
				"targetType": targetType,
			})
			continue
		}

		graphAssignment.SetTarget(target)
		policyAssignments = append(policyAssignments, graphAssignment)
	}

	tflog.Debug(ctx, "Completed assignment construction", map[string]any{
		"totalAssignments": len(policyAssignments),
	})

	requestBody.SetAssignments(policyAssignments)

	if err := constructors.DebugLogGraphObject(
		ctx,
		"Constructed assignment request body",
		requestBody,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructAssignmentTarget creates the appropriate target based on the target type
func constructAssignmentTarget(
	ctx context.Context,
	targetType string,
	assignment sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel,
) models.DeviceAndAppManagementAssignmentTargetable {
	var target models.DeviceAndAppManagementAssignmentTargetable

	switch targetType {
	case "allLicensedUsersAssignmentTarget":
		target = models.NewAllLicensedUsersAssignmentTarget()
	case "groupAssignmentTarget":
		groupTarget := models.NewGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() &&
			assignment.GroupId.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.GroupId, groupTarget.SetGroupId)
		} else {
			tflog.Error(ctx, "Group assignment target missing required group_id", map[string]any{
				"targetType": targetType,
			})
			return nil
		}
		target = groupTarget
	default:
		tflog.Error(ctx, "Unsupported target type", map[string]any{
			"targetType": targetType,
		})
		return nil
	}

	// Set filter if provided and meaningful (not default values)
	if !assignment.FilterId.IsNull() && !assignment.FilterId.IsUnknown() &&
		assignment.FilterId.ValueString() != "" &&
		assignment.FilterId.ValueString() != "00000000-0000-0000-0000-000000000000" {

		convert.FrameworkToGraphString(
			assignment.FilterId,
			target.SetDeviceAndAppManagementAssignmentFilterId,
		)

		if !assignment.FilterType.IsNull() && !assignment.FilterType.IsUnknown() &&
			assignment.FilterType.ValueString() != "" && assignment.FilterType.ValueString() != "none" {

			filterType := assignment.FilterType.ValueString()
			var filterTypeEnum models.DeviceAndAppManagementAssignmentFilterType
			switch filterType {
			case "include":
				filterTypeEnum = models.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			case "exclude":
				filterTypeEnum = models.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			default:
				tflog.Warn(ctx, "Unknown filter type, not setting filter", map[string]any{
					"filterType": filterType,
				})
				return target
			}
			target.SetDeviceAndAppManagementAssignmentFilterType(&filterTypeEnum)
		}
	}

	return target
}
