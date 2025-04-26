package graphBetaWindowsDriverUpdateProfileAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource creates an assign request body with assignments from the nested blocks
func constructResource(ctx context.Context, data *WindowsDriverUpdateProfileAssignmentResourceModel) (devicemanagement.WindowsDriverUpdateProfilesItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Creating assign request body from assignment blocks")

	// Create the assign request body
	assignRequest := devicemanagement.NewWindowsDriverUpdateProfilesItemAssignPostRequestBody()

	// Prepare a slice to hold all assignments
	var assignments []graphmodels.WindowsDriverUpdateProfileAssignmentable

	// Process each assignment block
	for i, assignmentBlock := range data.Assignments {
		if assignmentBlock.Target.IsNull() || assignmentBlock.Target.IsUnknown() {
			return nil, fmt.Errorf("assignment[%d]: target is required", i)
		}

		if assignmentBlock.GroupIds.IsNull() || assignmentBlock.GroupIds.IsUnknown() {
			return nil, fmt.Errorf("assignment[%d]: group_ids is required", i)
		}

		// Get the target type (include or exclude)
		targetType := assignmentBlock.Target.ValueString()

		// Get the group IDs from the set
		var groupIDs []string
		diags := assignmentBlock.GroupIds.ElementsAs(ctx, &groupIDs, false)
		if diags.HasError() {
			return nil, fmt.Errorf("assignment[%d]: error extracting group IDs: %s", i, diags.Errors())
		}

		// Skip if there are no group IDs
		if len(groupIDs) == 0 {
			continue
		}

		// Create assignments for each group ID
		for _, groupID := range groupIDs {
			// Create a new assignment
			assignment := graphmodels.NewWindowsDriverUpdateProfileAssignment()

			// Create the appropriate target based on the target type
			if targetType == "include" {
				// Regular group assignment target
				target := graphmodels.NewGroupAssignmentTarget()
				target.SetGroupId(&groupID)
				assignment.SetTarget(target)
				tflog.Debug(ctx, fmt.Sprintf("Added inclusion group assignment for group: %s", groupID))
			} else if targetType == "exclude" {
				// Exclusion group assignment target
				target := graphmodels.NewExclusionGroupAssignmentTarget()
				target.SetGroupId(&groupID)
				assignment.SetTarget(target)
				tflog.Debug(ctx, fmt.Sprintf("Added exclusion group assignment for group: %s", groupID))
			} else {
				return nil, fmt.Errorf("assignment[%d]: invalid target type: %s", i, targetType)
			}

			// Add the assignment to the list
			assignments = append(assignments, assignment)
		}
	}

	// Make sure we have at least one assignment
	if len(assignments) == 0 {
		return nil, fmt.Errorf("at least one assignment with group_ids is required")
	}

	// Set the assignments in the request
	assignRequest.SetAssignments(assignments)

	// Debug log the final request
	if err := constructors.DebugLogGraphObject(ctx, "Final assign request", assignRequest); err != nil {
		tflog.Error(ctx, "Failed to debug log assign request", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished creating assign request body with %d assignments", len(assignments)))
	return assignRequest, nil
}
