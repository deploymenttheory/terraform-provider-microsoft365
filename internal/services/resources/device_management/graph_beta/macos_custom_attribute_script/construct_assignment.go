package graphBetaMacOSCustomAttributeScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphsdkmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs and returns a ConfigurationPoliciesItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *DeviceCustomAttributeShellScriptResourceModel) (devicemanagement.DeviceManagementScriptsItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Starting assignment construction")

	// Create the request body with an empty assignments array
	requestBody := devicemanagement.NewDeviceShellScriptsItemAssignPostRequestBody()
	assignments := make([]graphsdkmodels.DeviceManagementScriptAssignmentable, 0)

	// If assignments is nil, return the request body with an empty assignments array
	// This will effectively remove all assignments when sent to the API
	if data.Assignments == nil {
		tflog.Debug(ctx, "Assignments is nil, creating empty assignments array to remove all assignments")
		requestBody.SetDeviceManagementScriptAssignments(assignments)

		if err := constructors.DebugLogGraphObject(ctx, "Constructed empty assignment request body", requestBody); err != nil {
			tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
				"error": err.Error(),
			})
		}

		return requestBody, nil
	}

	if err := validateAssignmentConfig(data.Assignments); err != nil {
		return nil, err
	}

	if !data.Assignments.AllDevices.IsNull() && data.Assignments.AllDevices.ValueBool() {
		assignments = append(assignments, constructAllDevicesAssignment())
	}

	if !data.Assignments.AllUsers.IsNull() && data.Assignments.AllUsers.ValueBool() {
		assignments = append(assignments, constructAllUsersAssignment())
	}

	if !data.Assignments.AllDevices.ValueBool() &&
		!data.Assignments.AllUsers.ValueBool() &&
		len(data.Assignments.IncludeGroupIds) > 0 {
		for _, id := range data.Assignments.IncludeGroupIds {
			if !id.IsNull() && !id.IsUnknown() && id.ValueString() != "" {
				assignments = append(assignments, constructGroupIncludeAssignments(data.Assignments)...)
				break
			}
		}
	}

	if len(data.Assignments.ExcludeGroupIds) > 0 {
		for _, id := range data.Assignments.ExcludeGroupIds {
			if !id.IsNull() && !id.IsUnknown() && id.ValueString() != "" {
				assignments = append(assignments, constructGroupExcludeAssignments(data.Assignments)...)
				break
			}
		}
	}

	// Always set assignments (will be empty array if no active assignments)
	// as update http method is a post not patch.
	requestBody.SetDeviceManagementScriptAssignments(assignments)

	// Debug log the final request body
	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructAllDevicesAssignment constructs and returns a DeviceManagementScriptAssignment object for all devices
func constructAllDevicesAssignment() graphsdkmodels.DeviceManagementScriptAssignmentable {
	assignment := graphsdkmodels.NewDeviceManagementScriptAssignment()
	target := graphsdkmodels.NewAllDevicesAssignmentTarget()

	assignment.SetTarget(target)
	return assignment
}

// constructAllUsersAssignment constructs and returns a DeviceManagementScriptAssignment object for all licensed users
func constructAllUsersAssignment() graphsdkmodels.DeviceManagementScriptAssignmentable {
	assignment := graphsdkmodels.NewDeviceManagementScriptAssignment()
	target := graphsdkmodels.NewAllLicensedUsersAssignmentTarget()

	assignment.SetTarget(target)
	return assignment
}

// constructGroupIncludeAssignments constructs and returns a list of DeviceManagementConfigurationPolicyAssignment objects for included groups
func constructGroupIncludeAssignments(config *sharedmodels.DeviceManagementScriptAssignmentResourceModel) []graphsdkmodels.DeviceManagementScriptAssignmentable {
	var assignments []graphsdkmodels.DeviceManagementScriptAssignmentable

	for _, groupId := range config.IncludeGroupIds {
		if !groupId.IsNull() && !groupId.IsUnknown() && groupId.ValueString() != "" {
			assignment := graphsdkmodels.NewDeviceManagementScriptAssignment()
			target := graphsdkmodels.NewGroupAssignmentTarget()

			convert.FrameworkToGraphString(groupId, target.SetGroupId)
			assignment.SetTarget(target)
			assignments = append(assignments, assignment)
		}
	}

	return assignments
}

func constructGroupExcludeAssignments(config *sharedmodels.DeviceManagementScriptAssignmentResourceModel) []graphsdkmodels.DeviceManagementScriptAssignmentable {
	var assignments []graphsdkmodels.DeviceManagementScriptAssignmentable

	// Check if we have any non-null, non-empty values
	hasValidExcludes := false
	for _, groupId := range config.ExcludeGroupIds {
		if !groupId.IsNull() && !groupId.IsUnknown() && groupId.ValueString() != "" {
			hasValidExcludes = true
			break
		}
	}

	// Only process if we have valid excludes
	if hasValidExcludes {
		for _, groupId := range config.ExcludeGroupIds {
			if !groupId.IsNull() && !groupId.IsUnknown() && groupId.ValueString() != "" {
				assignment := graphsdkmodels.NewDeviceManagementScriptAssignment()
				target := graphsdkmodels.NewExclusionGroupAssignmentTarget()

				convert.FrameworkToGraphString(groupId, target.SetGroupId)

				assignment.SetTarget(target)
				assignments = append(assignments, assignment)
			}
		}
	}

	return assignments
}
