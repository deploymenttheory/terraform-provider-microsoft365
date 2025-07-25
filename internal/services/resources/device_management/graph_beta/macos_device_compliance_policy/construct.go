package graphBetaDeviceCompliancePolicies

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *DeviceCompliancePolicyResourceModel) (graphmodels.DeviceCompliancePolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	odataType := data.Type.ValueString()
	var requestBody graphmodels.DeviceCompliancePolicyable

	switch odataType {
	case "aospDeviceOwnerCompliancePolicy":
		aospPolicy := graphmodels.NewAospDeviceOwnerCompliancePolicy()
		if err := constructAospDeviceOwnerCompliancePolicy(ctx, data, aospPolicy); err != nil {
			return nil, fmt.Errorf("failed to construct AOSP device owner policy: %s", err)
		}
		requestBody = aospPolicy
	case "androidDeviceOwnerCompliancePolicy":
		androidPolicy := graphmodels.NewAndroidDeviceOwnerCompliancePolicy()
		if err := constructAndroidDeviceOwnerPolicy(ctx, data, androidPolicy); err != nil {
			return nil, fmt.Errorf("failed to construct Android device owner policy: %s", err)
		}
		requestBody = androidPolicy
	case "iosCompliancePolicy":
		iosPolicy := graphmodels.NewIosCompliancePolicy()
		if err := constructIosCompliancePolicy(ctx, data, iosPolicy); err != nil {
			return nil, fmt.Errorf("failed to construct iOS policy: %s", err)
		}
		requestBody = iosPolicy
	case "windows10CompliancePolicy":
		windowsPolicy := graphmodels.NewWindows10CompliancePolicy()
		if err := constructWindows10CompliancePolicy(ctx, data, windowsPolicy); err != nil {
			return nil, fmt.Errorf("failed to construct Windows 10 policy: %s", err)
		}
		requestBody = windowsPolicy
	case "graph.macOSCompliancePolicy":
		macOsPolicy := graphmodels.NewMacOSCompliancePolicy()
		if err := constructMacOSCompliancePolicy(ctx, data, macOsPolicy); err != nil {
			return nil, fmt.Errorf("failed to construct macOS policy: %s", err)
		}
		requestBody = macOsPolicy

	default:
		return nil, fmt.Errorf("unsupported compliance policy type: %s", odataType)
	}

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if !data.ScheduledActionsForRule.IsNull() && !data.ScheduledActionsForRule.IsUnknown() {
		scheduledActions, err := constructScheduledActionsForRule(ctx, data.ScheduledActionsForRule)
		if err != nil {
			return nil, fmt.Errorf("failed to construct scheduled actions for rule: %s", err)
		}
		requestBody.SetScheduledActionsForRule(scheduledActions)
	}

	// NOTE: localActions does not have an SDK setter, it's a property that doesn't exist in the base SDK model
	// This property might be platform-specific or deprecated. Skipping for now.

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructAospDeviceOwnerCompliancePolicy handles AOSP Device Owner specific settings using SDK setters
func constructAospDeviceOwnerCompliancePolicy(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.AospDeviceOwnerCompliancePolicy) error {

	convert.FrameworkToGraphString(data.OsMinimumVersion, policy.SetOsMinimumVersion)
	convert.FrameworkToGraphString(data.OsMaximumVersion, policy.SetOsMaximumVersion)
	convert.FrameworkToGraphBool(data.PasswordRequired, policy.SetPasswordRequired)

	if err := convert.FrameworkToGraphEnum(data.PasswordRequiredType,
		graphmodels.ParseAndroidDeviceOwnerRequiredPasswordType, policy.SetPasswordRequiredType); err != nil {
		return fmt.Errorf("failed to set password required type: %s", err)
	}

	if !data.AospDeviceOwnerSettings.IsNull() && !data.AospDeviceOwnerSettings.IsUnknown() {
		var settings AospDeviceOwnerSettingsModel
		diags := data.AospDeviceOwnerSettings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fmt.Errorf("failed to convert AOSP device owner settings")
		}

		convert.FrameworkToGraphString(settings.MinAndroidSecurityPatchLevel, policy.SetMinAndroidSecurityPatchLevel)
		convert.FrameworkToGraphBool(settings.SecurityBlockJailbrokenDevices, policy.SetSecurityBlockJailbrokenDevices)
		convert.FrameworkToGraphBool(settings.StorageRequireEncryption, policy.SetStorageRequireEncryption)
		convert.FrameworkToGraphInt32(settings.PasswordMinimumLength, policy.SetPasswordMinimumLength)
		convert.FrameworkToGraphInt32(settings.PasswordMinutesOfInactivityBeforeLock, policy.SetPasswordMinutesOfInactivityBeforeLock)
	}

	return nil
}

