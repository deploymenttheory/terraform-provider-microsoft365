// Updated constructor functions to match the simplified model structure
package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs and returns a DeviceHealthScriptsItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *DeviceHealthScriptResourceModel) (devicemanagement.DeviceHealthScriptsItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Starting Device Health Script assignment construction")

	requestBody := devicemanagement.NewDeviceHealthScriptsItemAssignPostRequestBody()
	scriptAssignments := make([]graphmodels.DeviceHealthScriptAssignmentable, 0)

	if data.Assignments.IsNull() || data.Assignments.IsUnknown() {
		tflog.Debug(ctx, "Assignments is null or unknown, creating empty assignments array")
		requestBody.SetDeviceHealthScriptAssignments(scriptAssignments)
		return requestBody, nil
	}

	// Extract assignments using the proper struct types
	var terraformAssignments []WindowsRemediationScriptAssignmentModel
	diags := data.Assignments.ElementsAs(ctx, &terraformAssignments, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract assignments: %v", diags.Errors())
	}

	for idx, assignment := range terraformAssignments {
		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"index": idx,
		})

		// Create a new assignment
		graphAssignment := graphmodels.NewDeviceHealthScriptAssignment()

		// Process the assignment target based on type
		if assignment.Type.IsNull() || assignment.Type.IsUnknown() {
			tflog.Error(ctx, "Assignment target type is missing or invalid", map[string]interface{}{
				"index": idx,
			})
			continue
		}

		targetType := assignment.Type.ValueString()

		// Create the target based on the target type
		target := constructTarget(ctx, targetType, assignment)
		if target == nil {
			tflog.Error(ctx, "Failed to create target", map[string]interface{}{
				"index":      idx,
				"targetType": targetType,
			})
			continue
		}

		// Set the target on the assignment
		graphAssignment.SetTarget(target)

		// Set run schedule if provided
		runSchedule := constructRunSchedule(ctx, assignment)
		if runSchedule != nil {
			graphAssignment.SetRunSchedule(runSchedule)
		}

		scriptAssignments = append(scriptAssignments, graphAssignment)
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

// constructTarget creates the appropriate target based on the target type
func constructTarget(ctx context.Context, targetType string, assignment WindowsRemediationScriptAssignmentModel) graphmodels.DeviceAndAppManagementAssignmentTargetable {
	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch targetType {
	case "allDevicesAssignmentTarget":
		target = graphmodels.NewAllDevicesAssignmentTarget()
	case "allLicensedUsersAssignmentTarget":
		target = graphmodels.NewAllLicensedUsersAssignmentTarget()
	case "groupAssignmentTarget":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
			// Use helper to set group ID
			convert.FrameworkToGraphString(assignment.GroupId, groupTarget.SetGroupId)
		}
		target = groupTarget
	case "exclusionGroupAssignmentTarget":
		exclusionTarget := graphmodels.NewExclusionGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
			// Use helper to set group ID
			convert.FrameworkToGraphString(assignment.GroupId, exclusionTarget.SetGroupId)
		}
		target = exclusionTarget
	default:
		tflog.Error(ctx, "Unsupported target type", map[string]interface{}{
			"targetType": targetType,
		})
		return nil
	}

	// Set filter if provided
	if !assignment.FilterId.IsNull() && !assignment.FilterId.IsUnknown() && assignment.FilterId.ValueString() != "" {
		if targetWithFilter, ok := target.(graphmodels.DeviceAndAppManagementAssignmentTargetable); ok {
			// Use helper to set filter ID
			convert.FrameworkToGraphString(assignment.FilterId, targetWithFilter.SetDeviceAndAppManagementAssignmentFilterId)

			// Set filter type if provided
			if !assignment.FilterType.IsNull() && !assignment.FilterType.IsUnknown() && assignment.FilterType.ValueString() != "" {
				filterType := assignment.FilterType.ValueString()
				var filterTypeEnum graphmodels.DeviceAndAppManagementAssignmentFilterType
				switch filterType {
				case "include":
					filterTypeEnum = graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
				case "exclude":
					filterTypeEnum = graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
				default:
					return target
				}
				targetWithFilter.SetDeviceAndAppManagementAssignmentFilterType(&filterTypeEnum)
			}
		}
	}

	return target
}

