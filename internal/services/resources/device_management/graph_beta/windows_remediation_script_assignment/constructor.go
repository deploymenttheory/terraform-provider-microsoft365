package graphBetaWindowsRemediationScriptAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a DeviceHealthScriptAssignment
func constructResource(ctx context.Context, data DeviceHealthScriptAssignmentResourceModel) (graphmodels.DeviceHealthScriptAssignmentable, error) {
	tflog.Debug(ctx, "Starting device health script assignment construction")

	assignment := graphmodels.NewDeviceHealthScriptAssignment()

	convert.FrameworkToGraphBool(data.RunRemediationScript, assignment.SetRunRemediationScript)

	// Set Target
	target, err := constructAssignmentTarget(ctx, &data.Target)
	if err != nil {
		return nil, fmt.Errorf("error constructing device health script assignment target: %v", err)
	}
	assignment.SetTarget(target)

	// Set RunSchedule
	if data.RunSchedule != nil {
		runSchedule, err := constructRunSchedule(ctx, data.RunSchedule)
		if err != nil {
			return nil, fmt.Errorf("error constructing run schedule: %v", err)
		}
		if runSchedule != nil {
			assignment.SetRunSchedule(runSchedule)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed device health script assignment", assignment); err != nil {
		tflog.Error(ctx, "Failed to log device health script assignment", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return assignment, nil
}

// constructAssignmentTarget constructs the device health script assignment target
func constructAssignmentTarget(ctx context.Context, data *AssignmentTargetResourceModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	if data == nil {
		return nil, fmt.Errorf("assignment target data is required")
	}

	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch data.TargetType.ValueString() {
	case "allDevices":
		target = graphmodels.NewAllDevicesAssignmentTarget()
	case "allLicensedUsers":
		target = graphmodels.NewAllLicensedUsersAssignmentTarget()
	case "configurationManagerCollection":
		configManagerTarget := graphmodels.NewConfigurationManagerCollectionAssignmentTarget()
		convert.FrameworkToGraphString(data.CollectionId, configManagerTarget.SetCollectionId)
		target = configManagerTarget
	case "exclusionGroupAssignment":
		exclusionGroupTarget := graphmodels.NewExclusionGroupAssignmentTarget()
		convert.FrameworkToGraphString(data.GroupId, exclusionGroupTarget.SetGroupId)
		target = exclusionGroupTarget
	case "groupAssignment":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		convert.FrameworkToGraphString(data.GroupId, groupTarget.SetGroupId)
		target = groupTarget
	default:
		target = graphmodels.NewDeviceAndAppManagementAssignmentTarget()
	}

	// Set the filter properties using helpers
	convert.FrameworkToGraphString(data.DeviceAndAppManagementAssignmentFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

	// Set filter type enum property
	if !data.DeviceAndAppManagementAssignmentFilterType.IsNull() && !data.DeviceAndAppManagementAssignmentFilterType.IsUnknown() {
		err := convert.FrameworkToGraphEnum(
			data.DeviceAndAppManagementAssignmentFilterType,
			graphmodels.ParseDeviceAndAppManagementAssignmentFilterType,
			func(val *graphmodels.DeviceAndAppManagementAssignmentFilterType) {
				target.SetDeviceAndAppManagementAssignmentFilterType(val)
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error setting assignment filter type: %v", err)
		}
	}

	tflog.Debug(ctx, "Finished constructing assignment target")
	return target, nil
}

// constructRunSchedule constructs the device health script run schedule
func constructRunSchedule(ctx context.Context, data *RunScheduleResourceModel) (graphmodels.DeviceHealthScriptRunScheduleable, error) {
	if data == nil {
		return nil, nil
	}

	tflog.Debug(ctx, "Constructing run schedule")

	// Determine which schedule type is defined and construct the appropriate schedule
	if data.Daily != nil {
		dailySchedule := graphmodels.NewDeviceHealthScriptDailySchedule()
		convert.FrameworkToGraphTimeOnly(data.Daily.Time, dailySchedule.SetTime)
		convert.FrameworkToGraphInt32(data.Daily.Interval, dailySchedule.SetInterval)
		convert.FrameworkToGraphBool(data.Daily.UseUtc, dailySchedule.SetUseUtc)

		return dailySchedule, nil
	}

	if data.Hourly != nil {
		hourlySchedule := graphmodels.NewDeviceHealthScriptHourlySchedule()
		convert.FrameworkToGraphInt32(data.Hourly.Interval, hourlySchedule.SetInterval)

		return hourlySchedule, nil
	}

	if data.Once != nil {
		onceSchedule := graphmodels.NewDeviceHealthScriptRunOnceSchedule()

		convert.FrameworkToGraphBool(data.Once.UseUtc, onceSchedule.SetUseUtc)
		convert.FrameworkToGraphTimeOnly(data.Once.Time, onceSchedule.SetTime)
		convert.FrameworkToGraphDateOnly(data.Once.Date, onceSchedule.SetDate)

		return onceSchedule, nil
	}

	// No schedule defined
	return nil, fmt.Errorf("no valid schedule configuration provided, must define exactly one of: daily, hourly, or once")
}