// constructAndroidDeviceOwnerPolicy handles Android Device Owner specific settings using SDK setters
func constructAndroidDeviceOwnerPolicy(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.AndroidDeviceOwnerCompliancePolicy) error {

	convert.FrameworkToGraphString(data.OsMinimumVersion, policy.SetOsMinimumVersion)
	convert.FrameworkToGraphString(data.OsMaximumVersion, policy.SetOsMaximumVersion)
	convert.FrameworkToGraphBool(data.PasswordRequired, policy.SetPasswordRequired)

	if err := convert.FrameworkToGraphEnum(data.PasswordRequiredType,
		graphmodels.ParseAndroidDeviceOwnerRequiredPasswordType, policy.SetPasswordRequiredType); err != nil {
		return fmt.Errorf("failed to set password required type: %s", err)
	}

	if !data.AndroidDeviceOwnerSettings.IsNull() && !data.AndroidDeviceOwnerSettings.IsUnknown() {
		var settings AndroidDeviceOwnerSettingsModel
		diags := data.AndroidDeviceOwnerSettings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fmt.Errorf("failed to convert Android device owner settings")
		}

		convert.FrameworkToGraphString(settings.MinAndroidSecurityPatchLevel, policy.SetMinAndroidSecurityPatchLevel)
		convert.FrameworkToGraphBool(settings.SecurityBlockJailbrokenDevices, policy.SetSecurityBlockJailbrokenDevices)
		convert.FrameworkToGraphBool(settings.StorageRequireEncryption, policy.SetStorageRequireEncryption)
		convert.FrameworkToGraphInt32(settings.PasswordMinimumLength, policy.SetPasswordMinimumLength)
		convert.FrameworkToGraphInt32(settings.PasswordMinutesOfInactivityBeforeLock, policy.SetPasswordMinutesOfInactivityBeforeLock)

		if !settings.DeviceThreatProtectionRequiredSecurityLevel.IsNull() && !settings.DeviceThreatProtectionRequiredSecurityLevel.IsUnknown() {
			if err := convert.FrameworkToGraphEnum(settings.DeviceThreatProtectionRequiredSecurityLevel,
				graphmodels.ParseDeviceThreatProtectionLevel, policy.SetDeviceThreatProtectionRequiredSecurityLevel); err != nil {
				return fmt.Errorf("failed to set device threat protection required security level: %s", err)
			}
		}

		if err := convert.FrameworkToGraphEnum(settings.AdvancedThreatProtectionRequiredSecurityLevel,
			graphmodels.ParseDeviceThreatProtectionLevel, policy.SetAdvancedThreatProtectionRequiredSecurityLevel); err != nil {
			return fmt.Errorf("failed to set advanced threat protection required security level: %s", err)
		}

		convert.FrameworkToGraphInt32(settings.PasswordExpirationDays, policy.SetPasswordExpirationDays)
		convert.FrameworkToGraphInt32(settings.PasswordPreviousPasswordCountToBlock, policy.SetPasswordPreviousPasswordCountToBlock)

		if err := convert.FrameworkToGraphEnum(settings.SecurityRequiredAndroidSafetyNetEvaluationType,
			graphmodels.ParseAndroidSafetyNetEvaluationType, policy.SetSecurityRequiredAndroidSafetyNetEvaluationType); err != nil {
			return fmt.Errorf("failed to set security required android safety net evaluation type: %s", err)
		}

		convert.FrameworkToGraphBool(settings.SecurityRequireIntuneAppIntegrity, policy.SetSecurityRequireIntuneAppIntegrity)
		convert.FrameworkToGraphBool(settings.DeviceThreatProtectionEnabled, policy.SetDeviceThreatProtectionEnabled)
		convert.FrameworkToGraphBool(settings.SecurityRequireSafetyNetAttestationBasicIntegrity, policy.SetSecurityRequireSafetyNetAttestationBasicIntegrity)
		convert.FrameworkToGraphBool(settings.SecurityRequireSafetyNetAttestationCertifiedDevice, policy.SetSecurityRequireSafetyNetAttestationCertifiedDevice)
	}

	return nil
}

