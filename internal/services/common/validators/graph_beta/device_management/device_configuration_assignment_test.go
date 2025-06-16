package sharedValidators

import (
	"testing"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateDeviceConfiguationAssignmentSettings(t *testing.T) {
	tests := map[string]struct {
		config      *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel
		expectError bool
		errorMsg    string
	}{
		"valid_empty_config": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				AllDevices:           types.BoolValue(false),
				AllDevicesFilterType: types.StringNull(),
				AllDevicesFilterId:   types.StringNull(),
				AllUsers:             types.BoolValue(false),
				AllUsersFilterType:   types.StringNull(),
				AllUsersFilterId:     types.StringNull(),
				IncludeGroups:        []sharedmodels.IncludeGroup{},
				ExcludeGroupIds:      []types.String{},
			},
			expectError: false,
		},
		"invalid_all_devices_filter_type": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				AllDevices:           types.BoolValue(true),
				AllDevicesFilterType: types.StringValue("invalid"),
				AllDevicesFilterId:   types.StringNull(),
			},
			expectError: true,
			errorMsg:    "AllDevicesFilterType must be one of: include, exclude, none. Got: invalid",
		},
		"valid_all_devices_with_filter": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				AllDevices:           types.BoolValue(true),
				AllDevicesFilterType: types.StringValue("include"),
				AllDevicesFilterId:   types.StringValue("filter-id"),
			},
			expectError: false,
		},
		"invalid_all_devices_filter_without_id": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				AllDevices:           types.BoolValue(true),
				AllDevicesFilterType: types.StringValue("include"),
				AllDevicesFilterId:   types.StringNull(),
				AllUsers:             types.BoolValue(false),
				AllUsersFilterType:   types.StringNull(),
				AllUsersFilterId:     types.StringNull(),
				IncludeGroups:        []sharedmodels.IncludeGroup{},
				ExcludeGroupIds:      []types.String{},
			},
			expectError: true,
			errorMsg:    "all_devices_filter_type must be 'none' when all_devices is true and no filter ID is provided",
		},
		"valid_include_groups_alphanumeric": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				AllDevices:           types.BoolValue(false),
				AllUsers:             types.BoolValue(false),
				AllDevicesFilterType: types.StringNull(),
				AllDevicesFilterId:   types.StringNull(),
				AllUsersFilterType:   types.StringNull(),
				AllUsersFilterId:     types.StringNull(),
				IncludeGroups: []sharedmodels.IncludeGroup{
					{
						GroupId:                 types.StringValue("a1"),
						IncludeGroupsFilterType: types.StringValue("none"),
					},
					{
						GroupId:                 types.StringValue("b2"),
						IncludeGroupsFilterType: types.StringValue("none"),
					},
				},
			},
			expectError: false,
		},
		"invalid_include_groups_not_alphanumeric": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				IncludeGroups: []sharedmodels.IncludeGroup{
					{
						GroupId:                 types.StringValue("b2"),
						IncludeGroupsFilterType: types.StringValue("none"),
					},
					{
						GroupId:                 types.StringValue("a1"),
						IncludeGroupsFilterType: types.StringValue("none"),
					},
				},
			},
			expectError: true,
			errorMsg:    "include_groups must be in alphanumeric order by group_id. Found b2 before a1",
		},
		"invalid_duplicate_group_assignments": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				IncludeGroups: []sharedmodels.IncludeGroup{
					{
						GroupId:                 types.StringValue("same-id"),
						IncludeGroupsFilterType: types.StringValue("none"),
					},
				},
				ExcludeGroupIds: []types.String{
					types.StringValue("same-id"),
				},
			},
			expectError: true,
			errorMsg:    "group same-id is used in both include and exclude assignments. Each group assignment can only be used once across all assignment rules",
		},
		"invalid_all_devices_with_include_groups": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				AllDevices:           types.BoolValue(true),
				AllDevicesFilterType: types.StringValue("none"),
				AllDevicesFilterId:   types.StringNull(),
				IncludeGroups: []sharedmodels.IncludeGroup{
					{
						GroupId:                 types.StringValue("group1"),
						IncludeGroupsFilterType: types.StringValue("none"),
					},
				},
			},
			expectError: true,
			errorMsg:    "cannot assign to All Devices and Include Groups at the same time",
		},
		"valid_exclude_groups_alphanumeric": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				ExcludeGroupIds: []types.String{
					types.StringValue("a1"),
					types.StringValue("b2"),
					types.StringValue("c3"),
				},
			},
			expectError: false,
		},
		"invalid_exclude_groups_not_alphanumeric": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				ExcludeGroupIds: []types.String{
					types.StringValue("c3"),
					types.StringValue("a1"),
					types.StringValue("b2"),
				},
			},
			expectError: true,
			errorMsg:    "exclude_group_ids must be in alphanumeric order. Found c3 before a1",
		},
		"valid_include_groups_with_filter": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				IncludeGroups: []sharedmodels.IncludeGroup{
					{
						GroupId:                 types.StringValue("group1"),
						IncludeGroupsFilterType: types.StringValue("include"),
						IncludeGroupsFilterId:   types.StringValue("filter-id"),
					},
				},
			},
			expectError: false,
		},
		"invalid_include_groups_filter_type": {
			config: &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
				IncludeGroups: []sharedmodels.IncludeGroup{
					{
						GroupId:                 types.StringValue("group1"),
						IncludeGroupsFilterType: types.StringValue("invalid"),
					},
				},
			},
			expectError: true,
			errorMsg:    "IncludeGroupsFilterType must be one of: include, exclude, none. Got: invalid",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateDeviceConfiguationAssignmentSettings(tc.config)

			if tc.expectError && err == nil {
				t.Errorf("expected error containing %q, got no error", tc.errorMsg)
			}

			if !tc.expectError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.expectError && err != nil && tc.errorMsg != "" {
				if err.Error() != tc.errorMsg {
					t.Errorf("expected error message %q, got %q", tc.errorMsg, err.Error())
				}
			}
		})
	}
}
