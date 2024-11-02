package graphBetaSettingsCatalog

import (
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
)

// ValidateAssignmentConfiguration validates the assignment configuration
// - if all_devices is false, device filter settings should not be set
// - if all_users is false, user filter settings should not be set
// - if all_devices is true and all_devices_filter_id is nil, then all_devices_filter_type must set to "none"
// - if all_users is true and all_users_filter_id is nil, then all_users_filter_type must set to "none"
// - if all_devices is true and all_devices_filter_id is set, then all_devices_filter_type must be either "include" or "exclude"
// - if all_users is true and all_users_filter_id is set, then all_users_filter_type must be either "include" or "exclude"
// - if group_id is provided in include_groups, include_groups_filter_type must be set
// - if group_id is provided in include_groups and include_groups_filter_id is nil, then include_groups_filter_type must be "none"
// - if group_id is provided in include_groups and include_groups_filter_id is set, then include_groups_filter_type must be either "include" or "exclude"
// - AllDevicesFilterType must be one of: include, exclude, none
// - AllUsersFilterType must be one of: include, exclude, none
// - IncludeGroupsFilterType must be one of: include, exclude, none
// - ExcludeGroupIds must be in alphanumeric order by group_id
// - IncludeGroups must be in alphanumeric order by group_id
// - No group is used more than once across include and exclude assignments
// - AllDevices and IncludeGroups cannot be used at the same time
// - AllUsers and IncludeGroups cannot be used at the same time
func validateAssignmentConfig(config *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel) error {
	// Validate filter types have valid values
	validFilterTypes := map[string]bool{
		"include": true,
		"exclude": true,
		"none":    true,
	}

	// Validate filter type values
	if !config.AllDevicesFilterType.IsNull() && !validFilterTypes[config.AllDevicesFilterType.ValueString()] {
		return fmt.Errorf("AllDevicesFilterType must be one of: include, exclude, none. Got: %s",
			config.AllDevicesFilterType.ValueString())
	}

	if !config.AllUsersFilterType.IsNull() && !validFilterTypes[config.AllUsersFilterType.ValueString()] {
		return fmt.Errorf("AllUsersFilterType must be one of: include, exclude, none. Got: %s",
			config.AllUsersFilterType.ValueString())
	}

	// Validate all_devices related settings
	if !config.AllDevices.ValueBool() {
		if !config.AllDevicesFilterType.IsNull() || !config.AllDevicesFilterId.IsNull() {
			return fmt.Errorf("all_devices_filter_type and/or all_devices_filter_id cannot be set when all_devices is false")
		}
	} else {
		// all_devices is true
		if config.AllDevicesFilterId.IsNull() {
			// No filter ID provided, filter type MUST be "none"
			if config.AllDevicesFilterType.IsNull() || config.AllDevicesFilterType.ValueString() != "none" {
				return fmt.Errorf("all_devices_filter_type must be 'none' when all_devices is true and no filter ID is provided")
			}
		} else {
			// Filter ID is provided, filter type must be either include or exclude
			if config.AllDevicesFilterType.IsNull() || config.AllDevicesFilterType.ValueString() == "none" {
				return fmt.Errorf("all_devices_filter_type must be either 'include' or 'exclude' when filter ID is provided")
			}
		}
	}

	// Validate all_users related settings
	if !config.AllUsers.ValueBool() {
		if !config.AllUsersFilterType.IsNull() || !config.AllUsersFilterId.IsNull() {
			return fmt.Errorf("all_users_filter_type and/or all_users_filter_id cannot be set when all_users is false")
		}
	} else {
		// all_users is true
		if config.AllUsersFilterId.IsNull() {
			// No filter ID provided, filter type MUST be "none"
			if config.AllUsersFilterType.IsNull() || config.AllUsersFilterType.ValueString() != "none" {
				return fmt.Errorf("all_users_filter_type must be 'none' when all_users is true and no filter ID is provided")
			}
		} else {
			// Filter ID is provided, filter type must be either include or exclude
			if config.AllUsersFilterType.IsNull() || config.AllUsersFilterType.ValueString() == "none" {
				return fmt.Errorf("all_users_filter_type must be either 'include' or 'exclude' when filter ID is provided")
			}
		}
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

	// Validate include_groups settings
	for _, group := range config.IncludeGroups {
		// Validate that filter_type is set if group_id is provided
		if !group.GroupId.IsNull() && !group.GroupId.IsUnknown() && group.GroupId.ValueString() != "" {
			if group.IncludeGroupsFilterType.IsNull() {
				return fmt.Errorf("include_groups_filter_type must be set when group_id is provided for group %s",
					group.GroupId.ValueString())
			}

			// Validate filter type value
			if !validFilterTypes[group.IncludeGroupsFilterType.ValueString()] {
				return fmt.Errorf("IncludeGroupsFilterType must be one of: include, exclude, none. Got: %s",
					group.IncludeGroupsFilterType.ValueString())
			}

			// Apply filter type logic based on filter ID presence
			if group.IncludeGroupsFilterId.IsNull() {
				// No filter ID provided, filter type MUST be "none"
				if group.IncludeGroupsFilterType.ValueString() != "none" {
					return fmt.Errorf("include_groups_filter_type must be 'none' when no filter ID is provided for group %s",
						group.GroupId.ValueString())
				}
			} else {
				// Filter ID is provided, filter type must be either include or exclude
				if group.IncludeGroupsFilterType.ValueString() == "none" {
					return fmt.Errorf("include_groups_filter_type cannot be 'none' when filter ID is provided for group %s",
						group.GroupId.ValueString())
				}
			}
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
