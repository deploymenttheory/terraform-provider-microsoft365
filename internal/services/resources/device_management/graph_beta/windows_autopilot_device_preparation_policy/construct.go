package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	builders "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
)

// constructResource builds the resource model for the Windows Autopilot Device Preparation Policy.
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, planModel *WindowsAutopilotDevicePreparationPolicyResourceModel) (models.DeviceManagementConfigurationPolicyable, error) {

	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if diags := validateRequest(ctx, client, planModel); diags.HasError() {
		return nil, fmt.Errorf("validation failed: %v", diags.Errors())
	}

	// Determine policy mode based on deployment_type
	// deployment_type_0 = User-driven (uses user-driven template)
	// deployment_type_1 = Self-deploying/automatic (uses automatic template)
	var deploymentType string
	if planModel.DeploymentSettings != nil && !planModel.DeploymentSettings.DeploymentType.IsNull() {
		deploymentType = planModel.DeploymentSettings.DeploymentType.ValueString()
	}

	switch deploymentType {
	case DeploymentTypeUserDriven:
		return constructUserDrivenPolicy(ctx, planModel)
	case DeploymentTypeSelfDeploying:
		return constructAutomaticPolicy(ctx, planModel)
	default:
		return nil, fmt.Errorf("deployment_settings.deployment_type is required and must be '%s' (user-driven) or '%s' (self-deploying)", DeploymentTypeUserDriven, DeploymentTypeSelfDeploying)
	}
}

// constructAutomaticPolicy builds a self-deploying/automatic mode policy (simpler structure, no device security group)
func constructAutomaticPolicy(ctx context.Context, planModel *WindowsAutopilotDevicePreparationPolicyResourceModel) (models.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing self-deploying/automatic mode policy")

	// Validate that self-deploying mode doesn't have user-driven only fields
	if planModel.OOBESettings != nil {
		return nil, fmt.Errorf("oobe_settings cannot be set for self-deploying/automatic mode policies (deployment_type_1)")
	}
	if !planModel.DeviceSecurityGroup.IsNull() && !planModel.DeviceSecurityGroup.IsUnknown() {
		return nil, fmt.Errorf("device_security_group cannot be set for self-deploying/automatic mode policies (deployment_type_1)")
	}
	if !planModel.Assignments.IsNull() && !planModel.Assignments.IsUnknown() {
		return nil, fmt.Errorf("assignments cannot be set for automatic mode policies")
	}

	configurationPolicy := models.NewDeviceManagementConfigurationPolicy()

	convert.FrameworkToGraphString(planModel.Name, configurationPolicy.SetName)
	convert.FrameworkToGraphString(planModel.Description, configurationPolicy.SetDescription)

	// Set the template ID for automatic mode
	templateId := TemplateIDAutomatic
	templateFamily := TemplateFamily

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

	if err := convert.FrameworkToGraphStringSet(ctx, planModel.RoleScopeTagIds, configurationPolicy.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("%w: %w", sentinels.ErrSetRoleScopeTags, err)
	}

	// For automatic mode, only include allowed apps and scripts
	settings, err := constructAutomaticPolicySettings(ctx, planModel)
	if err != nil {
		return nil, err
	}
	configurationPolicy.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (automatic mode)", ResourceName), configurationPolicy); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing self-deploying/automatic mode policy")

	return configurationPolicy, nil
}

