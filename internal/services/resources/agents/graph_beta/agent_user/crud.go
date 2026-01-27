package graphBetaAgentUser

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/directory"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for agent user resources.
//
// Operation: Creates a new agent user
// API Calls:
//   - POST /users/microsoft.graph.agentUser
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentuser-post?view=graph-rest-beta
func (r *AgentUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AgentUserResourceModel

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
			"Error constructing agent user",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdUser, err := r.client.
		Users().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if createdUser.GetId() != nil {
		object.ID = types.StringValue(*createdUser.GetId())
	} else {
		resp.Diagnostics.AddError(
			"Error reading created resource ID",
			fmt.Sprintf("Could not extract ID from created resource: %s", ResourceName),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Waiting 5 seconds for eventual consistency after create")
	time.Sleep(5 * time.Second)

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

// Read handles the Read operation for agent user resources.
//
// Operation: Retrieves agent user details including sponsors
// API Calls:
//   - GET /users/microsoft.graph.agentUser/{userId}
//   - GET /users/{userId}/sponsors
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentuser-get?view=graph-rest-beta
func (r *AgentUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentUserResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := constants.TfOperationRead
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

	user, err := r.client.
		Users().
		ByUserId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, user)

	sponsors, err := r.client.
		Users().
		ByUserId(object.ID.ValueString()).
		Sponsors().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapSponsorsToTerraform(ctx, &object, sponsors)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for agent user resources.
//
// Operation: Updates agent user properties and manages sponsors
// API Calls:
//   - PATCH /users/microsoft.graph.agentUser/{userId}
//   - DELETE /users/{userId}/sponsors/{sponsorId}/$ref (for removed sponsors)
//   - POST /users/{userId}/sponsors/$ref (for new sponsors)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentuser-update?view=graph-rest-beta
func (r *AgentUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentUserResourceModel
	var state AgentUserResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructUpdateResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing agent user for update",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Users().
		ByUserId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// Handle sponsor changes
	sponsorsToAdd, sponsorsToRemove := ResolveSponsorChanges(ctx, &state, &plan)

	adapter := r.client.GetAdapter()

	for _, sponsorID := range sponsorsToRemove {
		if err := RemoveSponsor(ctx, adapter, state.ID.ValueString(), sponsorID); err != nil {
			resp.Diagnostics.AddError(
				"Error removing sponsor",
				fmt.Sprintf("Could not remove sponsor %s from %s: %s", sponsorID, ResourceName, err.Error()),
			)
			return
		}
	}

	for _, sponsorID := range sponsorsToAdd {
		if err := AddSponsor(ctx, adapter, state.ID.ValueString(), sponsorID); err != nil {
			resp.Diagnostics.AddError(
				"Error adding sponsor",
				fmt.Sprintf("Could not add sponsor %s to %s: %s", sponsorID, ResourceName, err.Error()),
			)
			return
		}
	}

	tflog.Debug(ctx, "Waiting 5 seconds for eventual consistency after update")
	time.Sleep(5 * time.Second)

	plan.ID = state.ID
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

// Delete handles the Delete operation for agent user resources.
//
// Operation: Deletes agent user with optional hard delete
// API Calls:
//   - DELETE /users/{userId}
//   - GET /directory/deletedItems/{id} (verification after soft delete)
//   - DELETE /directory/deletedItems/{id} (if hard_delete is true)
//   - GET /directory/deletedItems/{id} (verification after hard delete)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentuser-delete?view=graph-rest-beta
// Note: When hard_delete is true, performs soft delete then permanent deletion from directory deleted items after eventual consistency delay
func (r *AgentUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object AgentUserResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	agentUserId := object.ID.ValueString()

	// Define the soft delete function
	softDeleteFunc := func(ctx context.Context) error {
		return r.client.
			Users().
			ByUserId(agentUserId).
			Delete(ctx, nil)
	}

	// Define delete options
	deleteOpts := directory.DeleteOptions{
		MaxRetries:    10,
		RetryInterval: 5 * time.Second,
		ResourceType:  directory.ResourceTypeUser,
		ResourceID:    agentUserId,
		ResourceName:  object.DisplayName.ValueString(),
	}

	err := directory.ExecuteDeleteWithVerification(
		ctx,
		r.client,
		softDeleteFunc,
		object.HardDelete.ValueBool(),
		deleteOpts,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Failed",
			fmt.Sprintf("Failed to delete %s: %s", ResourceName, err.Error()),
		)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
