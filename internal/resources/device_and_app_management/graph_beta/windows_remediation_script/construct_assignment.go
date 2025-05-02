package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs and returns a DeviceHealthScriptsItemAssignPostRequestBody
func constructAssignment(ctx context.Context, assignments []WindowsRemediationScriptAssignmentResourceModel) (devicemanagement.DeviceHealthScriptsItemAssignPostRequestBodyable, error) {
	if assignments == nil {
		return nil, fmt.Errorf("assignments configuration block is required even if empty. Minimum config requires all_devices and all_users booleans to be set to false")
	}

	tflog.Debug(ctx, "Starting Device Health Script assignment construction", map[string]interface{}{
		"assignmentCount": len(assignments),
	})

	requestBody := devicemanagement.NewDeviceHealthScriptsItemAssignPostRequestBody()
	scriptAssignments := make([]graphmodels.DeviceHealthScriptAssignmentable, 0)

	// Process each assignment block
	for idx, assignment := range assignments {
		tflog.Debug(ctx, "Processing assignment block", map[string]interface{}{
			"index":             idx,
			"allDevices":        assignment.AllDevices.ValueBool(),
			"allUsers":          assignment.AllUsers.ValueBool(),
			"includeGroupsNull": assignment.IncludeGroups.IsNull(),
			"includeGroupsLen":  len(assignment.IncludeGroups.Elements()),
			"excludeGroupsNull": assignment.ExcludeGroupIds.IsNull(),
			"excludeGroupsLen":  len(assignment.ExcludeGroupIds.Elements()),
		})

		// Check All Devices
		if !assignment.AllDevices.IsNull() && assignment.AllDevices.ValueBool() {
			tflog.Debug(ctx, "Adding all devices assignment")
			scriptAssignments = append(scriptAssignments, constructAllDevicesAssignment(ctx, assignment))
		}

		// Check All Users
		if !assignment.AllUsers.IsNull() && assignment.AllUsers.ValueBool() {
			tflog.Debug(ctx, "Adding all users assignment")
			scriptAssignments = append(scriptAssignments, constructAllUsersAssignment(ctx, assignment))
		}

		// Check Include Groups
		if !assignment.IncludeGroups.IsNull() && len(assignment.IncludeGroups.Elements()) > 0 {
			tflog.Debug(ctx, "Processing include groups", map[string]interface{}{
				"count": len(assignment.IncludeGroups.Elements()),
			})
			includeAssignments := constructGroupIncludeAssignments(ctx, assignment)
			if len(includeAssignments) > 0 {
				scriptAssignments = append(scriptAssignments, includeAssignments...)
				tflog.Debug(ctx, "Added include group assignments", map[string]interface{}{
					"addedCount": len(includeAssignments),
				})
			} else {
				tflog.Warn(ctx, "No include group assignments were constructed")
			}
		}

		// Check Exclude Groups
		if !assignment.ExcludeGroupIds.IsNull() && len(assignment.ExcludeGroupIds.Elements()) > 0 {
			tflog.Debug(ctx, "Processing exclude groups", map[string]interface{}{
				"count": len(assignment.ExcludeGroupIds.Elements()),
			})
			excludeAssignments := constructGroupExcludeAssignments(ctx, assignment)
			if len(excludeAssignments) > 0 {
				scriptAssignments = append(scriptAssignments, excludeAssignments...)
				tflog.Debug(ctx, "Added exclude group assignments", map[string]interface{}{
					"addedCount": len(excludeAssignments),
				})
			}
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

// constructGroupIncludeAssignments constructs and returns a list of DeviceHealthScriptAssignment objects for included groups
func constructGroupIncludeAssignments(ctx context.Context, config WindowsRemediationScriptAssignmentResourceModel) []graphmodels.DeviceHealthScriptAssignmentable {
	var assignments []graphmodels.DeviceHealthScriptAssignmentable

	tflog.Debug(ctx, "Entering constructGroupIncludeAssignments", map[string]interface{}{
		"includeGroups.IsNull":    config.IncludeGroups.IsNull(),
		"includeGroups.IsUnknown": config.IncludeGroups.IsUnknown(),
	})

	if config.IncludeGroups.IsNull() {
		tflog.Debug(ctx, "IncludeGroups is null, returning empty assignments")
		return assignments
	}

	// Parse IncludeGroups set
	var includeGroups []IncludeGroupResourceModel
	diags := config.IncludeGroups.ElementsAs(ctx, &includeGroups, false)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to parse include groups", map[string]interface{}{
			"error": diags.Errors(),
		})
		return assignments
	}

	tflog.Debug(ctx, "Successfully parsed include groups", map[string]interface{}{
		"count": len(includeGroups),
	})

	for idx, group := range includeGroups {
		tflog.Debug(ctx, "Processing include group", map[string]interface{}{
			"index":   idx,
			"groupId": group.GroupId.ValueString(),
		})

		if group.GroupId.IsNull() || group.GroupId.IsUnknown() || group.GroupId.ValueString() == "" {
			tflog.Warn(ctx, "Skipping group with null/empty GroupId", map[string]interface{}{
				"index": idx,
			})
			continue
		}

		assignment := graphmodels.NewDeviceHealthScriptAssignment()
		target := graphmodels.NewGroupAssignmentTarget()

		constructors.SetStringProperty(group.GroupId, target.SetGroupId)

		// Set filter ID and type if present
		if !group.IncludeGroupsFilterId.IsNull() && !group.IncludeGroupsFilterId.IsUnknown() {
			constructors.SetStringProperty(group.IncludeGroupsFilterId,
				target.SetDeviceAndAppManagementAssignmentFilterId)

			if !group.IncludeGroupsFilterType.IsNull() && !group.IncludeGroupsFilterType.IsUnknown() {
				tflog.Debug(ctx, "Setting filter type", map[string]interface{}{
					"groupId":    group.GroupId.ValueString(),
					"filterType": group.IncludeGroupsFilterType.ValueString(),
				})

				err := constructors.SetEnumProperty(group.IncludeGroupsFilterType,
					graphmodels.ParseDeviceAndAppManagementAssignmentFilterType,
					target.SetDeviceAndAppManagementAssignmentFilterType)
				if err != nil {
					tflog.Warn(ctx, "Failed to parse include groups filter type", map[string]interface{}{
						"error":   err.Error(),
						"groupId": group.GroupId.ValueString(),
					})
				}
			}
		}

		assignment.SetTarget(target)

		// Set run remediation script if specified
		if !group.RunRemediationScript.IsNull() {
			runRemediation := group.RunRemediationScript.ValueBool()
			assignment.SetRunRemediationScript(&runRemediation)
			tflog.Debug(ctx, "Set RunRemediationScript", map[string]interface{}{
				"groupId":              group.GroupId.ValueString(),
				"runRemediationScript": runRemediation,
			})
		}

		if group.RunSchedule != nil {
			tflog.Debug(ctx, "Processing run schedule", map[string]interface{}{
				"groupId":      group.GroupId.ValueString(),
				"scheduleType": group.RunSchedule.ScheduleType.ValueString(),
			})

			schedule := constructRunSchedule(ctx, group.RunSchedule)
			if schedule != nil {
				assignment.SetRunSchedule(schedule)
				tflog.Debug(ctx, "Set run schedule for group", map[string]interface{}{
					"groupId": group.GroupId.ValueString(),
				})
			}
		}

		assignments = append(assignments, assignment)
		tflog.Debug(ctx, "Added include group assignment", map[string]interface{}{
			"groupId": group.GroupId.ValueString(),
			"index":   idx,
		})
	}

	tflog.Debug(ctx, "Completed constructGroupIncludeAssignments", map[string]interface{}{
		"totalAssignments": len(assignments),
	})

	return assignments
}

