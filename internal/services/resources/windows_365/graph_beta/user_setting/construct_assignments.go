package graphBetaCloudPcUserSetting

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignmentsRequestBody creates the request body for assigning a Cloud PC user setting to groups
func constructAssignmentsRequestBody(ctx context.Context, assignments []CloudPcUserSettingAssignmentModel) (*devicemanagement.VirtualEndpointUserSettingsItemAssignPostRequestBody, error) {
	tflog.Debug(ctx, "Constructing assignments request body")

	requestBody := devicemanagement.NewVirtualEndpointUserSettingsItemAssignPostRequestBody()

	// If no assignments are provided, return an empty assignments array
	// This is used for removing all assignments
	if len(assignments) == 0 {
		requestBody.SetAssignments([]models.CloudPcUserSettingAssignmentable{})
		tflog.Debug(ctx, "No assignments provided, setting empty assignments array")
		return requestBody, nil
	}

	graphAssignments := []models.CloudPcUserSettingAssignmentable{}

	for i, assignment := range assignments {
		if assignment.GroupId.IsNull() || assignment.GroupId.ValueString() == "" {
			tflog.Warn(ctx, fmt.Sprintf("Skipping assignment %d with empty group ID", i))
			continue
		}

		tflog.Debug(ctx, fmt.Sprintf("Creating assignment %d for group ID: %s", i, assignment.GroupId.ValueString()))

		graphAssignment := models.NewCloudPcUserSettingAssignment()
		if !assignment.ID.IsNull() && assignment.ID.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.ID, graphAssignment.SetId)
		}

		// Create the target with the proper @odata.type
		target := models.NewCloudPcManagementAssignmentTarget()

		// Always set the @odata.type for consistency
		odataType := "#microsoft.graph.cloudPcManagementGroupAssignmentTarget"
		target.SetOdataType(&odataType)

		// Set up additional data with the groupId
		additionalData := target.GetAdditionalData()

		// The groupId must be set directly in the additionalData map
		additionalData["groupId"] = assignment.GroupId.ValueString()

		tflog.Debug(ctx, fmt.Sprintf("Setting groupId in additionalData: %s", assignment.GroupId.ValueString()))

		graphAssignment.SetTarget(target)
		graphAssignments = append(graphAssignments, graphAssignment)

		tflog.Debug(ctx, fmt.Sprintf("Successfully created assignment %d", i))
	}

	requestBody.SetAssignments(graphAssignments)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s assignments", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructed assignments request body with %d assignments", len(graphAssignments)))
	return requestBody, nil
}
