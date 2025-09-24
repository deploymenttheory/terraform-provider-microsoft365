package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignments creates assignment objects for the Windows Autopilot deployment profile
func constructAssignments(ctx context.Context, assignments types.Set) ([]graphmodels.WindowsAutopilotDeploymentProfileAssignmentable, error) {
	if assignments.IsNull() || assignments.IsUnknown() {
		return nil, nil
	}

	var assignmentModels []AssignmentModel
	diags := assignments.ElementsAs(ctx, &assignmentModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert assignments: %v", diags)
	}

	var graphAssignments []graphmodels.WindowsAutopilotDeploymentProfileAssignmentable

	for _, assignment := range assignmentModels {
		graphAssignment := graphmodels.NewWindowsAutopilotDeploymentProfileAssignment()

		// Create the appropriate target based on assignment type
		switch assignment.Type.ValueString() {
		case "groupAssignmentTarget":
			if assignment.GroupID.IsNull() || assignment.GroupID.IsUnknown() {
				return nil, fmt.Errorf("group_id is required for groupAssignmentTarget")
			}
			target := graphmodels.NewGroupAssignmentTarget()
			groupId := assignment.GroupID.ValueString()
			target.SetGroupId(&groupId)
			graphAssignment.SetTarget(target)

		case "exclusionGroupAssignmentTarget":
			if assignment.GroupID.IsNull() || assignment.GroupID.IsUnknown() {
				return nil, fmt.Errorf("group_id is required for exclusionGroupAssignmentTarget")
			}
			target := graphmodels.NewExclusionGroupAssignmentTarget()
			groupId := assignment.GroupID.ValueString()
			target.SetGroupId(&groupId)
			graphAssignment.SetTarget(target)

		case "allDevicesAssignmentTarget":
			target := graphmodels.NewAllDevicesAssignmentTarget()
			graphAssignment.SetTarget(target)

		default:
			return nil, fmt.Errorf("unsupported assignment type: %s", assignment.Type.ValueString())
		}

		graphAssignments = append(graphAssignments, graphAssignment)
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructed %d assignments for Windows Autopilot deployment profile", len(graphAssignments)))
	return graphAssignments, nil
}
