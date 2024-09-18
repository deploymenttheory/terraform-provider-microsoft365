package graphBetaWinGetApp

import (
	"context"
	"fmt"

	graphBetaMobileAppAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/mobile_app_assignment"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *WinGetAppResource) createAssignments(ctx context.Context, appID string, assignments []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel) error {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to create")
		return nil
	}

	for _, assignment := range assignments {
		requestBody, err := graphBetaMobileAppAssignment.ConstructResource(ctx, &assignment)
		if err != nil {
			return fmt.Errorf("error constructing assignment: %v", err)
		}

		_, err = r.client.DeviceAppManagement().
			MobileApps().
			ByMobileAppId(appID).
			Assignments().
			Post(ctx, requestBody, nil)

		if err != nil {
			return fmt.Errorf("error creating assignment: %v", err)
		}
	}

	return nil
}

func (r *WinGetAppResource) readAssignments(ctx context.Context, appID string) ([]graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel, error) {
	tflog.Debug(ctx, fmt.Sprintf("Starting readAssignments for app ID: %s", appID))

	assignmentsResponse, err := r.client.DeviceAppManagement().
		MobileApps().
		ByMobileAppId(appID).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("error reading assignments: %v", err)
	}

	assignments := assignmentsResponse.GetValue()
	var result []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel

	for _, assignment := range assignments {
		var terraformAssignment graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel
		graphBetaMobileAppAssignment.MapRemoteStateToTerraform(ctx, &terraformAssignment, assignment)
		result = append(result, terraformAssignment)
	}

	return result, nil
}

func (r *WinGetAppResource) updateAssignments(ctx context.Context, appID string, newAssignments []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel) error {
	if len(newAssignments) == 0 {
		tflog.Debug(ctx, "No assignments to update")
		return nil
	}

	// First, get existing assignments
	existingAssignments, err := r.readAssignments(ctx, appID)
	if err != nil {
		return fmt.Errorf("error reading existing assignments: %v", err)
	}

	// Create a map of existing assignments for easy lookup
	existingMap := make(map[string]graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel)
	for _, assignment := range existingAssignments {
		existingMap[assignment.ID.ValueString()] = assignment
	}

	// Update or create assignments
	for _, newAssignment := range newAssignments {
		requestBody, err := graphBetaMobileAppAssignment.ConstructResource(ctx, &newAssignment)
		if err != nil {
			return fmt.Errorf("error constructing assignment: %v", err)
		}

		if _, exists := existingMap[newAssignment.ID.ValueString()]; exists {
			// Update existing assignment
			_, err = r.client.DeviceAppManagement().
				MobileApps().
				ByMobileAppId(appID).
				Assignments().
				ByMobileAppAssignmentId(newAssignment.ID.ValueString()).
				Patch(ctx, requestBody, nil)
		} else {
			// Create new assignment
			_, err = r.client.DeviceAppManagement().
				MobileApps().
				ByMobileAppId(appID).
				Assignments().
				Post(ctx, requestBody, nil)
		}

		if err != nil {
			return fmt.Errorf("error updating/creating assignment: %v", err)
		}
	}

	// Delete assignments that are not in the new set
	for _, existingAssignment := range existingAssignments {
		found := false
		for _, newAssignment := range newAssignments {
			if existingAssignment.ID == newAssignment.ID {
				found = true
				break
			}
		}
		if !found {
			err := r.client.DeviceAppManagement().
				MobileApps().
				ByMobileAppId(appID).
				Assignments().
				ByMobileAppAssignmentId(existingAssignment.ID.ValueString()).
				Delete(ctx, nil)
			if err != nil {
				return fmt.Errorf("error deleting assignment %s: %v", existingAssignment.ID.ValueString(), err)
			}
		}
	}

	return nil
}

func (r *WinGetAppResource) deleteAssignments(ctx context.Context, appID string) error {
	assignments, err := r.readAssignments(ctx, appID)
	if err != nil {
		return err
	}

	for _, assignment := range assignments {
		err := r.client.DeviceAppManagement().
			MobileApps().
			ByMobileAppId(appID).
			Assignments().
			ByMobileAppAssignmentId(assignment.ID.ValueString()).
			Delete(ctx, nil)

		if err != nil {
			return fmt.Errorf("error deleting assignment %s: %v", assignment.ID.ValueString(), err)
		}
	}

	return nil
}
