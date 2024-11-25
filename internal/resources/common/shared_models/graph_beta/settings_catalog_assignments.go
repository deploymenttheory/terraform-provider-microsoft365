package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

// SettingsCatalogSettingsAssignmentResourceModel struct to hold device configuation assignment configuration
type SettingsCatalogSettingsAssignmentResourceModel struct {
	AllDevices           types.Bool     `tfsdk:"all_devices"`
	AllDevicesFilterType types.String   `tfsdk:"all_devices_filter_type"`
	AllDevicesFilterId   types.String   `tfsdk:"all_devices_filter_id"`
	AllUsers             types.Bool     `tfsdk:"all_users"`
	AllUsersFilterType   types.String   `tfsdk:"all_users_filter_type"`
	AllUsersFilterId     types.String   `tfsdk:"all_users_filter_id"`
	IncludeGroups        []IncludeGroup `tfsdk:"include_groups"`
	ExcludeGroupIds      []types.String `tfsdk:"exclude_group_ids"`
}

// IncludeGroup represents a group with its corresponding filter type and filter group ID
type IncludeGroup struct {
	GroupId                 types.String `tfsdk:"group_id"`
	IncludeGroupsFilterType types.String `tfsdk:"include_groups_filter_type"`
	IncludeGroupsFilterId   types.String `tfsdk:"include_groups_filter_id"`
}

// DeviceManagementScriptAssignmentResourceModel struct to hold platform script assignment configuration
type DeviceManagementScriptAssignmentResourceModel struct {
	AllDevices      types.Bool     `tfsdk:"all_devices"`
	AllUsers        types.Bool     `tfsdk:"all_users"`
	IncludeGroupIds []types.String `tfsdk:"include_groups_ids"`
	ExcludeGroupIds []types.String `tfsdk:"exclude_group_ids"`
}