// constructUserDrivenPolicy builds a user-driven mode policy (full structure with device security group)
func constructUserDrivenPolicy(ctx context.Context, planModel *WindowsAutopilotDevicePreparationPolicyResourceModel) (models.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing user-driven mode policy")

	// Validate that user-driven mode has required fields
	if planModel.DeploymentSettings == nil {
		return nil, fmt.Errorf("deployment_settings is required for user_driven mode policies")
	}
	if planModel.OOBESettings == nil {
		return nil, fmt.Errorf("oobe_settings is required for user_driven mode policies")
	}
	if planModel.DeviceSecurityGroup.IsNull() || planModel.DeviceSecurityGroup.IsUnknown() {
		return nil, fmt.Errorf("device_security_group is required for user_driven mode policies")
	}

	configurationPolicy := models.NewDeviceManagementConfigurationPolicy()

	convert.FrameworkToGraphString(planModel.Name, configurationPolicy.SetName)
	convert.FrameworkToGraphString(planModel.Description, configurationPolicy.SetDescription)

	// Set the template ID for user-driven mode
	templateId := TemplateIDUserDriven
	templateFamily := TemplateFamily

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

	settings, err := constructUserDrivenPolicySettings(ctx, planModel)
	if err != nil {
		return nil, err
	}
	configurationPolicy.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (automatic mode)", ResourceName), configurationPolicy); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing user-driven mode policy")

	return configurationPolicy, nil
}

// constructAutomaticPolicySettings builds settings for automatic mode policies (apps and scripts).
func constructAutomaticPolicySettings(ctx context.Context, planModel *WindowsAutopilotDevicePreparationPolicyResourceModel) ([]models.DeviceManagementConfigurationSettingable, error) {
	var settings []models.DeviceManagementConfigurationSettingable

	settings = appendAllowedAppsAutomatic(settings, planModel.AllowedApps)
	settings = appendAllowedScriptsAutomatic(settings, planModel.AllowedScripts)

	return settings, nil
}

// constructUserDrivenPolicySettings builds the settings for user-driven mode policies (full settings).
//
// With the specific:
// - settingDefinitionId string extracted from graph x-ray testing.
// - value string from the HCL, e.g enrollment_autopilot_dpp_deploymentmode_0 or enrollment_autopilot_dpp_accountype_0 which is the settings catalog value. Extracted from graph x-ray testing.
// - settingInstanceTemplateId string extracted from graph x-ray testing.
// - settingValueTemplateId string extracted from graph x-ray testing.
//
// It then calls the generic helpers to construct the setting instance and value
// This func is unique per policy type and therefore per terraform resource type that uses templates.
func constructUserDrivenPolicySettings(ctx context.Context, planModel *WindowsAutopilotDevicePreparationPolicyResourceModel) ([]models.DeviceManagementConfigurationSettingable, error) {
	var settings []models.DeviceManagementConfigurationSettingable

	if planModel.DeploymentSettings != nil {
		settings = appendDeploymentSettings(settings, planModel.DeploymentSettings)
	}

	if planModel.OOBESettings != nil {
		settings = appendOOBESettings(settings, planModel.OOBESettings)
	}

	settings = appendDeviceSecurityGroupSetting(settings)
	settings = appendAllowedAppsUserDriven(settings, planModel.AllowedApps)
	settings = appendAllowedScriptsUserDriven(settings, planModel.AllowedScripts)

	return settings, nil
}

