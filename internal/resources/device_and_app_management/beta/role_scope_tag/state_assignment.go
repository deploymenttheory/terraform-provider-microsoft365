package graphBetaRoleScopeTag

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps the assignment remote state to the Terraform model
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

	var groupIDs []types.String
	for _, assignment := range assignments {
		target := assignment.GetTarget()
		if target == nil {
			continue
		}

		if target.GetOdataType() != nil && *target.GetOdataType() == "#microsoft.graph.groupAssignmentTarget" {
			groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable)
			if !ok {
				tflog.Debug(ctx, "Failed to cast target to GroupAssignmentTargetable")
				continue
			}

			if groupTarget.GetGroupId() != nil {
				groupIDs = append(groupIDs, types.StringValue(*groupTarget.GetGroupId()))
			}
		}
	}

	sort.Slice(groupIDs, func(i, j int) bool {
		return groupIDs[i].ValueString() < groupIDs[j].ValueString()
	})

	terraform.Assignments = groupIDs
}
