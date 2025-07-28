package graphBetaIosDeviceCompliancePolicy

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

	// This resource only handles iOS compliance policies
	if iosPolicy, ok := remoteResource.(*graphmodels.IosCompliancePolicy); ok {
		mapIosCompliancePolicyToState(ctx, data, iosPolicy)
	} else {
		tflog.Error(ctx, "Remote resource is not an iOS compliance policy")
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

	// Map restricted apps
	if iosPolicy, ok := remoteResource.(*graphmodels.IosCompliancePolicy); ok && iosPolicy.GetRestrictedApps() != nil {
		restrictedApps, err := mapRestrictedAppsToState(ctx, iosPolicy.GetRestrictedApps())
		if err != nil {
			tflog.Error(ctx, "Failed to map restricted apps", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			data.RestrictedApps = restrictedApps
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
		data.Assignments = types.SetNull(IosDeviceCompliancePolicyAssignmentType())
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

// IosDeviceCompliancePolicyAssignmentType returns the object type for IosDeviceCompliancePolicyAssignmentModel
func IosDeviceCompliancePolicyAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"group_id":    types.StringType,
			"filter_id":   types.StringType,
			"filter_type": types.StringType,
		},
	}
}

// mapIosCompliancePolicyToState is a responder function that maps iOS compliance policy properties.
func mapIosCompliancePolicyToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.IosCompliancePolicy) {
	// Map passcode properties
	data.PasscodeRequired = convert.GraphToFrameworkBool(policy.GetPasscodeRequired())
	data.PasscodeBlockSimple = convert.GraphToFrameworkBool(policy.GetPasscodeBlockSimple())
	data.PasscodeMinutesOfInactivityBeforeLock = convert.GraphToFrameworkInt32(policy.GetPasscodeMinutesOfInactivityBeforeLock())
	data.PasscodeMinutesOfInactivityBeforeScreenTimeout = convert.GraphToFrameworkInt32(policy.GetPasscodeMinutesOfInactivityBeforeScreenTimeout())
	data.PasscodeExpirationDays = convert.GraphToFrameworkInt32(policy.GetPasscodeExpirationDays())
	data.PasscodeMinimumLength = convert.GraphToFrameworkInt32(policy.GetPasscodeMinimumLength())
	data.PasscodeMinimumCharacterSetCount = convert.GraphToFrameworkInt32(policy.GetPasscodeMinimumCharacterSetCount())
	data.PasscodeRequiredType = convert.GraphToFrameworkEnum(policy.GetPasscodeRequiredType())
	data.PasscodePreviousPasscodeBlockCount = convert.GraphToFrameworkInt32(policy.GetPasscodePreviousPasscodeBlockCount())

	// Map OS version properties
	data.OsMinimumVersion = convert.GraphToFrameworkString(policy.GetOsMinimumVersion())
	data.OsMaximumVersion = convert.GraphToFrameworkString(policy.GetOsMaximumVersion())
	data.OsMinimumBuildVersion = convert.GraphToFrameworkString(policy.GetOsMinimumBuildVersion())
	data.OsMaximumBuildVersion = convert.GraphToFrameworkString(policy.GetOsMaximumBuildVersion())

	// Map iOS security properties
	data.DeviceThreatProtectionEnabled = convert.GraphToFrameworkBool(policy.GetDeviceThreatProtectionEnabled())
	data.DeviceThreatProtectionRequiredSecurityLevel = convert.GraphToFrameworkEnum(policy.GetDeviceThreatProtectionRequiredSecurityLevel())
	data.AdvancedThreatProtectionRequiredSecurityLevel = convert.GraphToFrameworkEnum(policy.GetAdvancedThreatProtectionRequiredSecurityLevel())
	data.ManagedEmailProfileRequired = convert.GraphToFrameworkBool(policy.GetManagedEmailProfileRequired())
	data.SecurityBlockJailbrokenDevices = convert.GraphToFrameworkBool(policy.GetSecurityBlockJailbrokenDevices())
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

// mapRestrictedAppsToState maps restricted apps from SDK to state
func mapRestrictedAppsToState(ctx context.Context, restrictedApps []graphmodels.AppListItemable) (types.Set, error) {
	restrictedAppType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":          types.StringType,
			"publisher":     types.StringType,
			"app_id":        types.StringType,
			"app_store_url": types.StringType,
		},
	}

	if len(restrictedApps) == 0 {
		return types.SetNull(restrictedAppType), nil
	}

	appValues := make([]attr.Value, 0, len(restrictedApps))

	for _, app := range restrictedApps {
		appAttrs := map[string]attr.Value{
			"name":          convert.GraphToFrameworkString(app.GetName()),
			"publisher":     convert.GraphToFrameworkString(app.GetPublisher()),
			"app_id":        convert.GraphToFrameworkString(app.GetAppId()),
			"app_store_url": convert.GraphToFrameworkString(app.GetAppStoreUrl()),
		}

		appValue, _ := types.ObjectValue(restrictedAppType.AttrTypes, appAttrs)
		appValues = append(appValues, appValue)
	}

	set, diags := types.SetValue(restrictedAppType, appValues)
	if diags.HasError() {
		return types.SetNull(restrictedAppType), fmt.Errorf("failed to create restricted apps set")
	}
	return set, nil
}

// MapAssignmentsToTerraform maps the remote DeviceCompliancePolicy assignments to Terraform state
func MapAssignmentsToTerraform(ctx context.Context, data *DeviceCompliancePolicyResourceModel, assignments []graphmodels.DeviceCompliancePolicyAssignmentable) {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process")
		data.Assignments = types.SetNull(IosDeviceCompliancePolicyAssignmentType())
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

		tflog.Debug(ctx, "Creating assignment object value", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		objValue, diags := types.ObjectValue(IosDeviceCompliancePolicyAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
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
		setVal, diags := types.SetValue(IosDeviceCompliancePolicyAssignmentType(), assignmentValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]interface{}{
				"errors":     diags.Errors(),
				"resourceId": data.ID.ValueString(),
			})
			data.Assignments = types.SetNull(IosDeviceCompliancePolicyAssignmentType())
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
		data.Assignments = types.SetNull(IosDeviceCompliancePolicyAssignmentType())
	}

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]interface{}{
		"finalAssignmentCount": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})
}
