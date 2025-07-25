package graphBetaDeviceCompliancePolicies

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote GraphServiceClient object to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *DeviceCompliancePolicyResourceModel, remoteResource graphmodels.DeviceCompliancePolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.Type = convert.GraphToFrameworkString(remoteResource.GetOdataType())

	odataType := data.Type.ValueString()
	switch odataType {
	case "aospDeviceOwnerCompliancePolicy":
		if aospPolicy, ok := remoteResource.(*graphmodels.AospDeviceOwnerCompliancePolicy); ok {
			mapAospDeviceOwnerCompliancePolicyToState(ctx, data, aospPolicy)
		}
	case "androidDeviceOwnerCompliancePolicy":
		if androidPolicy, ok := remoteResource.(*graphmodels.AndroidDeviceOwnerCompliancePolicy); ok {
			mapAndroidDeviceOwnerCompliancePolicyToState(ctx, data, androidPolicy)
		}
	case "#microsoft.graph.iosCompliancePolicy":
		if iosPolicy, ok := remoteResource.(*graphmodels.IosCompliancePolicy); ok {
			mapIosCompliancePolicyToState(ctx, data, iosPolicy)
		}
	case "#microsoft.graph.windows10CompliancePolicy":
		if windowsPolicy, ok := remoteResource.(*graphmodels.Windows10CompliancePolicy); ok {
			mapWindows10CompliancePolicyToState(ctx, data, windowsPolicy)
		}
	default:
		tflog.Warn(ctx, fmt.Sprintf("Unsupported compliance policy type: %s, using additionalData", odataType))
	}

	// Map scheduled actions using SDK getters
	if scheduledActions := remoteResource.GetScheduledActionsForRule(); scheduledActions != nil {
		mappedScheduledActions, err := mapScheduledActionsForRuleToState(ctx, scheduledActions)
		if err != nil {
			tflog.Error(ctx, "Failed to map scheduled actions for rule", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			data.ScheduledActionsForRule = mappedScheduledActions
		}
	}

	// Map local actions from additionalData only if no SDK getter exists
	if additionalData := remoteResource.GetAdditionalData(); additionalData != nil {
		mapLocalActionsToState(ctx, data, additionalData)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state for resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapAospDeviceOwnerCompliancePolicyToState is a responder function that maps AOSP Device Owner compliance policy properties.
func mapAospDeviceOwnerCompliancePolicyToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.AospDeviceOwnerCompliancePolicy) {
	// Map common properties
	data.OsMinimumVersion = convert.GraphToFrameworkString(policy.GetOsMinimumVersion())
	data.OsMaximumVersion = convert.GraphToFrameworkString(policy.GetOsMaximumVersion())
	data.PasswordRequired = convert.GraphToFrameworkBool(policy.GetPasswordRequired())
	data.PasswordRequiredType = convert.GraphToFrameworkEnum(policy.GetPasswordRequiredType())

	// Map AOSP-specific settings
	mapAospDeviceOwnerSettingsToState(ctx, data, policy)
}

// mapAndroidDeviceOwnerCompliancePolicyToState is a responder function that maps Android Device Owner compliance policy properties.
func mapAndroidDeviceOwnerCompliancePolicyToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.AndroidDeviceOwnerCompliancePolicy) {
	// Map common properties
	data.OsMinimumVersion = convert.GraphToFrameworkString(policy.GetOsMinimumVersion())
	data.OsMaximumVersion = convert.GraphToFrameworkString(policy.GetOsMaximumVersion())
	data.PasswordRequired = convert.GraphToFrameworkBool(policy.GetPasswordRequired())
	data.PasswordRequiredType = convert.GraphToFrameworkEnum(policy.GetPasswordRequiredType())

	// Map Android Device Owner-specific settings
	mapAndroidDeviceOwnerSettingsToState(ctx, data, policy)
}

