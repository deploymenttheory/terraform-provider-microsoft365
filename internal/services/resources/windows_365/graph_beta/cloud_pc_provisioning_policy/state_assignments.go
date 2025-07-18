package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapAssignmentsSliceToTerraform maps a slice of assignments directly to the Terraform model
func MapAssignmentsSliceToTerraform(ctx context.Context, assignments []models.CloudPcProvisioningPolicyAssignmentable) []CloudPcProvisioningPolicyAssignmentModel {
	if assignments == nil {
		return []CloudPcProvisioningPolicyAssignmentModel{}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d assignments from Graph API to Terraform model", len(assignments)))

	result := make([]CloudPcProvisioningPolicyAssignmentModel, 0, len(assignments))

	for i, assignment := range assignments {
		if assignment == nil {
			tflog.Debug(ctx, fmt.Sprintf("Assignment at index %d is nil, skipping", i))
			continue
		}

		id := safeGetStringPtr(assignment.GetId())
		tflog.Debug(ctx, fmt.Sprintf("Processing assignment %d with ID: %s", i, id))

		tfAssignment := CloudPcProvisioningPolicyAssignmentModel{}

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

		// Extract servicePlanId (only present for frontline)
		if v, ok := additionalData["servicePlanId"].(string); ok && v != "" && v != "null" {
			tfAssignment.ServicePlanId = types.StringValue(v)
		} else {
			tfAssignment.ServicePlanId = types.StringNull()
		}

		// Extract allotmentDisplayName (only present for frontline)
		if v, ok := additionalData["allotmentDisplayName"].(string); ok && v != "" && v != "null" {
			tfAssignment.AllotmentDisplayName = types.StringValue(v)
		} else {
			tfAssignment.AllotmentDisplayName = types.StringNull()
		}

		// Extract allotmentLicensesCount (only present for frontline)
		if v, ok := additionalData["allotmentLicensesCount"]; ok && v != nil && v != "null" {
			switch val := v.(type) {
			case int32:
				tfAssignment.AllotmentLicenseCount = types.Int64Value(int64(val))
			case int64:
				tfAssignment.AllotmentLicenseCount = types.Int64Value(val)
			case float64:
				tfAssignment.AllotmentLicenseCount = types.Int64Value(int64(val))
			default:
				tfAssignment.AllotmentLicenseCount = types.Int64Null()
			}
		} else {
			tfAssignment.AllotmentLicenseCount = types.Int64Null()
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

// MapAssignmentsToTerraform maps the assignments from the Graph API response to the Terraform model
// used for a read operation.
func MapAssignmentsToTerraform(ctx context.Context, assignmentsResponse models.CloudPcProvisioningPolicyAssignmentCollectionResponseable) []CloudPcProvisioningPolicyAssignmentModel {
	if assignmentsResponse == nil {
		return []CloudPcProvisioningPolicyAssignmentModel{}
	}

	assignments := assignmentsResponse.GetValue()
	return MapAssignmentsSliceToTerraform(ctx, assignments)
}
