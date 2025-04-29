package graphBetaDeviceEnrollmentConfiguration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignments builds the /assign request body from the unified resource model
func constructAssignments(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel) (devicemanagement.DeviceEnrollmentConfigurationsItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Creating assign request body from assignment blocks")

	assignRequest := devicemanagement.NewDeviceEnrollmentConfigurationsItemAssignPostRequestBody()
	var assignments []graphmodels.DeviceEnrollmentConfigurationAssignmentable

	for i, assignmentBlock := range data.Assignments {
		if assignmentBlock.Target.IsNull() || assignmentBlock.Target.IsUnknown() {
			return nil, fmt.Errorf("assignment[%d]: target is required", i)
		}

		if assignmentBlock.GroupIds.IsNull() || assignmentBlock.GroupIds.IsUnknown() {
			return nil, fmt.Errorf("assignment[%d]: group_ids is required", i)
		}

		targetType := assignmentBlock.Target.ValueString()

		var groupIDs []string
		diags := assignmentBlock.GroupIds.ElementsAs(ctx, &groupIDs, false)
		if diags.HasError() {
			return nil, fmt.Errorf("assignment[%d]: error extracting group IDs: %s", i, diags.Errors())
		}

		if len(groupIDs) == 0 {
			continue
		}

		for _, groupID := range groupIDs {
			assignment := graphmodels.NewDeviceEnrollmentConfigurationAssignment()

			switch targetType {
			case "include":
				target := graphmodels.NewGroupAssignmentTarget()
				target.SetGroupId(&groupID)
				assignment.SetTarget(target)
				tflog.Debug(ctx, fmt.Sprintf("Added inclusion group assignment for group: %s", groupID))

			case "exclude":
				target := graphmodels.NewExclusionGroupAssignmentTarget()
				target.SetGroupId(&groupID)
				assignment.SetTarget(target)
				tflog.Debug(ctx, fmt.Sprintf("Added exclusion group assignment for group: %s", groupID))

			default:
				return nil, fmt.Errorf("assignment[%d]: invalid target type: %s", i, targetType)
			}

			assignments = append(assignments, assignment)
		}
	}

	if len(assignments) == 0 {
		return nil, fmt.Errorf("at least one assignment with group_ids is required")
	}

	assignRequest.SetAssignments(assignments)

	if err := constructors.DebugLogGraphObject(ctx, "Final assign request", assignRequest); err != nil {
		tflog.Error(ctx, "Failed to debug log assign request", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished creating assign request body with %d assignments", len(assignments)))
	return assignRequest, nil
}
