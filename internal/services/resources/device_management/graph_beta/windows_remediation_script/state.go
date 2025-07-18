package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"
	"sort"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote DeviceHealthScript resource state to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *DeviceHealthScriptResourceModel, remoteResource graphmodels.DeviceHealthScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	// Map basic resource properties
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(remoteResource.GetPublisher())
	data.RunAs32Bit = convert.GraphToFrameworkBool(remoteResource.GetRunAs32Bit())
	data.EnforceSignatureCheck = convert.GraphToFrameworkBool(remoteResource.GetEnforceSignatureCheck())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())
	data.IsGlobalScript = convert.GraphToFrameworkBool(remoteResource.GetIsGlobalScript())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.HighestAvailableVersion = convert.GraphToFrameworkString(remoteResource.GetHighestAvailableVersion())
	data.RunAsAccount = convert.GraphToFrameworkEnum(remoteResource.GetRunAsAccount())
	data.DeviceHealthScriptType = convert.GraphToFrameworkEnum(remoteResource.GetDeviceHealthScriptType())
	data.DetectionScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetDetectionScriptContent())
	data.RemediationScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetRemediationScriptContent())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	// Map assignments
	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]interface{}{
		"assignmentCount": len(assignments),
	})

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, setting assignments to nil")
		data.Assignments = nil
	} else {
		MapAssignmentsToTerraform(ctx, data, assignments)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// MapAssignmentsToTerraform maps the remote DeviceHealthScript assignments to Terraform state
// Updated to use simplified model structure with unified Type and FilterId fields
func MapAssignmentsToTerraform(ctx context.Context, data *DeviceHealthScriptResourceModel, assignments []graphmodels.DeviceHealthScriptAssignmentable) {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process")
		data.Assignments = nil
		return
	}

	tflog.Debug(ctx, "Processing assignments from API response", map[string]interface{}{
		"assignmentCount": len(assignments),
	})

	// Create a single assignment model to aggregate all assignment types
	assignment := WindowsRemediationScriptAssignmentResourceModel{
		AllDevices: types.BoolValue(false),
		AllUsers:   types.BoolValue(false),
	}

	// Separate assignments by type
	var includeGroups []IncludeGroupResourceModel
	var excludeGroupIds []types.String

	// Process each assignment from the API
	for i, assignmentItem := range assignments {
		target := assignmentItem.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment has no target", map[string]interface{}{"index": i})
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			tflog.Warn(ctx, "Assignment target has no @odata.type", map[string]interface{}{"index": i})
			continue
		}

		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"index":      i,
			"targetType": *odataType,
		})

		switch *odataType {
		case "#microsoft.graph.allDevicesAssignmentTarget":
			assignment.AllDevices = types.BoolValue(true)
			mapAllDevicesTarget(ctx, &assignment, target)

		case "#microsoft.graph.allLicensedUsersAssignmentTarget":
			assignment.AllUsers = types.BoolValue(true)
			mapAllUsersTarget(ctx, &assignment, target)

		case "#microsoft.graph.groupAssignmentTarget":
			includeGroup := mapGroupAssignmentTarget(ctx, assignmentItem, target)
			if includeGroup != nil {
				includeGroups = append(includeGroups, *includeGroup)
			}

		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			excludeGroupId := mapExclusionGroupTarget(ctx, target)
			if excludeGroupId != nil {
				excludeGroupIds = append(excludeGroupIds, *excludeGroupId)
			}

		default:
			tflog.Warn(ctx, "Unknown assignment target type", map[string]interface{}{
				"targetType": *odataType,
			})
		}
	}

	// Set include_groups
	assignment.IncludeGroups = mapIncludeGroupsToSet(ctx, includeGroups)

	// Set exclude_group_ids
	assignment.ExcludeGroupIds = mapExcludeGroupIdsToSet(ctx, excludeGroupIds)

	// Create the assignments slice with our single aggregated assignment
	data.Assignments = []WindowsRemediationScriptAssignmentResourceModel{assignment}

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]interface{}{
		"assignmentCount":      len(data.Assignments),
		"includeGroupsCount":   len(includeGroups),
		"excludeGroupIdsCount": len(excludeGroupIds),
	})
}

// mapAllDevicesTarget maps AllDevicesAssignmentTarget properties
// Updated to use unified Type and FilterId fields
func mapAllDevicesTarget(ctx context.Context, assignment *WindowsRemediationScriptAssignmentResourceModel, target graphmodels.DeviceAndAppManagementAssignmentTargetable) {
	if allDevicesTarget, ok := target.(graphmodels.AllDevicesAssignmentTargetable); ok {
		if filterId := allDevicesTarget.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil {
			assignment.FilterId = types.StringValue(*filterId)
		}
		if filterType := allDevicesTarget.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
			assignment.Type = types.StringValue(filterType.String())
		}
	}
}

