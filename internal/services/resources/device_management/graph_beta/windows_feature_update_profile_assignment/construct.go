package graphBetaWindowsFeatureUpdateProfileAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a WindowsFeatureUpdateProfileAssignment
func constructResource(ctx context.Context, data WindowsFeatureUpdateProfileAssignmentResourceModel) (graphmodels.WindowsFeatureUpdateProfileAssignmentable, error) {
	tflog.Debug(ctx, "Starting windows feature update profile assignment construction")

	assignment := graphmodels.NewWindowsFeatureUpdateProfileAssignment()

	// Set Target
	target, err := constructAssignmentTarget(ctx, &data.Target)
	if err != nil {
		return nil, fmt.Errorf("error constructing windows feature update profile assignment target: %v", err)
	}
	assignment.SetTarget(target)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed windows feature update profile assignment", assignment); err != nil {
		tflog.Error(ctx, "Failed to log windows feature update profile assignment", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return assignment, nil
}

// constructAssignmentTarget constructs the windows feature update profile assignment target
func constructAssignmentTarget(ctx context.Context, data *AssignmentTargetResourceModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	if data == nil {
		return nil, fmt.Errorf("assignment target data is required")
	}

	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch data.TargetType.ValueString() {
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

	tflog.Debug(ctx, "Finished constructing assignment target")
	return target, nil
}
