package graphBetaWindowsQualityUpdateProfileAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a WindowsQualityUpdateProfileAssignment
func constructResource(ctx context.Context, data WindowsQualityUpdateProfileAssignmentResourceModel) (graphmodels.WindowsQualityUpdateProfileAssignmentable, error) {
	tflog.Debug(ctx, "Starting windows quality update profile assignment construction")

	assignment := graphmodels.NewWindowsQualityUpdateProfileAssignment()

	// Set Target
	target, err := constructAssignmentTarget(ctx, &data.Target)
	if err != nil {
		return nil, fmt.Errorf("error constructing windows quality update profile assignment target: %v", err)
	}
	assignment.SetTarget(target)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed windows quality update profile assignment", assignment); err != nil {
		tflog.Error(ctx, "Failed to log windows quality update profile assignment", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return assignment, nil
}

// constructAssignmentTarget constructs the windows quality update profile assignment target
func constructAssignmentTarget(ctx context.Context, data *AssignmentTargetResourceModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	if data == nil {
		return nil, fmt.Errorf("assignment target data is required")
	}

	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch data.TargetType.ValueString() {
	case "configurationManagerCollection":
		configManagerTarget := graphmodels.NewConfigurationManagerCollectionAssignmentTarget()
		constructors.SetStringProperty(data.CollectionId, configManagerTarget.SetCollectionId)
		target = configManagerTarget
	case "exclusionGroupAssignment":
		exclusionGroupTarget := graphmodels.NewExclusionGroupAssignmentTarget()
		constructors.SetStringProperty(data.GroupId, exclusionGroupTarget.SetGroupId)
		target = exclusionGroupTarget
	case "groupAssignment":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		constructors.SetStringProperty(data.GroupId, groupTarget.SetGroupId)
		target = groupTarget
	default:
		target = graphmodels.NewDeviceAndAppManagementAssignmentTarget()
	}

	tflog.Debug(ctx, "Finished constructing assignment target")
	return target, nil
}