// mapAllUsersTarget maps AllLicensedUsersAssignmentTarget properties
// Updated to use unified Type and FilterId fields
func mapAllUsersTarget(ctx context.Context, assignment *WindowsRemediationScriptAssignmentResourceModel, target graphmodels.DeviceAndAppManagementAssignmentTargetable) {
	if allUsersTarget, ok := target.(graphmodels.AllLicensedUsersAssignmentTargetable); ok {
		if filterId := allUsersTarget.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil {
			assignment.FilterId = types.StringValue(*filterId)
		}
		if filterType := allUsersTarget.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
			assignment.Type = types.StringValue(filterType.String())
		}
	}
}

// mapGroupAssignmentTarget maps GroupAssignmentTarget to IncludeGroupResourceModel
// Updated to use unified Type and FilterId fields
func mapGroupAssignmentTarget(ctx context.Context, assignmentItem graphmodels.DeviceHealthScriptAssignmentable, target graphmodels.DeviceAndAppManagementAssignmentTargetable) *IncludeGroupResourceModel {
	groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable)
	if !ok {
		return nil
	}

	groupId := groupTarget.GetGroupId()
	if groupId == nil {
		return nil
	}

	includeGroup := &IncludeGroupResourceModel{
		GroupId: types.StringValue(*groupId),
	}

	// Map filter properties using unified field names
	if filterId := groupTarget.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil {
		includeGroup.FilterId = types.StringValue(*filterId)
	} else {
		includeGroup.FilterId = types.StringValue("")
	}

	if filterType := groupTarget.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
		includeGroup.Type = types.StringValue(filterType.String())
	} else {
		includeGroup.Type = types.StringValue("")
	}

	// Handle API inconsistency with runRemediationScript
	runRemediationFromAPI := assignmentItem.GetRunRemediationScript()

	tflog.Debug(ctx, "API returned RunRemediationScript", map[string]interface{}{
		"groupId":       *groupId,
		"apiValue":      runRemediationFromAPI,
		"apiValueIsNil": runRemediationFromAPI == nil,
	})

	// Map run schedule with better error handling
	if runSchedule := assignmentItem.GetRunSchedule(); runSchedule != nil {
		scheduleModel := mapRunSchedule(ctx, runSchedule)
		if scheduleModel != nil {
			includeGroup.RunSchedule = scheduleModel
		} else {
			tflog.Warn(ctx, "Failed to map run schedule from API", map[string]interface{}{
				"groupId": *groupId,
			})
		}
	}

	tflog.Debug(ctx, "Successfully mapped group assignment target", map[string]interface{}{
		"groupId":     *groupId,
		"type":        includeGroup.Type.ValueString(),
		"filterId":    includeGroup.FilterId.ValueString(),
		"hasSchedule": includeGroup.RunSchedule != nil,
	})

	return includeGroup
}

// mapExclusionGroupTarget maps ExclusionGroupAssignmentTarget to string
// No changes needed as this only deals with group IDs
func mapExclusionGroupTarget(ctx context.Context, target graphmodels.DeviceAndAppManagementAssignmentTargetable) *types.String {
	exclusionTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable)
	if !ok {
		return nil
	}

	groupId := exclusionTarget.GetGroupId()
	if groupId == nil {
		return nil
	}

	result := types.StringValue(*groupId)
	tflog.Debug(ctx, "Mapped exclusion group target", map[string]interface{}{
		"groupId": *groupId,
	})

	return &result
}

