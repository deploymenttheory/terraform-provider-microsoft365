package graphBetaApplicationsAgentIdentityBlueprint

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/directory"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
//
//   - Retrieves the current state from the create request
//   - Validates the state data and timeout configuration
//   - Sends POST request to create the resource
//   - Reads the created resource state
//   - Cleans up by removing the resource from Terraform state
func (r *AgentIdentityBlueprintResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AgentIdentityBlueprintResourceModel

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

	if err := validateRequest(ctx, r.client, &object); err != nil {
		resp.Diagnostics.AddError(
			"Validation Error",
			fmt.Sprintf("Failed to validate %s: %s", ResourceName, err.Error()),
		)
		return
	}

	requestBody, err := constructResource(ctx, &object, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing agent identity blueprint",
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

// Read handles the Read operation.
//
//   - Retrieves the current state from the read request
//   - Validates the state data and timeout configuration
//   - Sends GET request to retrieve the resource
//   - Maps the remote resource state to the Terraform state
//   - Cleans up by removing the resource from Terraform state
func (r *AgentIdentityBlueprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentIdentityBlueprintResourceModel

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

	MapRemoteResourceStateToTerraform(ctx, &object, application)

	// Fetch sponsors using custom request (SDK doesn't support the cast endpoint)
	sponsorConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          fmt.Sprintf("/applications/%s/microsoft.graph.agentIdentityBlueprint/sponsors", object.ID.ValueString()),
		ResourceIDPattern: "",
		ResourceID:        "",
		EndpointSuffix:    "",
	}

	sponsorResponse, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), sponsorConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapSponsorIdsToTerraform(ctx, &object, sponsorResponse)

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

// Update handles the Update operation.
//
//   - Retrieves the current state from the update request
//   - Validates the state data and timeout configuration
//   - Sends PATCH request to update the resource
//   - Reads the updated resource state
//   - Cleans up by removing the resource from Terraform state
func (r *AgentIdentityBlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentIdentityBlueprintResourceModel
	var state AgentIdentityBlueprintResourceModel

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

	if err := validateRequest(ctx, r.client, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Validation Error",
			fmt.Sprintf("Failed to validate %s: %s", ResourceName, err.Error()),
		)
		return
	}

	requestBody, err := constructResource(ctx, &plan, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing agent identity blueprint",
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

	sponsorsToAdd, sponsorsToRemove, ownersToAdd, ownersToRemove := ResolveSponsorAndOwnerChanges(ctx, &state, &plan)

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

	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency ")
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

// Delete handles the Delete operation for resource Agent Identity Blueprint.
// When hard_delete is true, the blueprint is deleted in two steps:
// 1. Delete the application (soft delete - moves to deleted items)
// 2. Wait for the resource to appear in deleted items (handles eventual consistency)
// 3. Permanently delete from /directory/deleteditems/{id}
// 4. Verify deletion by confirming resource is gone
// When hard_delete is false (default), only the soft delete is performed.
// REF: https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-list?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-delete?view=graph-rest-beta
func (r *AgentIdentityBlueprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AgentIdentityBlueprintResourceModel

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

	blueprintId := data.ID.ValueString()

	// Define the soft delete function for applications
	softDeleteFunc := func(ctx context.Context) error {
		return r.client.
			Applications().
			ByApplicationId(blueprintId).
			Delete(ctx, nil)
	}

	deleteOpts := directory.DeleteOptions{
		MaxRetries:    10,
		RetryInterval: 5 * time.Second,
		ResourceType:  directory.ResourceTypeApplication,
		ResourceID:    blueprintId,
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
