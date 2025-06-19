package graphBetaWindowsRemediationScriptAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote assignment to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data DeviceHealthScriptAssignmentResourceModel, assignment graphmodels.DeviceHealthScriptAssignmentable) DeviceHealthScriptAssignmentResourceModel {
	if assignment == nil {
		tflog.Debug(ctx, "Remote assignment is nil")
		return data
	}

	data.ID = convert.GraphToFrameworkString(assignment.GetId())
	data.RunRemediationScript = convert.GraphToFrameworkBool(assignment.GetRunRemediationScript())

	if target := assignment.GetTarget(); target != nil {
		data.Target = mapRemoteTargetToTerraform(target)
	}

	if runSchedule := assignment.GetRunSchedule(); runSchedule != nil {
		data.RunSchedule = mapRemoteRunScheduleToTerraform(runSchedule)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}

// mapRemoteTargetToTerraform maps a remote assignment target to a Terraform assignment target
func mapRemoteTargetToTerraform(remoteTarget graphmodels.DeviceAndAppManagementAssignmentTargetable) AssignmentTargetResourceModel {
	target := AssignmentTargetResourceModel{
		DeviceAndAppManagementAssignmentFilterId:   convert.GraphToFrameworkString(remoteTarget.GetDeviceAndAppManagementAssignmentFilterId()),
		DeviceAndAppManagementAssignmentFilterType: convert.GraphToFrameworkEnum(remoteTarget.GetDeviceAndAppManagementAssignmentFilterType()),
	}

	switch v := remoteTarget.(type) {
	case *graphmodels.GroupAssignmentTarget:
		target.TargetType = types.StringValue("groupAssignment")
		target.GroupId = convert.GraphToFrameworkString(v.GetGroupId())
	case *graphmodels.ExclusionGroupAssignmentTarget:
		target.TargetType = types.StringValue("exclusionGroupAssignment")
		target.GroupId = convert.GraphToFrameworkString(v.GetGroupId())
	case *graphmodels.ConfigurationManagerCollectionAssignmentTarget:
		target.TargetType = types.StringValue("configurationManagerCollection")
		target.CollectionId = convert.GraphToFrameworkString(v.GetCollectionId())
	case *graphmodels.AllDevicesAssignmentTarget:
		target.TargetType = types.StringValue("allDevices")
	case *graphmodels.AllLicensedUsersAssignmentTarget:
		target.TargetType = types.StringValue("allLicensedUsers")
	}

	return target
}

// mapRemoteRunScheduleToTerraform maps a remote run schedule to the Terraform model
func mapRemoteRunScheduleToTerraform(remoteSchedule graphmodels.DeviceHealthScriptRunScheduleable) *RunScheduleResourceModel {
	if remoteSchedule == nil {
		return nil
	}

	switch schedule := remoteSchedule.(type) {
	case *graphmodels.DeviceHealthScriptDailySchedule:
		return &RunScheduleResourceModel{
			Daily: &DailyScheduleResourceModel{
				Interval: convert.GraphToFrameworkInt32(schedule.GetInterval()),
				UseUtc:   convert.GraphToFrameworkBool(schedule.GetUseUtc()),
				Time:     convert.GraphToFrameworkTimeOnly(schedule.GetTime()),
			},
		}

	case *graphmodels.DeviceHealthScriptHourlySchedule:
		return &RunScheduleResourceModel{
			Hourly: &HourlyScheduleResourceModel{
				Interval: convert.GraphToFrameworkInt32(schedule.GetInterval()),
			},
		}

	case *graphmodels.DeviceHealthScriptRunOnceSchedule:
		return &RunScheduleResourceModel{
			Once: &RunOnceScheduleResourceModel{
				Interval: convert.GraphToFrameworkInt32(schedule.GetInterval()),
				Date:     convert.GraphToFrameworkDateOnly(schedule.GetDate()),
				Time:     convert.GraphToFrameworkTimeOnly(schedule.GetTime()),
				UseUtc:   convert.GraphToFrameworkBool(schedule.GetUseUtc()),
			},
		}

	default:
		// Unknown schedule type: omit from state
		return nil
	}
}
