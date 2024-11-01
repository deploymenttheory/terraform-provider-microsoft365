package graphBetaWindowsSettingsCatalog

import (
	"context"
	"sort"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps the remote policy assignment state to the Terraform state
func MapRemoteAssignmentStateToTerraform(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel, assignmentsResponse models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable) {
	if assignmentsResponse == nil {
		tflog.Debug(ctx, "Assignments response is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map policy assignment to Terraform state")

	assignments := &SettingsCatalogSettingsAssignmentResourceModel{
		AllDevices: types.BoolValue(false),
		AllUsers:   types.BoolValue(false),
	}

	// Map All Devices assignments
	allDeviceAssignments := GetAllDeviceAssignments(assignmentsResponse)
	if len(allDeviceAssignments) > 0 {
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

	// Map All Users assignments
	allUserAssignments := GetAllUserAssignments(assignmentsResponse)
	if len(allUserAssignments) > 0 {
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

	// Map Include Group assignments
	includeGroupAssignments := GetIncludeGroupAssignments(assignmentsResponse)
	if len(includeGroupAssignments) > 0 {
		assignments.IncludeGroups = make([]IncludeGroup, 0, len(includeGroupAssignments))
		for _, assignment := range includeGroupAssignments {
			if target, ok := assignment.GetTarget().(models.GroupAssignmentTargetable); ok {
				includeGroup := IncludeGroup{
					GroupId: types.StringValue(state.StringPtrToString(target.GetGroupId())),
				}

				// Handle filter ID and type
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

	// Map Exclude Group assignments
	excludeGroupAssignments := GetExcludeGroupAssignments(assignmentsResponse)
	if len(excludeGroupAssignments) > 0 {
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
	} else {
		assignments.ExcludeGroupIds = nil
	}

	data.Assignments = assignments

	tflog.Debug(ctx, "Finished mapping assignment to Terraform state")
}

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
