package graphBetaCloudPcUserSetting

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// CloudPcUserSettingAssignmentType returns the object type for CloudPcUserSettingAssignmentModel
func CloudPcUserSettingAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":     types.StringType,
			"group_id": types.StringType,
		},
	}
}

// MapAssignmentsToTerraform maps the remote CloudPcUserSetting assignments to Terraform state
func MapAssignmentsToTerraform(ctx context.Context, data *CloudPcUserSettingResourceModel, assignments []models.CloudPcUserSettingAssignmentable) {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process")
		data.Assignments = types.SetNull(CloudPcUserSettingAssignmentType())
		return
	}

	tflog.Debug(ctx, "Starting assignment mapping process", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	assignmentValues := []attr.Value{}

	for i, assignment := range assignments {
		if assignment == nil {
			tflog.Debug(ctx, "Assignment is nil, skipping", map[string]interface{}{
				"assignmentIndex": i,
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"assignmentIndex": i,
			"resourceId":      data.ID.ValueString(),
		})

		assignmentObj := createAssignmentObject(ctx, assignment, i, data.ID.ValueString())
		if assignmentObj == nil {
			continue
		}

		objValue, diags := types.ObjectValue(CloudPcUserSettingAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
		if !diags.HasError() {
			tflog.Debug(ctx, "Successfully created assignment object", map[string]interface{}{
				"assignmentIndex": i,
				"resourceId":      data.ID.ValueString(),
			})
			assignmentValues = append(assignmentValues, objValue)
		} else {
			tflog.Error(ctx, "Failed to create assignment object value", map[string]interface{}{
				"assignmentIndex": i,
				"errors":          diags.Errors(),
				"resourceId":      data.ID.ValueString(),
			})
		}
	}

	tflog.Debug(ctx, "Creating assignments set", map[string]interface{}{
		"processedAssignments": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})

	if len(assignmentValues) > 0 {
		setVal, diags := types.SetValue(CloudPcUserSettingAssignmentType(), assignmentValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]interface{}{
				"errors":     diags.Errors(),
				"resourceId": data.ID.ValueString(),
			})
			data.Assignments = types.SetNull(CloudPcUserSettingAssignmentType())
		} else {
			tflog.Debug(ctx, "Successfully created assignments set", map[string]interface{}{
				"assignmentCount": len(assignmentValues),
				"resourceId":      data.ID.ValueString(),
			})
			data.Assignments = setVal
		}
	} else {
		tflog.Debug(ctx, "No valid assignments processed, setting assignments to null", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(CloudPcUserSettingAssignmentType())
	}

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]interface{}{
		"finalAssignmentCount": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})
}

// MapAssignmentsToTerraformSet maps assignments from Graph API response directly to a Terraform Set (deprecated, use MapAssignmentsToTerraform)
func MapAssignmentsToTerraformSet(ctx context.Context, assignments []models.CloudPcUserSettingAssignmentable) types.Set {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process, returning null set")
		return types.SetNull(CloudPcUserSettingAssignmentType())
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
		})

		assignmentObj := createAssignmentObject(ctx, assignment, i, "")
		if assignmentObj == nil {
			continue
		}

		objValue, diags := types.ObjectValue(CloudPcUserSettingAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
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
func createAssignmentObject(_ context.Context, assignment models.CloudPcUserSettingAssignmentable, _ int, _ string) map[string]attr.Value {
	assignmentObj := map[string]attr.Value{
		"type":     types.StringNull(),
		"group_id": types.StringNull(),
	}

	// Process target data
	target := assignment.GetTarget()
	if target == nil {
		return nil
	}

	// Get the target's OData type
	odataType := target.GetOdataType()
	if odataType == nil {
		return nil
	}

	// Map based on target type - only cloudPcManagementGroupAssignmentTarget is supported for Windows 365 user settings
	switch *odataType {
	case "#microsoft.graph.cloudPcManagementGroupAssignmentTarget":
		assignmentObj["type"] = types.StringValue("groupAssignmentTarget")

		if managementGroupTarget, ok := target.(models.CloudPcManagementGroupAssignmentTargetable); ok {
			groupId := managementGroupTarget.GetGroupId()
			if groupId != nil && *groupId != "" {
				assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
			} else {
				return nil
			}
		} else {
			return nil
		}

	default:
		return nil
	}

	return assignmentObj
}

// createAssignmentsSet creates the final Set from processed assignment values
func createAssignmentsSet(ctx context.Context, assignmentValues []attr.Value) types.Set {
	tflog.Debug(ctx, "Creating assignments set", map[string]interface{}{
		"processedAssignments": len(assignmentValues),
	})

	if len(assignmentValues) > 0 {
		setVal, diags := types.SetValue(CloudPcUserSettingAssignmentType(), assignmentValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]interface{}{
				"errors": diags.Errors(),
			})
			return types.SetNull(CloudPcUserSettingAssignmentType())
		} else {
			tflog.Debug(ctx, "Successfully created assignments set", map[string]interface{}{
				"assignmentCount": len(assignmentValues),
			})
			return setVal
		}
	}

	tflog.Debug(ctx, "No valid assignments processed, returning null set")
	return types.SetNull(CloudPcUserSettingAssignmentType())
}