// constructIosCompliancePolicy handles iOS specific settings using SDK setters
func constructIosCompliancePolicy(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.IosCompliancePolicy) error {

	convert.FrameworkToGraphString(data.OsMinimumVersion, policy.SetOsMinimumVersion)
	convert.FrameworkToGraphString(data.OsMaximumVersion, policy.SetOsMaximumVersion)
	convert.FrameworkToGraphBool(data.PasswordRequired, policy.SetPasscodeRequired)

	if !data.IosSettings.IsNull() && !data.IosSettings.IsUnknown() {
		var settings IosSettingsModel
		diags := data.IosSettings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fmt.Errorf("failed to convert iOS settings")
		}

		if err := convert.FrameworkToGraphEnum(settings.DeviceThreatProtectionRequiredSecurityLevel,
			graphmodels.ParseDeviceThreatProtectionLevel, policy.SetDeviceThreatProtectionRequiredSecurityLevel); err != nil {
			return fmt.Errorf("failed to set device threat protection required security level: %s", err)
		}

		if err := convert.FrameworkToGraphEnum(settings.AdvancedThreatProtectionRequiredSecurityLevel,
			graphmodels.ParseDeviceThreatProtectionLevel, policy.SetAdvancedThreatProtectionRequiredSecurityLevel); err != nil {
			return fmt.Errorf("failed to set advanced threat protection required security level: %s", err)
		}

		convert.FrameworkToGraphBool(settings.DeviceThreatProtectionEnabled, policy.SetDeviceThreatProtectionEnabled)

		if err := convert.FrameworkToGraphEnum(settings.PasscodeRequiredType,
			graphmodels.ParseRequiredPasswordType, policy.SetPasscodeRequiredType); err != nil {
			return fmt.Errorf("failed to set passcode required type: %s", err)
		}

		convert.FrameworkToGraphBool(settings.ManagedEmailProfileRequired, policy.SetManagedEmailProfileRequired)
		convert.FrameworkToGraphBool(settings.SecurityBlockJailbrokenDevices, policy.SetSecurityBlockJailbrokenDevices)
		convert.FrameworkToGraphString(settings.OsMinimumBuildVersion, policy.SetOsMinimumBuildVersion)
		convert.FrameworkToGraphString(settings.OsMaximumBuildVersion, policy.SetOsMaximumBuildVersion)
		convert.FrameworkToGraphInt32(settings.PasscodeMinimumCharacterSetCount, policy.SetPasscodeMinimumCharacterSetCount)
		convert.FrameworkToGraphInt32(settings.PasscodeMinutesOfInactivityBeforeLock, policy.SetPasscodeMinutesOfInactivityBeforeLock)
		convert.FrameworkToGraphInt32(settings.PasscodeMinutesOfInactivityBeforeScreenTimeout, policy.SetPasscodeMinutesOfInactivityBeforeScreenTimeout)
		convert.FrameworkToGraphInt32(settings.PasscodeExpirationDays, policy.SetPasscodeExpirationDays)
		convert.FrameworkToGraphInt32(settings.PasscodePreviousPasscodeBlockCount, policy.SetPasscodePreviousPasscodeBlockCount)

		if !settings.RestrictedApps.IsNull() && !settings.RestrictedApps.IsUnknown() {
			restrictedApps, err := constructRestrictedAppsSDK(ctx, settings.RestrictedApps)
			if err != nil {
				return fmt.Errorf("failed to construct restricted apps: %s", err)
			}
			policy.SetRestrictedApps(restrictedApps)
		}
	}

	return nil
}

