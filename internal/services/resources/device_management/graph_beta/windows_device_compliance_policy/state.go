package graphBetaWindowsDeviceCompliancePolicies

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

	// This resource only handles Windows 10 compliance policies
	if windowsPolicy, ok := remoteResource.(*graphmodels.Windows10CompliancePolicy); ok {
		mapWindows10CompliancePolicyToState(ctx, data, windowsPolicy)
	} else {
		tflog.Error(ctx, "Remote resource is not a Windows 10 compliance policy")
		return
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
		// We no longer have LocalActions in the model, so we don't need to map them
		// This comment is kept for reference
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state for resource %s with id %s", ResourceName, data.ID.ValueString()))
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

// mapWindows10SettingsToState maps Windows 10 specific settings using SDK getters.
func mapWindows10SettingsToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) {
	// Map all properties directly to the model

	// Password-related properties
	data.PasswordRequired = convert.GraphToFrameworkBool(policy.GetPasswordRequired())
	data.PasswordBlockSimple = convert.GraphToFrameworkBool(policy.GetPasswordBlockSimple())
	data.PasswordRequiredToUnlockFromIdle = convert.GraphToFrameworkBool(policy.GetPasswordRequiredToUnlockFromIdle())
	data.PasswordMinutesOfInactivityBeforeLock = convert.GraphToFrameworkInt32(policy.GetPasswordMinutesOfInactivityBeforeLock())
	data.PasswordExpirationDays = convert.GraphToFrameworkInt32(policy.GetPasswordExpirationDays())
	data.PasswordMinimumLength = convert.GraphToFrameworkInt32(policy.GetPasswordMinimumLength())
	data.PasswordMinimumCharacterSetCount = convert.GraphToFrameworkInt32(policy.GetPasswordMinimumCharacterSetCount())
	data.PasswordRequiredType = convert.GraphToFrameworkEnum(policy.GetPasswordRequiredType())
	data.PasswordPreviousPasswordBlockCount = convert.GraphToFrameworkInt32(policy.GetPasswordPreviousPasswordBlockCount())

	// Device health and attestation properties
	data.RequireHealthyDeviceReport = convert.GraphToFrameworkBool(policy.GetRequireHealthyDeviceReport())
	data.EarlyLaunchAntiMalwareDriverEnabled = convert.GraphToFrameworkBool(policy.GetEarlyLaunchAntiMalwareDriverEnabled())
	data.BitLockerEnabled = convert.GraphToFrameworkBool(policy.GetBitLockerEnabled())
	data.SecureBootEnabled = convert.GraphToFrameworkBool(policy.GetSecureBootEnabled())
	data.CodeIntegrityEnabled = convert.GraphToFrameworkBool(policy.GetCodeIntegrityEnabled())
	data.MemoryIntegrityEnabled = convert.GraphToFrameworkBool(policy.GetMemoryIntegrityEnabled())
	data.KernelDmaProtectionEnabled = convert.GraphToFrameworkBool(policy.GetKernelDmaProtectionEnabled())
	data.VirtualizationBasedSecurityEnabled = convert.GraphToFrameworkBool(policy.GetVirtualizationBasedSecurityEnabled())
	data.FirmwareProtectionEnabled = convert.GraphToFrameworkBool(policy.GetFirmwareProtectionEnabled())

	// Security and compliance properties
	data.StorageRequireEncryption = convert.GraphToFrameworkBool(policy.GetStorageRequireEncryption())
	data.ActiveFirewallRequired = convert.GraphToFrameworkBool(policy.GetActiveFirewallRequired())
	data.DefenderEnabled = convert.GraphToFrameworkBool(policy.GetDefenderEnabled())
	data.DefenderVersion = convert.GraphToFrameworkString(policy.GetDefenderVersion())
	data.SignatureOutOfDate = convert.GraphToFrameworkBool(policy.GetSignatureOutOfDate())
	data.RtpEnabled = convert.GraphToFrameworkBool(policy.GetRtpEnabled())
	data.AntivirusRequired = convert.GraphToFrameworkBool(policy.GetAntivirusRequired())
	data.AntiSpywareRequired = convert.GraphToFrameworkBool(policy.GetAntiSpywareRequired())
	data.DeviceThreatProtectionEnabled = convert.GraphToFrameworkBool(policy.GetDeviceThreatProtectionEnabled())
	data.DeviceThreatProtectionRequiredSecurityLevel = convert.GraphToFrameworkEnum(policy.GetDeviceThreatProtectionRequiredSecurityLevel())
	data.ConfigurationManagerComplianceRequired = convert.GraphToFrameworkBool(policy.GetConfigurationManagerComplianceRequired())
	data.TpmRequired = convert.GraphToFrameworkBool(policy.GetTpmRequired())

	// Version and OS properties - already mapped in mapWindows10CompliancePolicyToState
	data.MobileOsMinimumVersion = convert.GraphToFrameworkString(policy.GetMobileOsMinimumVersion())
	data.MobileOsMaximumVersion = convert.GraphToFrameworkString(policy.GetMobileOsMaximumVersion())

	// Map valid operating system build ranges
	data.ValidOperatingSystemBuildRanges = mapValidOperatingSystemVersionRange(ctx, policy.GetValidOperatingSystemBuildRanges())

	// Map WSL distributions
	data.WslDistributions = mapWslDistribution(ctx, policy.GetWslDistributions())

	// Map device compliance policy script
	// This would need special handling based on the actual structure
	// For now, we'll leave it as null/empty

	// Custom compliance required - this might need to be derived from the presence of a script
	data.CustomComplianceRequired = types.BoolValue(false)
	if policy.GetDeviceCompliancePolicyScript() != nil {
		data.CustomComplianceRequired = types.BoolValue(true)

		// Create the device compliance policy script object
		scriptType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"device_compliance_script_id": types.StringType,
				"rules_content":               types.StringType,
			},
		}

		// Get the rules content as a byte array and convert to string
		var rulesContentStr string
		rulesContent := policy.GetDeviceCompliancePolicyScript().GetRulesContent()
		if rulesContent != nil {
			// The SDK returns the rules content as a byte array
			// We convert it to a string for storage in the Terraform state
			rulesContentStr = string(rulesContent)
		}

		scriptAttrs := map[string]attr.Value{
			"device_compliance_script_id": convert.GraphToFrameworkString(policy.GetDeviceCompliancePolicyScript().GetDeviceComplianceScriptId()),
			"rules_content":               types.StringValue(rulesContentStr),
		}

		scriptObj, diags := types.ObjectValue(scriptType.AttrTypes, scriptAttrs)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create device compliance policy script object", map[string]interface{}{
				"error": diags.Errors(),
			})
		} else {
			data.DeviceCompliancePolicyScript = scriptObj
		}
	}
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

