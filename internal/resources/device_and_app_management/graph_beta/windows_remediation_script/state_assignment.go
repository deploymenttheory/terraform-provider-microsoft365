package graphBetaWindowsRemediationScript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps remote assignment state to Terraform state
func MapRemoteAssignmentStateToTerraform(ctx context.Context, state *DeviceHealthScriptResourceModel, remoteAssignments []graphmodels.DeviceHealthScriptAssignmentable) {
	if remoteAssignments == nil || len(remoteAssignments) == 0 {
		tflog.Debug(ctx, "No remote assignments found")
		state.Assignment = nil
		return
	}

	tflog.Debug(ctx, "Starting to map remote assignment state to Terraform state")

	assignment := &WindowsRemediationScriptAssignmentResourceModel{
		AllDevices: types.BoolValue(false),
		AllUsers:   types.BoolValue(false),
	}

	var includeGroups []IncludeGroupResourceModel
	var excludeGroups []string

	// Process each remote assignment
	for _, remoteAssignment := range remoteAssignments {
		if remoteAssignment == nil || remoteAssignment.GetTarget() == nil {
			continue
		}

		target := remoteAssignment.GetTarget()

		// Process assignment based on target type
		switch t := target.(type) {
		case *graphmodels.AllDevicesAssignmentTarget:
			processAllDevicesTarget(ctx, assignment, t)
		case *graphmodels.AllLicensedUsersAssignmentTarget:
			processAllUsersTarget(ctx, assignment, t)
		case *graphmodels.GroupAssignmentTarget:
			group := processGroupIncludeTarget(ctx, remoteAssignment, t)
			if group != nil {
				includeGroups = append(includeGroups, *group)
			}
		case *graphmodels.ExclusionGroupAssignmentTarget:
			groupId := processGroupExcludeTarget(ctx, t)
			if groupId != "" {
				excludeGroups = append(excludeGroups, groupId)
			}
		}
	}

	// Set include groups
	setIncludeGroups(ctx, assignment, includeGroups)

	// Set exclude groups
	setExcludeGroups(ctx, assignment, excludeGroups)

	state.Assignment = []WindowsRemediationScriptAssignmentResourceModel{*assignment}

	tflog.Debug(ctx, "Finished mapping remote assignment state to Terraform state")
}

// processAllDevicesTarget handles all devices assignment target
func processAllDevicesTarget(ctx context.Context, assignment *WindowsRemediationScriptAssignmentResourceModel, target *graphmodels.AllDevicesAssignmentTarget) {
	tflog.Debug(ctx, "Found all devices assignment")
	assignment.AllDevices = types.BoolValue(true)

	if filterId := target.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil {
		assignment.AllDevicesFilterId = types.StringValue(*filterId)
	}
	if filterType := target.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
		assignment.AllDevicesFilterType = types.StringValue(filterType.String())
	}
}

// processAllUsersTarget handles all users assignment target
func processAllUsersTarget(ctx context.Context, assignment *WindowsRemediationScriptAssignmentResourceModel, target *graphmodels.AllLicensedUsersAssignmentTarget) {
	tflog.Debug(ctx, "Found all users assignment")
	assignment.AllUsers = types.BoolValue(true)

	if filterId := target.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil {
		assignment.AllUsersFilterId = types.StringValue(*filterId)
	}
	if filterType := target.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
		assignment.AllUsersFilterType = types.StringValue(filterType.String())
	}
}

// processGroupIncludeTarget handles group include assignment target
func processGroupIncludeTarget(ctx context.Context, remoteAssignment graphmodels.DeviceHealthScriptAssignmentable, target *graphmodels.GroupAssignmentTarget) *IncludeGroupResourceModel {
	groupId := target.GetGroupId()
	if groupId == nil {
		return nil
	}

	tflog.Debug(ctx, "Found include group assignment", map[string]interface{}{
		"groupId": *groupId,
	})

	includeGroup := &IncludeGroupResourceModel{
		GroupId: types.StringValue(*groupId),
	}

	if filterId := target.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil {
		includeGroup.IncludeGroupsFilterId = types.StringValue(*filterId)
	}
	if filterType := target.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
		includeGroup.IncludeGroupsFilterType = types.StringValue(filterType.String())
	}

	// Map run remediation script
	if runRemediation := remoteAssignment.GetRunRemediationScript(); runRemediation != nil {
		includeGroup.RunRemediationScript = types.BoolValue(*runRemediation)
	}

	// Map run schedule
	// if schedule := remoteAssignment.GetRunSchedule(); schedule != nil {
	// 	runSchedule, err := mapRunScheduleToTerraform(ctx, schedule)
	// 	if err != nil {
	// 		tflog.Error(ctx, "Failed to map run schedule", map[string]interface{}{
	// 			"error": err.Error(),
	// 		})
	// 	} else {
	// 		includeGroup.RunSchedule = runSchedule
	// 	}
	// }

	return includeGroup
}