// constructRunSchedule constructs the appropriate schedule type based on the schedule model
func constructRunSchedule(ctx context.Context, schedule *RunScheduleResourceModel) graphmodels.DeviceHealthScriptRunScheduleable {
	if schedule == nil {
		tflog.Debug(ctx, "Schedule is nil, returning nil")
		return nil
	}

	tflog.Debug(ctx, "Constructing run schedule from model", map[string]interface{}{
		"scheduleType": schedule.ScheduleType.ValueString(),
		"interval":     schedule.Interval.ValueInt32(),
		"time":         schedule.Time.ValueString(),
		"date":         schedule.Date.ValueString(),
		"useUtc":       schedule.UseUtc.ValueBool(),
	})

	switch schedule.ScheduleType.ValueString() {
	case "daily":
		tflog.Debug(ctx, "Creating daily schedule")
		dailySchedule := graphmodels.NewDeviceHealthScriptDailySchedule()
		constructors.SetInt32Property(schedule.Interval, dailySchedule.SetInterval)
		constructors.StringToTimeOnly(schedule.Time, dailySchedule.SetTime)
		constructors.SetBoolProperty(schedule.UseUtc, dailySchedule.SetUseUtc)
		return dailySchedule

	case "hourly":
		tflog.Debug(ctx, "Creating hourly schedule")
		hourlySchedule := graphmodels.NewDeviceHealthScriptHourlySchedule()
		constructors.SetInt32Property(schedule.Interval, hourlySchedule.SetInterval)
		return hourlySchedule

	case "once":
		tflog.Debug(ctx, "Creating once schedule")
		onceSchedule := graphmodels.NewDeviceHealthScriptRunOnceSchedule()
		constructors.SetInt32Property(schedule.Interval, onceSchedule.SetInterval)
		constructors.StringToDateOnly(schedule.Date, onceSchedule.SetDate)
		constructors.StringToTimeOnly(schedule.Time, onceSchedule.SetTime)
		constructors.SetBoolProperty(schedule.UseUtc, onceSchedule.SetUseUtc)
		return onceSchedule

	default:
		tflog.Warn(ctx, "Unknown schedule type", map[string]interface{}{
			"scheduleType": schedule.ScheduleType.ValueString(),
		})
		return nil
	}
}