// mapWslDistribution maps WSL distributions from SDK to state.
func mapWslDistribution(ctx context.Context, wslDistributions []graphmodels.WslDistributionConfigurationable) types.Set {
	wslDistributionType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"distribution":       types.StringType,
			"minimum_os_version": types.StringType,
			"maximum_os_version": types.StringType,
		},
	}

	if len(wslDistributions) == 0 {
		return types.SetNull(wslDistributionType)
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

	set, diags := types.SetValue(wslDistributionType, wslDistributionValues)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create WSL distributions set from SDK", map[string]interface{}{
			"error": diags.Errors(),
		})
		return types.SetNull(wslDistributionType)
	}
	return set
}

// mapValidOperatingSystemVersionRange maps valid operating system build ranges from SDK to state.
func mapValidOperatingSystemVersionRange(ctx context.Context, buildRanges []graphmodels.OperatingSystemVersionRangeable) types.List {
	validOSBuildRangeType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"low_os_version":  types.StringType,
			"high_os_version": types.StringType,
		},
	}

	if len(buildRanges) == 0 {
		return types.ListNull(validOSBuildRangeType)
	}

	buildRangeValues := make([]attr.Value, 0, len(buildRanges))

	for _, buildRange := range buildRanges {
		buildRangeAttrs := map[string]attr.Value{
			"low_os_version":  convert.GraphToFrameworkString(buildRange.GetLowestVersion()),
			"high_os_version": convert.GraphToFrameworkString(buildRange.GetHighestVersion()),
		}

		buildRangeValue, _ := types.ObjectValue(validOSBuildRangeType.AttrTypes, buildRangeAttrs)
		buildRangeValues = append(buildRangeValues, buildRangeValue)
	}

	list, diags := types.ListValue(validOSBuildRangeType, buildRangeValues)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create valid operating system build ranges list from SDK", map[string]interface{}{
			"error": diags.Errors(),
		})
		return types.ListNull(validOSBuildRangeType)
	}
	return list
}
