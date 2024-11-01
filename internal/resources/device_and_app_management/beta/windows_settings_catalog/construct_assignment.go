package graphBetaWindowsSettingsCatalog

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignments constructs and returns a ConfigurationPoliciesItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel) (devicemanagement.ConfigurationPoliciesItemAssignPostRequestBodyable, error) {
	if data.Assignments == nil {
		return nil, fmt.Errorf("assignments configuration is required")
	}

	tflog.Debug(ctx, "Starting assignment construction")
	logAssignmentDetails(ctx, data.Assignments)

	if err := validateAssignmentConfig(data.Assignments); err != nil {
		return nil, err
	}

	requestBody := devicemanagement.NewConfigurationPoliciesItemAssignPostRequestBody()
	var assignments []graphmodels.DeviceManagementConfigurationPolicyAssignmentable

	// Handle All Devices assignment
	if !data.Assignments.AllDevices.IsNull() && data.Assignments.AllDevices.ValueBool() {
		assignments = append(assignments, constructAllDevicesAssignment(ctx, data.Assignments))
	}

	// Handle All Users assignment
	if !data.Assignments.AllUsers.IsNull() && data.Assignments.AllUsers.ValueBool() {
		assignments = append(assignments, constructAllUsersAssignment(ctx, data.Assignments))
	}

	// Handle Include Groups assignment - only if AllDevices and AllUsers are false
	if !data.Assignments.AllDevices.ValueBool() &&
		!data.Assignments.AllUsers.ValueBool() &&
		len(data.Assignments.IncludeGroups) > 0 {
		assignments = append(assignments, constructGroupIncludeAssignments(ctx, data.Assignments)...)
	}

	// Handle Exclude Groups assignment
	if len(data.Assignments.ExcludeGroupIds) > 0 {
		assignments = append(assignments, constructGroupExcludeAssignments(data.Assignments)...)
	}

	requestBody.SetAssignments(assignments)
	return requestBody, nil
}

