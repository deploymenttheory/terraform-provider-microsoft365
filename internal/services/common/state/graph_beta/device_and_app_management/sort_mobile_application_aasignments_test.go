package sharedStater

import (
	"testing"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestSortMobileAppAssignments tests the SortMobileAppAssignments function
func TestSortMobileAppAssignments(t *testing.T) {
	tests := []struct {
		name              string
		assignments       []sharedmodels.MobileAppAssignmentResourceModel
		expectedOrder     []string // Expected intent order after sort
		expectedGroupIds  []string // Expected group IDs after sort (for group assignments)
		validateFunc      func(t *testing.T, sorted []sharedmodels.MobileAppAssignmentResourceModel)
	}{
		{
			name:        "Empty assignments list",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{},
			validateFunc: func(t *testing.T, sorted []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, sorted, 0)
			},
		},
		{
			name: "Sort by intent priority (required, available, uninstall)",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{
				{
					Intent: types.StringValue("uninstall"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType: types.StringValue("allDevices"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType: types.StringValue("allDevices"),
					},
				},
				{
					Intent: types.StringValue("available"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType: types.StringValue("allDevices"),
					},
				},
			},
			expectedOrder: []string{"required", "available", "uninstall"},
			validateFunc: func(t *testing.T, sorted []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, sorted, 3)
				assert.Equal(t, "required", sorted[0].Intent.ValueString())
				assert.Equal(t, "available", sorted[1].Intent.ValueString())
				assert.Equal(t, "uninstall", sorted[2].Intent.ValueString())
			},
		},
		{
			name: "Sort by target type priority (groupAssignment, exclusionGroupAssignment, allLicensedUsers, allDevices)",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType: types.StringValue("allDevices"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType: types.StringValue("groupAssignment"),
						GroupId:    types.StringValue("group-123"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType: types.StringValue("allLicensedUsers"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType: types.StringValue("exclusionGroupAssignment"),
						GroupId:    types.StringValue("group-456"),
					},
				},
			},
			validateFunc: func(t *testing.T, sorted []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, sorted, 4)
				assert.Equal(t, "groupAssignment", sorted[0].Target.TargetType.ValueString())
				assert.Equal(t, "exclusionGroupAssignment", sorted[1].Target.TargetType.ValueString())
				assert.Equal(t, "allLicensedUsers", sorted[2].Target.TargetType.ValueString())
				assert.Equal(t, "allDevices", sorted[3].Target.TargetType.ValueString())
			},
		},
		{
			name: "Sort by filter type (exclude, include, none)",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringValue("group-123"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringValue("group-123"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("include"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringValue("group-123"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("exclude"),
					},
				},
			},
			validateFunc: func(t *testing.T, sorted []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, sorted, 3)
				assert.Equal(t, "exclude", sorted[0].Target.DeviceAndAppManagementAssignmentFilterType.ValueString())
				assert.Equal(t, "include", sorted[1].Target.DeviceAndAppManagementAssignmentFilterType.ValueString())
				assert.Equal(t, "none", sorted[2].Target.DeviceAndAppManagementAssignmentFilterType.ValueString())
			},
		},
		{
			name: "Sort group assignments by group ID",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringValue("group-zzz"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringValue("group-aaa"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringValue("group-mmm"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
					},
				},
			},
			expectedGroupIds: []string{"group-aaa", "group-mmm", "group-zzz"},
			validateFunc: func(t *testing.T, sorted []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, sorted, 3)
				assert.Equal(t, "group-aaa", sorted[0].Target.GroupId.ValueString())
				assert.Equal(t, "group-mmm", sorted[1].Target.GroupId.ValueString())
				assert.Equal(t, "group-zzz", sorted[2].Target.GroupId.ValueString())
			},
		},
		{
			name: "Complex multi-tier sort (intent > target type > filter type > group ID)",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{
				{
					Intent: types.StringValue("available"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("allDevices"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringValue("group-zzz"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("include"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringValue("group-aaa"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("exclude"),
					},
				},
				{
					Intent: types.StringValue("uninstall"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("allLicensedUsers"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("exclusionGroupAssignment"),
						GroupId:                                    types.StringValue("group-bbb"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
					},
				},
			},
			validateFunc: func(t *testing.T, sorted []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, sorted, 5)
				
				// First by intent: required comes first
				assert.Equal(t, "required", sorted[0].Intent.ValueString())
				assert.Equal(t, "required", sorted[1].Intent.ValueString())
				assert.Equal(t, "required", sorted[2].Intent.ValueString())
				assert.Equal(t, "available", sorted[3].Intent.ValueString())
				assert.Equal(t, "uninstall", sorted[4].Intent.ValueString())
				
				// Within required: groupAssignment, then exclusionGroupAssignment
				assert.Equal(t, "groupAssignment", sorted[0].Target.TargetType.ValueString())
				assert.Equal(t, "groupAssignment", sorted[1].Target.TargetType.ValueString())
				assert.Equal(t, "exclusionGroupAssignment", sorted[2].Target.TargetType.ValueString())
				
				// Within groupAssignments: exclude before include
				assert.Equal(t, "exclude", sorted[0].Target.DeviceAndAppManagementAssignmentFilterType.ValueString())
				assert.Equal(t, "include", sorted[1].Target.DeviceAndAppManagementAssignmentFilterType.ValueString())
				
				// Group IDs sorted alphabetically
				assert.Equal(t, "group-aaa", sorted[0].Target.GroupId.ValueString())
				assert.Equal(t, "group-zzz", sorted[1].Target.GroupId.ValueString())
			},
		},
		{
			name: "Single assignment",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType: types.StringValue("allDevices"),
					},
				},
			},
			validateFunc: func(t *testing.T, sorted []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, sorted, 1)
				assert.Equal(t, "required", sorted[0].Intent.ValueString())
			},
		},
		{
			name: "Assignments with null group IDs",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringNull(),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
					},
				},
				{
					Intent: types.StringValue("required"),
					Target: sharedmodels.AssignmentTargetResourceModel{
						TargetType:                                 types.StringValue("groupAssignment"),
						GroupId:                                    types.StringValue("group-123"),
						DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
					},
				},
			},
			validateFunc: func(t *testing.T, sorted []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, sorted, 2)
				// Should handle null gracefully without panicking
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the original
			assignmentsCopy := make([]sharedmodels.MobileAppAssignmentResourceModel, len(tt.assignments))
			copy(assignmentsCopy, tt.assignments)
			
			SortMobileAppAssignments(assignmentsCopy)
			
			if tt.validateFunc != nil {
				tt.validateFunc(t, assignmentsCopy)
			}
		})
	}
}