// mapIosCompliancePolicyToState is a responder function that maps iOS compliance policy properties.
func mapIosCompliancePolicyToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.IosCompliancePolicy) {
	// Map common properties
	data.OsMinimumVersion = convert.GraphToFrameworkString(policy.GetOsMinimumVersion())
	data.OsMaximumVersion = convert.GraphToFrameworkString(policy.GetOsMaximumVersion())
	// iOS uses passcodeRequired instead of passwordRequired
	data.PasswordRequired = convert.GraphToFrameworkBool(policy.GetPasscodeRequired())

	// Map iOS-specific settings
	mapIosSettingsToState(ctx, data, policy)
}

// mapWindows10CompliancePolicyToState is a responder function that maps Windows 10 compliance policy properties.
func mapWindows10CompliancePolicyToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) {
	// Map common properties
	data.OsMinimumVersion = convert.GraphToFrameworkString(policy.GetOsMinimumVersion())
	data.OsMaximumVersion = convert.GraphToFrameworkString(policy.GetOsMaximumVersion())
	data.PasswordRequired = convert.GraphToFrameworkBool(policy.GetPasswordRequired())
	data.PasswordRequiredType = convert.GraphToFrameworkEnum(policy.GetPasswordRequiredType())

	// Map Windows 10-specific settings
	mapWindows10SettingsToState(ctx, data, policy)
}

// mapAospDeviceOwnerSettingsToState maps AOSP Device Owner specific settings using SDK getters.
func mapAospDeviceOwnerSettingsToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.AospDeviceOwnerCompliancePolicy) {
	settingsType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_android_security_patch_level":           types.StringType,
			"security_block_jailbroken_devices":          types.BoolType,
			"storage_require_encryption":                 types.BoolType,
			"password_minimum_length":                    types.Int32Type,
			"password_minutes_of_inactivity_before_lock": types.Int32Type,
		},
	}

	settingsAttrs := map[string]attr.Value{
		"min_android_security_patch_level":           convert.GraphToFrameworkString(policy.GetMinAndroidSecurityPatchLevel()),
		"security_block_jailbroken_devices":          convert.GraphToFrameworkBool(policy.GetSecurityBlockJailbrokenDevices()),
		"storage_require_encryption":                 convert.GraphToFrameworkBool(policy.GetStorageRequireEncryption()),
		"password_minimum_length":                    convert.GraphToFrameworkInt32(policy.GetPasswordMinimumLength()),
		"password_minutes_of_inactivity_before_lock": convert.GraphToFrameworkInt32(policy.GetPasswordMinutesOfInactivityBeforeLock()),
	}

	data.AospDeviceOwnerSettings, _ = types.ObjectValue(settingsType.AttrTypes, settingsAttrs)
}

