package graphBetaWinGetApp

import (
	"context"
	"fmt"

	graphBetaMobileAppAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/mobile_app_assignment"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *WinGetAppResource) createAssignments(ctx context.Context, appID string, assignments []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel) error {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to create")
		return nil
	}

	assignmentResource := graphBetaMobileAppAssignment.NewMobileAppAssignmentResource()

	for _, assignment := range assignments {
		assignment.SourceID = types.StringValue(appID)

		req := resource.CreateRequest{
			Plan: assignment,
		}
		resp := &resource.CreateResponse{}

		assignmentResource.Create(ctx, req, resp)

		if resp.Diagnostics.HasError() {
			return fmt.Errorf("error creating assignment: %v", resp.Diagnostics)
		}
	}
	return nil
}

func (r *WinGetAppResource) readAssignments(ctx context.Context, appID string) ([]graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel, error) {
	tflog.Debug(ctx, fmt.Sprintf("Starting readAssignments for app ID: %s", appID))

	assignmentResource := graphBetaMobileAppAssignment.NewMobileAppAssignmentResource()

	req := resource.ReadRequest{
		State: graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel{
			SourceID: types.StringValue(appID),
		},
	}
	resp := &resource.ReadResponse{}

	assignmentResource.Read(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		return nil, fmt.Errorf("error reading assignments: %v", resp.Diagnostics)
	}

	var assignments []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel
	resp.State.Get(ctx, &assignments)

	return assignments, nil
}

func (r *WinGetAppResource) updateAssignments(ctx context.Context, appID string, newAssignments []graphBetaMobileAppAssignment.MobileAppAssignmentResourceModel) error {
	if len(newAssignments) == 0 {
		tflog.Debug(ctx, "No assignments to update")
		return nil
	}

	assignmentResource := graphBetaMobileAppAssignment.NewMobileAppAssignmentResource()

	for _, assignment := range newAssignments {
		assignment.SourceID = types.StringValue(appID)

		req := resource.UpdateRequest{
			Plan: assignment,
		}
		resp := &resource.UpdateResponse{}

		assignmentResource.Update(ctx, req, resp)

		if resp.Diagnostics.HasError() {
			return fmt.Errorf("error updating assignment: %v", resp.Diagnostics)
		}
	}
	return nil
}

func (r *WinGetAppResource) deleteAssignments(ctx context.Context, appID string) error {
	assignments, err := r.readAssignments(ctx, appID)
	if err != nil {
		return err
	}

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to delete")
		return nil
	}

	assignmentResource := graphBetaMobileAppAssignment.NewMobileAppAssignmentResource()

	for _, assignment := range assignments {
		req := resource.DeleteRequest{
			State: assignment,
		}
		resp := &resource.DeleteResponse{}

		assignmentResource.Delete(ctx, req, resp)

		if resp.Diagnostics.HasError() {
			return fmt.Errorf("error deleting assignment: %v", resp.Diagnostics)
		}
	}
	return nil
}
