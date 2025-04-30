package sharedStater

import (
	"context"
	"sort"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// StateConfigurationPolicyAssignment maps the remote policy assignment state to the Terraform state
func StateConfigurationPolicyAssignment(ctx context.Context, data *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel, assignmentsResponse models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable) {
	if assignmentsResponse == nil {
		tflog.Debug(ctx, "Assignments response is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map policy assignment to Terraform state")

	assignments := &sharedmodels.SettingsCatalogSettingsAssignmentResourceModel{
		AllDevices: types.BoolValue(false),
		AllUsers:   types.BoolValue(false),
	}

	allDeviceAssignments := GetAllDeviceAssignments(assignmentsResponse)
	MapAllDeviceAssignments(assignments, allDeviceAssignments)

	allUserAssignments := GetAllUserAssignments(assignmentsResponse)
	MapAllUserAssignments(assignments, allUserAssignments)

	includeGroupAssignments := GetIncludeGroupAssignments(assignmentsResponse)
	MapIncludeGroupAssignments(assignments, includeGroupAssignments)

	excludeGroupAssignments := GetExcludeGroupAssignments(assignmentsResponse)
	MapExcludeGroupAssignments(assignments, excludeGroupAssignments)

	//data.Assignments = assignments

	tflog.Debug(ctx, "Finished mapping assignment to Terraform state")
}

// MapAllDeviceAssignments maps the all devices assignment configuration to the assignments model
func MapAllDeviceAssignments(assignments *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel, allDeviceAssignments []models.DeviceManagementConfigurationPolicyAssignmentable) {
	if len(allDeviceAssignments) == 0 {
		return
	}

	assignments.AllDevices = types.BoolValue(true)

	if target := allDeviceAssignments[0].GetTarget(); target != nil {
		if filterId := target.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil && *filterId != "" {
			assignments.AllDevicesFilterId = types.StringValue(*filterId)
		} else {
			assignments.AllDevicesFilterId = types.StringNull()
		}

		if filterType := target.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil && assignments.AllDevicesFilterId.ValueString() != "" {
			assignments.AllDevicesFilterType = state.EnumPtrToTypeString(filterType)
		} else {
			assignments.AllDevicesFilterType = types.StringValue("none")
		}
	}
}

// MapAllUserAssignments maps the all users assignment configuration to the assignments model
func MapAllUserAssignments(assignments *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel, allUserAssignments []models.DeviceManagementConfigurationPolicyAssignmentable) {
	if len(allUserAssignments) == 0 {
		return
	}

	assignments.AllUsers = types.BoolValue(true)

	if target := allUserAssignments[0].GetTarget(); target != nil {
		if filterId := target.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil && *filterId != "" {
			assignments.AllUsersFilterId = types.StringValue(*filterId)
		} else {
			assignments.AllUsersFilterId = types.StringNull()
		}

		if filterType := target.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil && assignments.AllUsersFilterId.ValueString() != "" {
			assignments.AllUsersFilterType = state.EnumPtrToTypeString(filterType)
		} else {
			assignments.AllUsersFilterType = types.StringValue("none")
		}
	}
}

// MapIncludeGroupAssignments maps the include groups assignment configuration to the assignments model
func MapIncludeGroupAssignments(assignments *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel, includeGroupAssignments []models.DeviceManagementConfigurationPolicyAssignmentable) {
	if len(includeGroupAssignments) == 0 {
		return
	}

	assignments.IncludeGroups = make([]sharedmodels.IncludeGroup, 0, len(includeGroupAssignments))

	for _, assignment := range includeGroupAssignments {
		if target, ok := assignment.GetTarget().(models.GroupAssignmentTargetable); ok {
			includeGroup := sharedmodels.IncludeGroup{
				GroupId: types.StringValue(state.StringPtrToString(target.GetGroupId())),
			}

			if filterId := target.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil && *filterId != "" {
				includeGroup.IncludeGroupsFilterId = types.StringValue(*filterId)
				if filterType := target.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
					includeGroup.IncludeGroupsFilterType = state.EnumPtrToTypeString(filterType)
				}
			} else {
				includeGroup.IncludeGroupsFilterId = types.StringNull()
				includeGroup.IncludeGroupsFilterType = types.StringValue("none")
			}

			assignments.IncludeGroups = append(assignments.IncludeGroups, includeGroup)
		}
	}

	// Sort IncludeGroups by GroupId
	sort.Slice(assignments.IncludeGroups, func(i, j int) bool {
		return assignments.IncludeGroups[i].GroupId.ValueString() < assignments.IncludeGroups[j].GroupId.ValueString()
	})
}

// MapExcludeGroupAssignments maps the exclude groups assignment configuration to the assignments model
func MapExcludeGroupAssignments(assignments *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel, excludeGroupAssignments []models.DeviceManagementConfigurationPolicyAssignmentable) {
	if len(excludeGroupAssignments) == 0 {
		assignments.ExcludeGroupIds = nil
		return
	}

	excludeGroupIds := make([]types.String, 0, len(excludeGroupAssignments))
	for _, assignment := range excludeGroupAssignments {
		if target, ok := assignment.GetTarget().(models.GroupAssignmentTargetable); ok {
			if groupId := target.GetGroupId(); groupId != nil {
				excludeGroupIds = append(excludeGroupIds, types.StringValue(*groupId))
			}
		}
	}

	// Sort exclude group IDs alphanumerically
	sort.Slice(excludeGroupIds, func(i, j int) bool {
		return excludeGroupIds[i].ValueString() < excludeGroupIds[j].ValueString()
	})

	if len(excludeGroupIds) > 0 {
		assignments.ExcludeGroupIds = excludeGroupIds
	} else {
		assignments.ExcludeGroupIds = nil
	}
}

// Helpers

// GetAllDeviceAssignments retrieves all device assignments from the list of assignments
func GetAllDeviceAssignments(assignmentsResponse models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable) []models.DeviceManagementConfigurationPolicyAssignmentable {
	if assignmentsResponse == nil {
		return nil
	}

	var allDeviceAssignments []models.DeviceManagementConfigurationPolicyAssignmentable
	assignments := assignmentsResponse.GetValue()

	for _, assignment := range assignments {
		if target := assignment.GetTarget(); target != nil {
			if odataType := target.GetOdataType(); odataType != nil && *odataType == "#microsoft.graph.allDevicesAssignmentTarget" {
				allDeviceAssignments = append(allDeviceAssignments, assignment)
			}
		}
	}

	return allDeviceAssignments
}

// GetAllUserAssignments retrieves all user assignments from the list of assignments
func GetAllUserAssignments(assignmentsResponse models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable) []models.DeviceManagementConfigurationPolicyAssignmentable {
	if assignmentsResponse == nil {
		return nil
	}

	var allUserAssignments []models.DeviceManagementConfigurationPolicyAssignmentable
	assignments := assignmentsResponse.GetValue()

	for _, assignment := range assignments {
		if target := assignment.GetTarget(); target != nil {
			if odataType := target.GetOdataType(); odataType != nil && *odataType == "#microsoft.graph.allLicensedUsersAssignmentTarget" {
				allUserAssignments = append(allUserAssignments, assignment)
			}
		}
	}

	return allUserAssignments
}

// GetIncludeGroupAssignments retrieves include group assignments from the list of assignments
func GetIncludeGroupAssignments(assignmentsResponse models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable) []models.DeviceManagementConfigurationPolicyAssignmentable {
	if assignmentsResponse == nil {
		return nil
	}

	var includeGroupAssignments []models.DeviceManagementConfigurationPolicyAssignmentable
	assignments := assignmentsResponse.GetValue()

	for _, assignment := range assignments {
		if target := assignment.GetTarget(); target != nil {
			if odataType := target.GetOdataType(); odataType != nil && *odataType == "#microsoft.graph.groupAssignmentTarget" {
				includeGroupAssignments = append(includeGroupAssignments, assignment)
			}
		}
	}

	return includeGroupAssignments
}

// GetExcludeGroupAssignments retrieves exclude group assignments from the list of assignments
func GetExcludeGroupAssignments(assignmentsResponse models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable) []models.DeviceManagementConfigurationPolicyAssignmentable {
	if assignmentsResponse == nil {
		return nil
	}

	var excludeGroupAssignments []models.DeviceManagementConfigurationPolicyAssignmentable
	assignments := assignmentsResponse.GetValue()

	for _, assignment := range assignments {
		if target := assignment.GetTarget(); target != nil {
			if odataType := target.GetOdataType(); odataType != nil && *odataType == "#microsoft.graph.exclusionGroupAssignmentTarget" {
				excludeGroupAssignments = append(excludeGroupAssignments, assignment)
			}
		}
	}

	return excludeGroupAssignments
}
