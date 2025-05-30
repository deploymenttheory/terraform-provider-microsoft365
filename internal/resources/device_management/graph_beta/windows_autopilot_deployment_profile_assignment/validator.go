package graphBetaWindowsAutopilotDeploymentProfileAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ValidateWindowsAutopilotDeploymentProfileAssignment is the main validation function that orchestrates all validation checks
func ValidateWindowsAutopilotDeploymentProfileAssignment(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	object WindowsAutopilotDeploymentProfileAssignmentResourceModel,
	isUpdate bool,
) error {
	tflog.Debug(ctx, "Starting validation for Windows Autopilot Deployment Profile Assignment")

	// Get existing assignments for validation
	existingAssignments, err := getExistingAssignments(ctx, client, object.WindowsAutopilotDeploymentProfileId.ValueString())
	if err != nil {
		return fmt.Errorf("failed to retrieve existing assignments for validation: %w", err)
	}

	// Run validation checks
	if err := validateAllDevicesUniqueness(ctx, object, existingAssignments, isUpdate); err != nil {
		return err
	}

	// Validate that allDevices assignments are exclusive with other assignment types
	if err := validateAllDevicesExclusivity(ctx, object, existingAssignments, isUpdate); err != nil {
		return err
	}

	// Validate that group IDs are unique across all assignments
	if err := validateGroupIdUniqueness(ctx, object, existingAssignments, isUpdate); err != nil {
		return err
	}

	tflog.Debug(ctx, "All validation checks passed for Windows Autopilot Deployment Profile Assignment")
	return nil
}

// getExistingAssignments retrieves all current assignments for the Windows Autopilot Deployment Profile
func getExistingAssignments(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	profileId string,
) ([]graphmodels.WindowsAutopilotDeploymentProfileAssignmentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Retrieving existing assignments for profile: %s", profileId))

	assignments, err := client.
		DeviceManagement().
		WindowsAutopilotDeploymentProfiles().
		ByWindowsAutopilotDeploymentProfileId(profileId).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to get existing assignments: %w", err)
	}

	if assignments == nil || assignments.GetValue() == nil {
		tflog.Debug(ctx, "No existing assignments found")
		return []graphmodels.WindowsAutopilotDeploymentProfileAssignmentable{}, nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Found %d existing assignments", len(assignments.GetValue())))
	return assignments.GetValue(), nil
}

// validateAllDevicesUniqueness ensures that only one assignment can target "All Devices"
func validateAllDevicesUniqueness(
	ctx context.Context,
	object WindowsAutopilotDeploymentProfileAssignmentResourceModel,
	existingAssignments []graphmodels.WindowsAutopilotDeploymentProfileAssignmentable,
	isUpdate bool,
) error {
	tflog.Debug(ctx, "Validating All Devices uniqueness constraint")

	// Check if the current object targets all devices
	isCurrentAllDevices := object.Target.TargetType.ValueString() == "allDevices"

	if !isCurrentAllDevices {
		tflog.Debug(ctx, "Current assignment does not target all devices, skipping validation")
		return nil
	}

	// Count existing "All Devices" assignments
	allDevicesCount := 0
	var existingAllDevicesAssignment graphmodels.WindowsAutopilotDeploymentProfileAssignmentable

	for _, assignment := range existingAssignments {
		if assignment == nil {
			continue
		}

		target := assignment.GetTarget()
		if target == nil {
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			continue
		}

		// Check if this assignment targets all devices
		if *odataType == "#microsoft.graph.allDevicesAssignmentTarget" {
			allDevicesCount++
			existingAllDevicesAssignment = assignment

			// For updates, check if this is the same assignment being updated
			if isUpdate && assignment.GetId() != nil && !object.ID.IsNull() {
				if *assignment.GetId() == object.ID.ValueString() {
					tflog.Debug(ctx, "Found existing all devices assignment, but it's the same one being updated")
					allDevicesCount-- // Don't count the assignment being updated
				}
			}
		}
	}

	// If we found existing "All Devices" assignments that aren't the current one
	if allDevicesCount > 0 {
		assignmentId := "unknown"
		if existingAllDevicesAssignment != nil && existingAllDevicesAssignment.GetId() != nil {
			assignmentId = *existingAllDevicesAssignment.GetId()
		}

		return fmt.Errorf(
			"cannot assign Windows Autopilot Deployment Profile to 'All Devices' because another assignment (ID: %s) already targets 'All Devices'. "+
				"Only one assignment per profile can target 'All Devices'",
			assignmentId,
		)
	}

	tflog.Debug(ctx, "All Devices uniqueness validation passed")
	return nil
}

