// Updated constructor functions to match the simplified model structure
package graphBetaWindowsRemediationScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs and returns a DeviceHealthScriptsItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *DeviceHealthScriptResourceModel) (devicemanagement.DeviceHealthScriptsItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Starting Device Health Script assignment construction")

	// Create the request body - use the specific assign request body type
	requestBody := devicemanagement.NewDeviceHealthScriptsItemAssignPostRequestBody()
	scriptAssignments := make([]graphmodels.DeviceHealthScriptAssignmentable, 0)

	// If assignments is nil, return empty array to remove all assignments
	if data.Assignments == nil {
		tflog.Debug(ctx, "Assignments is nil, creating empty assignments array")
		requestBody.SetDeviceHealthScriptAssignments(scriptAssignments)
		return requestBody, nil
	}

	// Process each assignment block
	for idx, assignment := range data.Assignments {
		tflog.Debug(ctx, "Processing assignment block", map[string]interface{}{
			"index": idx,
		})

		// 1. Handle All Devices Assignment
		if !assignment.AllDevices.IsNull() && assignment.AllDevices.ValueBool() {
			allDevicesAssignment := constructAllDevicesAssignment(ctx, assignment)
			scriptAssignments = append(scriptAssignments, allDevicesAssignment)
			tflog.Debug(ctx, "Added all devices assignment")
		}

		// 2. Handle All Users Assignment
		if !assignment.AllUsers.IsNull() && assignment.AllUsers.ValueBool() {
			allUsersAssignment := constructAllUsersAssignment(ctx, assignment)
			scriptAssignments = append(scriptAssignments, allUsersAssignment)
			tflog.Debug(ctx, "Added all users assignment")
		}

		// 3. Handle Include Groups (GroupAssignmentTarget objects)
		if !assignment.IncludeGroups.IsNull() && len(assignment.IncludeGroups.Elements()) > 0 {
			includeAssignments := constructIncludeGroupAssignments(ctx, assignment)
			scriptAssignments = append(scriptAssignments, includeAssignments...)
			tflog.Debug(ctx, "Added include group assignments", map[string]interface{}{
				"count": len(includeAssignments),
			})
		}

		// 4. Handle Exclude Groups (ExclusionGroupAssignmentTarget objects)
		if !assignment.ExcludeGroupIds.IsNull() && len(assignment.ExcludeGroupIds.Elements()) > 0 {
			excludeAssignments := constructExcludeGroupAssignments(ctx, assignment)
			scriptAssignments = append(scriptAssignments, excludeAssignments...)
			tflog.Debug(ctx, "Added exclude group assignments", map[string]interface{}{
				"count": len(excludeAssignments),
			})
		}
	}

	tflog.Debug(ctx, "Completed assignment construction", map[string]interface{}{
		"totalAssignments": len(scriptAssignments),
	})

	requestBody.SetDeviceHealthScriptAssignments(scriptAssignments)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructAllDevicesAssignment creates an assignment targeting all devices
func constructAllDevicesAssignment(ctx context.Context, assignment WindowsRemediationScriptAssignmentResourceModel) graphmodels.DeviceHealthScriptAssignmentable {
	tflog.Debug(ctx, "Constructing all devices assignment")

	scriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()

	// Create AllDevicesAssignmentTarget
	target := graphmodels.NewAllDevicesAssignmentTarget()

	// Set filter properties if provided
	convert.FrameworkToGraphString(assignment.FilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

	// Use enum constants instead of parsing strings
	if !assignment.Type.IsNull() && assignment.Type.ValueString() != "" {
		switch assignment.Type.ValueString() {
		case "include":
			filterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			target.SetDeviceAndAppManagementAssignmentFilterType(&filterType)
		case "exclude":
			filterType := graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			target.SetDeviceAndAppManagementAssignmentFilterType(&filterType)
		default:
			tflog.Warn(ctx, "Unknown filter type for all devices assignment", map[string]interface{}{
				"type": assignment.Type.ValueString(),
			})
		}
	}

	scriptAssignment.SetTarget(target)

	// Set runRemediationScript to true (as per Microsoft example)
	runRemediation := true
	scriptAssignment.SetRunRemediationScript(&runRemediation)

	return scriptAssignment
}

// constructAllUsersAssignment creates an assignment targeting all licensed users
func constructAllUsersAssignment(ctx context.Context, assignment WindowsRemediationScriptAssignmentResourceModel) graphmodels.DeviceHealthScriptAssignmentable {
	tflog.Debug(ctx, "Constructing all users assignment")

	scriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()

	// Create AllLicensedUsersAssignmentTarget
	target := graphmodels.NewAllLicensedUsersAssignmentTarget()

	// Set filter properties if provided
	convert.FrameworkToGraphString(assignment.FilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

	// Use enum constants instead of parsing strings
	if !assignment.Type.IsNull() && assignment.Type.ValueString() != "" {
		switch assignment.Type.ValueString() {
		case "include":
			filterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			target.SetDeviceAndAppManagementAssignmentFilterType(&filterType)
		case "exclude":
			filterType := graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			target.SetDeviceAndAppManagementAssignmentFilterType(&filterType)
		default:
			tflog.Warn(ctx, "Unknown filter type for all users assignment", map[string]interface{}{
				"type": assignment.Type.ValueString(),
			})
		}
	}

	scriptAssignment.SetTarget(target)

	// Set runRemediationScript to true (as per Microsoft example)
	runRemediation := true
	scriptAssignment.SetRunRemediationScript(&runRemediation)

	return scriptAssignment
}

// constructIncludeGroupAssignments creates group assignment targets for included groups
func constructIncludeGroupAssignments(ctx context.Context, assignment WindowsRemediationScriptAssignmentResourceModel) []graphmodels.DeviceHealthScriptAssignmentable {
	tflog.Debug(ctx, "Constructing include group assignments")

	var assignments []graphmodels.DeviceHealthScriptAssignmentable

	// Convert the Set to slice of IncludeGroupResourceModel
	var includeGroups []IncludeGroupResourceModel
	if diags := assignment.IncludeGroups.ElementsAs(ctx, &includeGroups, false); diags.HasError() {
		tflog.Error(ctx, "Failed to convert include_groups set", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return assignments
	}

	for _, group := range includeGroups {
		scriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()

		// Create GroupAssignmentTarget
		target := graphmodels.NewGroupAssignmentTarget()
		groupId := group.GroupId.ValueString()
		target.SetGroupId(&groupId)

		// Set filter properties if provided
		convert.FrameworkToGraphString(group.FilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		// Use enum constants instead of parsing strings
		if !group.Type.IsNull() && group.Type.ValueString() != "" {
			switch group.Type.ValueString() {
			case "include":
				filterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
				target.SetDeviceAndAppManagementAssignmentFilterType(&filterType)
			case "exclude":
				filterType := graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
				target.SetDeviceAndAppManagementAssignmentFilterType(&filterType)
			default:
				tflog.Warn(ctx, "Unknown filter type for group assignment", map[string]interface{}{
					"type":    group.Type.ValueString(),
					"groupId": group.GroupId.ValueString(),
				})
			}
		}

		scriptAssignment.SetTarget(target)

		// hardcode runRemediationScript to true it's never returned by API
		// in a resp. but needed for a valid request.
		runRemediation := true
		scriptAssignment.SetRunRemediationScript(&runRemediation)

		// Set run schedule if provided
		if group.RunSchedule != nil {
			if schedule := constructRunSchedule(ctx, group.RunSchedule); schedule != nil {
				scriptAssignment.SetRunSchedule(schedule)
			}
		}

		assignments = append(assignments, scriptAssignment)

		tflog.Debug(ctx, "Added include group assignment", map[string]interface{}{
			"groupId":              group.GroupId.ValueString(),
			"runRemediationScript": runRemediation,
			"hasSchedule":          group.RunSchedule != nil,
		})
	}

	return assignments
}

// constructExcludeGroupAssignments creates exclusion group assignment targets
func constructExcludeGroupAssignments(ctx context.Context, assignment WindowsRemediationScriptAssignmentResourceModel) []graphmodels.DeviceHealthScriptAssignmentable {
	tflog.Debug(ctx, "Constructing exclude group assignments")

	var assignments []graphmodels.DeviceHealthScriptAssignmentable

	// Convert the Set to slice of strings
	var excludeGroupIds []string
	if diags := assignment.ExcludeGroupIds.ElementsAs(ctx, &excludeGroupIds, false); diags.HasError() {
		tflog.Error(ctx, "Failed to convert exclude_group_ids set", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return assignments
	}

	for _, groupId := range excludeGroupIds {
		scriptAssignment := graphmodels.NewDeviceHealthScriptAssignment()

		// Create ExclusionGroupAssignmentTarget
		target := graphmodels.NewExclusionGroupAssignmentTarget()
		target.SetGroupId(&groupId)

		scriptAssignment.SetTarget(target)

		// As per Microsoft example: exclusion assignments also have runRemediationScript = true
		runRemediation := true
		scriptAssignment.SetRunRemediationScript(&runRemediation)

		// As per Microsoft example: exclusion assignments have runSchedule = null
		// Don't set any schedule for exclusion assignments

		assignments = append(assignments, scriptAssignment)

		tflog.Debug(ctx, "Added exclude group assignment", map[string]interface{}{
			"groupId":              groupId,
			"runRemediationScript": runRemediation,
		})
	}

	return assignments
}

// constructRunSchedule creates a run schedule from the model
func constructRunSchedule(ctx context.Context, scheduleModel *RunScheduleResourceModel) graphmodels.DeviceHealthScriptRunScheduleable {
	if scheduleModel == nil {
		return nil
	}

	scheduleType := scheduleModel.ScheduleType.ValueString()

	switch scheduleType {
	case "daily":
		schedule := graphmodels.NewDeviceHealthScriptDailySchedule()

		convert.FrameworkToGraphInt32(scheduleModel.Interval, schedule.SetInterval)

		if !scheduleModel.Time.IsNull() && scheduleModel.Time.ValueString() != "" {
			if err := convert.FrameworkToGraphTimeOnly(scheduleModel.Time, schedule.SetTime); err != nil {
				tflog.Warn(ctx, "Failed to convert time for daily schedule", map[string]interface{}{
					"error": err.Error(),
					"time":  scheduleModel.Time.ValueString(),
				})
			}
		}

		convert.FrameworkToGraphBool(scheduleModel.UseUtc, schedule.SetUseUtc)

		tflog.Debug(ctx, "Created daily schedule", map[string]interface{}{
			"interval": scheduleModel.Interval.ValueInt32(),
			"time":     scheduleModel.Time.ValueString(),
			"useUtc":   scheduleModel.UseUtc.ValueBool(),
		})

		return schedule

	case "hourly":
		schedule := graphmodels.NewDeviceHealthScriptHourlySchedule()

		convert.FrameworkToGraphInt32(scheduleModel.Interval, schedule.SetInterval)

		tflog.Debug(ctx, "Created hourly schedule", map[string]interface{}{
			"interval": scheduleModel.Interval.ValueInt32(),
		})

		return schedule

	case "once":
		schedule := graphmodels.NewDeviceHealthScriptRunOnceSchedule()

		convert.FrameworkToGraphInt32(scheduleModel.Interval, schedule.SetInterval)

		if !scheduleModel.Date.IsNull() && scheduleModel.Date.ValueString() != "" {
			if err := convert.FrameworkToGraphDateOnly(scheduleModel.Date, schedule.SetDate); err != nil {
				tflog.Warn(ctx, "Failed to convert date for once schedule", map[string]interface{}{
					"error": err.Error(),
					"date":  scheduleModel.Date.ValueString(),
				})
			}
		}

		if !scheduleModel.Time.IsNull() && scheduleModel.Time.ValueString() != "" {
			if err := convert.FrameworkToGraphTimeOnly(scheduleModel.Time, schedule.SetTime); err != nil {
				tflog.Warn(ctx, "Failed to convert time for once schedule", map[string]interface{}{
					"error": err.Error(),
					"time":  scheduleModel.Time.ValueString(),
				})
			}
		}

		convert.FrameworkToGraphBool(scheduleModel.UseUtc, schedule.SetUseUtc)

		tflog.Debug(ctx, "Created once schedule", map[string]interface{}{
			"interval": scheduleModel.Interval.ValueInt32(),
			"date":     scheduleModel.Date.ValueString(),
			"time":     scheduleModel.Time.ValueString(),
			"useUtc":   scheduleModel.UseUtc.ValueBool(),
		})

		return schedule

	default:
		tflog.Warn(ctx, "Unknown schedule type", map[string]interface{}{
			"scheduleType": scheduleType,
		})
		return nil
	}
}
