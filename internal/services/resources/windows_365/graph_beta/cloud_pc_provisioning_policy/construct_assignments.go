package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignmentsRequestBody creates the request body for assigning a Cloud PC provisioning policy to groups
func constructAssignmentsRequestBody(ctx context.Context, assignments types.Set) (*devicemanagement.VirtualEndpointProvisioningPoliciesItemAssignPostRequestBody, error) {
	tflog.Debug(ctx, "Constructing assignments request body")

	requestBody := devicemanagement.NewVirtualEndpointProvisioningPoliciesItemAssignPostRequestBody()

	// If no assignments are provided, return an empty assignments array
	// This is used for removing all assignments
	if assignments.IsNull() || assignments.IsUnknown() {
		requestBody.SetAssignments([]models.CloudPcProvisioningPolicyAssignmentable{})
		tflog.Debug(ctx, "No assignments provided, setting empty assignments array")
		return requestBody, nil
	}

	var terraformAssignments []CloudPcProvisioningPolicyAssignmentModel
	diags := assignments.ElementsAs(ctx, &terraformAssignments, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract assignments: %v", diags.Errors())
	}

	graphAssignments := []models.CloudPcProvisioningPolicyAssignmentable{}

	for i, assignment := range terraformAssignments {
		if assignment.GroupId.IsNull() || assignment.GroupId.ValueString() == "" {
			tflog.Warn(ctx, fmt.Sprintf("Skipping assignment %d with empty group ID", i))
			continue
		}

		if assignment.Type.IsNull() || assignment.Type.IsUnknown() {
			tflog.Error(ctx, "Assignment target type is missing or invalid", map[string]interface{}{
				"index": i,
			})
			continue
		}

		targetType := assignment.Type.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Creating assignment %d for group ID: %s with target type: %s", i, assignment.GroupId.ValueString(), targetType))

		graphAssignment := models.NewCloudPcProvisioningPolicyAssignment()

		target := constructProvisioningPolicyTarget(ctx, targetType, assignment)
		if target == nil {
			tflog.Error(ctx, "Failed to create target", map[string]interface{}{
				"index":      i,
				"targetType": targetType,
			})
			continue
		}

		// Set up additional data with the groupId
		additionalData := target.GetAdditionalData()
		additionalData["groupId"] = assignment.GroupId.ValueString()

		tflog.Debug(ctx, fmt.Sprintf("Setting groupId in additionalData: %s", assignment.GroupId.ValueString()))

		// Handle Frontline-specific fields
		if !assignment.ServicePlanId.IsNull() && assignment.ServicePlanId.ValueString() != "" {
			// Frontline (dedicated/shared)
			tflog.Debug(ctx, fmt.Sprintf("Setting frontline-specific fields for assignment %d", i))
			additionalData["servicePlanId"] = assignment.ServicePlanId.ValueString()

			if !assignment.AllotmentDisplayName.IsNull() && assignment.AllotmentDisplayName.ValueString() != "" {
				additionalData["allotmentDisplayName"] = assignment.AllotmentDisplayName.ValueString()
			} else {
				additionalData["allotmentDisplayName"] = "Default Allotment"
			}

			if !assignment.AllotmentLicenseCount.IsNull() {
				additionalData["allotmentLicensesCount"] = assignment.AllotmentLicenseCount.ValueInt32()
			} else {
				additionalData["allotmentLicensesCount"] = int32(1)
			}
		}

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

// constructProvisioningPolicyTarget creates the appropriate target based on the target type
func constructProvisioningPolicyTarget(ctx context.Context, targetType string, assignment CloudPcProvisioningPolicyAssignmentModel) models.CloudPcManagementAssignmentTargetable {
	var target models.CloudPcManagementAssignmentTargetable
	var odataType string

	switch targetType {
	case "groupAssignmentTarget":
		target = models.NewCloudPcManagementAssignmentTarget()
		odataType = "#microsoft.graph.cloudPcManagementGroupAssignmentTarget"
		tflog.Debug(ctx, "Created CloudPcManagementGroupAssignmentTarget", map[string]interface{}{
			"groupId": assignment.GroupId.ValueString(),
		})
	case "exclusionGroupAssignmentTarget":
		target = models.NewCloudPcManagementAssignmentTarget()
		odataType = "#microsoft.graph.cloudPcManagementExclusionGroupAssignmentTarget"
		tflog.Debug(ctx, "Created CloudPcManagementExclusionGroupAssignmentTarget", map[string]interface{}{
			"groupId": assignment.GroupId.ValueString(),
		})
	default:
		tflog.Error(ctx, "Unsupported target type", map[string]interface{}{
			"targetType": targetType,
		})
		return nil
	}

	target.SetOdataType(&odataType)
	return target
}
