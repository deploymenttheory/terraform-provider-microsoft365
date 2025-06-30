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

	for _, assignment := range assignments {
		if assignment == nil {
			continue
		}

		tfAssignment := CloudPcProvisioningPolicyAssignmentModel{}
		id := assignment.GetId()
		if id != nil {
			tfAssignment.ID = types.StringValue(*id)
		} else {
			tfAssignment.ID = types.StringNull()
		}

		target := assignment.GetTarget()
		if target != nil {
			additionalData := target.GetAdditionalData()
			// groupId is always present
			if v, ok := additionalData["groupId"].(string); ok {
				tfAssignment.GroupId = types.StringValue(v)
			} else {
				tfAssignment.GroupId = types.StringNull()
			}
			// servicePlanId is only present for frontline
			if v, ok := additionalData["servicePlanId"].(string); ok && v != "" {
				tfAssignment.ServicePlanId = types.StringValue(v)
			} else {
				tfAssignment.ServicePlanId = types.StringNull()
			}
			// allotmentDisplayName is only present for frontline
			if v, ok := additionalData["allotmentDisplayName"].(string); ok && v != "" {
				tfAssignment.AllotmentDisplayName = types.StringValue(v)
			} else {
				tfAssignment.AllotmentDisplayName = types.StringNull()
			}
			// allotmentLicensesCount is only present for frontline
			if v, ok := additionalData["allotmentLicensesCount"]; ok {
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
		}

		result = append(result, tfAssignment)
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapped %d assignments to Terraform model", len(result)))
	return result
}

// MapAssignmentsToTerraform maps the assignments from the Graph API response to the Terraform model
// used for a read operation.
func MapAssignmentsToTerraform(ctx context.Context, assignmentsResponse models.CloudPcProvisioningPolicyAssignmentCollectionResponseable) []CloudPcProvisioningPolicyAssignmentModel {
	if assignmentsResponse == nil {
		return []CloudPcProvisioningPolicyAssignmentModel{}
	}

	assignments := assignmentsResponse.GetValue()

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d assignments from Graph API to Terraform model", len(assignments)))

	result := make([]CloudPcProvisioningPolicyAssignmentModel, 0, len(assignments))

	for _, assignment := range assignments {
		if assignment == nil {
			continue
		}

		tfAssignment := CloudPcProvisioningPolicyAssignmentModel{}
		id := assignment.GetId()
		if id != nil {
			tfAssignment.ID = types.StringValue(*id)
		} else {
			tfAssignment.ID = types.StringNull()
		}

		target := assignment.GetTarget()
		if target != nil {
			additionalData := target.GetAdditionalData()
			// groupId is always present
			if v, ok := additionalData["groupId"].(string); ok {
				tfAssignment.GroupId = types.StringValue(v)
			} else {
				tfAssignment.GroupId = types.StringNull()
			}
			// servicePlanId is only present for frontline
			if v, ok := additionalData["servicePlanId"].(string); ok && v != "" {
				tfAssignment.ServicePlanId = types.StringValue(v)
			} else {
				tfAssignment.ServicePlanId = types.StringNull()
			}
			// allotmentDisplayName is only present for frontline
			if v, ok := additionalData["allotmentDisplayName"].(string); ok && v != "" {
				tfAssignment.AllotmentDisplayName = types.StringValue(v)
			} else {
				tfAssignment.AllotmentDisplayName = types.StringNull()
			}
			// allotmentLicensesCount is only present for frontline
			if v, ok := additionalData["allotmentLicensesCount"]; ok {
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
		}

		result = append(result, tfAssignment)
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapped %d assignments to Terraform model", len(result)))
	return result
}
