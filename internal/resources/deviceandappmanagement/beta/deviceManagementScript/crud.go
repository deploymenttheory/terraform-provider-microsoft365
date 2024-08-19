package graphbetadevicemanagementscript

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
func (r *DeviceManagementScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceManagementScriptResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing device management script",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	script, err := r.client.DeviceManagement().DeviceManagementScripts().Post(ctx, requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device management script",
			fmt.Sprintf("Could not create device management script: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(*script.GetId())

	// Handle assignments
	if len(plan.Assignments) > 0 {
		assignments, err := constructAssignments(ctx, plan.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing resource assignments",
				fmt.Sprintf("Could not construct assignments: %s", err.Error()),
			)
			return
		}

		for _, assignment := range assignments {
			_, err := r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(*script.GetId()).Assignments().Post(ctx, assignment, nil)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating resource assignment",
					fmt.Sprintf("Could not create assignment for device management script %s: %s", *script.GetId(), err.Error()),
				)
				return
			}
		}
	}

	// Handle group assignments
	if len(plan.GroupAssignments) > 0 {
		groupAssignments, err := constructGroupAssignments(ctx, plan.GroupAssignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing resource group assignments",
				fmt.Sprintf("Could not construct group assignments: %s", err.Error()),
			)
			return
		}

		for _, groupAssignment := range groupAssignments {
			_, err := r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(*script.GetId()).GroupAssignments().Post(ctx, groupAssignment, nil)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating resource group assignment",
					fmt.Sprintf("Could not create group assignment for device management script %s: %s", *script.GetId(), err.Error()),
				)
				return
			}
		}
	}

	MapRemoteStateToTerraform(ctx, &plan, script)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *DeviceManagementScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceManagementScriptResourceModel
	tflog.Debug(ctx, "Starting Read method for device management script")

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading device management script with ID: %s", state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	assignmentFilter, err := r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(state.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		crud.HandleReadErrorIfNotFound(ctx, resp, r, &state, err)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, assignmentFilter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for the DeviceManagementScript resource.
// It independently updates the script, assignments, and group assignments,
// using separate mapping functions to update the Terraform state for each.
func (r *DeviceManagementScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DeviceManagementScriptResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.updateResource(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error updating script resource", err.Error())
		return
	}

	err = r.updateAssignments(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error updating assignments", err.Error())
		return
	}

	err = r.updateGroupAssignments(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error updating group assignments", err.Error())
		return
	}

	// Map the updated resource to Terraform state
	updatedResource, err := r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(plan.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving updated resource", err.Error())
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, updatedResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

func (r *DeviceManagementScriptResource) updateResource(ctx context.Context, plan, state *DeviceManagementScriptResourceModel) error {
	requestBody, err := constructResource(ctx, plan)
	if err != nil {
		return fmt.Errorf("could not construct resource: %v", err)
	}

	_, err = r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(plan.ID.ValueString()).Patch(ctx, requestBody, nil)
	if err != nil {
		return fmt.Errorf("error updating script: %v", err)
	}

	return nil
}

func (r *DeviceManagementScriptResource) updateAssignments(ctx context.Context, plan, state *DeviceManagementScriptResourceModel) error {
	// Delete assignments that are in state but not in plan
	for _, stateAssignment := range state.Assignments {
		if !assignmentExistsInPlan(stateAssignment, plan.Assignments) {
			err := r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(plan.ID.ValueString()).Assignments().ByDeviceManagementScriptAssignmentId(stateAssignment.Id.ValueString()).Delete(ctx, nil)
			if err != nil {
				return fmt.Errorf("error deleting assignment: %v", err)
			}
		}
	}

	// Create or update assignments in plan
	for _, planAssignment := range plan.Assignments {
		assignment, err := constructAssignment(planAssignment)
		if err != nil {
			return fmt.Errorf("error constructing assignment: %v", err)
		}

		if assignmentExistsInState(planAssignment, state.Assignments) {
			// Update existing assignment
			_, err = r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(plan.ID.ValueString()).Assignments().ByDeviceManagementScriptAssignmentId(planAssignment.Id.ValueString()).Patch(ctx, assignment, nil)
		} else {
			// Create new assignment
			_, err = r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(plan.ID.ValueString()).Assignments().Post(ctx, assignment, nil)
		}

		if err != nil {
			return fmt.Errorf("error creating/updating assignment: %v", err)
		}
	}

	return nil
}

func (r *DeviceManagementScriptResource) updateGroupAssignments(ctx context.Context, plan, state *DeviceManagementScriptResourceModel) error {
	// Implement group assignment update logic here, similar to updateAssignments
	return nil
}

func assignmentExistsInPlan(assignment DeviceManagementScriptAssignmentResourceModel, planAssignments []DeviceManagementScriptAssignmentResourceModel) bool {
	for _, planAssignment := range planAssignments {
		if assignment.Id == planAssignment.Id {
			return true
		}
	}
	return false
}

func assignmentExistsInState(assignment DeviceManagementScriptAssignmentResourceModel, stateAssignments []DeviceManagementScriptAssignmentResourceModel) bool {
	for _, stateAssignment := range stateAssignments {
		if assignment.Id == stateAssignment.Id {
			return true
		}
	}
	return false
}

func constructAssignment(assignment DeviceManagementScriptAssignmentResourceModel) (models.DeviceManagementScriptAssignmentable, error) {
	// Implement assignment construction logic here
	return nil, nil
}

// Delete handles the Delete operation for the DeviceManagementScript resource.
// It performs the following steps:
// 1. Retrieves the current state of the resource.
// 2. Deletes all associated assignments (both regular and group assignments).
// 3. If no assignments remain, it deletes the script itself and removes the resource from the state.
// 4. If assignments still exist, it updates the state to reflect the removal of assignments.
//
// This approach allows for selective deletion of assignments while potentially
// preserving the script itself if assignments remain. The resource is only fully
// removed from Terraform state when both the script and all its assignments are deleted.
func (r *DeviceManagementScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeviceManagementScriptResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	for _, assignment := range data.Assignments {
		err := r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(data.ID.ValueString()).Assignments().ByDeviceManagementScriptAssignmentId(assignment.Id.ValueString()).Delete(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError("Error deleting assignment", err.Error())
			return
		}
	}
	data.Assignments = []DeviceManagementScriptAssignmentResourceModel{}

	for _, groupAssignment := range data.GroupAssignments {
		err := r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(data.ID.ValueString()).GroupAssignments().ByDeviceManagementScriptGroupAssignmentId(groupAssignment.ID.ValueString()).Delete(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError("Error deleting group assignment", err.Error())
			return
		}
	}
	data.GroupAssignments = []DeviceManagementScriptGroupAssignmentResourceModel{}

	if len(data.Assignments) == 0 && len(data.GroupAssignments) == 0 {
		err := r.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(data.ID.ValueString()).Delete(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
		resp.State.RemoveResource(ctx)
	} else {
		// If there are still assignments, just update the state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