// ValidateAssignmentConfiguration validates the assignment configuration
// - if all_devices is false, device filter settings should not be set
// - if all_users is false, user filter settings should not be set
// - AllDevicesFilterType must be one of: include, exclude, none
// - AllUsersFilterType must be one of: include, exclude, none
// - IncludeGroupsFilterType must be one of: include, exclude, none
// - ExcludeGroupIds must be in alphanumeric order by group_id
// - IncludeGroups must be in alphanumeric order by group_id
// - No group is used more than once across include and exclude assignments
// - AllDevices and IncludeGroups cannot be used at the same time
// - AllUsers and IncludeGroups cannot be used at the same time
func validateAssignmentConfig(config *SettingsCatalogSettingsAssignmentResourceModel) error {
	// Validate filter types have valid values
	validFilterTypes := map[string]bool{
		"include": true,
		"exclude": true,
		"none":    true,
	}

	if !config.AllDevices.ValueBool() {
		if !config.AllDevicesFilterType.IsNull() || !config.AllDevicesFilterId.IsNull() {
			return fmt.Errorf("all_devices_filter_type and/or all_devices_filter_id cannot be set when all_devices is false")
		}
	}

	if !config.AllUsers.ValueBool() {
		if !config.AllUsersFilterType.IsNull() || !config.AllUsersFilterId.IsNull() {
			return fmt.Errorf("all_users_filter_type and/or all_users_filter_id cannot be set when all_users is false")
		}
	}

	if !config.AllDevicesFilterType.IsNull() && !validFilterTypes[config.AllDevicesFilterType.ValueString()] {
		return fmt.Errorf("AllDevicesFilterType must be one of: include, exclude, none. Got: %s",
			config.AllDevicesFilterType.ValueString())
	}

	if !config.AllUsersFilterType.IsNull() && !validFilterTypes[config.AllUsersFilterType.ValueString()] {
		return fmt.Errorf("AllUsersFilterType must be one of: include, exclude, none. Got: %s",
			config.AllUsersFilterType.ValueString())
	}

	// Validate include_groups filter types and alphanumeric order
	if len(config.IncludeGroups) > 1 {
		for i := 0; i < len(config.IncludeGroups)-1; i++ {
			current := config.IncludeGroups[i].GroupId.ValueString()
			next := config.IncludeGroups[i+1].GroupId.ValueString()
			if current > next {
				return fmt.Errorf("include_groups must be in alphanumeric order by group_id. Found %s before %s",
					current, next)
			}
		}
	}

	for _, group := range config.IncludeGroups {
		if !group.IncludeGroupsFilterType.IsNull() && !validFilterTypes[group.IncludeGroupsFilterType.ValueString()] {
			return fmt.Errorf("IncludeGroupsFilterType must be one of: include, exclude, none. Got: %s",
				group.IncludeGroupsFilterType.ValueString())
		}
	}

	// Validate exclude_group_ids are in alphanumeric order
	if len(config.ExcludeGroupIds) > 1 {
		for i := 0; i < len(config.ExcludeGroupIds)-1; i++ {
			current := config.ExcludeGroupIds[i].ValueString()
			next := config.ExcludeGroupIds[i+1].ValueString()
			if current > next {
				return fmt.Errorf("exclude_group_ids must be in alphanumeric order. Found %s before %s",
					current, next)
			}
		}
	}

	// Validate no group is used more than once across include and exclude assignments
	for _, includeGroup := range config.IncludeGroups {
		for _, excludeGroupId := range config.ExcludeGroupIds {
			if !includeGroup.GroupId.IsNull() && !excludeGroupId.IsNull() &&
				includeGroup.GroupId.ValueString() == excludeGroupId.ValueString() {
				return fmt.Errorf("group %s is used in both include and exclude assignments. Each group assignment can only be used once across all assignment rules",
					includeGroup.GroupId.ValueString())
			}
		}
	}

	// Validate AllDevices cannot be used with Include Groups
	if !config.AllDevices.IsNull() && config.AllDevices.ValueBool() && len(config.IncludeGroups) > 0 {
		return fmt.Errorf("cannot assign to All Devices and Include Groups at the same time")
	}

	// Validate AllUsers cannot be used with Include Groups
	if !config.AllUsers.IsNull() && config.AllUsers.ValueBool() && len(config.IncludeGroups) > 0 {
		return fmt.Errorf("cannot assign to All Users and Include Groups at the same time")
	}

	return nil
}

