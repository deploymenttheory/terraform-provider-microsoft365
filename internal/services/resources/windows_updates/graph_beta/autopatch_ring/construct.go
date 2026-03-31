package graphBetaWindowsUpdatesAutopatchRing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, data *WindowsUpdatesAutopatchRingResourceModel) (graphmodelswindowsupdates.QualityUpdateRingable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewQualityUpdateRing()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.IsPaused, requestBody.SetIsPaused)

	convert.FrameworkToGraphInt32(data.DeferralInDays, requestBody.SetDeferralInDays)

	if !data.IsHotpatchEnabled.IsNull() && !data.IsHotpatchEnabled.IsUnknown() {
		convert.FrameworkToGraphBool(data.IsHotpatchEnabled, requestBody.SetIsHotpatchEnabled)
	}

	requestBody.SetIncludedGroupAssignment(buildIncludedGroupAssignment(ctx, data.IncludedGroupAssignment))
	requestBody.SetExcludedGroupAssignment(buildExcludedGroupAssignment(ctx, data.ExcludedGroupAssignment))

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}

// constructUpdateResource builds a PATCH body with all mutable fields.
func constructUpdateResource(ctx context.Context, data *WindowsUpdatesAutopatchRingResourceModel) (graphmodelswindowsupdates.QualityUpdateRingable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing update request for %s resource", ResourceName))

	requestBody := graphmodelswindowsupdates.NewQualityUpdateRing()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.IsPaused, requestBody.SetIsPaused)

	convert.FrameworkToGraphInt32(data.DeferralInDays, requestBody.SetDeferralInDays)

	if !data.IsHotpatchEnabled.IsNull() && !data.IsHotpatchEnabled.IsUnknown() {
		convert.FrameworkToGraphBool(data.IsHotpatchEnabled, requestBody.SetIsHotpatchEnabled)
	}

	requestBody.SetIncludedGroupAssignment(buildIncludedGroupAssignment(ctx, data.IncludedGroupAssignment))
	requestBody.SetExcludedGroupAssignment(buildExcludedGroupAssignment(ctx, data.ExcludedGroupAssignment))

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing deployment settings for %s resource", ResourceName))
	return requestBody, nil
}

// buildIncludedGroupAssignment converts a types.Object to an SDK IncludedGroupAssignment.
// If the object is null/unknown, returns an empty assignment (sends [] to API).
func buildIncludedGroupAssignment(ctx context.Context, obj types.Object) graphmodelswindowsupdates.IncludedGroupAssignmentable {
	included := graphmodelswindowsupdates.NewIncludedGroupAssignment()
	included.SetAssignments(extractAssignedGroupsFromObject(ctx, obj))
	return included
}

// buildExcludedGroupAssignment converts a types.Object to an SDK ExcludedGroupAssignment.
// If the object is null/unknown, returns an empty assignment (sends [] to API).
func buildExcludedGroupAssignment(ctx context.Context, obj types.Object) graphmodelswindowsupdates.ExcludedGroupAssignmentable {
	excluded := graphmodelswindowsupdates.NewExcludedGroupAssignment()
	excluded.SetAssignments(extractAssignedGroupsFromObject(ctx, obj))
	return excluded
}

// extractAssignedGroupsFromObject extracts a slice of AssignedGroupable from a types.Object.
// Returns an empty slice if the object is null/unknown or has no assignments.
func extractAssignedGroupsFromObject(ctx context.Context, obj types.Object) []graphmodelswindowsupdates.AssignedGroupable {
	if obj.IsNull() || obj.IsUnknown() {
		return []graphmodelswindowsupdates.AssignedGroupable{}
	}

	var model GroupAssignmentModel
	diags := obj.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return []graphmodelswindowsupdates.AssignedGroupable{}
	}

	if model.Assignments.IsNull() || model.Assignments.IsUnknown() {
		return []graphmodelswindowsupdates.AssignedGroupable{}
	}

	var agModels []AssignedGroupModel
	diags = model.Assignments.ElementsAs(ctx, &agModels, false)
	if diags.HasError() {
		return []graphmodelswindowsupdates.AssignedGroupable{}
	}

	assignments := make([]graphmodelswindowsupdates.AssignedGroupable, len(agModels))
	for i, ag := range agModels {
		group := graphmodelswindowsupdates.NewAssignedGroup()
		groupId := ag.GroupId.ValueString()
		group.SetGroupId(&groupId)
		assignments[i] = group
	}
	return assignments
}
