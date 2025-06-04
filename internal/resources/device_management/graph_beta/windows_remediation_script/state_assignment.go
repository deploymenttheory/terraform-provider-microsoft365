package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps the remote assignment state to Terraform state
// following a similar structure to the constructAssignment function
func MapRemoteAssignmentStateToTerraform(ctx context.Context, tfState *DeviceHealthScriptResourceModel, remoteAssignments []graphmodels.DeviceHealthScriptAssignmentable) {
	if len(remoteAssignments) == 0 {
		tflog.Debug(ctx, "No remote assignments found")
		tfState.Assignment = nil
		return
	}

	tflog.Debug(ctx, "Starting to map remote assignment state to Terraform state", map[string]interface{}{
		"assignmentCount": len(remoteAssignments),
	})

	// Initialize assignment configuration
	assignment := &WindowsRemediationScriptAssignmentResourceModel{
		AllDevices: types.BoolValue(false),
		AllUsers:   types.BoolValue(false),
	}

	// Track include and exclude groups separately to build the sets
	var includeGroups []attr.Value
	var excludeGroupIds []attr.Value

	// Process each assignment
	for idx, remoteAssignment := range remoteAssignments {
		if remoteAssignment == nil || remoteAssignment.GetTarget() == nil {
			tflog.Debug(ctx, "Skipping nil assignment or target", map[string]interface{}{
				"index": idx,
			})
			continue
		}

		target := remoteAssignment.GetTarget()
		targetType := fmt.Sprintf("%T", target)

		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"index":      idx,
			"targetType": targetType,
		})

		// Handle different target types
		switch t := target.(type) {
		case graphmodels.AllLicensedUsersAssignmentTargetable:
			tflog.Debug(ctx, "Found AllLicensedUsersAssignmentTarget")
			assignment.AllUsers = types.BoolValue(true)
			assignment.AllUsersFilterId = state.StringPointerValue(t.GetDeviceAndAppManagementAssignmentFilterId())
			assignment.AllUsersFilterType = state.EnumPtrToTypeString(t.GetDeviceAndAppManagementAssignmentFilterType())

		case graphmodels.AllDevicesAssignmentTargetable:
			tflog.Debug(ctx, "Found AllDevicesAssignmentTarget")
			assignment.AllDevices = types.BoolValue(true)
			assignment.AllDevicesFilterId = state.StringPointerValue(t.GetDeviceAndAppManagementAssignmentFilterId())
			assignment.AllDevicesFilterType = state.EnumPtrToTypeString(t.GetDeviceAndAppManagementAssignmentFilterType())

		case graphmodels.GroupAssignmentTargetable:
			tflog.Debug(ctx, "Found GroupAssignmentTarget")
			groupId := t.GetGroupId()

			if groupId == nil || *groupId == "" {
				tflog.Warn(ctx, "Skipping group assignment with nil or empty group ID")
				continue
			}

			// Create include group attributes map
			attrs := map[string]attr.Value{
				"group_id":                   state.StringPointerValue(groupId),
				"include_groups_filter_type": state.EnumPtrToTypeString(t.GetDeviceAndAppManagementAssignmentFilterType()),
				"include_groups_filter_id":   state.StringPointerValue(t.GetDeviceAndAppManagementAssignmentFilterId()),
				"run_remediation_script":     state.BoolPointerValue(remoteAssignment.GetRunRemediationScript()),
				"run_schedule":               mapRunScheduleToTerraform(ctx, remoteAssignment.GetRunSchedule()),
			}

			// Create the object and add it to include groups
			includeGroupObj, diags := types.ObjectValue(getIncludeGroupAttrTypes(), attrs)
			if diags.HasError() {
				tflog.Error(ctx, "Failed to create include group object", map[string]interface{}{
					"errors":  diags.Errors(),
					"groupId": *groupId,
				})
				continue
			}

			includeGroups = append(includeGroups, includeGroupObj)
			tflog.Debug(ctx, "Added include group", map[string]interface{}{
				"groupId": *groupId,
			})

		case graphmodels.ExclusionGroupAssignmentTargetable:
			tflog.Debug(ctx, "Found ExclusionGroupAssignmentTarget")
			groupId := t.GetGroupId()

			if groupId == nil || *groupId == "" {
				tflog.Warn(ctx, "Skipping exclusion group with nil or empty group ID")
				continue
			}

			excludeGroupIds = append(excludeGroupIds, state.StringPointerValue(groupId))
			tflog.Debug(ctx, "Added exclude group", map[string]interface{}{
				"groupId": *groupId,
			})

		default:
			tflog.Warn(ctx, "Unknown assignment target type", map[string]interface{}{
				"type": targetType,
			})
		}
	}

	// Set include groups in assignment if any
	if len(includeGroups) > 0 {
		includeGroupsSet, diags := types.SetValue(
			types.ObjectType{AttrTypes: getIncludeGroupAttrTypes()},
			includeGroups,
		)

		if diags.HasError() {
			tflog.Error(ctx, "Failed to create include groups set", map[string]interface{}{
				"errors": diags.Errors(),
			})
			assignment.IncludeGroups = types.SetNull(types.ObjectType{AttrTypes: getIncludeGroupAttrTypes()})
		} else {
			assignment.IncludeGroups = includeGroupsSet
			tflog.Debug(ctx, "Set include groups", map[string]interface{}{
				"count": len(includeGroups),
			})
		}
	} else {
		assignment.IncludeGroups = types.SetNull(types.ObjectType{AttrTypes: getIncludeGroupAttrTypes()})
	}

	// Set exclude groups in assignment if any
	if len(excludeGroupIds) > 0 {
		excludeGroupsSet, diags := types.SetValue(types.StringType, excludeGroupIds)

		if diags.HasError() {
			tflog.Error(ctx, "Failed to create exclude groups set", map[string]interface{}{
				"errors": diags.Errors(),
			})
			assignment.ExcludeGroupIds = types.SetNull(types.StringType)
		} else {
			assignment.ExcludeGroupIds = excludeGroupsSet
			tflog.Debug(ctx, "Set exclude groups", map[string]interface{}{
				"count": len(excludeGroupIds),
			})
		}
	} else {
		assignment.ExcludeGroupIds = types.SetNull(types.StringType)
	}

	// Set the assignment in the state
	tfState.Assignment = []WindowsRemediationScriptAssignmentResourceModel{*assignment}

	tflog.Debug(ctx, "Completed mapping of remote assignment state to Terraform state", map[string]interface{}{
		"includeGroupsCount": len(includeGroups),
		"excludeGroupsCount": len(excludeGroupIds),
		"allDevices":         assignment.AllDevices.ValueBool(),
		"allUsers":           assignment.AllUsers.ValueBool(),
	})
}

