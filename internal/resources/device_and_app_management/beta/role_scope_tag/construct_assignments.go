package graphBetaRoleScopeTag

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructAssignment(ctx context.Context, data *RoleScopeTagResourceModel) (devicemanagement.RoleScopeTagsItemAssignPostRequestBodyable, error) {
	if data.Assignments == nil {
		return nil, fmt.Errorf("assignments configuration block is required even if empty")
	}

	tflog.Debug(ctx, "Starting assignment construction")

	requestBody := devicemanagement.NewRoleScopeTagsItemAssignPostRequestBody()
	var assignments []graphmodels.RoleScopeTagAutoAssignmentable

	for _, assignment := range data.Assignments {

		roleScopeTagAutoAssignment := graphmodels.NewRoleScopeTagAutoAssignment()

		target := graphmodels.NewGroupAssignmentTarget()

		odataType := "#microsoft.graph.groupAssignmentTarget"
		target.SetOdataType(&odataType)

		groupID := assignment.GroupID.ValueString()
		target.SetGroupId(&groupID)

		roleScopeTagAutoAssignment.SetTarget(target)

		assignments = append(assignments, roleScopeTagAutoAssignment)
	}

	requestBody.SetAssignments(assignments)

	if err := construct.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
