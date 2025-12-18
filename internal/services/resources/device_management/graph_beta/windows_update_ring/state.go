package graphBetaWindowsUpdateRing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the API response to the Terraform model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsUpdateRingResourceModel, remoteResource graphmodels.WindowsUpdateForBusinessConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.MicrosoftUpdateServiceAllowed = convert.GraphToFrameworkBool(remoteResource.GetMicrosoftUpdateServiceAllowed())
	data.DriversExcluded = convert.GraphToFrameworkBool(remoteResource.GetDriversExcluded())
	data.QualityUpdatesDeferralPeriodInDays = convert.GraphToFrameworkInt32(remoteResource.GetQualityUpdatesDeferralPeriodInDays())
	data.FeatureUpdatesDeferralPeriodInDays = convert.GraphToFrameworkInt32(remoteResource.GetFeatureUpdatesDeferralPeriodInDays())
	data.AllowWindows11Upgrade = convert.GraphToFrameworkBool(remoteResource.GetAllowWindows11Upgrade())
	data.QualityUpdatesPaused = convert.GraphToFrameworkBool(remoteResource.GetQualityUpdatesPaused())
	data.FeatureUpdatesPaused = convert.GraphToFrameworkBool(remoteResource.GetFeatureUpdatesPaused())

	data.FeatureUpdatesPauseStartDate = convert.GraphToFrameworkDateOnly(remoteResource.GetFeatureUpdatesPauseStartDate())
	data.FeatureUpdatesPauseExpiryDateTime = convert.GraphToFrameworkTime(remoteResource.GetFeatureUpdatesPauseExpiryDateTime())
	data.FeatureUpdatesRollbackStartDateTime = convert.GraphToFrameworkTime(remoteResource.GetFeatureUpdatesRollbackStartDateTime())
	data.QualityUpdatesPauseStartDate = convert.GraphToFrameworkDateOnly(remoteResource.GetQualityUpdatesPauseStartDate())
	data.QualityUpdatesPauseExpiryDateTime = convert.GraphToFrameworkTime(remoteResource.GetQualityUpdatesPauseExpiryDateTime())
	data.QualityUpdatesRollbackStartDateTime = convert.GraphToFrameworkTime(remoteResource.GetQualityUpdatesRollbackStartDateTime())

	data.SkipChecksBeforeRestart = convert.GraphToFrameworkBool(remoteResource.GetSkipChecksBeforeRestart())
	data.BusinessReadyUpdatesOnly = convert.GraphToFrameworkEnum(remoteResource.GetBusinessReadyUpdatesOnly())
	data.AutomaticUpdateMode = convert.GraphToFrameworkEnum(remoteResource.GetAutomaticUpdateMode())
	data.UpdateNotificationLevel = convert.GraphToFrameworkEnum(remoteResource.GetUpdateNotificationLevel())
	data.UpdateWeeks = convert.GraphToFrameworkEnum(remoteResource.GetUpdateWeeks())

	// Handle installation schedule (can be either active hours or scheduled install)
	if installationSchedule := remoteResource.GetInstallationSchedule(); installationSchedule != nil {
		if activeHoursInstall, ok := installationSchedule.(graphmodels.WindowsUpdateActiveHoursInstallable); ok {
			// Handle WindowsUpdateActiveHoursInstall
			if activeHoursInstall.GetActiveHoursStart() != nil {
				data.ActiveHoursStart = types.StringValue(activeHoursInstall.GetActiveHoursStart().String())
			} else {
				data.ActiveHoursStart = types.StringNull()
			}

			if activeHoursInstall.GetActiveHoursEnd() != nil {
				data.ActiveHoursEnd = types.StringValue(activeHoursInstall.GetActiveHoursEnd().String())
			} else {
				data.ActiveHoursEnd = types.StringNull()
			}

			// Clear scheduled install fields
			data.ScheduledInstallDay = types.StringNull()
			data.ScheduledInstallTime = types.StringNull()
		} else if scheduledInstall, ok := installationSchedule.(graphmodels.WindowsUpdateScheduledInstallable); ok {
			// Handle WindowsUpdateScheduledInstall
			data.ScheduledInstallDay = convert.GraphToFrameworkEnum(scheduledInstall.GetScheduledInstallDay())

			if scheduledInstall.GetScheduledInstallTime() != nil {
				data.ScheduledInstallTime = types.StringValue(scheduledInstall.GetScheduledInstallTime().String())
			} else {
				data.ScheduledInstallTime = types.StringNull()
			}

			// Clear active hours fields
			data.ActiveHoursStart = types.StringNull()
			data.ActiveHoursEnd = types.StringNull()
		} else {
			tflog.Warn(ctx, "Installation schedule is not of a recognized type")
			data.ActiveHoursStart = types.StringNull()
			data.ActiveHoursEnd = types.StringNull()
			data.ScheduledInstallDay = types.StringNull()
			data.ScheduledInstallTime = types.StringNull()
		}
	} else {
		// No installation schedule configured
		data.ActiveHoursStart = types.StringNull()
		data.ActiveHoursEnd = types.StringNull()
		data.ScheduledInstallDay = types.StringNull()
		data.ScheduledInstallTime = types.StringNull()
	}

	data.UserPauseAccess = convert.GraphToFrameworkEnum(remoteResource.GetUserPauseAccess())
	data.UserWindowsUpdateScanAccess = convert.GraphToFrameworkEnum(remoteResource.GetUserWindowsUpdateScanAccess())
	data.FeatureUpdatesRollbackWindowInDays = convert.GraphToFrameworkInt32(remoteResource.GetFeatureUpdatesRollbackWindowInDays())

	// Only set deadline settings if any deadline field has a non-null value from the API
	if remoteResource.GetDeadlineForFeatureUpdatesInDays() != nil ||
		remoteResource.GetDeadlineForQualityUpdatesInDays() != nil ||
		remoteResource.GetDeadlineGracePeriodInDays() != nil ||
		remoteResource.GetPostponeRebootUntilAfterDeadline() != nil {
		deadlineSettings := DeadlineSettingsModel{
			DeadlineForFeatureUpdatesInDays:  convert.GraphToFrameworkInt32(remoteResource.GetDeadlineForFeatureUpdatesInDays()),
			DeadlineForQualityUpdatesInDays:  convert.GraphToFrameworkInt32(remoteResource.GetDeadlineForQualityUpdatesInDays()),
			DeadlineGracePeriodInDays:        convert.GraphToFrameworkInt32(remoteResource.GetDeadlineGracePeriodInDays()),
			PostponeRebootUntilAfterDeadline: convert.GraphToFrameworkBool(remoteResource.GetPostponeRebootUntilAfterDeadline()),
		}
		deadlineSettingsObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"deadline_for_feature_updates_in_days": types.Int32Type,
			"deadline_for_quality_updates_in_days": types.Int32Type,
			"deadline_grace_period_in_days":        types.Int32Type,
			"postpone_reboot_until_after_deadline": types.BoolType,
		}, deadlineSettings)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create deadline settings object", map[string]any{
				"errors": diags.Errors(),
			})
			data.DeadlineSettings = types.ObjectNull(map[string]attr.Type{
				"deadline_for_feature_updates_in_days": types.Int32Type,
				"deadline_for_quality_updates_in_days": types.Int32Type,
				"deadline_grace_period_in_days":        types.Int32Type,
				"postpone_reboot_until_after_deadline": types.BoolType,
			})
		} else {
			data.DeadlineSettings = deadlineSettingsObj
		}
	} else {
		data.DeadlineSettings = types.ObjectNull(map[string]attr.Type{
			"deadline_for_feature_updates_in_days": types.Int32Type,
			"deadline_for_quality_updates_in_days": types.Int32Type,
			"deadline_grace_period_in_days":        types.Int32Type,
			"postpone_reboot_until_after_deadline": types.BoolType,
		})
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
		data.Assignments = types.SetNull(WindowsUpdateRingAssignmentType())
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

// WindowsUpdateRingAssignmentType returns the object type for WindowsRemediationScriptAssignmentModel
func WindowsUpdateRingAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"group_id":    types.StringType,
			"filter_id":   types.StringType,
			"filter_type": types.StringType,
		},
	}
}

