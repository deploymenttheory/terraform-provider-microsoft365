package graphBetaLinuxPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphsdkmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs and returns a ConfigurationPoliciesItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *LinuxPlatformScriptResourceModel) (devicemanagement.ConfigurationPoliciesItemAssignPostRequestBodyable, error) {
	if data.Assignments == nil {
		return nil, fmt.Errorf("assignments configuration block is required even if empty. Minimum config requires all_devices and all_users booleans to be set to false")
	}

	tflog.Debug(ctx, "Starting assignment construction")

	if err := validateAssignmentConfig(data.Assignments); err != nil {
		return nil, err
	}

	requestBody := devicemanagement.NewConfigurationPoliciesItemAssignPostRequestBody()
	assignments := make([]graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable, 0)

	// Check All Devices
	if !data.Assignments.AllDevices.IsNull() && data.Assignments.AllDevices.ValueBool() {
		assignments = append(assignments, constructAllDevicesAssignment(ctx, data.Assignments))
	}

	// Check All Users
	if !data.Assignments.AllUsers.IsNull() && data.Assignments.AllUsers.ValueBool() {
		assignments = append(assignments, constructAllUsersAssignment(ctx, data.Assignments))
	}

	// Check Include Groups
	if !data.Assignments.AllDevices.ValueBool() &&
		!data.Assignments.AllUsers.ValueBool() &&
		len(data.Assignments.IncludeGroups) > 0 {
		for _, group := range data.Assignments.IncludeGroups {
			if !group.GroupId.IsNull() && !group.GroupId.IsUnknown() && group.GroupId.ValueString() != "" {
				assignments = append(assignments, constructGroupIncludeAssignments(ctx, data.Assignments)...)
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
	requestBody.SetAssignments(assignments)

	// Debug log the final request body
	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructAllDevicesAssignment constructs and returns a DeviceManagementConfigurationPolicyAssignment object for all devices
func constructAllDevicesAssignment(ctx context.Context, config *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel) graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
	target := graphsdkmodels.NewAllDevicesAssignmentTarget()

	assignment.SetTarget(target)

	return assignment
}

// constructAllUsersAssignment constructs and returns a DeviceManagementConfigurationPolicyAssignment object for all licensed users
func constructAllUsersAssignment(ctx context.Context, config *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel) graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
	target := graphsdkmodels.NewAllLicensedUsersAssignmentTarget()

	assignment.SetTarget(target)

	return assignment
}

// constructGroupIncludeAssignments constructs and returns a list of DeviceManagementConfigurationPolicyAssignment objects for included groups
func constructGroupIncludeAssignments(ctx context.Context, config *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel) []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	var assignments []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable

	for _, groupFilter := range config.IncludeGroups {
		assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
		target := graphsdkmodels.NewGroupAssignmentTarget()

		constructors.SetStringProperty(groupFilter.GroupId, target.SetGroupId)

		assignment.SetTarget(target)
		assignments = append(assignments, assignment)
	}

	return assignments
}

// constructGroupExcludeAssignments constructs and returns a list of DeviceManagementConfigurationPolicyAssignment objects for excluded groups
func constructGroupExcludeAssignments(config *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel) []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	var assignments []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable

	// Check if we have any non-null, non-empty values
	hasValidExcludes := false
	for _, groupId := range config.ExcludeGroupIds {
		if !groupId.IsNull() && !groupId.IsUnknown() && groupId.ValueString() != "" {
			hasValidExcludes = true
			break
		}
	}

	if hasValidExcludes {
		for _, groupId := range config.ExcludeGroupIds {
			if !groupId.IsNull() && !groupId.IsUnknown() && groupId.ValueString() != "" {
				assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
				target := graphsdkmodels.NewExclusionGroupAssignmentTarget()

				constructors.SetStringProperty(groupId, target.SetGroupId)

				assignment.SetTarget(target)
				assignments = append(assignments, assignment)
			}
		}
	}

	return assignments
}
