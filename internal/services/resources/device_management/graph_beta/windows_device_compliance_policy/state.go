package graphBetaWindowsDeviceCompliancePolicy

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

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
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

	if scheduledActions := remoteResource.GetScheduledActionsForRule(); scheduledActions != nil {
		mappedScheduledActions, err := mapScheduledActionsForRuleToState(ctx, scheduledActions)
		if err != nil {
			tflog.Error(ctx, "Failed to map scheduled actions for rule", map[string]any{
				"error": err.Error(),
			})
		} else {
			data.ScheduledActionsForRule = mappedScheduledActions
		}
	}

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, setting assignments to null", map[string]any{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(WindowsDeviceCompliancePolicyAssignmentType())
	} else {
		tflog.Debug(ctx, "Starting assignment mapping process", map[string]any{
			"resourceId":      data.ID.ValueString(),
			"assignmentCount": len(assignments),
		})
		MapAssignmentsToTerraform(ctx, data, assignments)
		tflog.Debug(ctx, "Completed assignment mapping process", map[string]any{
			"resourceId": data.ID.ValueString(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapWindows10SettingsToState maps Windows 10 specific settings using SDK getters.
func mapWindows10CompliancePolicyToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) {

	// all of these fields are now deprecated.
	// data.RequireHealthyDeviceReport = convert.GraphToFrameworkBool(policy.GetRequireHealthyDeviceReport())
	// data.EarlyLaunchAntiMalwareDriverEnabled = convert.GraphToFrameworkBool(policy.GetEarlyLaunchAntiMalwareDriverEnabled())
	// data.MemoryIntegrityEnabled = convert.GraphToFrameworkBool(policy.GetMemoryIntegrityEnabled())
	// data.KernelDmaProtectionEnabled = convert.GraphToFrameworkBool(policy.GetKernelDmaProtectionEnabled())
	// data.VirtualizationBasedSecurityEnabled = convert.GraphToFrameworkBool(policy.GetVirtualizationBasedSecurityEnabled())
	// data.FirmwareProtectionEnabled = convert.GraphToFrameworkBool(policy.GetFirmwareProtectionEnabled())

	// Map device_health object only if it was configured by the user
	if !data.DeviceHealth.IsNull() && !data.DeviceHealth.IsUnknown() {
		mapDeviceHealthToState(ctx, data, policy)
	}

	// Map system_security object only if it was configured by the user
	if !data.SystemSecurity.IsNull() && !data.SystemSecurity.IsUnknown() {
		mapSystemSecurityToState(ctx, data, policy)
	}

	// Map microsoft_defender_for_endpoint object only if it was configured by the user
	if !data.MicrosoftDefenderForEndpoint.IsNull() && !data.MicrosoftDefenderForEndpoint.IsUnknown() {
		mapMicrosoftDefenderForEndpointToState(ctx, data, policy)
	}

	// Map device_properties object only if it was configured by the user
	if !data.DeviceProperties.IsNull() && !data.DeviceProperties.IsUnknown() {
		mapDevicePropertiesToState(ctx, data, policy)
	}

	data.WslDistributions = mapWslDistribution(ctx, policy.GetWslDistributions())

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
			rulesContentStr = string(rulesContent)
		}

		scriptAttrs := map[string]attr.Value{
			"device_compliance_script_id": convert.GraphToFrameworkString(policy.GetDeviceCompliancePolicyScript().GetDeviceComplianceScriptId()),
			"rules_content":               types.StringValue(rulesContentStr),
		}

		scriptObj, diags := types.ObjectValue(scriptType.AttrTypes, scriptAttrs)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create device compliance policy script object", map[string]any{
				"error": diags.Errors(),
			})
			data.DeviceCompliancePolicyScript = types.ObjectNull(scriptType.AttrTypes)
		} else {
			data.DeviceCompliancePolicyScript = scriptObj
		}
	} else {
		// Set to null when no script is present
		scriptType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"device_compliance_script_id": types.StringType,
				"rules_content":               types.StringType,
			},
		}
		data.DeviceCompliancePolicyScript = types.ObjectNull(scriptType.AttrTypes)
	}
}

