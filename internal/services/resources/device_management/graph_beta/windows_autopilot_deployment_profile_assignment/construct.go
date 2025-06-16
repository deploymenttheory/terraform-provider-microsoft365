package graphBetaWindowsAutopilotDeploymentProfileAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructWindowsAutopilotDeploymentProfileAssignment constructs a Windows Autopilot Deployment Profile Assignment object for API requests
func ConstructWindowsAutopilotDeploymentProfileAssignment(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	data WindowsAutopilotDeploymentProfileAssignmentResourceModel,
	isUpdate bool,
) (graphmodels.WindowsAutopilotDeploymentProfileAssignmentable, error) {

	if err := ValidateWindowsAutopilotDeploymentProfileAssignment(ctx, client, data, isUpdate); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	assignment := graphmodels.NewWindowsAutopilotDeploymentProfileAssignment()

	if err := constructors.SetEnumProperty(data.Source, graphmodels.ParseDeviceAndAppManagementAssignmentSource, assignment.SetSource); err != nil {
		return nil, fmt.Errorf("error setting source: %v", err)
	}

	constructors.SetStringProperty(data.SourceId, assignment.SetSourceId)

	target, err := constructTarget(data.Target)
	if err != nil {
		return nil, fmt.Errorf("error constructing target: %v", err)
	}
	assignment.SetTarget(target)

	return assignment, nil
}

// constructTarget creates an assignment target based on the target type
func constructTarget(targetData AssignmentTargetResourceModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	targetType := targetData.TargetType.ValueString()

	switch targetType {
	case "allDevices":
		target := graphmodels.NewAllDevicesAssignmentTarget()
		setCommonTargetProperties(target, targetData)
		return target, nil

	case "groupAssignment":
		target := graphmodels.NewGroupAssignmentTarget()
		if !targetData.GroupId.IsNull() && !targetData.GroupId.IsUnknown() {
			groupId := targetData.GroupId.ValueString()
			target.SetGroupId(&groupId)
		}
		setCommonTargetProperties(target, targetData)
		return target, nil

	case "exclusionGroupAssignment":
		target := graphmodels.NewExclusionGroupAssignmentTarget()
		if !targetData.GroupId.IsNull() && !targetData.GroupId.IsUnknown() {
			groupId := targetData.GroupId.ValueString()
			target.SetGroupId(&groupId)
		}
		setCommonTargetProperties(target, targetData)
		return target, nil

	default:
		return nil, fmt.Errorf("unsupported target type: %s", targetType)
	}
}

// setCommonTargetProperties sets properties common to all target types
func setCommonTargetProperties(target graphmodels.DeviceAndAppManagementAssignmentTargetable, targetData AssignmentTargetResourceModel) {
	if !targetData.DeviceAndAppManagementAssignmentFilterId.IsNull() && !targetData.DeviceAndAppManagementAssignmentFilterId.IsUnknown() {
		filterId := targetData.DeviceAndAppManagementAssignmentFilterId.ValueString()
		target.SetDeviceAndAppManagementAssignmentFilterId(&filterId)
	}

	if !targetData.DeviceAndAppManagementAssignmentFilterType.IsNull() && !targetData.DeviceAndAppManagementAssignmentFilterType.IsUnknown() {
		filterType := targetData.DeviceAndAppManagementAssignmentFilterType.ValueString()
		if parsedFilterType, err := graphmodels.ParseDeviceAndAppManagementAssignmentFilterType(filterType); err == nil {
			target.SetDeviceAndAppManagementAssignmentFilterType(parsedFilterType.(*graphmodels.DeviceAndAppManagementAssignmentFilterType))
		}
	}
}
