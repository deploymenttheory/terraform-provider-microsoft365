package graphBetaMacOSDeviceEnrollmentPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	builders "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
)

// constructResource builds the DeviceManagementConfigurationPolicy request body for the macOS ADE
// enrollment policy. creationSource (built from depOnboardingSettingsId) is sent on both Create
// and Update - confirmed against live Intune admin center traffic, which resends it unchanged on
// every PUT.
func constructResource(ctx context.Context, planModel *MacOSDeviceEnrollmentPolicyResourceModel, depOnboardingSettingsId string) (models.DeviceManagementConfigurationPolicyable, error) {
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

	settings, err := constructSettings(planModel)
	if err != nil {
		return nil, err
	}
	configurationPolicy.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), configurationPolicy); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing resource")

	return configurationPolicy, nil
}

// constructSettings builds the full settings catalog tree for the macOS ADE enrollment policy.
func constructSettings(planModel *MacOSDeviceEnrollmentPolicyResourceModel) ([]models.DeviceManagementConfigurationSettingable, error) {
	var settings []models.DeviceManagementConfigurationSettingable

	settings = append(settings, constructUserAffinitySetting(planModel))

	awaitSetting, err := constructAwaitConfigurationSetting(planModel)
	if err != nil {
		return nil, err
	}
	settings = append(settings, awaitSetting)

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

	return settings, nil
}

// constructUserAffinitySetting builds ade_macos_useraffinity, with its ade_macos_authenticationmethod
// child present only when requires_user_authentication is true.
func constructUserAffinitySetting(planModel *MacOSDeviceEnrollmentPolicyResourceModel) models.DeviceManagementConfigurationSettingable {
	requiresAuth := planModel.RequiresUserAuthentication.ValueBool()

	var children []models.DeviceManagementConfigurationSettingInstanceable
	if requiresAuth {
		authValue := AuthenticationMethodBasic
		if planModel.EnableAuthenticationViaCompanyPortal.ValueBool() {
			authValue = AuthenticationMethodCompanyPortal
		}
		if planModel.RequireCompanyPortalOnSetupAssistantEnrolledDevices.ValueBool() {
			authValue = AuthenticationMethodCompanyPortalOnSetupAssistant
		}
		authSetting := builders.ConstructChoiceSettingInstance(SettingDefAuthenticationMethod, authValue, "", "")
		children = append(children, authSetting.GetSettingInstance())
	}

	value := SettingDefUserAffinity + "_0"
	if requiresAuth {
		value = SettingDefUserAffinity + "_1"
	}

	return builders.ConstructChoiceSettingWithChildren(
		SettingDefUserAffinity,
		value,
		SettingInstanceTemplateUserAffinity,
		SettingValueTemplateUserAffinity,
		children,
	)
}

// constructAwaitConfigurationSetting builds ade_macos_awaitconfiguration. Confirmed against live
// Intune admin center traffic: this parent setting is always sent as "_1", with its
// ade_accountsettings_createlocaladmin child always present - even when no admin_account is
// configured, in which case the child (and its own createlocalprimary child) are sent as "_0"
// rather than omitted. await_device_configured therefore no longer varies this value; it remains
// a config-time gate on whether admin_account may/must be set (see requireAdminAccountWhenAwaitConfigured).
func constructAwaitConfigurationSetting(planModel *MacOSDeviceEnrollmentPolicyResourceModel) (models.DeviceManagementConfigurationSettingable, error) {
	admin := planModel.AdminAccount
	if admin == nil {
		admin = &AdminAccountModel{
			CreateLocalAdminAccount:   types.BoolValue(false),
			CreateLocalPrimaryAccount: types.BoolValue(false),
		}
	}

	adminInstance, err := constructCreateLocalAdminInstance(admin)
	if err != nil {
		return nil, err
	}

	return builders.ConstructChoiceSettingWithChildren(
		SettingDefAwaitConfiguration,
		SettingDefAwaitConfiguration+"_1",
		SettingInstanceTemplateAwaitConfiguration,
		SettingValueTemplateAwaitConfiguration,
		[]models.DeviceManagementConfigurationSettingInstanceable{adminInstance},
	), nil
}