// validateAllDevicesExclusivity ensures that allDevices assignments are exclusive with other assignment types
func validateAllDevicesExclusivity(
	ctx context.Context,
	object WindowsAutopilotDeploymentProfileAssignmentResourceModel,
	existingAssignments []graphmodels.WindowsAutopilotDeploymentProfileAssignmentable,
	isUpdate bool,
) error {
	tflog.Debug(ctx, "Validating All Devices exclusivity constraint")

	// Check if the current object targets all devices
	isCurrentAllDevices := object.Target.TargetType.ValueString() == "allDevices"

	// Count existing assignments by type
	allDevicesCount := 0
	otherAssignmentsCount := 0
	var conflictingAssignment graphmodels.WindowsAutopilotDeploymentProfileAssignmentable

	for _, assignment := range existingAssignments {
		if assignment == nil {
			continue
		}

		// Skip the assignment being updated to avoid false conflicts
		if isUpdate && assignment.GetId() != nil && !object.ID.IsNull() {
			if *assignment.GetId() == object.ID.ValueString() {
				continue
			}
		}

		target := assignment.GetTarget()
		if target == nil {
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			continue
		}

		// Count assignment types
		if *odataType == "#microsoft.graph.allDevicesAssignmentTarget" {
			allDevicesCount++
			conflictingAssignment = assignment
		} else {
			otherAssignmentsCount++
			conflictingAssignment = assignment
		}
	}

	// If we're trying to create an allDevices assignment but other assignment types exist
	if isCurrentAllDevices && otherAssignmentsCount > 0 {
		assignmentId := "unknown"
		if conflictingAssignment != nil && conflictingAssignment.GetId() != nil {
			assignmentId = *conflictingAssignment.GetId()
		}

		return fmt.Errorf(
			"cannot assign Windows Autopilot Deployment Profile to 'All Devices' because other assignment types already exist (conflicting assignment ID: %s). "+
				"'All Devices' assignments must be the only assignment for the profile",
			assignmentId,
		)
	}

	// If we're trying to create a non-allDevices assignment but allDevices assignment exists
	if !isCurrentAllDevices && allDevicesCount > 0 {
		assignmentId := "unknown"
		if conflictingAssignment != nil && conflictingAssignment.GetId() != nil {
			assignmentId = *conflictingAssignment.GetId()
		}

		return fmt.Errorf(
			"cannot create assignment with target type '%s' because an 'All Devices' assignment already exists (ID: %s). "+
				"'All Devices' assignments must be the only assignment for the profile",
			object.Target.TargetType.ValueString(),
			assignmentId,
		)
	}

	tflog.Debug(ctx, "All Devices exclusivity validation passed")
	return nil
}

// validateGroupIdUniqueness ensures that group IDs are unique across all assignments for the profile
func validateGroupIdUniqueness(
	ctx context.Context,
	object WindowsAutopilotDeploymentProfileAssignmentResourceModel,
	existingAssignments []graphmodels.WindowsAutopilotDeploymentProfileAssignmentable,
	isUpdate bool,
) error {
	tflog.Debug(ctx, "Validating group ID uniqueness constraint")

	// Only validate if the current assignment uses a group ID
	currentTargetType := object.Target.TargetType.ValueString()
	if currentTargetType != "groupAssignment" && currentTargetType != "exclusionGroupAssignment" {
		tflog.Debug(ctx, "Current assignment does not use group ID, skipping validation")
		return nil
	}

	// Check if current assignment has a group ID
	if object.Target.GroupId.IsNull() || object.Target.GroupId.IsUnknown() {
		tflog.Debug(ctx, "Current assignment has no group ID, skipping validation")
		return nil
	}

	currentGroupId := object.Target.GroupId.ValueString()

	// Check existing assignments for group ID conflicts
	for _, assignment := range existingAssignments {
		if assignment == nil {
			continue
		}

		// Skip the assignment being updated to avoid false conflicts
		if isUpdate && assignment.GetId() != nil && !object.ID.IsNull() {
			if *assignment.GetId() == object.ID.ValueString() {
				continue
			}
		}

		target := assignment.GetTarget()
		if target == nil {
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			continue
		}

		// Check group-based assignments for group ID conflicts
		var existingGroupId *string
		switch *odataType {
		case "#microsoft.graph.groupAssignmentTarget":
			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				existingGroupId = groupTarget.GetGroupId()
			}
		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			if exclusionTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				existingGroupId = exclusionTarget.GetGroupId()
			}
		}

		// If we found a group ID conflict
		if existingGroupId != nil && *existingGroupId == currentGroupId {
			assignmentId := "unknown"
			if assignment.GetId() != nil {
				assignmentId = *assignment.GetId()
			}

			return fmt.Errorf(
				"cannot create assignment with group ID '%s' because another assignment (ID: %s) already uses the same group ID. "+
					"Group IDs must be unique across all assignments for the profile",
				currentGroupId,
				assignmentId,
			)
		}
	}

	tflog.Debug(ctx, "Group ID uniqueness validation passed")
	return nil
}