// constructRunSchedule creates a run schedule from the assignment
func constructRunSchedule(ctx context.Context, assignment WindowsRemediationScriptAssignmentModel) graphmodels.DeviceHealthScriptRunScheduleable {
	// Check for daily schedule
	if assignment.DailySchedule != nil {
		dailySchedule := graphmodels.NewDeviceHealthScriptDailySchedule()

		// Set interval using Int32 helper
		if !assignment.DailySchedule.Interval.IsNull() && !assignment.DailySchedule.Interval.IsUnknown() {
			// Convert int32 to int32 for the API
			interval := int32(assignment.DailySchedule.Interval.ValueInt32())
			dailySchedule.SetInterval(&interval)
		}

		// Set time using TimeOnly helper
		if !assignment.DailySchedule.Time.IsNull() && !assignment.DailySchedule.Time.IsUnknown() && assignment.DailySchedule.Time.ValueString() != "" {
			// Use FrameworkToGraphTimeOnlyWithPrecision helper
			err := convert.FrameworkToGraphTimeOnlyWithPrecision(assignment.DailySchedule.Time, 0, dailySchedule.SetTime)
			if err != nil {
				tflog.Error(ctx, "Failed to parse daily schedule time", map[string]interface{}{
					"time":  assignment.DailySchedule.Time.ValueString(),
					"error": err.Error(),
				})
			}
		}

		// Set useUtc using Bool helper
		if !assignment.DailySchedule.UseUtc.IsNull() && !assignment.DailySchedule.UseUtc.IsUnknown() {
			convert.FrameworkToGraphBool(assignment.DailySchedule.UseUtc, dailySchedule.SetUseUtc)
		}

		return dailySchedule
	}

	// Check for hourly schedule
	if assignment.HourlySchedule != nil {
		hourlySchedule := graphmodels.NewDeviceHealthScriptHourlySchedule()

		// Set interval using Int32 helper
		if !assignment.HourlySchedule.Interval.IsNull() && !assignment.HourlySchedule.Interval.IsUnknown() {
			// Convert int32 to int32 for the API
			interval := int32(assignment.HourlySchedule.Interval.ValueInt32())
			hourlySchedule.SetInterval(&interval)
		}

		return hourlySchedule
	}

	// Check for run once schedule
	if assignment.RunOnceSchedule != nil {
		runOnceSchedule := graphmodels.NewDeviceHealthScriptRunOnceSchedule()

		// Set date using DateOnly helper
		if !assignment.RunOnceSchedule.Date.IsNull() && !assignment.RunOnceSchedule.Date.IsUnknown() && assignment.RunOnceSchedule.Date.ValueString() != "" {
			// Use DateOnly helper
			err := convert.FrameworkToGraphDateOnly(assignment.RunOnceSchedule.Date, runOnceSchedule.SetDate)
			if err != nil {
				tflog.Error(ctx, "Failed to parse run once schedule date", map[string]interface{}{
					"date":  assignment.RunOnceSchedule.Date.ValueString(),
					"error": err.Error(),
				})
			}
		}

		// Set time using TimeOnly helper
		if !assignment.RunOnceSchedule.Time.IsNull() && !assignment.RunOnceSchedule.Time.IsUnknown() && assignment.RunOnceSchedule.Time.ValueString() != "" {
			// Use FrameworkToGraphTimeOnlyWithPrecision helper directly
			err := convert.FrameworkToGraphTimeOnlyWithPrecision(assignment.RunOnceSchedule.Time, 0, runOnceSchedule.SetTime)
			if err != nil {
				tflog.Error(ctx, "Failed to parse run once schedule time", map[string]interface{}{
					"time":  assignment.RunOnceSchedule.Time.ValueString(),
					"error": err.Error(),
				})
			}
		}

		// Set useUtc using Bool helper
		if !assignment.RunOnceSchedule.UseUtc.IsNull() && !assignment.RunOnceSchedule.UseUtc.IsUnknown() {
			convert.FrameworkToGraphBool(assignment.RunOnceSchedule.UseUtc, runOnceSchedule.SetUseUtc)
		}

		return runOnceSchedule
	}

	return nil
}