// constructAllDevicesAssignment constructs and returns a DeviceHealthScriptAssignment object for all devices
func constructAllDevicesAssignment(ctx context.Context, config WindowsRemediationScriptAssignmentResourceModel) graphmodels.DeviceHealthScriptAssignmentable {
	assignment := graphmodels.NewDeviceHealthScriptAssignment()
	target := graphmodels.NewAllDevicesAssignmentTarget()

	if !config.AllDevicesFilterId.IsNull() && !config.AllDevicesFilterId.IsUnknown() &&
		config.AllDevicesFilterId.ValueString() != "" {
		constructors.SetStringProperty(config.AllDevicesFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		if !config.AllDevicesFilterType.IsNull() && !config.AllDevicesFilterType.IsUnknown() {
			err := constructors.SetEnumProperty(config.AllDevicesFilterType,
				graphmodels.ParseDeviceAndAppManagementAssignmentFilterType,
				target.SetDeviceAndAppManagementAssignmentFilterType)
			if err != nil {
				tflog.Warn(ctx, "Failed to parse all devices filter type", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}

	assignment.SetTarget(target)
	return assignment
}

// constructAllUsersAssignment constructs and returns a DeviceHealthScriptAssignment object for all licensed users
func constructAllUsersAssignment(ctx context.Context, config WindowsRemediationScriptAssignmentResourceModel) graphmodels.DeviceHealthScriptAssignmentable {
	assignment := graphmodels.NewDeviceHealthScriptAssignment()
	target := graphmodels.NewAllLicensedUsersAssignmentTarget()

	if !config.AllUsersFilterId.IsNull() && !config.AllUsersFilterId.IsUnknown() &&
		config.AllUsersFilterId.ValueString() != "" {
		constructors.SetStringProperty(config.AllUsersFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		if !config.AllUsersFilterType.IsNull() && !config.AllUsersFilterType.IsUnknown() {
			err := constructors.SetEnumProperty(config.AllUsersFilterType,
				graphmodels.ParseDeviceAndAppManagementAssignmentFilterType,
				target.SetDeviceAndAppManagementAssignmentFilterType)
			if err != nil {
				tflog.Warn(ctx, "Failed to parse all users filter type", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}

	assignment.SetTarget(target)
	return assignment
}

// constructGroupExcludeAssignments constructs and returns a list of DeviceHealthScriptAssignment objects for excluded groups
func constructGroupExcludeAssignments(ctx context.Context, config WindowsRemediationScriptAssignmentResourceModel) []graphmodels.DeviceHealthScriptAssignmentable {
	var assignments []graphmodels.DeviceHealthScriptAssignmentable

	for _, elem := range config.ExcludeGroupIds.Elements() {
		if !elem.IsNull() && !elem.IsUnknown() && elem.String() != "" {
			assignment := graphmodels.NewDeviceHealthScriptAssignment()
			target := graphmodels.NewExclusionGroupAssignmentTarget()

			groupId := strings.Trim(elem.String(), "\"")
			target.SetGroupId(&groupId)

			assignment.SetTarget(target)
			assignments = append(assignments, assignment)
		}
	}

	return assignments
}
