package graphBetaCloudPcUserSetting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapAssignmentsToTerraform maps the assignments from the Graph API response to the Terraform model
// used for a read operation.
func MapAssignmentsToTerraform(ctx context.Context, assignmentsResponse models.CloudPcUserSettingAssignmentCollectionResponseable) []CloudPcUserSettingAssignmentModel {
	if assignmentsResponse == nil {
		return []CloudPcUserSettingAssignmentModel{}
	}

	assignments := assignmentsResponse.GetValue()
	return MapAssignmentsSliceToTerraform(ctx, assignments)
}

// MapAssignmentsSliceToTerraform maps a slice of assignments directly to the Terraform model
func MapAssignmentsSliceToTerraform(ctx context.Context, assignments []models.CloudPcUserSettingAssignmentable) []CloudPcUserSettingAssignmentModel {
	if assignments == nil {
		return []CloudPcUserSettingAssignmentModel{}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d assignments from Graph API to Terraform model", len(assignments)))

	result := make([]CloudPcUserSettingAssignmentModel, 0, len(assignments))

	for i, assignment := range assignments {
		if assignment == nil {
			tflog.Debug(ctx, fmt.Sprintf("Assignment at index %d is nil, skipping", i))
			continue
		}

		id := safeGetStringPtr(assignment.GetId())
		tflog.Debug(ctx, fmt.Sprintf("Processing assignment %d with ID: %s", i, id))

		tfAssignment := CloudPcUserSettingAssignmentModel{}

		// Set ID
		if id != "<nil>" {
			tfAssignment.ID = types.StringValue(id)
		} else {
			tfAssignment.ID = types.StringNull()
		}

		// Process target data
		target := assignment.GetTarget()
		if target == nil {
			tflog.Debug(ctx, fmt.Sprintf("Target for assignment %d is nil", i))

			// In some API responses, the assignment ID itself is actually the group ID
			if id != "<nil>" {
				tflog.Debug(ctx, fmt.Sprintf("Using assignment ID as groupId: %s", id))
				tfAssignment.GroupId = types.StringValue(id)
			} else {
				tfAssignment.GroupId = types.StringNull()
			}

			result = append(result, tfAssignment)
			continue
		}

		// Debug log the target data
		additionalData := target.GetAdditionalData()
		tflog.Debug(ctx, "Target additional data", map[string]interface{}{
			"assignmentId": id,
			"data":         fmt.Sprintf("%v", additionalData),
			"odataType":    target.GetOdataType(),
		})

		// Based on the actual API response, the groupId is directly in the target's additionalData
		if groupId, ok := additionalData["groupId"].(string); ok && groupId != "" {
			tflog.Debug(ctx, fmt.Sprintf("Found groupId in target additionalData: %s", groupId))
			tfAssignment.GroupId = types.StringValue(groupId)
		} else {
			// Fallback to using the assignment ID as the group ID
			tflog.Debug(ctx, fmt.Sprintf("Fallback: Using assignment ID as groupId: %s", id))
			tfAssignment.GroupId = types.StringValue(id)
		}

		result = append(result, tfAssignment)
		tflog.Debug(ctx, fmt.Sprintf("Successfully mapped assignment %d", i), map[string]interface{}{
			"id":      tfAssignment.ID.ValueString(),
			"groupId": tfAssignment.GroupId.ValueString(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapped %d assignments to Terraform model", len(result)))
	return result
}

// Helper function to safely get string pointer value
func safeGetStringPtr(ptr *string) string {
	if ptr == nil {
		return "<nil>"
	}
	return *ptr
}
