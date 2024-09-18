package graphBetaWinGetApp

import (
	"context"
	"fmt"
	"reflect"

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

	// Get existing assignments
	existingAssignments, err := r.readAssignments(ctx, appID)
	if err != nil {
		return fmt.Errorf("error reading existing assignments: %v", err)
	}

	// Create map of existing assignments
	existingMap := make(map[string]graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel)
	for _, assignment := range existingAssignments {
		existingMap[assignment.ID.ValueString()] = assignment
	}

	// Check if there are any changes
	changesExist := false
	if len(existingAssignments) != len(newAssignments) {
		changesExist = true
	} else {
		for _, newAssignment := range newAssignments {
			existingAssignment, exists := existingMap[newAssignment.ID.ValueString()]
			if !exists || !reflect.DeepEqual(existingAssignment, newAssignment) {
				changesExist = true
				break
			}
		}
	}

	if !changesExist {
		tflog.Debug(ctx, "No changes in assignments")
		return nil
	}

	// If there are changes, delete all existing assignments
	err = r.deleteAssignments(ctx, appID)
	if err != nil {
		return fmt.Errorf("error deleting existing assignments: %v", err)
	}

	// Create all new assignments
	err = r.createAssignments(ctx, appID, newAssignments)
	if err != nil {
		return fmt.Errorf("error creating new assignments: %v", err)
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