// constructRestrictedAppsSDK converts Terraform List to Graph SDK model using proper SDK types
func constructRestrictedAppsSDK(ctx context.Context, restrictedAppsData types.List) ([]graphmodels.AppListItemable, error) {
	var restrictedAppsModels []RestrictedAppModel
	diags := restrictedAppsData.ElementsAs(ctx, &restrictedAppsModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert restricted apps")
	}

	restrictedApps := make([]graphmodels.AppListItemable, 0, len(restrictedAppsModels))
	for _, app := range restrictedAppsModels {
		appItem := graphmodels.NewAppListItem()

		convert.FrameworkToGraphString(app.Name, appItem.SetName)
		convert.FrameworkToGraphString(app.AppId, appItem.SetAppId)

		restrictedApps = append(restrictedApps, appItem)
	}

	return restrictedApps, nil
}

// constructMacOSCompliancePolicy handles macOS specific settings using SDK setters
func constructMacOSCompliancePolicy(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.MacOSCompliancePolicy) error {
	convert.FrameworkToGraphString(data.OsMinimumVersion, policy.SetOsMinimumVersion)
	convert.FrameworkToGraphString(data.OsMaximumVersion, policy.SetOsMaximumVersion)
	convert.FrameworkToGraphBool(data.PasswordRequired, policy.SetPasswordRequired)

	if err := convert.FrameworkToGraphEnum(data.PasswordRequiredType,
		graphmodels.ParseRequiredPasswordType, policy.SetPasswordRequiredType); err != nil {
		return fmt.Errorf("failed to set password required type: %s", err)
	}

	if !data.MacOsSettings.IsNull() && !data.MacOsSettings.IsUnknown() {
		var settings MacOsSettingsModel
		diags := data.MacOsSettings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fmt.Errorf("failed to convert macOS settings")
		}

		if err := convert.FrameworkToGraphEnum(settings.GatekeeperAllowedAppSource,
			graphmodels.ParseMacOSGatekeeperAppSources, policy.SetGatekeeperAllowedAppSource); err != nil {
			return fmt.Errorf("failed to set gatekeeper allowed app source: %s", err)
		}

		convert.FrameworkToGraphBool(settings.SystemIntegrityProtectionEnabled, policy.SetSystemIntegrityProtectionEnabled)
		convert.FrameworkToGraphString(settings.OsMinimumBuildVersion, policy.SetOsMinimumBuildVersion)
		convert.FrameworkToGraphString(settings.OsMaximumBuildVersion, policy.SetOsMaximumBuildVersion)
		convert.FrameworkToGraphBool(settings.PasswordBlockSimple, policy.SetPasswordBlockSimple)
		convert.FrameworkToGraphInt32(settings.PasswordMinimumCharacterSetCount, policy.SetPasswordMinimumCharacterSetCount)
		convert.FrameworkToGraphInt32(settings.PasswordMinutesOfInactivityBeforeLock, policy.SetPasswordMinutesOfInactivityBeforeLock)
		convert.FrameworkToGraphBool(settings.StorageRequireEncryption, policy.SetStorageRequireEncryption)
		convert.FrameworkToGraphBool(settings.FirewallEnabled, policy.SetFirewallEnabled)
		convert.FrameworkToGraphBool(settings.FirewallBlockAllIncoming, policy.SetFirewallBlockAllIncoming)
		convert.FrameworkToGraphBool(settings.FirewallEnableStealthMode, policy.SetFirewallEnableStealthMode)
	}

	return nil
}