// mapScheduledActionsForRuleToState maps scheduled actions for rule from SDK to state.
func mapScheduledActionsForRuleToState(ctx context.Context, scheduledActions []graphmodels.DeviceComplianceScheduledActionForRuleable) (types.List, error) {
	scheduledActionType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"scheduled_action_configurations": types.SetType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"action_type":                  types.StringType,
					"grace_period_hours":           types.Int32Type,
					"notification_template_id":     types.StringType,
					"notification_message_cc_list": types.ListType{ElemType: types.StringType},
				},
			}},
		},
	}

	if len(scheduledActions) == 0 {
		return types.ListNull(scheduledActionType), nil
	}

	actionValues := make([]attr.Value, 0, len(scheduledActions))

	for _, action := range scheduledActions {
		var mappedConfigs types.Set
		if configs := action.GetScheduledActionConfigurations(); configs != nil {
			var err error
			mappedConfigs, err = mapScheduledActionConfigurationsToState(ctx, configs)
			if err != nil {
				return types.ListNull(scheduledActionType), err
			}
		} else {
			mappedConfigs = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"action_type":                  types.StringType,
					"grace_period_hours":           types.Int32Type,
					"notification_template_id":     types.StringType,
					"notification_message_cc_list": types.ListType{ElemType: types.StringType},
				},
			})
		}

		actionAttrs := map[string]attr.Value{
			"scheduled_action_configurations": mappedConfigs,
		}

		actionValue, _ := types.ObjectValue(scheduledActionType.AttrTypes, actionAttrs)
		actionValues = append(actionValues, actionValue)
	}

	list, diags := types.ListValue(scheduledActionType, actionValues)
	if diags.HasError() {
		return types.ListNull(scheduledActionType), fmt.Errorf("failed to create scheduled actions list")
	}
	return list, nil
}

// mapScheduledActionConfigurationsToState maps scheduled action configurations from SDK to state.
func mapScheduledActionConfigurationsToState(ctx context.Context, configurations []graphmodels.DeviceComplianceActionItemable) (types.Set, error) {
	configurationType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"action_type":                  types.StringType,
			"grace_period_hours":           types.Int32Type,
			"notification_template_id":     types.StringType,
			"notification_message_cc_list": types.ListType{ElemType: types.StringType},
		},
	}

	configValues := make([]attr.Value, 0, len(configurations))

	for _, config := range configurations {
		configAttrs := map[string]attr.Value{
			"action_type":                  convert.GraphToFrameworkEnum(config.GetActionType()),
			"grace_period_hours":           convert.GraphToFrameworkInt32(config.GetGracePeriodHours()),
			"notification_template_id":     convert.GraphToFrameworkString(config.GetNotificationTemplateId()),
			"notification_message_cc_list": convert.GraphToFrameworkStringList(config.GetNotificationMessageCCList()),
		}

		configValue, _ := types.ObjectValue(configurationType.AttrTypes, configAttrs)
		configValues = append(configValues, configValue)
	}

	set, diags := types.SetValue(configurationType, configValues)
	if diags.HasError() {
		return types.SetNull(configurationType), fmt.Errorf("failed to create scheduled action configurations set")
	}
	return set, nil
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
		tflog.Error(ctx, "Failed to create WSL distributions set from SDK", map[string]any{
			"error": diags.Errors(),
		})
		return types.SetNull(wslDistributionType)
	}
	return set
}

// mapValidOperatingSystemVersionRange maps valid operating system build ranges from SDK to state.
func mapValidOperatingSystemVersionRange(ctx context.Context, buildRanges []graphmodels.OperatingSystemVersionRangeable) types.Set {
	validOSBuildRangeType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"low_os_version":  types.StringType,
			"high_os_version": types.StringType,
			"description":     types.StringType,
		},
	}

	if len(buildRanges) == 0 {
		return types.SetNull(validOSBuildRangeType)
	}

	buildRangeValues := make([]attr.Value, 0, len(buildRanges))

	for _, buildRange := range buildRanges {
		buildRangeAttrs := map[string]attr.Value{
			"low_os_version":  convert.GraphToFrameworkString(buildRange.GetLowestVersion()),
			"high_os_version": convert.GraphToFrameworkString(buildRange.GetHighestVersion()),
			"description":     convert.GraphToFrameworkString(buildRange.GetDescription()),
		}

		buildRangeValue, _ := types.ObjectValue(validOSBuildRangeType.AttrTypes, buildRangeAttrs)
		buildRangeValues = append(buildRangeValues, buildRangeValue)
	}

	set, diags := types.SetValue(validOSBuildRangeType, buildRangeValues)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create valid operating system build ranges set from SDK", map[string]any{
			"error": diags.Errors(),
		})
		return types.SetNull(validOSBuildRangeType)
	}
	return set
}

// mapDeviceHealthToState maps device health properties from SDK to state
func mapDeviceHealthToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) {
	data.DeviceHealth, _ = types.ObjectValueFrom(ctx, map[string]attr.Type{
		"bit_locker_enabled":     types.BoolType,
		"secure_boot_enabled":    types.BoolType,
		"code_integrity_enabled": types.BoolType,
	}, DeviceHealthModel{
		BitLockerEnabled:     convert.GraphToFrameworkBool(policy.GetBitLockerEnabled()),
		SecureBootEnabled:    convert.GraphToFrameworkBool(policy.GetSecureBootEnabled()),
		CodeIntegrityEnabled: convert.GraphToFrameworkBool(policy.GetCodeIntegrityEnabled()),
	})
}

