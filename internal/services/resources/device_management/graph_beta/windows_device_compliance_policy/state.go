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

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, setting assignments to null", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(WindowsDeviceCompliancePolicyAssignmentType())
	} else {
		tflog.Debug(ctx, "Starting assignment mapping process", map[string]interface{}{
			"resourceId":      data.ID.ValueString(),
			"assignmentCount": len(assignments),
		})
		MapAssignmentsToTerraform(ctx, data, assignments)
		tflog.Debug(ctx, "Completed assignment mapping process", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// WindowsRemediationScriptAssignmentType returns the object type for WindowsRemediationScriptAssignmentModel
func WindowsDeviceCompliancePolicyAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"group_id":    types.StringType,
			"filter_id":   types.StringType,
			"filter_type": types.StringType,
		},
	}
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
	data.PasswordRequired = convert.GraphToFrameworkBool(policy.GetPasswordRequired())
	data.PasswordBlockSimple = convert.GraphToFrameworkBool(policy.GetPasswordBlockSimple())
	data.PasswordRequiredToUnlockFromIdle = convert.GraphToFrameworkBool(policy.GetPasswordRequiredToUnlockFromIdle())
	data.PasswordMinutesOfInactivityBeforeLock = convert.GraphToFrameworkInt32(policy.GetPasswordMinutesOfInactivityBeforeLock())
	data.PasswordExpirationDays = convert.GraphToFrameworkInt32(policy.GetPasswordExpirationDays())
	data.PasswordMinimumLength = convert.GraphToFrameworkInt32(policy.GetPasswordMinimumLength())
	data.PasswordMinimumCharacterSetCount = convert.GraphToFrameworkInt32(policy.GetPasswordMinimumCharacterSetCount())
	data.PasswordRequiredType = convert.GraphToFrameworkEnum(policy.GetPasswordRequiredType())
	data.PasswordPreviousPasswordBlockCount = convert.GraphToFrameworkInt32(policy.GetPasswordPreviousPasswordBlockCount())
	data.RequireHealthyDeviceReport = convert.GraphToFrameworkBool(policy.GetRequireHealthyDeviceReport())
	data.EarlyLaunchAntiMalwareDriverEnabled = convert.GraphToFrameworkBool(policy.GetEarlyLaunchAntiMalwareDriverEnabled())
	data.BitLockerEnabled = convert.GraphToFrameworkBool(policy.GetBitLockerEnabled())
	data.SecureBootEnabled = convert.GraphToFrameworkBool(policy.GetSecureBootEnabled())
	data.CodeIntegrityEnabled = convert.GraphToFrameworkBool(policy.GetCodeIntegrityEnabled())
	data.MemoryIntegrityEnabled = convert.GraphToFrameworkBool(policy.GetMemoryIntegrityEnabled())
	data.KernelDmaProtectionEnabled = convert.GraphToFrameworkBool(policy.GetKernelDmaProtectionEnabled())
	data.VirtualizationBasedSecurityEnabled = convert.GraphToFrameworkBool(policy.GetVirtualizationBasedSecurityEnabled())
	data.FirmwareProtectionEnabled = convert.GraphToFrameworkBool(policy.GetFirmwareProtectionEnabled())
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
			tflog.Error(ctx, "Failed to create device compliance policy script object", map[string]interface{}{
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
			"rule_name": types.StringType,
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
			"rule_name":                       convert.GraphToFrameworkString(action.GetRuleName()),
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
			"notification_message_cc_list": types.SetType{ElemType: types.StringType},
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

// MapAssignmentsToTerraform maps the remote DeviceHealthScript assignments to Terraform state
func MapAssignmentsToTerraform(ctx context.Context, data *DeviceCompliancePolicyResourceModel, assignments []graphmodels.DeviceCompliancePolicyAssignmentable) {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process")
		data.Assignments = types.SetNull(WindowsDeviceCompliancePolicyAssignmentType())
		return
	}

	tflog.Debug(ctx, "Starting assignment mapping process", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	assignmentValues := []attr.Value{}

	for i, assignment := range assignments {
		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		target := assignment.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment target is nil, skipping assignment", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			tflog.Warn(ctx, "Assignment target OData type is nil, skipping assignment", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		tflog.Debug(ctx, "Processing assignment target", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"targetType":      *odataType,
			"resourceId":      data.ID.ValueString(),
		})

		assignmentObj := map[string]attr.Value{
			"type":        types.StringNull(),
			"group_id":    types.StringNull(),
			"filter_id":   types.StringNull(),
			"filter_type": types.StringNull(),
		}

		switch *odataType {
		case "#microsoft.graph.allDevicesAssignmentTarget":
			tflog.Debug(ctx, "Mapping allDevicesAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("allDevicesAssignmentTarget")
			assignmentObj["group_id"] = types.StringNull()

		case "#microsoft.graph.allLicensedUsersAssignmentTarget":
			tflog.Debug(ctx, "Mapping allLicensedUsersAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("allLicensedUsersAssignmentTarget")
			assignmentObj["group_id"] = types.StringNull()

		case "#microsoft.graph.groupAssignmentTarget":
			tflog.Debug(ctx, "Mapping groupAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("groupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					tflog.Debug(ctx, "Setting group ID for group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"groupId":         *groupId,
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(ctx, "Group ID is nil/empty for group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = types.StringNull()
				}
			} else {
				tflog.Error(ctx, "Failed to cast target to GroupAssignmentTargetable", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["group_id"] = types.StringNull()
			}

		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			tflog.Debug(ctx, "Mapping exclusionGroupAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("exclusionGroupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					tflog.Debug(ctx, "Setting group ID for exclusion group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"groupId":         *groupId,
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(ctx, "Group ID is nil/empty for exclusion group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = types.StringNull()
				}
			} else {
				tflog.Error(ctx, "Failed to cast target to ExclusionGroupAssignmentTargetable", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["group_id"] = types.StringNull()
			}

		default:
			tflog.Warn(ctx, "Unknown target type encountered", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"targetType":      *odataType,
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["group_id"] = types.StringNull()
		}

		tflog.Debug(ctx, "Processing assignment filters", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		if filterID != nil && *filterID != "" && *filterID != "00000000-0000-0000-0000-000000000000" {
			tflog.Debug(ctx, "Assignment has meaningful filter ID", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"filterId":        *filterID,
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["filter_id"] = convert.GraphToFrameworkString(filterID)
		} else {
			tflog.Debug(ctx, "Assignment has no meaningful filter ID, using schema default", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["filter_id"] = types.StringValue("00000000-0000-0000-0000-000000000000")
		}

		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()
		if filterType != nil {
			tflog.Debug(ctx, "Processing filter type", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"filterType":      *filterType,
				"resourceId":      data.ID.ValueString(),
			})

			switch *filterType {
			case graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				tflog.Debug(ctx, "Setting filter type to include", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("include")
			case graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				tflog.Debug(ctx, "Setting filter type to exclude", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("exclude")
			case graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				tflog.Debug(ctx, "Setting filter type to none", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("none")
			default:
				tflog.Debug(ctx, "Unknown filter type, using schema default", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"filterType":      *filterType,
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("none")
			}
		} else {
			tflog.Debug(ctx, "No filter type specified, using schema default", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["filter_type"] = types.StringValue("none")
		}

		tflog.Debug(ctx, "Processing assignment schedule", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		tflog.Debug(ctx, "Creating assignment object value", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		objValue, diags := types.ObjectValue(WindowsDeviceCompliancePolicyAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
		if !diags.HasError() {
			tflog.Debug(ctx, "Successfully created assignment object", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentValues = append(assignmentValues, objValue)
		} else {
			tflog.Error(ctx, "Failed to create assignment object value", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"errors":          diags.Errors(),
				"resourceId":      data.ID.ValueString(),
			})
		}
	}

	tflog.Debug(ctx, "Creating assignments set", map[string]interface{}{
		"processedAssignments": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})

	if len(assignmentValues) > 0 {
		setVal, diags := types.SetValue(WindowsDeviceCompliancePolicyAssignmentType(), assignmentValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]interface{}{
				"errors":     diags.Errors(),
				"resourceId": data.ID.ValueString(),
			})
			data.Assignments = types.SetNull(WindowsDeviceCompliancePolicyAssignmentType())
		} else {
			tflog.Debug(ctx, "Successfully created assignments set", map[string]interface{}{
				"assignmentCount": len(assignmentValues),
				"resourceId":      data.ID.ValueString(),
			})
			data.Assignments = setVal
		}
	} else {
		tflog.Debug(ctx, "No valid assignments processed, setting assignments to null", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(WindowsDeviceCompliancePolicyAssignmentType())
	}

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]interface{}{
		"finalAssignmentCount": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})
}