// constructWindows10CompliancePolicy handles Windows 10 specific settings using SDK setters
func constructWindows10CompliancePolicy(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) error {
	convert.FrameworkToGraphString(data.OsMinimumVersion, policy.SetOsMinimumVersion)
	convert.FrameworkToGraphString(data.OsMaximumVersion, policy.SetOsMaximumVersion)
	convert.FrameworkToGraphBool(data.PasswordRequired, policy.SetPasswordRequired)

	if err := convert.FrameworkToGraphEnum(data.PasswordRequiredType,
		graphmodels.ParseRequiredPasswordType, policy.SetPasswordRequiredType); err != nil {
		return fmt.Errorf("failed to set password required type: %s", err)
	}

	if !data.Windows10Settings.IsNull() && !data.Windows10Settings.IsUnknown() {
		var settings Windows10SettingsModel
		diags := data.Windows10Settings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fmt.Errorf("failed to convert Windows 10 settings")
		}

		if err := convert.FrameworkToGraphEnum(settings.DeviceThreatProtectionRequiredSecurityLevel,
			graphmodels.ParseDeviceThreatProtectionLevel, policy.SetDeviceThreatProtectionRequiredSecurityLevel); err != nil {
			return fmt.Errorf("failed to set device threat protection required security level: %s", err)
		}

		convert.FrameworkToGraphBool(settings.PasswordBlockSimple, policy.SetPasswordBlockSimple)
		convert.FrameworkToGraphBool(settings.PasswordRequiredToUnlockFromIdle, policy.SetPasswordRequiredToUnlockFromIdle)
		convert.FrameworkToGraphBool(settings.StorageRequireEncryption, policy.SetStorageRequireEncryption)
		convert.FrameworkToGraphInt32(settings.PasswordMinutesOfInactivityBeforeLock, policy.SetPasswordMinutesOfInactivityBeforeLock)
		convert.FrameworkToGraphInt32(settings.PasswordMinimumCharacterSetCount, policy.SetPasswordMinimumCharacterSetCount)
		convert.FrameworkToGraphBool(settings.ActiveFirewallRequired, policy.SetActiveFirewallRequired)
		convert.FrameworkToGraphBool(settings.TpmRequired, policy.SetTpmRequired)
		convert.FrameworkToGraphBool(settings.AntivirusRequired, policy.SetAntivirusRequired)
		convert.FrameworkToGraphBool(settings.AntiSpywareRequired, policy.SetAntiSpywareRequired)
		convert.FrameworkToGraphBool(settings.DefenderEnabled, policy.SetDefenderEnabled)
		convert.FrameworkToGraphBool(settings.SignatureOutOfDate, policy.SetSignatureOutOfDate)
		convert.FrameworkToGraphBool(settings.RtpEnabled, policy.SetRtpEnabled)
		convert.FrameworkToGraphString(settings.DefenderVersion, policy.SetDefenderVersion)
		convert.FrameworkToGraphBool(settings.ConfigurationManagerComplianceRequired, policy.SetConfigurationManagerComplianceRequired)
		convert.FrameworkToGraphString(settings.OsMinimumVersion, policy.SetOsMinimumVersion)
		convert.FrameworkToGraphString(settings.OsMaximumVersion, policy.SetOsMaximumVersion)
		convert.FrameworkToGraphString(settings.MobileOsMinimumVersion, policy.SetMobileOsMinimumVersion)
		convert.FrameworkToGraphString(settings.MobileOsMaximumVersion, policy.SetMobileOsMaximumVersion)
		convert.FrameworkToGraphBool(settings.SecureBootEnabled, policy.SetSecureBootEnabled)
		convert.FrameworkToGraphBool(settings.BitLockerEnabled, policy.SetBitLockerEnabled)
		convert.FrameworkToGraphBool(settings.CodeIntegrityEnabled, policy.SetCodeIntegrityEnabled)
		convert.FrameworkToGraphBool(settings.DeviceThreatProtectionEnabled, policy.SetDeviceThreatProtectionEnabled)

		if !settings.WslDistributions.IsNull() && !settings.WslDistributions.IsUnknown() {
			wslDistributions, err := constructWslDistributions(ctx, settings.WslDistributions)
			if err != nil {
				return fmt.Errorf("failed to construct WSL distributions: %s", err)
			}
			policy.SetWslDistributions(wslDistributions)
		}
	}

	return nil
}

