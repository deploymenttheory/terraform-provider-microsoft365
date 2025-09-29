package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapAssignmentsToState maps the assignments response to the state model.
func mapAssignmentsToState(ctx context.Context, stateModel *WindowsAutopilotDevicePreparationPolicyResourceModel, assignmentsResponse models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable) {
	if assignmentsResponse == nil {
		return
	}

	assignments := assignmentsResponse.GetValue()
	if len(assignments) == 0 {
		return
	}

	// Use the shared updateStateWithAssignments function to set the assignments
	updateStateWithAssignments(ctx, stateModel, assignments)
}

// updateStateWithAssignments updates the assignments in the state model
func updateStateWithAssignments(ctx context.Context, stateModel *WindowsAutopilotDevicePreparationPolicyResourceModel, assignments []models.DeviceManagementConfigurationPolicyAssignmentable) {
	if len(assignments) == 0 {
		return
	}

	var includeGroupIds []types.String

	for _, assignment := range assignments {
		assignmentIncludeGroupIds, _, _, err := extractAssignmentData(assignment)
		if err != nil {
			tflog.Error(ctx, "Failed to extract assignment data", map[string]any{
				"error": err.Error(),
			})
			continue
		}

		includeGroupIds = append(includeGroupIds, assignmentIncludeGroupIds...)
	}

	// Initialize the Assignments if nil
	if stateModel.Assignments == nil {
		stateModel.Assignments = &WindowsAutopilotDevicePreparationAssignment{}
	}

	stateModel.Assignments.IncludeGroupIds = includeGroupIds
}

// extractAssignmentData extracts assignment data from the policy response
func extractAssignmentData(assignment models.DeviceManagementConfigurationPolicyAssignmentable) ([]types.String, []types.String, []types.String, error) {
	var includeGroupIds []types.String
	var excludeGroupIds []types.String
	var includeAllUsers []types.String

	if assignment == nil {
		return includeGroupIds, excludeGroupIds, includeAllUsers, fmt.Errorf("assignment is nil")
	}

	assignmentTarget := assignment.GetTarget()
	if assignmentTarget == nil {
		return includeGroupIds, excludeGroupIds, includeAllUsers, fmt.Errorf("assignment target is nil")
	}

	// Get the OData type to determine the assignment target type
	targetType := ""
	if assignmentTarget.GetOdataType() != nil {
		targetType = *assignmentTarget.GetOdataType()
	}

	// Handle different assignment target types based on OData type
	switch targetType {
	case "#microsoft.graph.groupAssignmentTarget":
		// Handle group assignment target
		groupTarget, ok := assignmentTarget.(*models.GroupAssignmentTarget)
		if ok && groupTarget.GetGroupId() != nil {
			includeGroupIds = append(includeGroupIds, types.StringValue(*groupTarget.GetGroupId()))
		}
	case "#microsoft.graph.exclusionGroupAssignmentTarget":
		// Handle exclusion group assignment target
		exclusionTarget, ok := assignmentTarget.(*models.ExclusionGroupAssignmentTarget)
		if ok && exclusionTarget.GetGroupId() != nil {
			excludeGroupIds = append(excludeGroupIds, types.StringValue(*exclusionTarget.GetGroupId()))
		}
	case "#microsoft.graph.allLicensedUsersAssignmentTarget":
		// Handle all licensed users assignment target
		includeAllUsers = append(includeAllUsers, types.StringValue("All Users"))
	case "#microsoft.graph.allDevicesAssignmentTarget":
		// Handle all devices assignment target
		includeAllUsers = append(includeAllUsers, types.StringValue("All Devices"))
	default:
		// Try to extract group ID from additionalData if available
		additionalData := assignmentTarget.GetAdditionalData()
		if additionalData != nil {
			if groupId, ok := additionalData["groupId"].(string); ok && groupId != "" {
				includeGroupIds = append(includeGroupIds, types.StringValue(groupId))
			}
		}
	}

	return includeGroupIds, excludeGroupIds, includeAllUsers, nil
}
