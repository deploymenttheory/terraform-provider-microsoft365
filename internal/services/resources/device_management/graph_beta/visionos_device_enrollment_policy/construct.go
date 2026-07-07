package graphBetaVisionOSDeviceEnrollmentPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	builders "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
)

// constructResource builds the DeviceManagementConfigurationPolicy request body for the visionOS
// ADE enrollment policy. creationSource (built from depOnboardingSettingsId) is sent on both
// Create and Update - matching the macOS/iOS ADE policies, where live Intune admin center traffic
// resends it unchanged on every PUT.
func constructResource(ctx context.Context, planModel *VisionOSDeviceEnrollmentPolicyResourceModel, depOnboardingSettingsId string) (models.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	configurationPolicy := models.NewDeviceManagementConfigurationPolicy()

	convert.FrameworkToGraphString(planModel.Name, configurationPolicy.SetName)
	convert.FrameworkToGraphString(planModel.Description, configurationPolicy.SetDescription)

	templateId := TemplateID
	templateReference := models.NewDeviceManagementConfigurationPolicyTemplateReference()
	templateReference.SetTemplateId(&templateId)

	if parsedFamily, err := models.ParseDeviceManagementConfigurationTemplateFamily(TemplateFamily); err == nil {
		if family, ok := parsedFamily.(*models.DeviceManagementConfigurationTemplateFamily); ok && family != nil {
			templateReference.SetTemplateFamily(family)
		}
	}
	configurationPolicy.SetTemplateReference(templateReference)

	if parsedPlatform, err := models.ParseDeviceManagementConfigurationPlatforms(Platforms); err == nil {
		if platform, ok := parsedPlatform.(*models.DeviceManagementConfigurationPlatforms); ok && platform != nil {
			configurationPolicy.SetPlatforms(platform)
		}
	}

	if parsedTech, err := models.ParseDeviceManagementConfigurationTechnologies(Technologies); err == nil {
		if tech, ok := parsedTech.(*models.DeviceManagementConfigurationTechnologies); ok && tech != nil {
			configurationPolicy.SetTechnologies(tech)
		}
	}

	if err := convert.FrameworkToGraphStringSet(ctx, planModel.RoleScopeTagIds, configurationPolicy.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("%w: %w", sentinels.ErrSetRoleScopeTags, err)
	}

	if depOnboardingSettingsId != "" {
		creationSource := CreationSourcePrefix + depOnboardingSettingsId
		configurationPolicy.SetAdditionalData(map[string]any{
			"creationSource": creationSource,
		})
	}

	settings := constructSettings(planModel)
	configurationPolicy.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), configurationPolicy); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing resource")

	return configurationPolicy, nil
}

// constructSettings builds the full settings catalog tree for the visionOS ADE enrollment policy.
// Every setting is a flat, childless choice or string - visionOS uses the "basic" user affinity
// and await configuration variants with no nested subtrees. RequiresUserAuthentication defaults
// to (and is expected to stay) false: Graph only accepts ade_useraffinitybasic_0 for visionOS,
// since ADE enrollment for this platform is userless-only.
func constructSettings(planModel *VisionOSDeviceEnrollmentPolicyResourceModel) []models.DeviceManagementConfigurationSettingable {
	var settings []models.DeviceManagementConfigurationSettingable

	settings = append(settings, builders.ConstructBoolChoiceSettingInstance(
		SettingDefUserAffinity,
		planModel.UserAffinity.ValueBool(),
		SettingInstanceTemplateUserAffinity,
		SettingValueTemplateUserAffinity,
	))

	settings = append(settings, builders.ConstructBoolChoiceSettingInstance(
		SettingDefAwaitConfiguration,
		planModel.AwaitDeviceConfigured.ValueBool(),
		SettingInstanceTemplateAwaitConfiguration,
		SettingValueTemplateAwaitConfiguration,
	))

	settings = append(settings, builders.ConstructBoolChoiceSettingInstance(
		SettingDefLockedEnrollment,
		planModel.LockedEnrollmentEnabled.ValueBool(),
		SettingInstanceTemplateLockedEnrollment,
		SettingValueTemplateLockedEnrollment,
	))

	settings = append(settings, builders.ConstructStringSimpleSettingInstance(
		SettingDefDepartment,
		planModel.SupportDepartment.ValueString(),
		SettingInstanceTemplateDepartment,
		SettingValueTemplateDepartment,
	))
	settings = append(settings, builders.ConstructStringSimpleSettingInstance(
		SettingDefDepartmentPhone,
		planModel.SupportPhoneNumber.ValueString(),
		SettingInstanceTemplateDepartmentPhone,
		SettingValueTemplateDepartmentPhone,
	))

	settings = append(settings, constructSetupAssistantSettings(planModel)...)

	return settings
}

