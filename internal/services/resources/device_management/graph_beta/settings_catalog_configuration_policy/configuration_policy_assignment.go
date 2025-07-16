package graphBetaSettingsCatalogConfigurationPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphsdkmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructConfigurationPolicyAssignmen constructs and returns a ConfigurationPoliciesItemAssignPostRequestBody
func ConstructConfigurationPolicyAssignment(ctx context.Context, data *SettingsCatalogSettingsAssignmentResourceModel) (devicemanagement.ConfigurationPoliciesItemAssignPostRequestBodyable, error) {
	if data == nil {
		return nil, fmt.Errorf("assignments configuration block is required even if empty. Minimum config requires all_devices and all_users booleans to be set to false")
	}

	tflog.Debug(ctx, "Starting Configuartion Policy (settings catalog) assignment construction")

	// if err := validators.ValidateDeviceConfiguationAssignmentSettings(data); err != nil {
	// 	return nil, err
	// }

	requestBody := devicemanagement.NewConfigurationPoliciesItemAssignPostRequestBody()
	assignments := make([]graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable, 0)

	// Check All Devices
	if !data.AllDevices.IsNull() && data.AllDevices.ValueBool() {
		assignments = append(assignments, constructAllDevicesAssignment(ctx, data))
	}

	// Check All Users
	if !data.AllUsers.IsNull() && data.AllUsers.ValueBool() {
		assignments = append(assignments, constructAllUsersAssignment(ctx, data))
	}

	// Check Include Groups
	if !data.AllDevices.ValueBool() &&
		!data.AllUsers.ValueBool() &&
		len(data.IncludeGroups) > 0 {
		for _, group := range data.IncludeGroups {
			if !group.GroupId.IsNull() && !group.GroupId.IsUnknown() && group.GroupId.ValueString() != "" {
				assignments = append(assignments, constructGroupIncludeAssignments(ctx, data)...)
				break
			}
		}
	}

	// Check Exclude Groups
	if len(data.ExcludeGroupIds) > 0 {
		for _, id := range data.ExcludeGroupIds {
			if !id.IsNull() && !id.IsUnknown() && id.ValueString() != "" {
				assignments = append(assignments, constructGroupExcludeAssignments(data)...)
				break
			}
		}
	}

	// Always set assignments (will be empty array if no active assignments)
	// as update http method is a post not patch.
	requestBody.SetAssignments(assignments)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructAllDevicesAssignment constructs and returns a DeviceManagementConfigurationPolicyAssignment object for all devices
func constructAllDevicesAssignment(ctx context.Context, config *SettingsCatalogSettingsAssignmentResourceModel) graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
	target := graphsdkmodels.NewAllDevicesAssignmentTarget()

	if !config.AllDevicesFilterId.IsNull() && !config.AllDevicesFilterId.IsUnknown() &&
		config.AllDevicesFilterId.ValueString() != "" {
		convert.FrameworkToGraphString(config.AllDevicesFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		if !config.AllDevicesFilterType.IsNull() && !config.AllDevicesFilterType.IsUnknown() {
			err := convert.FrameworkToGraphEnum(config.AllDevicesFilterType,
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
func constructAllUsersAssignment(ctx context.Context, config *SettingsCatalogSettingsAssignmentResourceModel) graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
	target := graphsdkmodels.NewAllLicensedUsersAssignmentTarget()

	if !config.AllUsersFilterId.IsNull() && !config.AllUsersFilterId.IsUnknown() &&
		config.AllUsersFilterId.ValueString() != "" {
		convert.FrameworkToGraphString(config.AllUsersFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		if !config.AllUsersFilterType.IsNull() && !config.AllUsersFilterType.IsUnknown() {
			err := convert.FrameworkToGraphEnum(config.AllUsersFilterType,
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
func constructGroupIncludeAssignments(ctx context.Context, config *SettingsCatalogSettingsAssignmentResourceModel) []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	var assignments []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable

	for _, groupFilter := range config.IncludeGroups {
		assignment := graphsdkmodels.NewDeviceManagementConfigurationPolicyAssignment()
		target := graphsdkmodels.NewGroupAssignmentTarget()

		convert.FrameworkToGraphString(groupFilter.GroupId, target.SetGroupId)

		if !groupFilter.IncludeGroupsFilterId.IsNull() && !groupFilter.IncludeGroupsFilterType.IsNull() {
			convert.FrameworkToGraphString(groupFilter.IncludeGroupsFilterId,
				target.SetDeviceAndAppManagementAssignmentFilterId)

			err := convert.FrameworkToGraphEnum(groupFilter.IncludeGroupsFilterType,
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

func constructGroupExcludeAssignments(config *SettingsCatalogSettingsAssignmentResourceModel) []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable {
	var assignments []graphsdkmodels.DeviceManagementConfigurationPolicyAssignmentable

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
				convert.FrameworkToGraphString(groupId, target.SetGroupId)

				assignment.SetTarget(target)
				assignments = append(assignments, assignment)
			}
		}
	}

	return assignments
}