// mapIncludeGroupsToSet converts IncludeGroupResourceModel slice to types.Set
// Updated to use unified Type and FilterId fields
func mapIncludeGroupsToSet(ctx context.Context, includeGroups []IncludeGroupResourceModel) types.Set {
	if len(includeGroups) == 0 {
		return types.SetNull(getIncludeGroupObjectType())
	}

	// Sort for consistent ordering
	sort.Slice(includeGroups, func(i, j int) bool {
		return includeGroups[i].GroupId.ValueString() < includeGroups[j].GroupId.ValueString()
	})

	includeGroupValues := make([]attr.Value, 0, len(includeGroups))

	for _, group := range includeGroups {
		tflog.Debug(ctx, "Processing include group for set", map[string]interface{}{
			"groupId":        group.GroupId.ValueString(),
			"type":           group.Type.ValueString(),
			"filterId":       group.FilterId.ValueString(),
			"hasRunSchedule": group.RunSchedule != nil,
		})

		// Handle type: preserve the exact value from API or empty if not set
		var filterType types.String
		if group.Type.IsNull() || group.Type.ValueString() == "" {
			filterType = types.StringValue("")
		} else {
			filterType = group.Type
		}

		// Handle filter ID: preserve the exact value from API or empty if not set
		var filterId types.String
		if group.FilterId.IsNull() || group.FilterId.ValueString() == "" {
			filterId = types.StringValue("")
		} else {
			filterId = group.FilterId
		}

		// Create run_schedule object with exact value patterns
		var runScheduleObj attr.Value
		if group.RunSchedule != nil {
			scheduleAttrs := map[string]attr.Value{
				"schedule_type": group.RunSchedule.ScheduleType,
				"interval":      group.RunSchedule.Interval,
				"time":          group.RunSchedule.Time,
				"date":          group.RunSchedule.Date,
				"use_utc":       group.RunSchedule.UseUtc,
			}

			var diags diag.Diagnostics
			runScheduleObj, diags = types.ObjectValue(getRunScheduleObjectType().AttrTypes, scheduleAttrs)
			if diags.HasError() {
				tflog.Error(ctx, "Failed to create run_schedule object", map[string]interface{}{
					"errors": diags.Errors(),
					"group":  group.GroupId.ValueString(),
				})
				continue // Skip this group if we can't create the schedule object
			}
		} else {
			runScheduleObj = types.ObjectNull(getRunScheduleObjectType().AttrTypes)
		}

		groupAttrs := map[string]attr.Value{
			"group_id":     group.GroupId,
			"type":         filterType,
			"filter_id":    filterId,
			"run_schedule": runScheduleObj,
		}

		groupObj, diags := types.ObjectValue(getIncludeGroupObjectType().AttrTypes, groupAttrs)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create include group object", map[string]interface{}{
				"errors": diags.Errors(),
				"group":  group.GroupId.ValueString(),
			})
			continue
		}

		includeGroupValues = append(includeGroupValues, groupObj)

		tflog.Debug(ctx, "Successfully created include group object", map[string]interface{}{
			"groupId":        group.GroupId.ValueString(),
			"type":           filterType.ValueString(),
			"filterId":       filterId.ValueString(),
			"hasRunSchedule": group.RunSchedule != nil,
		})
	}

	if len(includeGroupValues) == 0 {
		return types.SetNull(getIncludeGroupObjectType())
	}

	includeGroupsSet, diags := types.SetValue(getIncludeGroupObjectType(), includeGroupValues)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create include_groups set", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return types.SetNull(getIncludeGroupObjectType())
	}

	return includeGroupsSet
}

// mapExcludeGroupIdsToSet converts string slice to types.Set
// No changes needed as this only deals with string IDs
func mapExcludeGroupIdsToSet(ctx context.Context, excludeGroupIds []types.String) types.Set {
	if len(excludeGroupIds) == 0 {
		return types.SetNull(types.StringType)
	}

	// Sort for consistent ordering
	sort.Slice(excludeGroupIds, func(i, j int) bool {
		return excludeGroupIds[i].ValueString() < excludeGroupIds[j].ValueString()
	})

	excludeGroupsSet, diags := types.SetValueFrom(ctx, types.StringType, excludeGroupIds)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create exclude_group_ids set", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return types.SetNull(types.StringType)
	}

	return excludeGroupsSet
}