// setupAssistantBoolSetting describes one Setup Assistant screen toggle.
type setupAssistantBoolSetting struct {
	settingDefinitionId       string
	enabled                   bool
	settingInstanceTemplateId string
	settingValueTemplateId    string
}

// constructSetupAssistantSettings builds every Setup Assistant screen toggle. Each schema
// attribute is named `<screen>_disabled`, so the value passed to Graph is the logical negation.
func constructSetupAssistantSettings(planModel *VisionOSDeviceEnrollmentPolicyResourceModel) []models.DeviceManagementConfigurationSettingable {
	specs := []setupAssistantBoolSetting{
		{SettingDefAppleId, !planModel.AppleIdDisabled.ValueBool(), SettingInstanceTemplateAppleId, SettingValueTemplateAppleId},
		{SettingDefApplePay, !planModel.ApplePayDisabled.ValueBool(), SettingInstanceTemplateApplePay, SettingValueTemplateApplePay},
		{SettingDefDiagnosticsData, !planModel.DiagnosticsDisabled.ValueBool(), SettingInstanceTemplateDiagnosticsData, SettingValueTemplateDiagnosticsData},
		{SettingDefGetStarted, !planModel.GetStartedScreenDisabled.ValueBool(), SettingInstanceTemplateGetStarted, SettingValueTemplateGetStarted},
		{SettingDefIntelligence, !planModel.AppleIntelligenceDisabled.ValueBool(), SettingInstanceTemplateIntelligence, SettingValueTemplateIntelligence},
		{SettingDefLocationServices, !planModel.LocationServicesDisabled.ValueBool(), SettingInstanceTemplateLocationServices, SettingValueTemplateLocationServices},
		{SettingDefPasscode, !planModel.PasscodeDisabled.ValueBool(), SettingInstanceTemplatePasscode, SettingValueTemplatePasscode},
		{SettingDefPrivacy, !planModel.PrivacyPaneDisabled.ValueBool(), SettingInstanceTemplatePrivacy, SettingValueTemplatePrivacy},
		{SettingDefScreenTime, !planModel.ScreenTimeScreenDisabled.ValueBool(), SettingInstanceTemplateScreenTime, SettingValueTemplateScreenTime},
		{SettingDefSiri, !planModel.SiriDisabled.ValueBool(), SettingInstanceTemplateSiri, SettingValueTemplateSiri},
		{SettingDefSoftwareUpdate, !planModel.SoftwareUpdateScreenDisabled.ValueBool(), SettingInstanceTemplateSoftwareUpdate, SettingValueTemplateSoftwareUpdate},
		{SettingDefTermsAndConditions, !planModel.TermsAndConditionsDisabled.ValueBool(), SettingInstanceTemplateTermsAndConditions, SettingValueTemplateTermsAndConditions},
		{SettingDefTips, !planModel.TipsScreenDisabled.ValueBool(), SettingInstanceTemplateTips, SettingValueTemplateTips},
		{SettingDefTouchFaceId, !planModel.TouchIdDisabled.ValueBool(), SettingInstanceTemplateTouchFaceId, SettingValueTemplateTouchFaceId},
	}

	settings := make([]models.DeviceManagementConfigurationSettingable, 0, len(specs))
	for _, spec := range specs {
		settings = append(settings, builders.ConstructBoolChoiceSettingInstance(
			spec.settingDefinitionId, spec.enabled, spec.settingInstanceTemplateId, spec.settingValueTemplateId,
		))
	}
	return settings
}
