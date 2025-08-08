package graphBetaCloudPcUserSetting

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignmentsRequestBody creates the request body for assigning a Cloud PC user setting to groups
func constructAssignmentsRequestBody(ctx context.Context, assignments types.Set) (*devicemanagement.VirtualEndpointUserSettingsItemAssignPostRequestBody, error) {
	tflog.Debug(ctx, "Constructing assignments request body")

	requestBody := devicemanagement.NewVirtualEndpointUserSettingsItemAssignPostRequestBody()

	// If no assignments are provided, return an empty assignments array
	// This is used for removing all assignments
	if assignments.IsNull() || assignments.IsUnknown() {
		requestBody.SetAssignments([]models.CloudPcUserSettingAssignmentable{})
		tflog.Debug(ctx, "No assignments provided, setting empty assignments array")
		return requestBody, nil
	}

	var terraformAssignments []CloudPcUserSettingAssignmentModel
	diags := assignments.ElementsAs(ctx, &terraformAssignments, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract assignments: %v", diags.Errors())
	}

	graphAssignments := []models.CloudPcUserSettingAssignmentable{}

	for i, assignment := range terraformAssignments {
		if assignment.Type.IsNull() || assignment.Type.IsUnknown() {
			tflog.Error(ctx, "Assignment target type is missing or invalid", map[string]interface{}{
				"index": i,
			})
			continue
		}

		targetType := assignment.Type.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Creating assignment %d with target type: %s", i, targetType))

		graphAssignment := models.NewCloudPcUserSettingAssignment()

		target := constructTarget(ctx, targetType, assignment)
		if target == nil {
			tflog.Error(ctx, "Failed to create target", map[string]interface{}{
				"index":      i,
				"targetType": targetType,
			})
			continue
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

// constructTarget creates the appropriate target based on the target type
func constructTarget(ctx context.Context, targetType string, assignment CloudPcUserSettingAssignmentModel) models.CloudPcManagementAssignmentTargetable {
	var target models.CloudPcManagementAssignmentTargetable

	switch targetType {
	case "allDevicesAssignmentTarget":
		allDevicesTarget := models.NewAllDevicesAssignmentTarget()
		target = allDevicesTarget
		tflog.Debug(ctx, "Created AllDevicesAssignmentTarget")
	case "allLicensedUsersAssignmentTarget":
		allUsersTarget := models.NewAllLicensedUsersAssignmentTarget()
		target = allUsersTarget
		tflog.Debug(ctx, "Created AllLicensedUsersAssignmentTarget")
	case "groupAssignmentTarget":
		groupTarget := models.NewGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.GroupId, groupTarget.SetGroupId)
			tflog.Debug(ctx, "Created GroupAssignmentTarget", map[string]interface{}{
				"groupId": assignment.GroupId.ValueString(),
			})
		} else {
			tflog.Error(ctx, "Group assignment target missing required group_id", map[string]interface{}{
				"targetType": targetType,
			})
			return nil
		}
		target = groupTarget
	case "exclusionGroupAssignmentTarget":
		exclusionTarget := models.NewExclusionGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.GroupId, exclusionTarget.SetGroupId)
			tflog.Debug(ctx, "Created ExclusionGroupAssignmentTarget", map[string]interface{}{
				"groupId": assignment.GroupId.ValueString(),
			})
		} else {
			tflog.Error(ctx, "Exclusion group assignment target missing required group_id", map[string]interface{}{
				"targetType": targetType,
			})
			return nil
		}
		target = exclusionTarget
	default:
		tflog.Error(ctx, "Unsupported target type", map[string]interface{}{
			"targetType": targetType,
		})
		return nil
	}

	return target
}