// appendDeploymentSettings adds deployment-related settings catalog template settings.
func appendDeploymentSettings(
	settings []models.DeviceManagementConfigurationSettingable,
	ds *DeploymentSettingsModel,
) []models.DeviceManagementConfigurationSettingable {
	if !ds.DeploymentMode.IsNull() && !ds.DeploymentMode.IsUnknown() {
		settings = append(settings, builders.ConstructChoiceSettingInstance(
			SettingDefDeploymentMode,
			ds.DeploymentMode.ValueString(),
			SettingInstanceTemplateDeploymentMode,
			SettingValueTemplateDeploymentMode,
		))
	}

	if !ds.DeploymentType.IsNull() && !ds.DeploymentType.IsUnknown() {
		settings = append(settings, builders.ConstructChoiceSettingInstance(
			SettingDefDeploymentType,
			ds.DeploymentType.ValueString(),
			SettingInstanceTemplateDeploymentType,
			SettingValueTemplateDeploymentType,
		))
	}

	// Join type is always Entra ID joined (enrollment_autopilot_dpp_jointype_0)
	// Hybrid join is not supported for Windows Autopilot Device Preparation policies
	settings = append(settings, builders.ConstructChoiceSettingInstance(
		SettingDefJoinType,
		"enrollment_autopilot_dpp_jointype_0",
		SettingInstanceTemplateJoinType,
		SettingValueTemplateJoinType,
	))

	if !ds.AccountType.IsNull() && !ds.AccountType.IsUnknown() {
		settings = append(settings, builders.ConstructChoiceSettingInstance(
			SettingDefAccountType,
			ds.AccountType.ValueString(),
			SettingInstanceTemplateAccountType,
			SettingValueTemplateAccountType,
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
			SettingDefTimeout,
			oobe.TimeoutInMinutes.ValueInt64(),
			SettingInstanceTemplateTimeout,
			SettingValueTemplateTimeout,
		))
	}

	if !oobe.CustomErrorMessage.IsNull() && !oobe.CustomErrorMessage.IsUnknown() {
		settings = append(settings, builders.ConstructStringSimpleSettingInstance(
			SettingDefCustomErrorMessage,
			oobe.CustomErrorMessage.ValueString(),
			SettingInstanceTemplateCustomErrorMessage,
			SettingValueTemplateCustomErrorMessage,
		))
	}

	if !oobe.AllowSkip.IsNull() && !oobe.AllowSkip.IsUnknown() {
		settings = append(settings, builders.ConstructBoolChoiceSettingInstance(
			SettingDefAllowSkip,
			oobe.AllowSkip.ValueBool(),
			SettingInstanceTemplateAllowSkip,
			SettingValueTemplateAllowSkip,
		))
	}

	if !oobe.AllowDiagnostics.IsNull() && !oobe.AllowDiagnostics.IsUnknown() {
		settings = append(settings, builders.ConstructBoolChoiceSettingInstance(
			SettingDefAllowDiagnostics,
			oobe.AllowDiagnostics.ValueBool(),
			SettingInstanceTemplateAllowDiagnostics,
			SettingValueTemplateAllowDiagnostics,
		))
	}

	return settings
}

// appendDeviceSecurityGroupSetting adds the device security group setting with empty value.
// The actual group ID is set via the separate setEnrollmentTimeDeviceMembershipTarget API call.
func appendDeviceSecurityGroupSetting(
	settings []models.DeviceManagementConfigurationSettingable,
) []models.DeviceManagementConfigurationSettingable {
	settings = append(settings, builders.ConstructStringSimpleSettingInstance(
		SettingDefDeviceSecurityGroupIDs,
		"",
		SettingInstanceTemplateDeviceSecurityGroup,
		SettingValueTemplateDeviceSecurityGroup,
	))
	return settings
}

// appendAllowedAppsAutomatic adds allowed app settings for automatic mode policies.
func appendAllowedAppsAutomatic(settings []models.DeviceManagementConfigurationSettingable,
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
			SettingDefAllowedAppIDs,
			appIds,
			SettingInstanceTemplateAllowedAppsAutomatic,
		))
	}

	return settings
}

// appendAllowedAppsUserDriven adds allowed app settings for user-driven mode policies.
func appendAllowedAppsUserDriven(settings []models.DeviceManagementConfigurationSettingable,
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
			SettingDefAllowedAppIDs,
			appIds,
			SettingInstanceTemplateAllowedAppsUserDriven,
		))
	}

	return settings
}

// appendAllowedScriptsAutomatic adds allowed script settings for automatic mode policies.
func appendAllowedScriptsAutomatic(
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
			SettingDefAllowedScriptIDs,
			scriptIds,
			SettingInstanceTemplateAllowedScriptsAutomatic,
		))
	}

	return settings
}

// appendAllowedScriptsUserDriven adds allowed script settings for user-driven mode policies.
func appendAllowedScriptsUserDriven(
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
			SettingDefAllowedScriptIDs,
			scriptIds,
			SettingInstanceTemplateAllowedScriptsUserDriven,
		))
	}

	return settings
}
