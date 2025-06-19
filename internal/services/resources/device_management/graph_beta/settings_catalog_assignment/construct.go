package graphBetaDeviceManagementConfigurationPolicyAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a DeviceManagementConfigurationPolicyAssignment
func constructResource(ctx context.Context, data DeviceManagementConfigurationPolicyAssignmentResourceModel) (graphmodels.DeviceManagementConfigurationPolicyAssignmentable, error) {
	tflog.Debug(ctx, "Starting device management configuration policy assignment construction")

	assignment := graphmodels.NewDeviceManagementConfigurationPolicyAssignment()

	// Set source
	err := convert.FrameworkToGraphEnum(
		data.Source,
		graphmodels.ParseDeviceAndAppManagementAssignmentSource,
		func(val *graphmodels.DeviceAndAppManagementAssignmentSource) {
			assignment.SetSource(val)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error setting assignment source: %v", err)
	}

	convert.FrameworkToGraphString(data.SourceId, assignment.SetSourceId)

	target, err := constructAssignmentTarget(ctx, &data.Target)
	if err != nil {
		return nil, fmt.Errorf("error constructing configuration policy assignment target: %v", err)
	}
	assignment.SetTarget(target)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed configuration policy assignment", assignment); err != nil {
		tflog.Error(ctx, "Failed to log configuration policy assignment", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return assignment, nil
}

// constructAssignmentTarget constructs the configuration policy assignment target
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

	tflog.Debug(ctx, "Finished constructing assignment target")
	return target, nil
}