// mapDevicePropertiesToState maps device properties from SDK to state
func mapDevicePropertiesToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) {
	data.DeviceProperties, _ = types.ObjectValueFrom(ctx, map[string]attr.Type{
		"os_minimum_version":        types.StringType,
		"os_maximum_version":        types.StringType,
		"mobile_os_minimum_version": types.StringType,
		"mobile_os_maximum_version": types.StringType,
		"valid_operating_system_build_ranges": types.SetType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"low_os_version":  types.StringType,
				"high_os_version": types.StringType,
				"description":     types.StringType,
			},
		}},
	}, DevicePropertiesModel{
		OsMinimumVersion:                convert.GraphToFrameworkString(policy.GetOsMinimumVersion()),
		OsMaximumVersion:                convert.GraphToFrameworkString(policy.GetOsMaximumVersion()),
		MobileOsMinimumVersion:          convert.GraphToFrameworkString(policy.GetMobileOsMinimumVersion()),
		MobileOsMaximumVersion:          convert.GraphToFrameworkString(policy.GetMobileOsMaximumVersion()),
		ValidOperatingSystemBuildRanges: mapValidOperatingSystemVersionRange(ctx, policy.GetValidOperatingSystemBuildRanges()),
	})
}

// mapSystemSecurityToState maps system security settings to state
func mapSystemSecurityToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) {
	systemSecurityType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"active_firewall_required":                  types.BoolType,
			"anti_spyware_required":                     types.BoolType,
			"antivirus_required":                        types.BoolType,
			"configuration_manager_compliance_required": types.BoolType,
			"defender_enabled":                          types.BoolType,
			"defender_version":                          types.StringType,
			"password_block_simple":                     types.BoolType,
			"password_minimum_character_set_count":      types.Int32Type,
			"password_required":                         types.BoolType,
			"password_required_to_unlock_from_idle":     types.BoolType,
			"password_required_type":                    types.StringType,
			"rtp_enabled":                               types.BoolType,
			"signature_out_of_date":                     types.BoolType,
			"storage_require_encryption":                types.BoolType,
			"tpm_required":                              types.BoolType,
		},
	}

	data.SystemSecurity, _ = types.ObjectValue(systemSecurityType.AttrTypes, map[string]attr.Value{
		"active_firewall_required":                  convert.GraphToFrameworkBool(policy.GetActiveFirewallRequired()),
		"anti_spyware_required":                     convert.GraphToFrameworkBool(policy.GetAntiSpywareRequired()),
		"antivirus_required":                        convert.GraphToFrameworkBool(policy.GetAntivirusRequired()),
		"configuration_manager_compliance_required": convert.GraphToFrameworkBool(policy.GetConfigurationManagerComplianceRequired()),
		"defender_enabled":                          convert.GraphToFrameworkBool(policy.GetDefenderEnabled()),
		"defender_version":                          convert.GraphToFrameworkString(policy.GetDefenderVersion()),
		"password_block_simple":                     convert.GraphToFrameworkBool(policy.GetPasswordBlockSimple()),
		"password_minimum_character_set_count":      convert.GraphToFrameworkInt32(policy.GetPasswordMinimumCharacterSetCount()),
		"password_required":                         convert.GraphToFrameworkBool(policy.GetPasswordRequired()),
		"password_required_to_unlock_from_idle":     convert.GraphToFrameworkBool(policy.GetPasswordRequiredToUnlockFromIdle()),
		"password_required_type":                    convert.GraphToFrameworkEnum(policy.GetPasswordRequiredType()),
		"rtp_enabled":                               convert.GraphToFrameworkBool(policy.GetRtpEnabled()),
		"signature_out_of_date":                     convert.GraphToFrameworkBool(policy.GetSignatureOutOfDate()),
		"storage_require_encryption":                convert.GraphToFrameworkBool(policy.GetStorageRequireEncryption()),
		"tpm_required":                              convert.GraphToFrameworkBool(policy.GetTpmRequired()),
	})
}

// mapMicrosoftDefenderForEndpointToState maps Microsoft Defender for Endpoint settings to state
func mapMicrosoftDefenderForEndpointToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) {
	microsoftDefenderType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"device_threat_protection_enabled":                 types.BoolType,
			"device_threat_protection_required_security_level": types.StringType,
		},
	}

	data.MicrosoftDefenderForEndpoint, _ = types.ObjectValue(microsoftDefenderType.AttrTypes, map[string]attr.Value{
		"device_threat_protection_enabled":                 convert.GraphToFrameworkBool(policy.GetDeviceThreatProtectionEnabled()),
		"device_threat_protection_required_security_level": convert.GraphToFrameworkEnum(policy.GetDeviceThreatProtectionRequiredSecurityLevel()),
	})
}