// processGroupExcludeTarget handles group exclude assignment target
func processGroupExcludeTarget(ctx context.Context, target *graphmodels.ExclusionGroupAssignmentTarget) string {
	groupId := target.GetGroupId()
	if groupId == nil {
		return ""
	}

	tflog.Debug(ctx, "Found exclude group assignment", map[string]interface{}{
		"groupId": *groupId,
	})
	return *groupId
}

// setIncludeGroups sets the include groups on the assignment
func setIncludeGroups(ctx context.Context, assignment *WindowsRemediationScriptAssignmentResourceModel, includeGroups []IncludeGroupResourceModel) {
	if len(includeGroups) == 0 {
		assignment.IncludeGroups = types.SetNull(getIncludeGroupObjectType())
		return
	}

	includeGroupElements := make([]attr.Value, 0, len(includeGroups))
	includeGroupObjType := getIncludeGroupObjectType()

	for _, group := range includeGroups {
		groupObj := createIncludeGroupObject(ctx, group, includeGroupObjType)
		if groupObj != nil {
			includeGroupElements = append(includeGroupElements, groupObj)
		}
	}

	if len(includeGroupElements) > 0 {
		includeGroupsSet, diags := types.SetValue(includeGroupObjType, includeGroupElements)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create include groups set", map[string]interface{}{
				"errors": diags.Errors(),
			})
			assignment.IncludeGroups = types.SetNull(includeGroupObjType)
		} else {
			assignment.IncludeGroups = includeGroupsSet
		}
	} else {
		assignment.IncludeGroups = types.SetNull(includeGroupObjType)
	}
}

// setExcludeGroups sets the exclude groups on the assignment
func setExcludeGroups(ctx context.Context, assignment *WindowsRemediationScriptAssignmentResourceModel, excludeGroups []string) {
	if len(excludeGroups) > 0 {
		excludeSet, diags := types.SetValueFrom(ctx, types.StringType, excludeGroups)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create exclude groups set", map[string]interface{}{
				"errors": diags.Errors(),
			})
		} else {
			assignment.ExcludeGroupIds = excludeSet
		}
	} else {
		assignment.ExcludeGroupIds = types.SetNull(types.StringType)
	}
}

// getIncludeGroupObjectType returns the object type for include groups
func getIncludeGroupObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"group_id":                   types.StringType,
			"include_groups_filter_type": types.StringType,
			"include_groups_filter_id":   types.StringType,
			"run_remediation_script":     types.BoolType,
			"run_schedule": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: getRunScheduleObjectTypes(),
				},
			},
		},
	}
}

// getRunScheduleObjectTypes returns the attribute types for run schedule
func getRunScheduleObjectTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"schedule_type": types.StringType,
		"interval":      types.Int32Type,
		"time":          types.StringType,
		"date":          types.StringType,
		"use_utc":       types.BoolType,
	}
}

// createIncludeGroupObject creates an include group object
func createIncludeGroupObject(ctx context.Context, group IncludeGroupResourceModel, objType types.ObjectType) attr.Value {
	attrs := map[string]attr.Value{
		"group_id":                   group.GroupId,
		"include_groups_filter_type": group.IncludeGroupsFilterType,
		"include_groups_filter_id":   group.IncludeGroupsFilterId,
		"run_remediation_script":     group.RunRemediationScript,
		"run_schedule":               group.RunSchedule,
	}

	groupObj, diags := types.ObjectValue(objType.AttrTypes, attrs)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create include group object", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return nil
	}
	return groupObj
}