// constructCreateLocalAdminInstance builds ade_accountsettings_createlocaladmin and its children.
// Confirmed against live Intune admin center traffic: ade_accountsettings_createlocalprimary is
// always present as a child, even when create_local_admin_account is false, in which case it (and
// its own subtree) is sent as "_0" with no further children.
func constructCreateLocalAdminInstance(admin *AdminAccountModel) (models.DeviceManagementConfigurationSettingInstanceable, error) {
	createAdmin := admin.CreateLocalAdminAccount.ValueBool()
	createPrimary := createAdmin && admin.CreateLocalPrimaryAccount.ValueBool()

	if !createAdmin && (admin.PrimaryAccount != nil || admin.CreateLocalPrimaryAccount.ValueBool()) {
		return nil, fmt.Errorf("admin_account sub-fields must not be set when create_local_admin_account is false")
	}
	if createAdmin && !createPrimary && admin.PrimaryAccount != nil {
		return nil, fmt.Errorf("admin_account.primary_account must not be set when create_local_primary_account is false")
	}

	var primaryChildren []models.DeviceManagementConfigurationSettingInstanceable
	if createPrimary {
		primaryInstance, err := constructPrefillAccountInfoInstance(admin.PrimaryAccount)
		if err != nil {
			return nil, err
		}
		primaryChildren = append(primaryChildren, primaryInstance)
	}

	primaryValue := SettingDefCreateLocalPrimary + "_0"
	if createPrimary {
		primaryValue = SettingDefCreateLocalPrimary + "_1"
	}
	createPrimarySetting := builders.ConstructChoiceSettingWithChildren(
		SettingDefCreateLocalPrimary, primaryValue, "", SettingValueTemplateAccountSettings, primaryChildren,
	)

	children := []models.DeviceManagementConfigurationSettingInstanceable{createPrimarySetting.GetSettingInstance()}
	if createAdmin {
		nameInstance := builders.ConstructStringSimpleSettingInstance(
			SettingDefAdminAccountName, admin.UserName.ValueString(), "", SettingValueTemplateAccountSettings,
		).GetSettingInstance()
		fullNameInstance := builders.ConstructStringSimpleSettingInstance(
			SettingDefAdminAccountFullName, admin.FullName.ValueString(), "", SettingValueTemplateAccountSettings,
		).GetSettingInstance()
		hideInstance := builders.ConstructBoolChoiceSettingInstance(
			SettingDefHideUsersGroups, admin.HideAccount.ValueBool(), "", SettingValueTemplateAccountSettings,
		).GetSettingInstance()
		rotationInstance := builders.ConstructIntSimpleSettingInstance(
			SettingDefAdminAccountPasswordRotation, admin.PasswordRotationInDays.ValueInt64(), "", SettingValueTemplateAccountSettings,
		).GetSettingInstance()

		children = append([]models.DeviceManagementConfigurationSettingInstanceable{
			nameInstance, fullNameInstance, hideInstance, rotationInstance,
		}, children...)
	}

	value := SettingDefCreateLocalAdmin + "_0"
	if createAdmin {
		value = SettingDefCreateLocalAdmin + "_1"
	}

	setting := builders.ConstructChoiceSettingWithChildren(
		SettingDefCreateLocalAdmin, value, "", SettingValueTemplateAccountSettings, children,
	)
	return setting.GetSettingInstance(), nil
}

