package graphBetaTargetedManagedAppConfigurations

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignments constructs assignments for the mobile app configuration
// this is the same for both create and update requests
func constructAssignments(ctx context.Context, assignmentsSet types.Set) ([]graphmodels.TargetedManagedAppPolicyAssignmentable, error) {
	if assignmentsSet.IsNull() || assignmentsSet.IsUnknown() {
		// Return empty slice instead of nil to ensure the field is present in JSON
		return []graphmodels.TargetedManagedAppPolicyAssignmentable{}, nil
	}

	var assignmentModels []sharedmodels.InclusionGroupAndExclusionGroupAssignmentModel
	diags := assignmentsSet.ElementsAs(ctx, &assignmentModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert assignments set: %v", diags)
	}

	assignments := make([]graphmodels.TargetedManagedAppPolicyAssignmentable, len(assignmentModels))
	for i, assignmentModel := range assignmentModels {
		assignment := graphmodels.NewTargetedManagedAppPolicyAssignment()

		target, err := constructAssignmentTarget(ctx, &assignmentModel)
		if err != nil {
			return nil, fmt.Errorf("failed to construct assignment target: %s", err)
		}
		assignment.SetTarget(target)

		assignments[i] = assignment
	}

	return assignments, nil
}

// constructAssignmentTarget constructs an assignment target
func constructAssignmentTarget(_ context.Context, assignmentModel *sharedmodels.InclusionGroupAndExclusionGroupAssignmentModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	assignmentType := assignmentModel.Type.ValueString()

	switch assignmentType {
	case "groupAssignmentTarget":
		target := graphmodels.NewGroupAssignmentTarget()
		groupOdataType := "#microsoft.graph.groupAssignmentTarget"
		target.SetOdataType(&groupOdataType)
		if !assignmentModel.GroupId.IsNull() && !assignmentModel.GroupId.IsUnknown() {
			target.SetGroupId(assignmentModel.GroupId.ValueStringPointer())
		}
		return target, nil
	case "exclusionGroupAssignmentTarget":
		target := graphmodels.NewExclusionGroupAssignmentTarget()
		exclusionOdataType := "#microsoft.graph.exclusionGroupAssignmentTarget"
		target.SetOdataType(&exclusionOdataType)
		if !assignmentModel.GroupId.IsNull() && !assignmentModel.GroupId.IsUnknown() {
			target.SetGroupId(assignmentModel.GroupId.ValueStringPointer())
		}
		return target, nil
	default:
		return nil, fmt.Errorf("unsupported assignment target type: %s. Valid types are: 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'", assignmentType)
	}
}
