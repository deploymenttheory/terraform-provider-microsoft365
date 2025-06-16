package graphBetaRoleScopeTag

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs and returns a RoleScopeTagsItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *RoleScopeTagResourceModel) (devicemanagement.RoleScopeTagsItemAssignPostRequestBodyable, error) {
	if data.Assignments == nil {
		return nil, fmt.Errorf("assignments configuration block is required even if empty")
	}

	tflog.Debug(ctx, "Starting assignment construction")
	tflog.Debug(ctx, fmt.Sprintf("Number of assignments to process: %d", len(data.Assignments)))

	hasValidAssignments := false
	for _, groupID := range data.Assignments {
		if !groupID.IsNull() && !groupID.IsUnknown() && groupID.ValueString() != "" {
			hasValidAssignments = true
			break
		}
	}

	requestBody := devicemanagement.NewRoleScopeTagsItemAssignPostRequestBody()
	var assignments []graphmodels.RoleScopeTagAutoAssignmentable

	if hasValidAssignments {
		for _, groupID := range data.Assignments {
			if !groupID.IsNull() && !groupID.IsUnknown() && groupID.ValueString() != "" {
				tflog.Debug(ctx, fmt.Sprintf("Processing group ID: %s", groupID.ValueString()))
				assignment := constructGroupAssignment(groupID.ValueString())
				assignments = append(assignments, assignment)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Final number of assignments created: %d", len(assignments)))
	requestBody.SetAssignments(assignments)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructGroupAssignment constructs and returns a RoleScopeTagAutoAssignment for a group
func constructGroupAssignment(groupID string) graphmodels.RoleScopeTagAutoAssignmentable {
	assignment := graphmodels.NewRoleScopeTagAutoAssignment()
	target := graphmodels.NewGroupAssignmentTarget()
	target.SetGroupId(&groupID)
	assignment.SetTarget(target)
	return assignment
}