// constructWslDistributions converts Terraform List to Graph SDK model using proper SDK types
func constructWslDistributions(ctx context.Context, wslDistributionsData types.List) ([]graphmodels.WslDistributionConfigurationable, error) {
	var wslDistributionModels []WslDistributionModel
	diags := wslDistributionsData.ElementsAs(ctx, &wslDistributionModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert WSL distributions")
	}

	wslDistributions := make([]graphmodels.WslDistributionConfigurationable, 0, len(wslDistributionModels))
	for _, wslDist := range wslDistributionModels {
		distribution := graphmodels.NewWslDistributionConfiguration()

		convert.FrameworkToGraphString(wslDist.Distribution, distribution.SetDistribution)
		convert.FrameworkToGraphString(wslDist.MinimumOSVersion, distribution.SetMinimumOSVersion)
		convert.FrameworkToGraphString(wslDist.MaximumOSVersion, distribution.SetMaximumOSVersion)

		wslDistributions = append(wslDistributions, distribution)
	}

	return wslDistributions, nil
}

// constructScheduledActionsForRule converts Terraform Set to Graph SDK model using proper SDK types
func constructScheduledActionsForRule(ctx context.Context, scheduledActionsData types.Set) ([]graphmodels.DeviceComplianceScheduledActionForRuleable, error) {
	var scheduledActionsModels []ScheduledActionForRuleModel
	diags := scheduledActionsData.ElementsAs(ctx, &scheduledActionsModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert scheduled actions for rule")
	}

	scheduledActions := make([]graphmodels.DeviceComplianceScheduledActionForRuleable, 0, len(scheduledActionsModels))
	for _, action := range scheduledActionsModels {
		scheduledAction := graphmodels.NewDeviceComplianceScheduledActionForRule()

		convert.FrameworkToGraphString(action.RuleName, scheduledAction.SetRuleName)

		if !action.ScheduledActionConfigurations.IsNull() && !action.ScheduledActionConfigurations.IsUnknown() {
			configs, err := constructScheduledActionConfigurationsSDK(ctx, action.ScheduledActionConfigurations)
			if err != nil {
				return nil, fmt.Errorf("failed to construct scheduled action configurations: %s", err)
			}
			scheduledAction.SetScheduledActionConfigurations(configs)
		}

		scheduledActions = append(scheduledActions, scheduledAction)
	}

	return scheduledActions, nil
}

// constructScheduledActionConfigurationsSDK converts Terraform List to Graph SDK model using proper SDK types
func constructScheduledActionConfigurationsSDK(ctx context.Context, configurationsData types.List) ([]graphmodels.DeviceComplianceActionItemable, error) {
	var configModels []ScheduledActionConfigurationModel
	diags := configurationsData.ElementsAs(ctx, &configModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert scheduled action configurations")
	}

	configurations := make([]graphmodels.DeviceComplianceActionItemable, 0, len(configModels))
	for _, config := range configModels {
		actionItem := graphmodels.NewDeviceComplianceActionItem()

		if err := convert.FrameworkToGraphEnum(config.ActionType,
			graphmodels.ParseDeviceComplianceActionType, actionItem.SetActionType); err != nil {
			return nil, fmt.Errorf("failed to set action type: %s", err)
		}

		convert.FrameworkToGraphInt32(config.GracePeriodHours, actionItem.SetGracePeriodHours)
		convert.FrameworkToGraphString(config.NotificationTemplateId, actionItem.SetNotificationTemplateId)

		if err := convert.FrameworkToGraphStringSet(ctx, config.NotificationMessageCcList, actionItem.SetNotificationMessageCCList); err != nil {
			return nil, fmt.Errorf("failed to set notification message CC list: %s", err)
		}

		configurations = append(configurations, actionItem)
	}

	return configurations, nil
}
