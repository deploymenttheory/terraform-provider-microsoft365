package graphBetaDeviceManagementTemplate

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/validators"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphsdkmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs and returns a ConfigurationPoliciesItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *DeviceManagementTemplateResourceModel) (devicemanagement.ConfigurationPoliciesItemAssignPostRequestBodyable, error) {
	if data.Assignments == nil {
		return nil, fmt.Errorf("assignments configuration block is required even if empty. Minimum config requires all_devices and all_users booleans to be set to false")
	}

	tflog.Debug(ctx, "Starting assignment construction")

	if err := validators.ValidateDeviceConfiguationAssignmentSettings(data.Assignments); err != nil {
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

	if err := construct.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
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

	if !config.AllDevicesFilterId.IsNull() && !config.AllDevicesFilterId.IsUnknown() &&
		config.AllDevicesFilterId.ValueString() != "" {
		construct.SetStringProperty(config.AllDevicesFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		if !config.AllDevicesFilterType.IsNull() && !config.AllDevicesFilterType.IsUnknown() {
			err := construct.ParseEnum(config.AllDevicesFilterType,
				graphsdkmodels.ParseDeviceAndAppManagementAssignmentFilterType,
				target.SetDeviceAndAppManagementAssignmentFilterType)
			if err != nil {
				tflog.Warn(ctx, "Failed to parse all devices filter type", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}

	assignment.SetTarget(target)
	return assignment
}

// constructAllUsersAssignment constructs and returns a DeviceManagementConfigurationPolicyAssignment object for all licensed users
func constructAllUsersAssignment(ctx context.Context, config *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel) graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
	target := graphsdkmodels.NewAllLicensedUsersAssignmentTarget()

	if !config.AllUsersFilterId.IsNull() && !config.AllUsersFilterId.IsUnknown() &&
		config.AllUsersFilterId.ValueString() != "" {
		construct.SetStringProperty(config.AllUsersFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		if !config.AllUsersFilterType.IsNull() && !config.AllUsersFilterType.IsUnknown() {
			err := construct.ParseEnum(config.AllUsersFilterType,
				graphsdkmodels.ParseDeviceAndAppManagementAssignmentFilterType,
				target.SetDeviceAndAppManagementAssignmentFilterType)
			if err != nil {
				tflog.Warn(ctx, "Failed to parse all users filter type", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}

	assignment.SetTarget(target)
	return assignment
}

// constructGroupIncludeAssignments constructs and returns a list of DeviceManagementConfigurationPolicyAssignment objects for included groups
func constructGroupIncludeAssignments(ctx context.Context, config *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel) []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	var assignments []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable

	for _, groupFilter := range config.IncludeGroups {
		assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
		target := graphsdkmodels.NewGroupAssignmentTarget()

		construct.SetStringProperty(groupFilter.GroupId, target.SetGroupId)

		if !groupFilter.IncludeGroupsFilterId.IsNull() && !groupFilter.IncludeGroupsFilterType.IsNull() {
			construct.SetStringProperty(groupFilter.IncludeGroupsFilterId,
				target.SetDeviceAndAppManagementAssignmentFilterId)

			err := construct.ParseEnum(groupFilter.IncludeGroupsFilterType,
				graphsdkmodels.ParseDeviceAndAppManagementAssignmentFilterType,
				target.SetDeviceAndAppManagementAssignmentFilterType)
			if err != nil {
				tflog.Warn(ctx, "Failed to parse include groups filter type", map[string]interface{}{
					"error":   err.Error(),
					"groupId": groupFilter.GroupId.ValueString(),
				})
			}
		}

		assignment.SetTarget(target)
		assignments = append(assignments, assignment)
	}

	return assignments
}

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

	// Only process if we have valid excludes
	if hasValidExcludes {
		for _, groupId := range config.ExcludeGroupIds {
			if !groupId.IsNull() && !groupId.IsUnknown() && groupId.ValueString() != "" {
				assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
				target := graphsdkmodels.NewExclusionGroupAssignmentTarget()
				construct.SetStringProperty(groupId, target.SetGroupId)

				assignment.SetTarget(target)
				assignments = append(assignments, assignment)
			}
		}
	}

	return assignments
}