// constructAllDevicesAssignment constructs and returns a DeviceManagementConfigurationPolicyAssignment object for all devices
func constructAllDevicesAssignment(ctx context.Context, config *SettingsCatalogSettingsAssignmentResourceModel) graphmodels.DeviceManagementConfigurationPolicyAssignmentable {
	assignment := graphmodels.NewDeviceManagementConfigurationPolicyAssignment()
	target := graphmodels.NewAllDevicesAssignmentTarget()

	if !config.AllDevicesFilterId.IsNull() && !config.AllDevicesFilterType.IsNull() {
		construct.SetStringProperty(config.AllDevicesFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		err := construct.ParseEnum(config.AllDevicesFilterType,
			graphmodels.ParseDeviceAndAppManagementAssignmentFilterType,
			target.SetDeviceAndAppManagementAssignmentFilterType)
		if err != nil {
			tflog.Warn(ctx, "Failed to parse all devices filter type", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	assignment.SetTarget(target)
	return assignment
}

// constructAllUsersAssignment constructs and returns a DeviceManagementConfigurationPolicyAssignment object for all licensed users
func constructAllUsersAssignment(ctx context.Context, config *SettingsCatalogSettingsAssignmentResourceModel) graphmodels.DeviceManagementConfigurationPolicyAssignmentable {
	assignment := graphmodels.NewDeviceManagementConfigurationPolicyAssignment()
	target := graphmodels.NewAllLicensedUsersAssignmentTarget()

	if !config.AllUsersFilterId.IsNull() && !config.AllUsersFilterType.IsNull() {
		construct.SetStringProperty(config.AllUsersFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		err := construct.ParseEnum(config.AllUsersFilterType,
			graphmodels.ParseDeviceAndAppManagementAssignmentFilterType,
			target.SetDeviceAndAppManagementAssignmentFilterType)
		if err != nil {
			tflog.Warn(ctx, "Failed to parse all users filter type", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	assignment.SetTarget(target)
	return assignment
}

// constructGroupIncludeAssignments constructs and returns a list of DeviceManagementConfigurationPolicyAssignment objects for included groups
func constructGroupIncludeAssignments(ctx context.Context, config *SettingsCatalogSettingsAssignmentResourceModel) []graphmodels.DeviceManagementConfigurationPolicyAssignmentable {
	var assignments []graphmodels.DeviceManagementConfigurationPolicyAssignmentable

	for _, groupFilter := range config.IncludeGroups {
		assignment := graphmodels.NewDeviceManagementConfigurationPolicyAssignment()
		target := graphmodels.NewGroupAssignmentTarget()

		construct.SetStringProperty(groupFilter.GroupId, target.SetGroupId)

		if !groupFilter.IncludeGroupsFilterId.IsNull() && !groupFilter.IncludeGroupsFilterType.IsNull() {
			construct.SetStringProperty(groupFilter.IncludeGroupsFilterId,
				target.SetDeviceAndAppManagementAssignmentFilterId)

			err := construct.ParseEnum(groupFilter.IncludeGroupsFilterType,
				graphmodels.ParseDeviceAndAppManagementAssignmentFilterType,
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

func constructGroupExcludeAssignments(config *SettingsCatalogSettingsAssignmentResourceModel) []graphmodels.DeviceManagementConfigurationPolicyAssignmentable {
	var assignments []graphmodels.DeviceManagementConfigurationPolicyAssignmentable

	for _, groupId := range config.ExcludeGroupIds {
		// Only process if the group ID is not null or unknown
		if !groupId.IsNull() && !groupId.IsUnknown() {
			assignment := graphmodels.NewDeviceManagementConfigurationPolicyAssignment()
			target := graphmodels.NewExclusionGroupAssignmentTarget()

			// Use construct helper for setting the group ID
			construct.SetStringProperty(groupId, target.SetGroupId)

			assignment.SetTarget(target)
			assignments = append(assignments, assignment)
		}
	}

	return assignments
}

// logAssignmentDetails logs detailed information about the assignments configuration
func logAssignmentDetails(ctx context.Context, assignments *SettingsCatalogSettingsAssignmentResourceModel) {
	tflog.Debug(ctx, "Policy Assignment Configuration Details", map[string]interface{}{
		// All Devices fields
		"all_devices":             assignments.AllDevices.ValueBool(),
		"all_devices_filter_id":   assignments.AllDevicesFilterId.ValueString(),
		"all_devices_filter_type": assignments.AllDevicesFilterType.ValueString(),

		// All Users fields
		"all_users":             assignments.AllUsers.ValueBool(),
		"all_users_filter_id":   assignments.AllUsersFilterId.ValueString(),
		"all_users_filter_type": assignments.AllUsersFilterType.ValueString(),

		// Include Groups count
		"include_groups_count": len(assignments.IncludeGroups),

		// Exclude Groups count
		"exclude_groups_count": len(assignments.ExcludeGroupIds),
	})

	// Log each include group separately
	if len(assignments.IncludeGroups) > 0 {
		for i, group := range assignments.IncludeGroups {
			tflog.Debug(ctx, "Include Group Details", map[string]interface{}{
				"index":                      i,
				"group_id":                   group.GroupId.ValueString(),
				"include_groups_filter_id":   group.IncludeGroupsFilterId.ValueString(),
				"include_groups_filter_type": group.IncludeGroupsFilterType.ValueString(),
			})
		}
	}

	// Log exclude groups
	if len(assignments.ExcludeGroupIds) > 0 {
		excludeIds := make([]string, 0, len(assignments.ExcludeGroupIds))
		for _, id := range assignments.ExcludeGroupIds {
			excludeIds = append(excludeIds, id.ValueString())
		}
		tflog.Debug(ctx, "Exclude Group Details", map[string]interface{}{
			"exclude_group_ids": excludeIds,
		})
	}
}
