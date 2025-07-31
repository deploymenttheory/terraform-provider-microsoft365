package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// CloudPcProvisioningPolicyAssignmentType returns the object type for CloudPcProvisioningPolicyAssignmentModel
func CloudPcProvisioningPolicyAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":                    types.StringType,
			"group_id":                types.StringType,
			"service_plan_id":         types.StringType,
			"allotment_license_count": types.Int32Type,
			"allotment_display_name":  types.StringType,
		},
	}
}

// MapAssignmentsToTerraformSet maps assignments from Graph API response directly to a Terraform Set
func MapAssignmentsToTerraformSet(ctx context.Context, assignments []models.CloudPcProvisioningPolicyAssignmentable) types.Set {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process, returning null set")
		return types.SetNull(CloudPcProvisioningPolicyAssignmentType())
	}

	tflog.Debug(ctx, "Starting assignment mapping process", map[string]interface{}{
		"assignmentCount": len(assignments),
	})

	assignmentValues := []attr.Value{}

	for i, assignment := range assignments {
		if assignment == nil {
			tflog.Debug(ctx, "Assignment is nil, skipping", map[string]interface{}{
				"assignmentIndex": i,
			})
			continue
		}

		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    convert.GraphToFrameworkString(assignment.GetId()).ValueString(),
		})

		assignmentObj := createAssignmentObject(ctx, assignment, i)
		if assignmentObj == nil {
			continue
		}

		objValue, diags := types.ObjectValue(CloudPcProvisioningPolicyAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
		if !diags.HasError() {
			tflog.Debug(ctx, "Successfully created assignment object", map[string]interface{}{
				"assignmentIndex": i,
			})
			assignmentValues = append(assignmentValues, objValue)
		} else {
			tflog.Error(ctx, "Failed to create assignment object value", map[string]interface{}{
				"assignmentIndex": i,
				"errors":          diags.Errors(),
			})
		}
	}

	return createAssignmentsSet(ctx, assignmentValues)
}

// createAssignmentObject creates a single assignment object from Graph API data
func createAssignmentObject(ctx context.Context, assignment models.CloudPcProvisioningPolicyAssignmentable, index int) map[string]attr.Value {
	assignmentObj := map[string]attr.Value{
		"type":                    types.StringNull(),
		"group_id":                types.StringNull(),
		"service_plan_id":         types.StringNull(),
		"allotment_license_count": types.Int32Null(),
		"allotment_display_name":  types.StringNull(),
	}

	// For Cloud PC assignments, the API returns the group ID as the assignment ID
	// The target additionalData is empty when reading back from the API
	assignmentID := convert.GraphToFrameworkString(assignment.GetId())
	if !assignmentID.IsNull() {
		assignmentObj["group_id"] = assignmentID
		tflog.Debug(ctx, "Set group_id from assignment ID (Cloud PC API behavior)", map[string]interface{}{
			"assignmentIndex": index,
			"assignmentId":    assignmentID.ValueString(),
		})
	}

	// Process target data for any additional fields (like frontline service plan info)
	target := assignment.GetTarget()
	if target != nil {
		mapTargetData(ctx, target, assignmentObj, index)
	}

	// Debug log the final assignment object
	tflog.Debug(ctx, "Final assignment object", map[string]interface{}{
		"assignmentIndex":         index,
		"type":                    assignmentObj["type"],
		"group_id":                assignmentObj["group_id"],
		"service_plan_id":         assignmentObj["service_plan_id"],
		"allotment_license_count": assignmentObj["allotment_license_count"],
		"allotment_display_name":  assignmentObj["allotment_display_name"],
	})

	return assignmentObj
}

// mapTargetData extracts data from the assignment target
func mapTargetData(ctx context.Context, target models.CloudPcManagementAssignmentTargetable, assignmentObj map[string]attr.Value, index int) {
	// Map the target type from OData type
	odataType := ""
	if target.GetOdataType() != nil {
		odataType = *target.GetOdataType()
	}

	switch odataType {
	case "#microsoft.graph.cloudPcManagementGroupAssignmentTarget":
		assignmentObj["type"] = types.StringValue("groupAssignmentTarget")
	case "#microsoft.graph.cloudPcManagementExclusionGroupAssignmentTarget":
		assignmentObj["type"] = types.StringValue("exclusionGroupAssignmentTarget")
	default:
		// Default to groupAssignmentTarget for backwards compatibility
		assignmentObj["type"] = types.StringValue("groupAssignmentTarget")
		tflog.Debug(ctx, "Unknown or null OData type, defaulting to groupAssignmentTarget", map[string]interface{}{
			"assignmentIndex": index,
			"odataType":       odataType,
		})
	}
	additionalData := target.GetAdditionalData()
	if additionalData == nil {
		tflog.Debug(ctx, "Target additionalData is nil", map[string]interface{}{
			"assignmentIndex": index,
		})
		return
	}

	tflog.Debug(ctx, "Processing target additionalData", map[string]interface{}{
		"assignmentIndex": index,
		"additionalData":  fmt.Sprintf("%+v", additionalData),
	})

	// Extract groupId
	if groupId, ok := additionalData["groupId"].(string); ok && groupId != "" {
		assignmentObj["group_id"] = types.StringValue(groupId)
		tflog.Debug(ctx, "Set groupId from additionalData", map[string]interface{}{
			"assignmentIndex": index,
			"groupId":         groupId,
		})
	} else {
		tflog.Debug(ctx, "No groupId found in additionalData", map[string]interface{}{
			"assignmentIndex": index,
		})
	}

	// Extract servicePlanId (for frontline)
	if servicePlanId, ok := additionalData["servicePlanId"].(string); ok && servicePlanId != "" {
		assignmentObj["service_plan_id"] = types.StringValue(servicePlanId)
	}

	// Extract allotmentDisplayName (for frontline)
	if allotmentDisplayName, ok := additionalData["allotmentDisplayName"].(string); ok && allotmentDisplayName != "" && allotmentDisplayName != "null" {
		assignmentObj["allotment_display_name"] = types.StringValue(allotmentDisplayName)
	}

	// Extract allotmentLicensesCount (for frontline) - convert to int32
	if allotmentLicensesCount, ok := additionalData["allotmentLicensesCount"]; ok && allotmentLicensesCount != nil && allotmentLicensesCount != "null" {
		switch val := allotmentLicensesCount.(type) {
		case int32:
			assignmentObj["allotment_license_count"] = types.Int32Value(val)
		case int64:
			assignmentObj["allotment_license_count"] = types.Int32Value(int32(val))
		case float64:
			assignmentObj["allotment_license_count"] = types.Int32Value(int32(val))
		}
	}
}

// createAssignmentsSet creates the final Set from processed assignment values
func createAssignmentsSet(ctx context.Context, assignmentValues []attr.Value) types.Set {
	tflog.Debug(ctx, "Creating assignments set", map[string]interface{}{
		"processedAssignments": len(assignmentValues),
	})

	if len(assignmentValues) > 0 {
		setVal, diags := types.SetValue(CloudPcProvisioningPolicyAssignmentType(), assignmentValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]interface{}{
				"errors": diags.Errors(),
			})
			return types.SetNull(CloudPcProvisioningPolicyAssignmentType())
		}

		tflog.Debug(ctx, "Successfully created assignments set", map[string]interface{}{
			"assignmentCount": len(assignmentValues),
		})
		return setVal
	}

	tflog.Debug(ctx, "No valid assignments processed, returning null set")
	return types.SetNull(CloudPcProvisioningPolicyAssignmentType())
}
