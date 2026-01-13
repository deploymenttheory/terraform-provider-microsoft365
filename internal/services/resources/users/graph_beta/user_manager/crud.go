package graphBetaUsersUserManager

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for the user manager relationship.
//
//   - Retrieves the planned configuration from the create request
//   - Constructs the manager reference from the plan
//   - Sends PUT request to assign the manager
//   - Captures the resource ID (user_id)
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API
//   - Updates the final state with the fresh data from the API
//
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-post-manager?view=graph-rest-beta
func (r *UserManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserManagerResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	userId := plan.UserID.ValueString()
	managerId := plan.ManagerID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Assigning manager %s to user %s", managerId, userId))

	requestBody := constructManagerReference(managerId)

	err := r.client.
		Users().
		ByUserId(userId).
		Manager().
		Ref().
		Put(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(userId)
	tflog.Debug(ctx, fmt.Sprintf("Successfully assigned manager %s to user %s", managerId, userId))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName

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

// Read handles the Read operation for the user manager relationship.
//
//   - Retrieves the current state from the read request
//   - Gets the manager details from the API
//   - Maps the manager details to Terraform state
//
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-list-manager?view=graph-rest-beta
func (r *UserManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserManagerResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := constants.TfOperationRead
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading manager for user %s", state.UserID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	userId := state.UserID.ValueString()

	manager, err := r.client.
		Users().
		ByUserId(userId).
		Manager().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, manager)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for the user manager relationship.
//
// To update the manager, we must first remove the existing manager and then add the new one.
//
//   - Retrieves the planned changes from the update request
//   - Removes the existing manager reference
//   - Adds the new manager reference
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API
//   - Updates the final state with the fresh data from the API
//
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-delete-manager?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-post-manager?view=graph-rest-beta
func (r *UserManagerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserManagerResourceModel
	var state UserManagerResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	userId := plan.UserID.ValueString()
	newManagerId := plan.ManagerID.ValueString()
	oldManagerId := state.ManagerID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Updating manager for user %s from %s to %s", userId, oldManagerId, newManagerId))

	// Step 1: Remove the existing manager
	tflog.Debug(ctx, fmt.Sprintf("Removing existing manager %s from user %s", oldManagerId, userId))

	err := r.client.
		Users().
		ByUserId(userId).
		Manager().
		Ref().
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// Brief delay to allow for eventual consistency
	time.Sleep(2 * time.Second)

	// Step 2: Add the new manager
	tflog.Debug(ctx, fmt.Sprintf("Adding new manager %s to user %s", newManagerId, userId))

	requestBody := constructManagerReference(newManagerId)

	err = r.client.
		Users().
		ByUserId(userId).
		Manager().
		Ref().
		Put(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(userId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for the user manager relationship.
//
//   - Retrieves the current state from the delete request
//   - Sends DELETE request to remove the manager reference
//   - Cleans up by removing the resource from Terraform state
//
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-delete-manager?view=graph-rest-beta
func (r *UserManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserManagerResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	userId := state.UserID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Removing manager from user %s", userId))

	err := r.client.
		Users().
		ByUserId(userId).
		Manager().
		Ref().
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