// constructPrefillAccountInfoInstance builds ade_accountsettings_prefillaccountinfo and its children.
func constructPrefillAccountInfoInstance(primary *PrimaryAccountModel) (models.DeviceManagementConfigurationSettingInstanceable, error) {
	if primary == nil {
		primary = &PrimaryAccountModel{}
	}

	prefill := primary.PrefillAccountInfo.ValueBool()

	var children []models.DeviceManagementConfigurationSettingInstanceable
	if prefill {
		restrictInstance := builders.ConstructBoolChoiceSettingInstance(
			SettingDefRestrictEditing, primary.RestrictEditing.ValueBool(), "", "",
		).GetSettingInstance()
		fullNameInstance := builders.ConstructStringSimpleSettingInstance(
			SettingDefPrimaryAccountFullName, primary.FullName.ValueString(), "", SettingValueTemplateAccountSettings,
		).GetSettingInstance()
		userNameInstance := builders.ConstructStringSimpleSettingInstance(
			SettingDefPrimaryAccountName, primary.UserName.ValueString(), "", SettingValueTemplateAccountSettings,
		).GetSettingInstance()

		children = append(children, restrictInstance, fullNameInstance, userNameInstance)
	}

	value := SettingDefPrefillAccountInfo + "_0"
	if prefill {
		value = SettingDefPrefillAccountInfo + "_1"
	}

	setting := builders.ConstructChoiceSettingWithChildren(
		SettingDefPrefillAccountInfo, value, "", SettingValueTemplateAccountSettings, children,
	)
	return setting.GetSettingInstance(), nil
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
func constructSetupAssistantSettings(planModel *MacOSDeviceEnrollmentPolicyResourceModel) []models.DeviceManagementConfigurationSettingable {
	specs := []setupAssistantBoolSetting{
		{SettingDefLocationServices, !planModel.LocationServicesDisabled.ValueBool(), SettingInstanceTemplateLocationServices, SettingValueTemplateLocationServices},
		{SettingDefRestore, !planModel.RestoreDisabled.ValueBool(), SettingInstanceTemplateRestore, SettingValueTemplateRestore},
		{SettingDefAppleId, !planModel.AppleIdDisabled.ValueBool(), SettingInstanceTemplateAppleId, SettingValueTemplateAppleId},
		{SettingDefTermsAndConditions, !planModel.TermsAndConditionsDisabled.ValueBool(), SettingInstanceTemplateTermsAndConditions, SettingValueTemplateTermsAndConditions},
		{SettingDefTouchFaceId, !planModel.TouchIdDisabled.ValueBool(), SettingInstanceTemplateTouchFaceId, SettingValueTemplateTouchFaceId},
		{SettingDefApplePay, !planModel.ApplePayDisabled.ValueBool(), SettingInstanceTemplateApplePay, SettingValueTemplateApplePay},
		{SettingDefSiri, !planModel.SiriDisabled.ValueBool(), SettingInstanceTemplateSiri, SettingValueTemplateSiri},
		{SettingDefDiagnosticsData, !planModel.DiagnosticsDisabled.ValueBool(), SettingInstanceTemplateDiagnosticsData, SettingValueTemplateDiagnosticsData},
		{SettingDefFileVault, !planModel.FileVaultDisabled.ValueBool(), SettingInstanceTemplateFileVault, SettingValueTemplateFileVault},
		{SettingDefICloudDiagnostics, !planModel.ICloudDiagnosticsDisabled.ValueBool(), SettingInstanceTemplateICloudDiagnostics, SettingValueTemplateICloudDiagnostics},
		{SettingDefICloudStorage, !planModel.ICloudStorageDisabled.ValueBool(), SettingInstanceTemplateICloudStorage, SettingValueTemplateICloudStorage},
		{SettingDefAppearance, !planModel.DisplayToneSetupDisabled.ValueBool(), SettingInstanceTemplateAppearance, SettingValueTemplateAppearance},
		{SettingDefScreenTime, !planModel.ScreenTimeScreenDisabled.ValueBool(), SettingInstanceTemplateScreenTime, SettingValueTemplateScreenTime},
		{SettingDefPrivacy, !planModel.PrivacyPaneDisabled.ValueBool(), SettingInstanceTemplatePrivacy, SettingValueTemplatePrivacy},
		{SettingDefAccessibility, !planModel.AccessibilityScreenDisabled.ValueBool(), SettingInstanceTemplateAccessibility, SettingValueTemplateAccessibility},
		{SettingDefUnlockWithWatch, !planModel.AutoUnlockWithWatchDisabled.ValueBool(), SettingInstanceTemplateUnlockWithWatch, SettingValueTemplateUnlockWithWatch},
		{SettingDefEnableLockdownMode, !planModel.LockdownModeDisabled.ValueBool(), SettingInstanceTemplateEnableLockdownMode, SettingValueTemplateEnableLockdownMode},
		{SettingDefSoftwareUpdate, !planModel.SoftwareUpdateScreenDisabled.ValueBool(), SettingInstanceTemplateSoftwareUpdate, SettingValueTemplateSoftwareUpdate},
		{SettingDefSoftwareUpdateCompleted, !planModel.SoftwareUpdateCompletedScreenDisabled.ValueBool(), SettingInstanceTemplateSoftwareUpdateCompleted, SettingValueTemplateSoftwareUpdateCompleted},
		{SettingDefTermsOfAddress, !planModel.TermsOfAddressScreenDisabled.ValueBool(), SettingInstanceTemplateTermsOfAddress, SettingValueTemplateTermsOfAddress},
		{SettingDefIntelligence, !planModel.AppleIntelligenceDisabled.ValueBool(), SettingInstanceTemplateIntelligence, SettingValueTemplateIntelligence},
		{SettingDefOSShowcase, !planModel.OsShowcaseScreenDisabled.ValueBool(), SettingInstanceTemplateOSShowcase, SettingValueTemplateOSShowcase},
		{SettingDefAppStore, !planModel.AppStoreDisabled.ValueBool(), SettingInstanceTemplateAppStore, SettingValueTemplateAppStore},
	}

	settings := make([]models.DeviceManagementConfigurationSettingable, 0, len(specs))
	for _, spec := range specs {
		settings = append(settings, builders.ConstructBoolChoiceSettingInstance(
			spec.settingDefinitionId, spec.enabled, spec.settingInstanceTemplateId, spec.settingValueTemplateId,
		))
	}
	return settings
}
