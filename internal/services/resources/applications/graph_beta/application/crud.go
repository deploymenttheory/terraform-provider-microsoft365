package graphBetaApplication

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

// Create handles the Create operation for application resources.
//
// Operation: Creates a new Microsoft Entra application
// API Calls:
//   - POST /applications
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-post-applications?view=graph-rest-beta
func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ApplicationResourceModel

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

	requestBody, err := constructResource(ctx, &object, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing application",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdApplication, err := r.client.
		Applications().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdApplication.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Waiting 25 seconds for eventual consistency after create")
	time.Sleep(25 * time.Second)

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

// Read handles the Read operation for application resources.
//
// Operation: Retrieves application details including owners
// API Calls:
//   - GET /applications/{id}
//   - GET /applications/{id}/owners
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta
func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ApplicationResourceModel

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

	application, err := r.client.
		Applications().
		ByApplicationId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, application)

	owners, err := r.client.
		Applications().
		ByApplicationId(object.ID.ValueString()).
		Owners().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteOwnersToTerraform(ctx, &object, owners)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for application resources.
//
// Operation: Updates application properties and manages owners
// API Calls:
//   - PATCH /applications/{id}
//   - DELETE /applications/{id}/owners/{ownerId}/$ref (for removed owners)
//   - POST /applications/{id}/owners/$ref (for new owners)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationResourceModel
	var state ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting update of resource: %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing application",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Applications().
		ByApplicationId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	ownersToAdd, ownersToRemove := ResolveOwnerChanges(ctx, &state, &plan)

	adapter := r.client.GetAdapter()

	for _, ownerID := range ownersToRemove {
		if err := RemoveOwner(ctx, adapter, state.ID.ValueString(), ownerID); err != nil {
			resp.Diagnostics.AddError(
				"Error removing owner",
				fmt.Sprintf("Could not remove owner %s from %s: %s", ownerID, ResourceName, err.Error()),
			)
			return
		}
	}

	for _, ownerID := range ownersToAdd {
		if err := AddOwner(ctx, adapter, state.ID.ValueString(), ownerID); err != nil {
			resp.Diagnostics.AddError(
				"Error adding owner",
				fmt.Sprintf("Could not add owner %s to %s: %s", ownerID, ResourceName, err.Error()),
			)
			return
		}
	}

	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency")
	time.Sleep(15 * time.Second)

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

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for application resources.
//
// Operation: Deletes application with optional hard delete
// API Calls:
//   - DELETE /applications/{id}
//   - GET /directory/deletedItems/microsoft.graph.application (if hard_delete is true)
//   - DELETE /directory/deletedItems/{id} (if hard_delete is true)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-delete?view=graph-rest-beta
// Note: When hard_delete is true, performs soft delete then permanent deletion from directory deleted items after eventual consistency delay
func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	applicationId := data.ID.ValueString()

	// Define the soft delete function for applications
	softDeleteFunc := func(ctx context.Context) error {
		return r.client.
			Applications().
			ByApplicationId(applicationId).
			Delete(ctx, nil)
	}

	deleteOpts := directory.DeleteOptions{
		MaxRetries:    10,
		RetryInterval: 5 * time.Second,
		ResourceType:  directory.ResourceTypeApplication,
		ResourceID:    applicationId,
		ResourceName:  data.DisplayName.ValueString(),
	}

	err := directory.ExecuteDeleteWithVerification(
		ctx,
		r.client,
		softDeleteFunc,
		data.HardDelete.ValueBool(),
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