// mapAndroidDeviceOwnerSettingsToState maps Android Device Owner specific settings using SDK getters.
func mapAndroidDeviceOwnerSettingsToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.AndroidDeviceOwnerCompliancePolicy) {
	settingsType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_android_security_patch_level":                         types.StringType,
			"security_block_jailbroken_devices":                        types.BoolType,
			"storage_require_encryption":                               types.BoolType,
			"password_minimum_length":                                  types.Int32Type,
			"password_minutes_of_inactivity_before_lock":               types.Int32Type,
			"device_threat_protection_required_security_level":         types.StringType,
			"advanced_threat_protection_required_security_level":       types.StringType,
			"password_expiration_days":                                 types.Int32Type,
			"password_previous_password_count_to_block":                types.Int32Type,
			"security_required_android_safety_net_evaluation_type":     types.StringType,
			"security_require_intune_app_integrity":                    types.BoolType,
			"device_threat_protection_enabled":                         types.BoolType,
			"security_require_safety_net_attestation_basic_integrity":  types.BoolType,
			"security_require_safety_net_attestation_certified_device": types.BoolType,
		},
	}

	settingsAttrs := map[string]attr.Value{
		"min_android_security_patch_level":                         convert.GraphToFrameworkString(policy.GetMinAndroidSecurityPatchLevel()),
		"security_block_jailbroken_devices":                        convert.GraphToFrameworkBool(policy.GetSecurityBlockJailbrokenDevices()),
		"storage_require_encryption":                               convert.GraphToFrameworkBool(policy.GetStorageRequireEncryption()),
		"password_minimum_length":                                  convert.GraphToFrameworkInt32(policy.GetPasswordMinimumLength()),
		"password_minutes_of_inactivity_before_lock":               convert.GraphToFrameworkInt32(policy.GetPasswordMinutesOfInactivityBeforeLock()),
		"device_threat_protection_required_security_level":         convert.GraphToFrameworkEnum(policy.GetDeviceThreatProtectionRequiredSecurityLevel()),
		"advanced_threat_protection_required_security_level":       convert.GraphToFrameworkEnum(policy.GetAdvancedThreatProtectionRequiredSecurityLevel()),
		"password_expiration_days":                                 convert.GraphToFrameworkInt32(policy.GetPasswordExpirationDays()),
		"password_previous_password_count_to_block":                convert.GraphToFrameworkInt32(policy.GetPasswordPreviousPasswordCountToBlock()),
		"security_required_android_safety_net_evaluation_type":     convert.GraphToFrameworkEnum(policy.GetSecurityRequiredAndroidSafetyNetEvaluationType()),
		"security_require_intune_app_integrity":                    convert.GraphToFrameworkBool(policy.GetSecurityRequireIntuneAppIntegrity()),
		"device_threat_protection_enabled":                         convert.GraphToFrameworkBool(policy.GetDeviceThreatProtectionEnabled()),
		"security_require_safety_net_attestation_basic_integrity":  convert.GraphToFrameworkBool(policy.GetSecurityRequireSafetyNetAttestationBasicIntegrity()),
		"security_require_safety_net_attestation_certified_device": convert.GraphToFrameworkBool(policy.GetSecurityRequireSafetyNetAttestationCertifiedDevice()),
	}

	data.AndroidDeviceOwnerSettings, _ = types.ObjectValue(settingsType.AttrTypes, settingsAttrs)
}

// mapIosSettingsToState maps iOS specific settings using SDK getters.
func mapIosSettingsToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.IosCompliancePolicy) {
	restrictedAppType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":   types.StringType,
			"app_id": types.StringType,
		},
	}

	settingsType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"device_threat_protection_required_security_level":     types.StringType,
			"advanced_threat_protection_required_security_level":   types.StringType,
			"device_threat_protection_enabled":                     types.BoolType,
			"passcode_required_type":                               types.StringType,
			"managed_email_profile_required":                       types.BoolType,
			"security_block_jailbroken_devices":                    types.BoolType,
			"os_minimum_build_version":                             types.StringType,
			"os_maximum_build_version":                             types.StringType,
			"passcode_minimum_character_set_count":                 types.Int32Type,
			"passcode_minutes_of_inactivity_before_lock":           types.Int32Type,
			"passcode_minutes_of_inactivity_before_screen_timeout": types.Int32Type,
			"passcode_expiration_days":                             types.Int32Type,
			"passcode_previous_passcode_block_count":               types.Int32Type,
			"restricted_apps":                                      types.ListType{ElemType: restrictedAppType},
		},
	}

	settingsAttrs := map[string]attr.Value{
		"device_threat_protection_required_security_level":     convert.GraphToFrameworkEnum(policy.GetDeviceThreatProtectionRequiredSecurityLevel()),
		"advanced_threat_protection_required_security_level":   convert.GraphToFrameworkEnum(policy.GetAdvancedThreatProtectionRequiredSecurityLevel()),
		"device_threat_protection_enabled":                     convert.GraphToFrameworkBool(policy.GetDeviceThreatProtectionEnabled()),
		"passcode_required_type":                               convert.GraphToFrameworkEnum(policy.GetPasscodeRequiredType()),
		"managed_email_profile_required":                       convert.GraphToFrameworkBool(policy.GetManagedEmailProfileRequired()),
		"security_block_jailbroken_devices":                    convert.GraphToFrameworkBool(policy.GetSecurityBlockJailbrokenDevices()),
		"os_minimum_build_version":                             convert.GraphToFrameworkString(policy.GetOsMinimumBuildVersion()),
		"os_maximum_build_version":                             convert.GraphToFrameworkString(policy.GetOsMaximumBuildVersion()),
		"passcode_minimum_character_set_count":                 convert.GraphToFrameworkInt32(policy.GetPasscodeMinimumCharacterSetCount()),
		"passcode_minutes_of_inactivity_before_lock":           convert.GraphToFrameworkInt32(policy.GetPasscodeMinutesOfInactivityBeforeLock()),
		"passcode_minutes_of_inactivity_before_screen_timeout": convert.GraphToFrameworkInt32(policy.GetPasscodeMinutesOfInactivityBeforeScreenTimeout()),
		"passcode_expiration_days":                             convert.GraphToFrameworkInt32(policy.GetPasscodeExpirationDays()),
		"passcode_previous_passcode_block_count":               convert.GraphToFrameworkInt32(policy.GetPasscodePreviousPasscodeBlockCount()),
		"restricted_apps":                                      mapRestrictedAppsFromSDK(ctx, policy.GetRestrictedApps()),
	}

	data.IosSettings, _ = types.ObjectValue(settingsType.AttrTypes, settingsAttrs)
}

