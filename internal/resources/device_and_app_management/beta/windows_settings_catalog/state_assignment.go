package graphBetaWindowsSettingsCatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteAssignmentStateToTerraform(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel, assignmentsResponse graphmodels.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable) {
	if assignmentsResponse == nil {
		tflog.Debug(ctx, "Assignments response is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map assignment state to Terraform state")

	assignments := &SettingsCatalogSettingsAssignmentResourceModel{
		AllDevices: types.BoolValue(false),
		AllUsers:   types.BoolValue(false),
	}

	// Map All Devices assignments
	allDeviceAssignments := GetAllDeviceAssignments(assignmentsResponse)
	if len(allDeviceAssignments) > 0 {
		assignments.AllDevices = types.BoolValue(true)
		if target := allDeviceAssignments[0].GetTarget(); target != nil {
			assignments.AllDevicesFilterId = types.StringValue(state.StringPtrToString(target.GetDeviceAndAppManagementAssignmentFilterId()))
			assignments.AllDevicesFilterType = state.EnumPtrToTypeString(target.GetDeviceAndAppManagementAssignmentFilterType())
		}
	}

	// Map All Users assignments
	allUserAssignments := GetAllUserAssignments(assignmentsResponse)
	if len(allUserAssignments) > 0 {
		assignments.AllUsers = types.BoolValue(true)
		if target := allUserAssignments[0].GetTarget(); target != nil {
			assignments.AllUsersFilterId = types.StringValue(state.StringPtrToString(target.GetDeviceAndAppManagementAssignmentFilterId()))
			assignments.AllUsersFilterType = state.EnumPtrToTypeString(target.GetDeviceAndAppManagementAssignmentFilterType())
		}
	}

	// Map Include Group assignments
	includeGroupAssignments := GetIncludeGroupAssignments(assignmentsResponse)
	if len(includeGroupAssignments) > 0 {
		assignments.IncludeGroups = make([]IncludeGroup, 0, len(includeGroupAssignments))
		for _, assignment := range includeGroupAssignments {
			if target, ok := assignment.GetTarget().(graphmodels.GroupAssignmentTargetable); ok {
				includeGroup := IncludeGroup{
					GroupId: types.StringValue(state.StringPtrToString(target.GetGroupId())),
				}

				if filterId := target.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil {
					includeGroup.IncludeGroupsFilterId = types.StringValue(state.StringPtrToString(filterId))
					includeGroup.IncludeGroupsFilterType = state.EnumPtrToTypeString(target.GetDeviceAndAppManagementAssignmentFilterType())
				}

				assignments.IncludeGroups = append(assignments.IncludeGroups, includeGroup)
			}
		}
	}

	// Map Exclude Group assignments
	excludeGroupAssignments := GetExcludeGroupAssignments(assignmentsResponse)
	if len(excludeGroupAssignments) > 0 {
		assignments.ExcludeGroupIds = make([]types.String, 0, len(excludeGroupAssignments))
		for _, assignment := range excludeGroupAssignments {
			if target, ok := assignment.GetTarget().(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				if groupId := target.GetGroupId(); groupId != nil {
					assignments.ExcludeGroupIds = append(assignments.ExcludeGroupIds, types.StringValue(*groupId))
				}
			}
		}
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
