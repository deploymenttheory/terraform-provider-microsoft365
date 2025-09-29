package graphDeviceConfigurationAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *DeviceConfigurationAssignmentResourceModel) (*models.DeviceConfigurationAssignment, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewDeviceConfigurationAssignment()

	// Construct target based on target type
	target, err := constructAssignmentTarget(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("error constructing assignment target: %v", err)
	}

	// Set the target to the request body
	if target != nil {
		requestBody.SetTarget(target)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructAssignmentTarget constructs the device configuration assignment target
func constructAssignmentTarget(ctx context.Context, data *DeviceConfigurationAssignmentResourceModel) (models.DeviceAndAppManagementAssignmentTargetable, error) {
	if data == nil {
		return nil, fmt.Errorf("assignment data is required")
	}

	var target models.DeviceAndAppManagementAssignmentTargetable

	switch data.TargetType.ValueString() {
	case "allDevices":
		target = models.NewAllDevicesAssignmentTarget()

	case "allLicensedUsers":
		target = models.NewAllLicensedUsersAssignmentTarget()

	case "groupAssignment":
		groupTarget := models.NewGroupAssignmentTarget()
		convert.FrameworkToGraphString(data.GroupId, groupTarget.SetGroupId)
		target = groupTarget

	case "exclusionGroupAssignment":
		exclusionTarget := models.NewExclusionGroupAssignmentTarget()
		convert.FrameworkToGraphString(data.GroupId, exclusionTarget.SetGroupId)
		target = exclusionTarget

	case "configurationManagerCollection":
		configManagerTarget := models.NewConfigurationManagerCollectionAssignmentTarget()
		convert.FrameworkToGraphString(data.GroupId, configManagerTarget.SetCollectionId)
		target = configManagerTarget

	default:
		return nil, fmt.Errorf("unsupported target type: %s", data.TargetType.ValueString())
	}

	tflog.Debug(ctx, "Finished constructing assignment target")
	return target, nil
}
