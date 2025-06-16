package graphBetaMacOSPlatformScript

import (
	"context"
	"sort"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps the remote policy assignment state to the Terraform state
func MapRemoteAssignmentStateToTerraform(ctx context.Context, data *MacOSPlatformScriptResourceModel, assignmentsResponse models.DeviceManagementScriptAssignmentCollectionResponseable) {
	if assignmentsResponse == nil {
		tflog.Debug(ctx, "Assignments response is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map policy assignment to Terraform state")

	assignments := &sharedmodels.DeviceManagementScriptAssignmentResourceModel{
		AllDevices: types.BoolValue(false),
		AllUsers:   types.BoolValue(false),
	}

	// Map All Devices assignments
	allDeviceAssignments := GetAllDeviceAssignments(assignmentsResponse)
	MapAllDeviceAssignments(assignments, allDeviceAssignments)

	// Map All Users assignments
	allUserAssignments := GetAllUserAssignments(assignmentsResponse)
	MapAllUserAssignments(assignments, allUserAssignments)

	// Map Include Group assignments
	includeGroupAssignments := GetIncludeGroupAssignments(assignmentsResponse)
	MapIncludeGroupAssignments(assignments, includeGroupAssignments)

	// Map Exclude Group assignments
	excludeGroupAssignments := GetExcludeGroupAssignments(assignmentsResponse)
	MapExcludeGroupAssignments(assignments, excludeGroupAssignments)

	data.Assignments = assignments

	tflog.Debug(ctx, "Finished mapping assignment to Terraform state")
}

// MapAllDeviceAssignments maps the all devices assignment configuration to the assignments model
func MapAllDeviceAssignments(assignments *sharedmodels.DeviceManagementScriptAssignmentResourceModel, allDeviceAssignments []models.DeviceManagementScriptAssignmentable) {
	if len(allDeviceAssignments) > 0 {
		assignments.AllDevices = types.BoolValue(true)
	}
}

// MapAllUserAssignments maps the all users assignment configuration to the assignments model
func MapAllUserAssignments(assignments *sharedmodels.DeviceManagementScriptAssignmentResourceModel, allUserAssignments []models.DeviceManagementScriptAssignmentable) {
	if len(allUserAssignments) > 0 {
		assignments.AllUsers = types.BoolValue(true)
	}
}

// MapIncludeGroupAssignments maps the include groups assignment configuration to the assignments model
func MapIncludeGroupAssignments(assignments *sharedmodels.DeviceManagementScriptAssignmentResourceModel, includeGroupAssignments []models.DeviceManagementScriptAssignmentable) {
	if len(includeGroupAssignments) == 0 {
		return
	}

	includeGroupIds := make([]types.String, 0)
	for _, assignment := range includeGroupAssignments {
		if target, ok := assignment.GetTarget().(models.GroupAssignmentTargetable); ok {
			if groupId := target.GetGroupId(); groupId != nil {
				includeGroupIds = append(includeGroupIds, types.StringValue(*groupId))
			}
		}
	}

	// Sort include group IDs alphanumerically
	sort.Slice(includeGroupIds, func(i, j int) bool {
		return includeGroupIds[i].ValueString() < includeGroupIds[j].ValueString()
	})

	assignments.IncludeGroupIds = includeGroupIds
}

// MapExcludeGroupAssignments maps the exclude groups assignment configuration to the assignments model
func MapExcludeGroupAssignments(assignments *sharedmodels.DeviceManagementScriptAssignmentResourceModel, excludeGroupAssignments []models.DeviceManagementScriptAssignmentable) {
	if len(excludeGroupAssignments) == 0 {
		return
	}

	excludeGroupIds := make([]types.String, 0)
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

	assignments.ExcludeGroupIds = excludeGroupIds
}

// Helpers

// GetAllDeviceAssignments retrieves all device assignments from the list of assignments
func GetAllDeviceAssignments(assignmentsResponse models.DeviceManagementScriptAssignmentCollectionResponseable) []models.DeviceManagementScriptAssignmentable {
	if assignmentsResponse == nil {
		return nil
	}

	var allDeviceAssignments []models.DeviceManagementScriptAssignmentable
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
func GetAllUserAssignments(assignmentsResponse models.DeviceManagementScriptAssignmentCollectionResponseable) []models.DeviceManagementScriptAssignmentable {
	if assignmentsResponse == nil {
		return nil
	}

	var allUserAssignments []models.DeviceManagementScriptAssignmentable
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
func GetIncludeGroupAssignments(assignmentsResponse models.DeviceManagementScriptAssignmentCollectionResponseable) []models.DeviceManagementScriptAssignmentable {
	if assignmentsResponse == nil {
		return nil
	}

	var includeGroupAssignments []models.DeviceManagementScriptAssignmentable
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
func GetExcludeGroupAssignments(assignmentsResponse models.DeviceManagementScriptAssignmentCollectionResponseable) []models.DeviceManagementScriptAssignmentable {
	if assignmentsResponse == nil {
		return nil
	}

	var excludeGroupAssignments []models.DeviceManagementScriptAssignmentable
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
