package graphBetaWinGetApp

import (
	"context"
	"fmt"

	graphBetaMobileAppAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/mobile_app_assignment"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Helper functions for assignments
func (r *WinGetAppResource) createAssignments(ctx context.Context, appID string, assignments []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel) error {
	for _, assignment := range assignments {
		// Create a new MobileAppAssignmentResource
		assignmentResource := graphBetaMobileAppAssignment.NewMobileAppAssignmentResource()

		// Set the SourceID for the assignment
		assignment.SourceID = types.StringValue(appID)

		// Create a new CreateRequest
		req := resource.CreateRequest{}

		// Set the Plan in the CreateRequest
		diags := req.Plan.Set(ctx, assignment)
		if diags.HasError() {
			return fmt.Errorf("error setting plan for assignment: %v", diags)
		}

		// Create a new CreateResponse
		resp := &resource.CreateResponse{}

		// Call the Create method of the MobileAppAssignmentResource
		assignmentResource.Create(ctx, req, resp)

		if resp.Diagnostics.HasError() {
			return fmt.Errorf("error creating assignment: %v", resp.Diagnostics)
		}
	}
	return nil
}

func (r *WinGetAppResource) readAssignments(ctx context.Context, appID string) ([]graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel, error) {
	tflog.Debug(ctx, fmt.Sprintf("Starting readAssignments for app ID: %s", appID))

	// Create a new MobileAppAssignmentResource
	assignmentResource := graphBetaMobileAppAssignment.NewMobileAppAssignmentResource()

	// Ensure the resource is not nil
	if assignmentResource == nil {
		return nil, fmt.Errorf("MobileAppAssignmentResource is nil")
	}

	// Create a new ReadRequest
	req := resource.ReadRequest{}

	// Create a state to hold the appID
	state := struct {
		SourceID types.String `tfsdk:"source_id"`
	}{
		SourceID: types.StringValue(appID),
	}

	// Set the State in the ReadRequest
	diags := req.State.Set(ctx, &state)
	if diags.HasError() {
		return nil, fmt.Errorf("error setting state for reading assignments: %v", diags)
	}

	// Create a new ReadResponse
	resp := &resource.ReadResponse{}

	// Call the Read method of the MobileAppAssignmentResource
	assignmentResource.Read(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		return nil, fmt.Errorf("error reading assignments: %v", resp.Diagnostics)
	}

	// Create a variable to hold the assignments
	var assignments []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel

	// Get the assignments from the response state
	diags = resp.State.Get(ctx, &assignments)
	if diags.HasError() {
		tflog.Warn(ctx, fmt.Sprintf("No assignments found or error retrieving assignments for app ID %s: %v", appID, diags))
		return []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel{}, nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully read %d assignments for app ID: %s", len(assignments), appID))
	return assignments, nil
}

func (r *WinGetAppResource) updateAssignments(ctx context.Context, appID string, newAssignments []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel) error {
	// First, delete all existing assignments
	err := r.deleteAssignments(ctx, appID)
	if err != nil {
		return err
	}

	// Then, create all new assignments
	return r.createAssignments(ctx, appID, newAssignments)
}

func (r *WinGetAppResource) deleteAssignments(ctx context.Context, appID string) error {
	assignments, err := r.readAssignments(ctx, appID)
	if err != nil {
		return err
	}

	for _, assignment := range assignments {
		// Create a new MobileAppAssignmentResource for each assignment
		assignmentResource := graphBetaMobileAppAssignment.NewMobileAppAssignmentResource()

		// Create a new DeleteRequest
		req := resource.DeleteRequest{}

		// Set the State in the DeleteRequest
		diags := req.State.Set(ctx, assignment)
		if diags.HasError() {
			return fmt.Errorf("error setting state for assignment deletion: %v", diags)
		}

		// Create a new DeleteResponse
		resp := &resource.DeleteResponse{}

		// Call the Delete method of the MobileAppAssignmentResource
		assignmentResource.Delete(ctx, req, resp)

		if resp.Diagnostics.HasError() {
			return fmt.Errorf("error deleting assignment: %v", resp.Diagnostics)
		}
	}
	return nil
}