func mapRunSchedule(ctx context.Context, schedule graphmodels.DeviceHealthScriptRunScheduleable) *RunScheduleResourceModel {
	if schedule == nil {
		return nil
	}

	// Initialize with schema defaults
	result := &RunScheduleResourceModel{
		Interval: types.Int32Value(1),    // Schema default
		UseUtc:   types.BoolValue(false), // Schema default
	}

	odataType := schedule.GetOdataType()
	if odataType == nil {
		tflog.Warn(ctx, "Schedule missing @odata.type, defaulting to hourly")
		result.ScheduleType = types.StringValue("hourly")
		result.Time = types.StringNull()
		result.Date = types.StringNull()
		return result
	}

	switch *odataType {
	case "#microsoft.graph.deviceHealthScriptDailySchedule":
		if dailySchedule, ok := schedule.(graphmodels.DeviceHealthScriptDailyScheduleable); ok {
			result.ScheduleType = types.StringValue("daily")

			if interval := dailySchedule.GetInterval(); interval != nil {
				result.Interval = types.Int32Value(*interval)
			}

			// CRITICAL FIX: Handle API returning empty/null time values
			if time := dailySchedule.GetTime(); time != nil {
				timeStr := convert.GraphToFrameworkTimeOnly(time)
				if !timeStr.IsNull() && timeStr.ValueString() != "" {
					result.Time = timeStr
				} else {
					// API returned empty time - this is an API inconsistency
					tflog.Warn(ctx, "API returned empty time for daily schedule")
					result.Time = types.StringValue("") // Preserve as empty string
				}
			} else {
				// API didn't return time field at all
				tflog.Warn(ctx, "API didn't return time field for daily schedule")
				result.Time = types.StringValue("") // Set as empty string
			}

			// Date is always null for daily schedules
			result.Date = types.StringNull()

			if useUtc := dailySchedule.GetUseUtc(); useUtc != nil {
				result.UseUtc = types.BoolValue(*useUtc)
			}

			tflog.Debug(ctx, "Mapped daily schedule", map[string]interface{}{
				"interval":   result.Interval.ValueInt32(),
				"time":       result.Time.ValueString(),
				"timeIsNull": result.Time.IsNull(),
				"useUtc":     result.UseUtc.ValueBool(),
			})
		}

	case "#microsoft.graph.deviceHealthScriptHourlySchedule":
		if hourlySchedule, ok := schedule.(graphmodels.DeviceHealthScriptHourlyScheduleable); ok {
			result.ScheduleType = types.StringValue("hourly")

			if interval := hourlySchedule.GetInterval(); interval != nil {
				result.Interval = types.Int32Value(*interval)
			}

			// For hourly schedule: time and date should be null
			result.Time = types.StringNull()
			result.Date = types.StringNull()
			result.UseUtc = types.BoolValue(false) // Always false for hourly

			tflog.Debug(ctx, "Mapped hourly schedule", map[string]interface{}{
				"interval": result.Interval.ValueInt32(),
			})
		}

	case "#microsoft.graph.deviceHealthScriptRunOnceSchedule":
		if onceSchedule, ok := schedule.(graphmodels.DeviceHealthScriptRunOnceScheduleable); ok {
			result.ScheduleType = types.StringValue("once")

			if interval := onceSchedule.GetInterval(); interval != nil {
				result.Interval = types.Int32Value(*interval)
			}

			if date := onceSchedule.GetDate(); date != nil {
				dateStr := convert.GraphToFrameworkDateOnly(date)
				if !dateStr.IsNull() && dateStr.ValueString() != "" {
					result.Date = dateStr
				} else {
					tflog.Warn(ctx, "API returned empty date for once schedule")
					result.Date = types.StringValue("")
				}
			} else {
				tflog.Warn(ctx, "API didn't return date field for once schedule")
				result.Date = types.StringValue("")
			}

			if time := onceSchedule.GetTime(); time != nil {
				timeStr := convert.GraphToFrameworkTimeOnly(time)
				if !timeStr.IsNull() && timeStr.ValueString() != "" {
					result.Time = timeStr
				} else {
					tflog.Warn(ctx, "API returned empty time for once schedule")
					result.Time = types.StringValue("")
				}
			} else {
				tflog.Warn(ctx, "API didn't return time field for once schedule")
				result.Time = types.StringValue("")
			}

			if useUtc := onceSchedule.GetUseUtc(); useUtc != nil {
				result.UseUtc = types.BoolValue(*useUtc)
			}

			tflog.Debug(ctx, "Mapped once schedule", map[string]interface{}{
				"interval":   result.Interval.ValueInt32(),
				"date":       result.Date.ValueString(),
				"dateIsNull": result.Date.IsNull(),
				"time":       result.Time.ValueString(),
				"timeIsNull": result.Time.IsNull(),
				"useUtc":     result.UseUtc.ValueBool(),
			})
		}

	default:
		tflog.Warn(ctx, "Unknown schedule type", map[string]interface{}{
			"scheduleType": *odataType,
		})
		// Default to hourly for unknown types
		result.ScheduleType = types.StringValue("hourly")
		result.Time = types.StringNull()
		result.Date = types.StringNull()
		result.UseUtc = types.BoolValue(false)
	}

	return result
}

// Helper functions to get object types - Updated for simplified model
func getIncludeGroupObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"group_id":               types.StringType,
			"type":                   types.StringType,
			"filter_id":              types.StringType,
			"run_remediation_script": types.BoolType,
			"run_schedule":           getRunScheduleObjectType(),
		},
	}
}

// getRunScheduleObjectType remains unchanged
func getRunScheduleObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"schedule_type": types.StringType,
			"interval":      types.Int32Type,
			"time":          types.StringType,
			"date":          types.StringType,
			"use_utc":       types.BoolType,
		},
	}
}
