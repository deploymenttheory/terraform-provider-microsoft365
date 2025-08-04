package graphBetaMacOSSoftwareUpdateConfiguration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteResourceStateToTerraform maps the Graph API model to the Terraform state model.
func mapRemoteResourceStateToTerraform(ctx context.Context, data *MacOSSoftwareUpdateConfigurationResourceModel, remoteResource graphmodels.DeviceConfigurationable) {
	if remoteResource == nil {
		return
	}

	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]interface{}{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	resource, ok := remoteResource.(graphmodels.MacOSSoftwareUpdateConfigurationable)
	if !ok {
		return
	}

	data.ID = convert.GraphToFrameworkString(resource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(resource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(resource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, resource.GetRoleScopeTagIds())
	data.CriticalUpdateBehavior = convert.GraphToFrameworkEnum(resource.GetCriticalUpdateBehavior())
	data.ConfigDataUpdateBehavior = convert.GraphToFrameworkEnum(resource.GetConfigDataUpdateBehavior())
	data.FirmwareUpdateBehavior = convert.GraphToFrameworkEnum(resource.GetFirmwareUpdateBehavior())
	data.AllOtherUpdateBehavior = convert.GraphToFrameworkEnum(resource.GetAllOtherUpdateBehavior())
	data.UpdateScheduleType = convert.GraphToFrameworkEnum(resource.GetUpdateScheduleType())
	data.Priority = convert.GraphToFrameworkEnum(resource.GetPriority())
	data.UpdateTimeWindowUtcOffsetInMinutes = convert.GraphToFrameworkInt32(resource.GetUpdateTimeWindowUtcOffsetInMinutes())
	data.MaxUserDeferralsCount = convert.GraphToFrameworkInt32(resource.GetMaxUserDeferralsCount())

	// Custom update time windows
	objType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"start_day":  types.StringType,
			"end_day":    types.StringType,
			"start_time": types.StringType,
			"end_time":   types.StringType,
		},
	}
	var windows []attr.Value
	if remoteResourceWindows := resource.GetCustomUpdateTimeWindows(); remoteResourceWindows != nil {
		for _, win := range remoteResourceWindows {
			if win == nil {
				continue
			}
			obj, _ := types.ObjectValue(
				objType.AttrTypes,
				map[string]attr.Value{
					"start_day":  convert.GraphToFrameworkEnum(win.GetStartDay()),
					"end_day":    convert.GraphToFrameworkEnum(win.GetEndDay()),
					"start_time": convert.GraphToFrameworkTimeOnly(win.GetStartTime()),
					"end_time":   convert.GraphToFrameworkTimeOnly(win.GetEndTime()),
				},
			)
			windows = append(windows, obj)
		}
	}
	if len(windows) == 0 {
		data.CustomUpdateTimeWindows = types.ListNull(objType)
	} else {
		data.CustomUpdateTimeWindows = types.ListValueMust(objType, windows)
	}

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		data.Assignments = types.SetNull(MacOSSoftwareUpdateAssignmentType())
	} else {
		MapAssignmentsToTerraform(ctx, data, assignments)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