// mapWindows10SettingsToState maps Windows 10 specific settings using SDK getters.
func mapWindows10SettingsToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) {
	wslDistributionType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"distribution":       types.StringType,
			"minimum_os_version": types.StringType,
			"maximum_os_version": types.StringType,
		},
	}

	settingsType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"device_threat_protection_required_security_level": types.StringType,
			"device_compliance_policy_script":                  types.StringType,
			"password_required_type":                           types.StringType,
			"wsl_distributions":                                types.ListType{ElemType: wslDistributionType},
			"password_required":                                types.BoolType,
			"password_block_simple":                            types.BoolType,
			"password_required_to_unlock_from_idle":            types.BoolType,
			"storage_require_encryption":                       types.BoolType,
			"password_minutes_of_inactivity_before_lock":       types.Int32Type,
			"password_minimum_character_set_count":             types.Int32Type,
			"active_firewall_required":                         types.BoolType,
			"tpm_required":                                     types.BoolType,
			"antivirus_required":                               types.BoolType,
			"anti_spyware_required":                            types.BoolType,
			"defender_enabled":                                 types.BoolType,
			"signature_out_of_date":                            types.BoolType,
			"rtp_enabled":                                      types.BoolType,
			"defender_version":                                 types.StringType,
			"configuration_manager_compliance_required":        types.BoolType,
			"os_minimum_version":                               types.StringType,
			"os_maximum_version":                               types.StringType,
			"mobile_os_minimum_version":                        types.StringType,
			"mobile_os_maximum_version":                        types.StringType,
			"secure_boot_enabled":                              types.BoolType,
			"bit_locker_enabled":                               types.BoolType,
			"code_integrity_enabled":                           types.BoolType,
			"device_threat_protection_enabled":                 types.BoolType,
		},
	}

	settingsAttrs := map[string]attr.Value{
		"device_threat_protection_required_security_level": convert.GraphToFrameworkEnum(policy.GetDeviceThreatProtectionRequiredSecurityLevel()),
		"device_compliance_policy_script":                  convert.GraphToFrameworkString(nil), // This might need special handling
		"password_required_type":                           convert.GraphToFrameworkEnum(policy.GetPasswordRequiredType()),
		"wsl_distributions":                                mapWslDistributionsFromSDK(ctx, policy.GetWslDistributions()),
		"password_required":                                convert.GraphToFrameworkBool(policy.GetPasswordRequired()),
		"password_block_simple":                            convert.GraphToFrameworkBool(policy.GetPasswordBlockSimple()),
		"password_required_to_unlock_from_idle":            convert.GraphToFrameworkBool(policy.GetPasswordRequiredToUnlockFromIdle()),
		"storage_require_encryption":                       convert.GraphToFrameworkBool(policy.GetStorageRequireEncryption()),
		"password_minutes_of_inactivity_before_lock":       convert.GraphToFrameworkInt32(policy.GetPasswordMinutesOfInactivityBeforeLock()),
		"password_minimum_character_set_count":             convert.GraphToFrameworkInt32(policy.GetPasswordMinimumCharacterSetCount()),
		"active_firewall_required":                         convert.GraphToFrameworkBool(policy.GetActiveFirewallRequired()),
		"tpm_required":                                     convert.GraphToFrameworkBool(policy.GetTpmRequired()),
		"antivirus_required":                               convert.GraphToFrameworkBool(policy.GetAntivirusRequired()),
		"anti_spyware_required":                            convert.GraphToFrameworkBool(policy.GetAntiSpywareRequired()),
		"defender_enabled":                                 convert.GraphToFrameworkBool(policy.GetDefenderEnabled()),
		"signature_out_of_date":                            convert.GraphToFrameworkBool(policy.GetSignatureOutOfDate()),
		"rtp_enabled":                                      convert.GraphToFrameworkBool(policy.GetRtpEnabled()),
		"defender_version":                                 convert.GraphToFrameworkString(policy.GetDefenderVersion()),
		"configuration_manager_compliance_required":        convert.GraphToFrameworkBool(policy.GetConfigurationManagerComplianceRequired()),
		"os_minimum_version":                               convert.GraphToFrameworkString(policy.GetOsMinimumVersion()),
		"os_maximum_version":                               convert.GraphToFrameworkString(policy.GetOsMaximumVersion()),
		"mobile_os_minimum_version":                        convert.GraphToFrameworkString(policy.GetMobileOsMinimumVersion()),
		"mobile_os_maximum_version":                        convert.GraphToFrameworkString(policy.GetMobileOsMaximumVersion()),
		"secure_boot_enabled":                              convert.GraphToFrameworkBool(policy.GetSecureBootEnabled()),
		"bit_locker_enabled":                               convert.GraphToFrameworkBool(policy.GetBitLockerEnabled()),
		"code_integrity_enabled":                           convert.GraphToFrameworkBool(policy.GetCodeIntegrityEnabled()),
		"device_threat_protection_enabled":                 convert.GraphToFrameworkBool(policy.GetDeviceThreatProtectionEnabled()),
	}

	data.Windows10Settings, _ = types.ObjectValue(settingsType.AttrTypes, settingsAttrs)
}