// MapAssignmentsToTerraform maps the remote DeviceHealthScript assignments to Terraform state
func MapAssignmentsToTerraform(ctx context.Context, data *WindowsUpdateRingResourceModel, assignments []graphmodels.DeviceConfigurationAssignmentable) {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process")
		data.Assignments = types.SetNull(WindowsUpdateRingAssignmentType())
		return
	}

	tflog.Debug(ctx, "Starting assignment mapping process", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	assignmentValues := []attr.Value{}

	for i, assignment := range assignments {
		tflog.Debug(ctx, "Processing assignment", map[string]any{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		target := assignment.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment target is nil, skipping assignment", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			tflog.Warn(ctx, "Assignment target OData type is nil, skipping assignment", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		tflog.Debug(ctx, "Processing assignment target", map[string]any{
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
			tflog.Debug(ctx, "Mapping allDevicesAssignmentTarget", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("allDevicesAssignmentTarget")
			assignmentObj["group_id"] = types.StringNull()

		case "#microsoft.graph.allLicensedUsersAssignmentTarget":
			tflog.Debug(ctx, "Mapping allLicensedUsersAssignmentTarget", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("allLicensedUsersAssignmentTarget")
			assignmentObj["group_id"] = types.StringNull()

		case "#microsoft.graph.groupAssignmentTarget":
			tflog.Debug(ctx, "Mapping groupAssignmentTarget", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("groupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					tflog.Debug(ctx, "Setting group ID for group assignment target", map[string]any{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"groupId":         *groupId,
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(ctx, "Group ID is nil/empty for group assignment target", map[string]any{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = types.StringNull()
				}
			} else {
				tflog.Error(ctx, "Failed to cast target to GroupAssignmentTargetable", map[string]any{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["group_id"] = types.StringNull()
			}

		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			tflog.Debug(ctx, "Mapping exclusionGroupAssignmentTarget", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("exclusionGroupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					tflog.Debug(ctx, "Setting group ID for exclusion group assignment target", map[string]any{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"groupId":         *groupId,
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(ctx, "Group ID is nil/empty for exclusion group assignment target", map[string]any{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = types.StringNull()
				}
			} else {
				tflog.Error(ctx, "Failed to cast target to ExclusionGroupAssignmentTargetable", map[string]any{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["group_id"] = types.StringNull()
			}

		default:
			tflog.Warn(ctx, "Unknown target type encountered", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"targetType":      *odataType,
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["group_id"] = types.StringNull()
		}

		tflog.Debug(ctx, "Processing assignment filters", map[string]any{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		if filterID != nil && *filterID != "" && *filterID != "00000000-0000-0000-0000-000000000000" {
			tflog.Debug(ctx, "Assignment has meaningful filter ID", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"filterId":        *filterID,
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["filter_id"] = convert.GraphToFrameworkString(filterID)
		} else {
			tflog.Debug(ctx, "Assignment has no meaningful filter ID, using schema default", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["filter_id"] = types.StringValue("00000000-0000-0000-0000-000000000000")
		}

		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()
		if filterType != nil {
			tflog.Debug(ctx, "Processing filter type", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"filterType":      *filterType,
				"resourceId":      data.ID.ValueString(),
			})

			switch *filterType {
			case graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				tflog.Debug(ctx, "Setting filter type to include", map[string]any{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("include")
			case graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				tflog.Debug(ctx, "Setting filter type to exclude", map[string]any{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("exclude")
			case graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				tflog.Debug(ctx, "Setting filter type to none", map[string]any{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("none")
			default:
				tflog.Debug(ctx, "Unknown filter type, using schema default", map[string]any{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"filterType":      *filterType,
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("none")
			}
		} else {
			tflog.Debug(ctx, "No filter type specified, using schema default", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["filter_type"] = types.StringValue("none")
		}

		tflog.Debug(ctx, "Processing assignment schedule", map[string]any{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		tflog.Debug(ctx, "Creating assignment object value", map[string]any{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		objValue, diags := types.ObjectValue(WindowsUpdateRingAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
		if !diags.HasError() {
			tflog.Debug(ctx, "Successfully created assignment object", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentValues = append(assignmentValues, objValue)
		} else {
			tflog.Error(ctx, "Failed to create assignment object value", map[string]any{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"errors":          diags.Errors(),
				"resourceId":      data.ID.ValueString(),
			})
		}
	}

	tflog.Debug(ctx, "Creating assignments set", map[string]any{
		"processedAssignments": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})

	if len(assignmentValues) > 0 {
		setVal, diags := types.SetValue(WindowsUpdateRingAssignmentType(), assignmentValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]any{
				"errors":     diags.Errors(),
				"resourceId": data.ID.ValueString(),
			})
			data.Assignments = types.SetNull(WindowsUpdateRingAssignmentType())
		} else {
			tflog.Debug(ctx, "Successfully created assignments set", map[string]any{
				"assignmentCount": len(assignmentValues),
				"resourceId":      data.ID.ValueString(),
			})
			data.Assignments = setVal
		}
	} else {
		tflog.Debug(ctx, "No valid assignments processed, setting assignments to null", map[string]any{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(WindowsUpdateRingAssignmentType())
	}

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]any{
		"finalAssignmentCount": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})
}
