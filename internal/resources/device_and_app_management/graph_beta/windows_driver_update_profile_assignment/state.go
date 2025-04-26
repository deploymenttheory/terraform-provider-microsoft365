package graphBetaWindowsDriverUpdateProfileAssignment

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps properties from multiple assignments to Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsDriverUpdateProfileAssignmentResourceModel, assignments []graphmodels.WindowsDriverUpdateProfileAssignmentable) {
	if assignments == nil || len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found")
		return
	}

	tflog.Debug(ctx, "Mapping assignments to Terraform state", map[string]interface{}{
		"assignmentCount": len(assignments),
	})

	// If we don't have an ID yet, use the first assignment's ID
	if data.ID.IsNull() || data.ID.IsUnknown() || data.ID.ValueString() == "" {
		if len(assignments) > 0 && assignments[0] != nil && assignments[0].GetId() != nil {
			data.ID = types.StringPointerValue(assignments[0].GetId())
		}
	}

	// Group assignments by type (include or exclude)
	includeGroupIDs := make(map[string]bool)
	excludeGroupIDs := make(map[string]bool)

	// Process all assignments to categorize them
	for _, assignment := range assignments {
		if assignment == nil || assignment.GetTarget() == nil {
			continue
		}

		target := assignment.GetTarget()

		// Check for GroupAssignmentTarget (either regular or exclusion)
		if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok && groupTarget != nil && groupTarget.GetGroupId() != nil {
			groupID := *groupTarget.GetGroupId()

			// Get the odata.type directly from the target using GetOdataType()
			if odataType := groupTarget.GetOdataType(); odataType != nil {
				if strings.Contains(*odataType, "exclusion") || strings.Contains(*odataType, "ExclusionGroupAssignmentTarget") {
					// This is an exclusion group
					excludeGroupIDs[groupID] = true
					tflog.Debug(ctx, fmt.Sprintf("Found exclusion group ID: %s", groupID))
				} else {
					// This is a regular inclusion group
					includeGroupIDs[groupID] = true
					tflog.Debug(ctx, fmt.Sprintf("Found inclusion group ID: %s", groupID))
				}
			} else {
				// Fallback - if no odata type, assume it's an include
				includeGroupIDs[groupID] = true
				tflog.Debug(ctx, fmt.Sprintf("Found inclusion group ID (no odata type): %s", groupID))
			}
		}
	}

	// Convert maps to slices for inclusion and exclusion groups
	var inclusionGroupIDList []string
	for id := range includeGroupIDs {
		inclusionGroupIDList = append(inclusionGroupIDList, id)
	}

	var exclusionGroupIDList []string
	for id := range excludeGroupIDs {
		exclusionGroupIDList = append(exclusionGroupIDList, id)
	}

	// Create assignment blocks for state
	var assignmentBlocks []AssignmentResourceModel

	// Add inclusion block if there are any inclusion groups
	if len(inclusionGroupIDList) > 0 {
		// Create a set from the inclusion group IDs
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

	// Add exclusion block if there are any exclusion groups
	if len(exclusionGroupIDList) > 0 {
		// Create a set from the exclusion group IDs
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

	// Set the assignment blocks in the state
	data.Assignments = assignmentBlocks

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]interface{}{
		"assignmentBlockCount": len(assignmentBlocks),
	})
}