// mapRunScheduleToTerraform maps a GraphAPI run schedule to Terraform types
func mapRunScheduleToTerraform(ctx context.Context, schedule graphmodels.DeviceHealthScriptRunScheduleable) types.List {
	if schedule == nil {
		tflog.Debug(ctx, "Run schedule is nil")
		return types.ListNull(types.ObjectType{AttrTypes: getRunScheduleAttrTypes()})
	}

	tflog.Debug(ctx, "Mapping run schedule", map[string]interface{}{
		"scheduleType": fmt.Sprintf("%T", schedule),
	})

	scheduleAttrs := map[string]attr.Value{
		"schedule_type": types.StringNull(),
		"interval":      types.Int32Null(),
		"time":          types.StringNull(),
		"date":          types.StringNull(),
		"use_utc":       types.BoolNull(),
	}

	// Map based on schedule type
	switch s := schedule.(type) {
	case *graphmodels.DeviceHealthScriptDailySchedule:
		scheduleAttrs["schedule_type"] = types.StringValue("daily")
		scheduleAttrs["interval"] = state.Int32PtrToTypeInt32(s.GetInterval())
		scheduleAttrs["time"] = state.TimeOnlyPtrToString(s.GetTime())
		scheduleAttrs["use_utc"] = state.BoolPointerValue(s.GetUseUtc())

	case *graphmodels.DeviceHealthScriptHourlySchedule:
		scheduleAttrs["schedule_type"] = types.StringValue("hourly")
		scheduleAttrs["interval"] = state.Int32PtrToTypeInt32(s.GetInterval())

	case *graphmodels.DeviceHealthScriptRunOnceSchedule:
		scheduleAttrs["schedule_type"] = types.StringValue("once")
		scheduleAttrs["interval"] = state.Int32PtrToTypeInt32(s.GetInterval())
		scheduleAttrs["time"] = state.TimeOnlyPtrToString(s.GetTime())
		scheduleAttrs["date"] = state.DateOnlyPtrToString(s.GetDate())
		scheduleAttrs["use_utc"] = state.BoolPointerValue(s.GetUseUtc())

	default:
		tflog.Warn(ctx, "Unknown schedule type", map[string]interface{}{
			"type": fmt.Sprintf("%T", schedule),
		})
		return types.ListNull(types.ObjectType{AttrTypes: getRunScheduleAttrTypes()})
	}

	// Create schedule object
	scheduleObj, diags := types.ObjectValue(getRunScheduleAttrTypes(), scheduleAttrs)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create schedule object", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return types.ListNull(types.ObjectType{AttrTypes: getRunScheduleAttrTypes()})
	}

	// Create and return list with single schedule
	scheduleList, diags := types.ListValue(
		types.ObjectType{AttrTypes: getRunScheduleAttrTypes()},
		[]attr.Value{scheduleObj},
	)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to create schedule list", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return types.ListNull(types.ObjectType{AttrTypes: getRunScheduleAttrTypes()})
	}

	return scheduleList
}

// getIncludeGroupAttrTypes returns the attribute types for an include group
func getIncludeGroupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"group_id":                   types.StringType,
		"include_groups_filter_type": types.StringType,
		"include_groups_filter_id":   types.StringType,
		"run_remediation_script":     types.BoolType,
		"run_schedule": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: getRunScheduleAttrTypes(),
			},
		},
	}
}

// getRunScheduleAttrTypes returns the attribute types for a run schedule
func getRunScheduleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"schedule_type": types.StringType,
		"interval":      types.Int32Type,
		"time":          types.StringType,
		"date":          types.StringType,
		"use_utc":       types.BoolType,
	}
}
