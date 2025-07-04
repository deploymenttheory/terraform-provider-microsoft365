package graphBetaMacosCustomAttributeScriptAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a DeviceManagementScriptAssignment
func constructResource(ctx context.Context, data MacosCustomAttributeScriptAssignmentResourceModel) (graphmodels.DeviceManagementScriptAssignmentable, error) {
	tflog.Debug(ctx, "Starting device management script assignment construction")

	assignment := graphmodels.NewDeviceManagementScriptAssignment()

	// Set Target
	target, err := constructAssignmentTarget(ctx, &data.Target)
	if err != nil {
		return nil, fmt.Errorf("error constructing device management script assignment target: %v", err)
	}
	assignment.SetTarget(target)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed device management script assignment", assignment); err != nil {
		tflog.Error(ctx, "Failed to log device management script assignment", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return assignment, nil
}

// constructAssignmentTarget constructs the assignment target
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
	case "groupAssignment":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		convert.FrameworkToGraphString(data.GroupId, groupTarget.SetGroupId)
		target = groupTarget
	default:
		target = graphmodels.NewDeviceAndAppManagementAssignmentTarget()
	}

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
