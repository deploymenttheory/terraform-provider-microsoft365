package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignmentsRequestBody creates the request body for assigning a Cloud PC provisioning policy to groups
func constructAssignmentsRequestBody(ctx context.Context, assignments []CloudPcProvisioningPolicyAssignmentModel) (*devicemanagement.VirtualEndpointProvisioningPoliciesItemAssignPostRequestBody, error) {
	tflog.Debug(ctx, "Constructing assignments request body")

	requestBody := devicemanagement.NewVirtualEndpointProvisioningPoliciesItemAssignPostRequestBody()

	// If no assignments are provided, return an empty assignments array
	// This is used for removing all assignments
	if len(assignments) == 0 {
		requestBody.SetAssignments([]models.CloudPcProvisioningPolicyAssignmentable{})
		tflog.Debug(ctx, "No assignments provided, setting empty assignments array")
		return requestBody, nil
	}

	graphAssignments := []models.CloudPcProvisioningPolicyAssignmentable{}

	for _, assignment := range assignments {
		if assignment.GroupId.IsNull() || assignment.GroupId.ValueString() == "" {
			tflog.Warn(ctx, "Skipping assignment with empty group ID")
			continue
		}

		graphAssignment := models.NewCloudPcProvisioningPolicyAssignment()
		if !assignment.ID.IsNull() && assignment.ID.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.ID, graphAssignment.SetId)
		}

		target := models.NewCloudPcManagementAssignmentTarget()
		additionalData := target.GetAdditionalData()
		additionalData["groupId"] = assignment.GroupId.ValueString()

		if !assignment.ServicePlanId.IsNull() && assignment.ServicePlanId.ValueString() != "" {
			// Frontline (dedicated/shared)
			odataType := "#microsoft.graph.cloudPcManagementServicePlanAssignmentTarget"
			target.SetOdataType(&odataType)
			additionalData["servicePlanId"] = assignment.ServicePlanId.ValueString()
			if !assignment.AllotmentDisplayName.IsNull() && assignment.AllotmentDisplayName.ValueString() != "" {
				additionalData["allotmentDisplayName"] = assignment.AllotmentDisplayName.ValueString()
			} else {
				additionalData["allotmentDisplayName"] = "Default Allotment"
			}
			if !assignment.AllotmentLicenseCount.IsNull() {
				additionalData["allotmentLicensesCount"] = int32(assignment.AllotmentLicenseCount.ValueInt64())
			} else {
				additionalData["allotmentLicensesCount"] = int32(1)
			}
		} else {
			// Dedicated
			odataType := "#microsoft.graph.cloudPcManagementGroupAssignmentTarget"
			target.SetOdataType(&odataType)
		}

		graphAssignment.SetTarget(target)
		graphAssignments = append(graphAssignments, graphAssignment)
	}

	requestBody.SetAssignments(graphAssignments)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructed assignments request body with %d assignments", len(graphAssignments)))
	return requestBody, nil
}
