package graphBetaDeviceManagementScript

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	// mutex needed to lock Create requests during parallel runs to avoid overwhelming api and resulting in stating issues
	mu sync.Mutex

	// object is the resource model for the Endpoint Privilege Management resource
	object DeviceManagementScriptResourceModel
)

// Create handles the Create operation.
func (r *DeviceManagementScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	mu.Lock()
	defer mu.Unlock()

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	// create resource
	requestResource, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*requestResource.GetId())

	// create assignments
	if object.Assignments != nil {
		requestAssignment, err := constructAssignment(ctx, &object)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for create method",
				fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
			)
			return
		}

		_, err = r.client.
			DeviceManagement().
			DeviceManagementScripts().
			ByDeviceManagementScriptId(object.ID.ValueString()).
			Assignments().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
	}

	// Get the resource to update the state
	respResource, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Get(context.Background(), nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}
	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	respAssignments, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create - Assignments Fetch", r.ReadPermissions)
		return
	}
	MapRemoteAssignmentStateToTerraform(ctx, &object, respAssignments)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *DeviceManagementScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Debug(ctx, "Starting Read method for device management script")

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading device management script with ID: %s", object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	script, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	// Retrieve assignments
	assignments, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error reading assignments", fmt.Sprintf("Could not read assignments for device management script %s: %s", state.ID.ValueString(), err.Error()))
		return
	}

	// Retrieve group assignments
	groupAssignments, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		GroupAssignments().
		Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error reading group assignments", fmt.Sprintf("Could not read group assignments for device management script %s: %s", state.ID.ValueString(), err.Error()))
		return
	}

	// Map assignments to state
	if assignments != nil && len(assignments.GetValue()) > 0 {
		state.Assignments = make([]DeviceManagementScriptAssignmentResourceModel, len(assignments.GetValue()))
		for i, assignment := range assignments.GetValue() {
			state.Assignments[i] = MapAssignmentsRemoteStateToTerraform(assignment)
		}
	}
	if groupAssignments != nil && len(groupAssignments.GetValue()) > 0 {
		state.GroupAssignments = make([]DeviceManagementScriptGroupAssignmentResourceModel, len(groupAssignments.GetValue()))
		for i, groupAssignment := range groupAssignments.GetValue() {
			state.GroupAssignments[i] = MapGroupAssignmentsRemoteStateToTerraform(groupAssignment)
		}
	}

	MapRemoteStateToTerraform(ctx, &object, script)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for the DeviceManagementScript resource.
// It independently updates the script, assignments, and group assignments,
// using separate mapping functions to update the Terraform state for each.
func (r *DeviceManagementScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.ReadPermissions)
		return
	}

	// Update assignments
	err = r.updateAssignments(ctx, &object, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error updating assignments", err.Error())
		return
	}

	err = r.updateGroupAssignments(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error updating group assignments", err.Error())
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, requestBody)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
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

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.ReadPermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
