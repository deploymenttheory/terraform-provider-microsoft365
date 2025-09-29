package graphBetaRoleScopeTag

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
)

// constructAssignment constructs and returns a RoleScopeTagsItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *RoleScopeTagResourceModel) (devicemanagement.RoleScopeTagsItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Starting role scope tag assignment construction")

	requestBody := devicemanagement.NewRoleScopeTagsItemAssignPostRequestBody()
	roleScopeTagAssignments := make([]graphmodels.RoleScopeTagAutoAssignmentable, 0)

	if data.Assignments.IsNull() || data.Assignments.IsUnknown() {
		tflog.Debug(ctx, "Assignments is null or unknown, creating empty assignments array")
		requestBody.SetAssignments(roleScopeTagAssignments)
		return requestBody, nil
	}

	var terraformAssignments []sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithoutGroupFilterModel
	diags := data.Assignments.ElementsAs(ctx, &terraformAssignments, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract assignments: %v", diags.Errors())
	}

	for idx, assignment := range terraformAssignments {
		tflog.Debug(ctx, "Processing assignment", map[string]any{
			"index": idx,
		})

		// Create RoleScopeTagAutoAssignment
		roleScopeTagAssignment := graphmodels.NewRoleScopeTagAutoAssignment()

		if assignment.Type.IsNull() || assignment.Type.IsUnknown() {
			tflog.Error(ctx, "Assignment target type is missing or invalid", map[string]any{
				"index": idx,
			})
			continue
		}

		targetType := assignment.Type.ValueString()

		target := constructTarget(ctx, targetType, assignment)
		if target == nil {
			tflog.Error(ctx, "Failed to create target", map[string]any{
				"index":      idx,
				"targetType": targetType,
			})
			continue
		}

		roleScopeTagAssignment.SetTarget(target)
		roleScopeTagAssignments = append(roleScopeTagAssignments, roleScopeTagAssignment)
	}

	tflog.Debug(ctx, "Completed assignment construction", map[string]any{
		"totalAssignments": len(roleScopeTagAssignments),
	})

	requestBody.SetAssignments(roleScopeTagAssignments)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructTarget creates the appropriate target based on the target type
func constructTarget(ctx context.Context, targetType string, assignment sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithoutGroupFilterModel) graphmodels.DeviceAndAppManagementAssignmentTargetable {
	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch targetType {
	case "groupAssignmentTarget":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.GroupId, groupTarget.SetGroupId)
		} else {
			tflog.Error(ctx, "Group assignment target missing required group_id", map[string]any{
				"targetType": targetType,
			})
			return nil
		}
		target = groupTarget
	default:
		tflog.Error(ctx, "Unsupported target type", map[string]any{
			"targetType": targetType,
		})
		return nil
	}

	return target
}
