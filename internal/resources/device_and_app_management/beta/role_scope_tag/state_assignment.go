package graphBetaRoleScopeTag

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteAssignmentStateToTerraform(ctx context.Context, terraform *RoleScopeTagResourceModel, assignmentsResponse graphmodels.RoleScopeTagAutoAssignmentCollectionResponseable) {
	if assignmentsResponse == nil {
		terraform.Assignments = nil
		return
	}

	assignments := assignmentsResponse.GetValue()
	if assignments == nil {
		terraform.Assignments = nil
		return
	}

	var terraformAssignments []AssignmentModel
	for _, assignment := range assignments {
		target := assignment.GetTarget()
		if target == nil {
			continue
		}

		// Check if it's a group assignment target
		if target.GetOdataType() != nil && *target.GetOdataType() == "#microsoft.graph.groupAssignmentTarget" {
			groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable)
			if !ok {
				tflog.Debug(ctx, "Failed to cast target to GroupAssignmentTargetable")
				continue
			}

			if groupTarget.GetGroupId() != nil {
				terraformAssignments = append(terraformAssignments, AssignmentModel{
					GroupID: types.StringValue(*groupTarget.GetGroupId()),
				})
			}
		}
	}

	// Sort assignments by group ID to ensure consistent ordering
	sort.Slice(terraformAssignments, func(i, j int) bool {
		return terraformAssignments[i].GroupID.ValueString() < terraformAssignments[j].GroupID.ValueString()
	})

	terraform.Assignments = terraformAssignments
}
