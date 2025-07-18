package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Create handles the Create operation.
func (r *DeviceHealthScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object DeviceHealthScriptResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceManagement().
		DeviceHealthScripts().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	requestAssignment, err := constructAssignment(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Create Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		DeviceManagement().
		DeviceHealthScripts().
		ByDeviceHealthScriptId(object.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create - Assignments", r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for macos platform scripts resources.
//
//   - Retrieves the current state from the read request
//   - Gets the resource details including assignments from the API using expand
//   - Maps both resource and assignment details to Terraform state
//
// The function ensures all components are properly read and mapped into the
// Terraform state in a single API call, providing a complete view of the
// resource's current configuration on the server.
func (r *DeviceHealthScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object DeviceHealthScriptResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := "Read"
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	respResource, err := r.client.
		DeviceManagement().
		DeviceHealthScripts().
		ByDeviceHealthScriptId(object.ID.ValueString()).
		Get(ctx, &devicemanagement.DeviceHealthScriptsDeviceHealthScriptItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.DeviceHealthScriptsDeviceHealthScriptItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
			},
		})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for macos platform scripts resources.
//
// The function performs the following operations:
//   - Patches the existing script resource with updated settings using PATCH
//   - Updates assignments using POST or removes all assignments if nil
//   - Retrieves the updated resource with expanded assignments
//   - Maps the remote state back to Terraform
//
// The Microsoft Graph Beta API supports direct updates of device shell script resources
// through PATCH operations for the base resource, while assignments are handled through
// a separate POST operation to the assign endpoint. This allows for atomic updates
// of both the script properties and its assignments.
func (r *DeviceHealthScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceHealthScriptResourceModel
	var state DeviceHealthScriptResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		DeviceHealthScripts().
		ByDeviceHealthScriptId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Always handle assignments - either update with new assignments or remove all assignments if nil
	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for update method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		DeviceManagement().
		DeviceHealthScripts().
		ByDeviceHealthScriptId(state.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update - Assignments", r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for Device Management Script resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
//
// All assignments and settings associated with the resource are automatically removed as part of the deletion.
func (r *DeviceHealthScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object DeviceHealthScriptResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		DeviceHealthScripts().
		ByDeviceHealthScriptId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
