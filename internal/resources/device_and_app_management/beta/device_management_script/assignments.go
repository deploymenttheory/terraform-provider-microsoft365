package graphBetaDeviceManagementScript

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure DeviceManagementScriptAssignmentResourceModel implements the Assignment interface
var _ crud.Assignment = (*DeviceManagementScriptAssignmentResourceModel)(nil)

// Implement the necessary methods for DeviceManagementScriptAssignmentResourceModel
func (d DeviceManagementScriptAssignmentResourceModel) GetID() types.String {
	return d.ID
}

func (d DeviceManagementScriptAssignmentResourceModel) GetTarget() interface{} {
	return d.Target
}

// Custom comparison function for DeviceManagementScriptAssignmentResourceModel
func compareDeviceManagementScriptAssignments(a, b interface{}) bool {
	aAssign, aOk := a.(DeviceManagementScriptAssignmentResourceModel)
	bAssign, bOk := b.(DeviceManagementScriptAssignmentResourceModel)
	if !aOk || !bOk {
		return false
	}

	// Compare relevant fields of the target
	return aAssign.Target.DeviceAndAppManagementAssignmentFilterId == bAssign.Target.DeviceAndAppManagementAssignmentFilterId &&
		aAssign.Target.DeviceAndAppManagementAssignmentFilterType == bAssign.Target.DeviceAndAppManagementAssignmentFilterType &&
		aAssign.Target.TargetType == bAssign.Target.TargetType &&
		aAssign.Target.EntraObjectId == bAssign.Target.EntraObjectId
}

// Helper functions for assignments
func assignmentExistsInPlan(assignment DeviceManagementScriptAssignmentResourceModel, planAssignments []DeviceManagementScriptAssignmentResourceModel) bool {
	return crud.AssignmentExistsInSlice(assignment, planAssignments, compareDeviceManagementScriptAssignments)
}

func assignmentExistsInState(assignment DeviceManagementScriptAssignmentResourceModel, stateAssignments []DeviceManagementScriptAssignmentResourceModel) bool {
	return crud.AssignmentExistsInSlice(assignment, stateAssignments, compareDeviceManagementScriptAssignments)
}

// Helper functions for group assignments
func groupAssignmentExistsInPlan(groupAssignment DeviceManagementScriptGroupAssignmentResourceModel, planGroupAssignments []DeviceManagementScriptGroupAssignmentResourceModel) bool {
	return crud.ExistsInSlice(groupAssignment, planGroupAssignments, func(a, b interface{}) bool {
		aAssign := a.(DeviceManagementScriptGroupAssignmentResourceModel)
		bAssign := b.(DeviceManagementScriptGroupAssignmentResourceModel)
		return aAssign.ID == bAssign.ID
	})
}

func groupAssignmentExistsInState(groupAssignment DeviceManagementScriptGroupAssignmentResourceModel, stateGroupAssignments []DeviceManagementScriptGroupAssignmentResourceModel) bool {
	return crud.ExistsInSlice(groupAssignment, stateGroupAssignments, func(a, b interface{}) bool {
		aAssign := a.(DeviceManagementScriptGroupAssignmentResourceModel)
		bAssign := b.(DeviceManagementScriptGroupAssignmentResourceModel)
		return aAssign.ID == bAssign.ID
	})
}
