package graphBetaWindowsPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphsdkmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs and returns a ConfigurationPoliciesItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *WindowsPlatformScriptResourceModel) (devicemanagement.DeviceManagementScriptsItemAssignPostRequestBodyable, error) {
	if data.Assignments == nil {
		return nil, fmt.Errorf("assignments configuration block is required even if empty. Minimum config requires all_devices and all_users booleans to be set to false")
	}

	tflog.Debug(ctx, "Starting assignment construction")

	if err := validateAssignmentConfig(data.Assignments); err != nil {
		return nil, err
	}

	requestBody := devicemanagement.NewDeviceManagementScriptsItemAssignPostRequestBody()
	assignments := make([]graphsdkmodels.DeviceManagementScriptAssignmentable, 0)

	// Check All Devices
	if !data.Assignments.AllDevices.IsNull() && data.Assignments.AllDevices.ValueBool() {
		assignments = append(assignments, constructAllDevicesAssignment())
	}

	// Check All Users
	if !data.Assignments.AllUsers.IsNull() && data.Assignments.AllUsers.ValueBool() {
		assignments = append(assignments, constructAllUsersAssignment())
	}

	// Check Include Groups
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

	// Check Exclude Groups
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

			constructors.SetStringProperty(groupId, target.SetGroupId)
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

				constructors.SetStringProperty(groupId, target.SetGroupId)

				assignment.SetTarget(target)
				assignments = append(assignments, assignment)
			}
		}
	}

	return assignments
}
