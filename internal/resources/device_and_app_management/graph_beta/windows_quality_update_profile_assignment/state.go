package graphBetaWindowsQualityUpdateProfileAssignment

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps properties from multiple assignments to Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsQualityUpdateProfileAssignmentResourceModel, assignments []graphmodels.WindowsQualityUpdateProfileAssignmentable) {
	if assignments == nil || len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found")
		return
	}

	tflog.Debug(ctx, "Mapping assignments to Terraform state", map[string]interface{}{
		"assignmentCount": len(assignments),
	})

	if data.ID.IsNull() || data.ID.IsUnknown() || data.ID.ValueString() == "" {
		if len(assignments) > 0 && assignments[0] != nil && assignments[0].GetId() != nil {
			data.ID = types.StringPointerValue(assignments[0].GetId())
		}
	}

	includeGroupIDs := make(map[string]bool)
	excludeGroupIDs := make(map[string]bool)

	for _, assignment := range assignments {
		if assignment == nil || assignment.GetTarget() == nil {
			continue
		}

		target := assignment.GetTarget()

		if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok && groupTarget != nil && groupTarget.GetGroupId() != nil {
			groupID := *groupTarget.GetGroupId()

			if odataType := groupTarget.GetOdataType(); odataType != nil {
				if strings.Contains(*odataType, "exclusion") || strings.Contains(*odataType, "ExclusionGroupAssignmentTarget") {
					excludeGroupIDs[groupID] = true
					tflog.Debug(ctx, fmt.Sprintf("Found exclusion group ID: %s", groupID))
				} else {
					includeGroupIDs[groupID] = true
					tflog.Debug(ctx, fmt.Sprintf("Found inclusion group ID: %s", groupID))
				}
			}
		}
	}

	var inclusionGroupIDList []string
	for id := range includeGroupIDs {
		inclusionGroupIDList = append(inclusionGroupIDList, id)
	}

	var exclusionGroupIDList []string
	for id := range excludeGroupIDs {
		exclusionGroupIDList = append(exclusionGroupIDList, id)
	}

	var assignmentBlocks []AssignmentResourceModel

	if len(inclusionGroupIDList) > 0 {
		inclusionSet, diags := types.SetValueFrom(ctx, types.StringType, inclusionGroupIDList)
		if !diags.HasError() {
			assignmentBlocks = append(assignmentBlocks, AssignmentResourceModel{
				Target:   types.StringValue("include"),
				GroupIds: inclusionSet,
			})
			tflog.Debug(ctx, fmt.Sprintf("Added inclusion assignment block with %d groups", len(inclusionGroupIDList)))
		} else {
			tflog.Error(ctx, "Failed to create set for inclusion group IDs", map[string]interface{}{
				"error": diags.Errors(),
			})
		}
	}

	if len(exclusionGroupIDList) > 0 {
		exclusionSet, diags := types.SetValueFrom(ctx, types.StringType, exclusionGroupIDList)
		if !diags.HasError() {
			assignmentBlocks = append(assignmentBlocks, AssignmentResourceModel{
				Target:   types.StringValue("exclude"),
				GroupIds: exclusionSet,
			})
			tflog.Debug(ctx, fmt.Sprintf("Added exclusion assignment block with %d groups", len(exclusionGroupIDList)))
		} else {
			tflog.Error(ctx, "Failed to create set for exclusion group IDs", map[string]interface{}{
				"error": diags.Errors(),
			})
		}
	}

	data.Assignments = assignmentBlocks

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]interface{}{
		"assignmentBlockCount": len(assignmentBlocks),
	})
}
