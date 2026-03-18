package graphBetaWindowsUpdatesAutopatchRing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	commonstate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

// assignedGroupAttrTypes defines the attribute types for the AssignedGroupModel set elements.
var assignedGroupAttrTypes = map[string]attr.Type{
	"group_id": types.StringType,
}

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchRingResourceModel, remoteResource graphmodelswindowsupdates.Ringable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.IsPaused = convert.GraphToFrameworkBool(remoteResource.GetIsPaused())
	data.DeferralInDays = convert.GraphToFrameworkInt32(remoteResource.GetDeferralInDays())

	// Map isHotpatchEnabled from QualityUpdateRing subtype
	if qur, ok := remoteResource.(graphmodelswindowsupdates.QualityUpdateRingable); ok {
		data.IsHotpatchEnabled = convert.GraphToFrameworkBool(qur.GetIsHotpatchEnabled())
	}

	data.IncludedGroupAssignment = mapGroupAssignment(ctx, remoteResource.GetIncludedGroupAssignment())
	data.ExcludedGroupAssignment = mapGroupAssignment(ctx, remoteResource.GetExcludedGroupAssignment())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}

// mapGroupAssignment converts the SDK group assignment to a Terraform GroupAssignmentModel.
// Returns a model with an empty set if the assignment is nil or has no assignments.
func mapGroupAssignment(ctx context.Context, assignment interface {
	GetAssignments() []graphmodelswindowsupdates.AssignedGroupable
}) *GroupAssignmentModel {
	if assignment == nil {
		set := commonstate.BuildObjectSetFromSlice(ctx, assignedGroupAttrTypes, nil, 0)
		return &GroupAssignmentModel{Assignments: set}
	}

	assignments := assignment.GetAssignments()
	set := commonstate.BuildObjectSetFromSlice(
		ctx,
		assignedGroupAttrTypes,
		func(i int) map[string]attr.Value {
			return map[string]attr.Value{
				"group_id": convert.GraphToFrameworkString(assignments[i].GetGroupId()),
			}
		},
		len(assignments),
	)
	return &GroupAssignmentModel{Assignments: set}
}
