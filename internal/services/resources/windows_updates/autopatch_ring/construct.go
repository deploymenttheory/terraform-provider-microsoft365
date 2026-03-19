package graphBetaWindowsUpdatesAutopatchRing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
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

// buildIncludedGroupAssignment converts the Terraform model to an SDK IncludedGroupAssignment.
// If the model is nil or has no assignments, returns an empty assignment (sends [] to API).
func buildIncludedGroupAssignment(ctx context.Context, model *GroupAssignmentModel) graphmodelswindowsupdates.IncludedGroupAssignmentable {
	included := graphmodelswindowsupdates.NewIncludedGroupAssignment()
	included.SetAssignments(extractAssignedGroups(ctx, model))
	return included
}

// buildExcludedGroupAssignment converts the Terraform model to an SDK ExcludedGroupAssignment.
// If the model is nil or has no assignments, returns an empty assignment (sends [] to API).
func buildExcludedGroupAssignment(ctx context.Context, model *GroupAssignmentModel) graphmodelswindowsupdates.ExcludedGroupAssignmentable {
	excluded := graphmodelswindowsupdates.NewExcludedGroupAssignment()
	excluded.SetAssignments(extractAssignedGroups(ctx, model))
	return excluded
}

// extractAssignedGroups extracts a slice of AssignedGroupable from the Terraform model.
// Returns an empty slice if the model is nil or the set is null/unknown/empty.
func extractAssignedGroups(ctx context.Context, model *GroupAssignmentModel) []graphmodelswindowsupdates.AssignedGroupable {
	if model == nil || model.Assignments.IsNull() || model.Assignments.IsUnknown() {
		return []graphmodelswindowsupdates.AssignedGroupable{}
	}

	var agModels []AssignedGroupModel
	diags := model.Assignments.ElementsAs(ctx, &agModels, false)
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