// mapScheduledActionsForRuleToState maps scheduled actions for rule from SDK to state.
func mapScheduledActionsForRuleToState(ctx context.Context, scheduledActions []graphmodels.DeviceComplianceScheduledActionForRuleable) (types.Set, error) {
	scheduledActionsType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"rule_name": types.StringType,
			"scheduled_action_configurations": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"action_type":                  types.StringType,
						"grace_period_hours":           types.Int32Type,
						"notification_template_id":     types.StringType,
						"notification_message_cc_list": types.SetType{ElemType: types.StringType},
					},
				},
			},
		},
	}

	scheduledActionsValues := make([]attr.Value, 0, len(scheduledActions))

	for _, action := range scheduledActions {
		actionAttrs := map[string]attr.Value{
			"rule_name": convert.GraphToFrameworkString(action.GetRuleName()),
			"scheduled_action_configurations": types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"action_type":                  types.StringType,
					"grace_period_hours":           types.Int32Type,
					"notification_template_id":     types.StringType,
					"notification_message_cc_list": types.SetType{ElemType: types.StringType},
				},
			}),
		}

		if configs := action.GetScheduledActionConfigurations(); configs != nil {
			mappedConfigs, err := mapScheduledActionConfigurationsToState(ctx, configs)
			if err != nil {
				return types.SetNull(scheduledActionsType), err
			}
			actionAttrs["scheduled_action_configurations"] = mappedConfigs
		}

		actionValue, _ := types.ObjectValue(scheduledActionsType.AttrTypes, actionAttrs)
		scheduledActionsValues = append(scheduledActionsValues, actionValue)
	}

	set, diags := types.SetValue(scheduledActionsType, scheduledActionsValues)
	if diags.HasError() {
		return types.SetNull(scheduledActionsType), fmt.Errorf("failed to create scheduled actions set")
	}
	return set, nil
}

