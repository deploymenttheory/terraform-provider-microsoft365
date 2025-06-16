package graphBetaTermsAndConditionsAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a TermsAndConditionsAssignment
func constructResource(ctx context.Context, data TermsAndConditionsAssignmentResourceModel) (graphmodels.TermsAndConditionsAssignmentable, error) {
	tflog.Debug(ctx, "Starting terms and conditions assignment construction")

	assignment := graphmodels.NewTermsAndConditionsAssignment()

	// Set Target
	target, err := constructAssignmentTarget(ctx, &data.Target)
	if err != nil {
		return nil, fmt.Errorf("error constructing terms and conditions assignment target: %v", err)
	}
	assignment.SetTarget(target)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed terms and conditions assignment", assignment); err != nil {
		tflog.Error(ctx, "Failed to log terms and conditions assignment", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return assignment, nil
}

// constructAssignmentTarget constructs the terms and conditions assignment target
func constructAssignmentTarget(ctx context.Context, data *AssignmentTargetResourceModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	if data == nil {
		return nil, fmt.Errorf("assignment target data is required")
	}

	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch data.TargetType.ValueString() {
	case "allLicensedUsers":
		target = graphmodels.NewAllLicensedUsersAssignmentTarget()
	case "groupAssignment":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		constructors.SetStringProperty(data.GroupId, groupTarget.SetGroupId)
		target = groupTarget
	case "configurationManagerCollection":
		configManagerTarget := graphmodels.NewConfigurationManagerCollectionAssignmentTarget()
		constructors.SetStringProperty(data.CollectionId, configManagerTarget.SetCollectionId)
		target = configManagerTarget
	default:
		target = graphmodels.NewDeviceAndAppManagementAssignmentTarget()
	}

	tflog.Debug(ctx, "Finished constructing assignment target")
	return target, nil
}