// mapScheduledActionConfigurationsToState maps scheduled action configurations from SDK to state.
func mapScheduledActionConfigurationsToState(ctx context.Context, configurations []graphmodels.DeviceComplianceActionItemable) (types.List, error) {
	configurationType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"action_type":                  types.StringType,
			"grace_period_hours":           types.Int32Type,
			"notification_template_id":     types.StringType,
			"notification_message_cc_list": types.SetType{ElemType: types.StringType},
		},
	}

	configValues := make([]attr.Value, 0, len(configurations))

	for _, config := range configurations {
		configAttrs := map[string]attr.Value{
			"action_type":                  convert.GraphToFrameworkEnum(config.GetActionType()),
			"grace_period_hours":           convert.GraphToFrameworkInt32(config.GetGracePeriodHours()),
			"notification_template_id":     convert.GraphToFrameworkString(config.GetNotificationTemplateId()),
			"notification_message_cc_list": convert.GraphToFrameworkStringSet(ctx, config.GetNotificationMessageCCList()),
		}

		configValue, _ := types.ObjectValue(configurationType.AttrTypes, configAttrs)
		configValues = append(configValues, configValue)
	}

	list, diags := types.ListValue(configurationType, configValues)
	if diags.HasError() {
		return types.ListNull(configurationType), fmt.Errorf("failed to create scheduled action configurations list")
	}
	return list, nil
}

// mapLocalActionsToState maps local actions from additional data to state.
func mapLocalActionsToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, additionalData map[string]interface{}) {
	if localActionsData, ok := additionalData["localActions"].([]interface{}); ok {
		localActions := make([]attr.Value, 0, len(localActionsData))
		for _, action := range localActionsData {
			if actionStr, ok := action.(string); ok {
				localActions = append(localActions, convert.GraphToFrameworkString(&actionStr))
			}
		}
		data.LocalActions, _ = types.ListValue(types.StringType, localActions)
	}
}

// mapRestrictedAppsFromSDK maps restricted apps from SDK to state.
func mapRestrictedAppsFromSDK(ctx context.Context, restrictedApps []graphmodels.AppListItemable) types.List {
	restrictedAppType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":   types.StringType,
			"app_id": types.StringType,
		},
	}

	if restrictedApps == nil || len(restrictedApps) == 0 {
		return types.ListNull(restrictedAppType)
	}

	restrictedAppsValues := make([]attr.Value, 0, len(restrictedApps))

	for _, app := range restrictedApps {
		appAttrs := map[string]attr.Value{
			"name":   convert.GraphToFrameworkString(app.GetName()),
			"app_id": convert.GraphToFrameworkString(app.GetAppId()),
		}

		appValue, _ := types.ObjectValue(restrictedAppType.AttrTypes, appAttrs)
		restrictedAppsValues = append(restrictedAppsValues, appValue)
	}

	list, diags := types.ListValue(restrictedAppType, restrictedAppsValues)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create restricted apps list from SDK", map[string]interface{}{
			"error": diags.Errors(),
		})
		return types.ListNull(restrictedAppType)
	}
	return list
}

// mapWslDistributionsFromSDK maps WSL distributions from SDK to state.
func mapWslDistributionsFromSDK(ctx context.Context, wslDistributions []graphmodels.WslDistributionConfigurationable) types.List {
	wslDistributionType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"distribution":       types.StringType,
			"minimum_os_version": types.StringType,
			"maximum_os_version": types.StringType,
		},
	}

	if len(wslDistributions) == 0 {
		return types.ListNull(wslDistributionType)
	}

	wslDistributionValues := make([]attr.Value, 0, len(wslDistributions))

	for _, wslDist := range wslDistributions {
		wslDistAttrs := map[string]attr.Value{
			"distribution":       convert.GraphToFrameworkString(wslDist.GetDistribution()),
			"minimum_os_version": convert.GraphToFrameworkString(wslDist.GetMinimumOSVersion()),
			"maximum_os_version": convert.GraphToFrameworkString(wslDist.GetMaximumOSVersion()),
		}

		wslDistValue, _ := types.ObjectValue(wslDistributionType.AttrTypes, wslDistAttrs)
		wslDistributionValues = append(wslDistributionValues, wslDistValue)
	}

	list, diags := types.ListValue(wslDistributionType, wslDistributionValues)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create WSL distributions list from SDK", map[string]interface{}{
			"error": diags.Errors(),
		})
		return types.ListNull(wslDistributionType)
	}
	return list
}
