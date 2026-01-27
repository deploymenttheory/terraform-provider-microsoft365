package graphBetaAgentInstance

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

// Create handles the Create operation for agent instance resources.
//
// Operation: Creates a new agent instance with optional agent card manifest
// API Calls:
//   - POST /agentRegistry/agentInstances
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentregistry-post-agentinstances?view=graph-rest-beta
func (r *AgentInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AgentInstanceResourceModel

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
			"Error constructing agent instance",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		AgentRegistry().
		AgentInstances().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if createdResource.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error reading created resource ID",
			fmt.Sprintf("Could not extract ID from created resource: %s", ResourceName),
		)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	if createdResource.GetAgentUserId() != nil {
		object.AgentUserId = types.StringValue(*createdResource.GetAgentUserId())
	}

	// Capture agentCardManifest ID from create response so Read can fetch full details
	if manifest := createdResource.GetAgentCardManifest(); manifest != nil && manifest.GetId() != nil {
		if object.AgentCardManifest != nil {
			object.AgentCardManifest.ID = types.StringValue(*manifest.GetId())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Waiting for eventual consistency after create")
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

// Read handles the Read operation for agent instance resources.
//
// Operation: Retrieves agent instance details including agent card manifest
// API Calls:
//   - GET /agentRegistry/agentInstances/{agentInstanceId}
//   - GET /agentRegistry/agentInstances/{agentInstanceId}/agentCardManifest
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentinstance-get?view=graph-rest-beta
func (r *AgentInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentInstanceResourceModel

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

	// Get the agent instance
	agentInstance, err := r.client.
		AgentRegistry().
		AgentInstances().
		ByAgentInstanceId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, agentInstance)

	// Get the agent card manifest (separate API call - not included in main GET)
	manifest, err := r.client.
		AgentRegistry().
		AgentInstances().
		ByAgentInstanceId(object.ID.ValueString()).
		AgentCardManifest().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if manifest != nil {
		MapAgentCardManifestToTerraform(ctx, &object, manifest)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for agent instance resources.
//
// Operation: Updates agent instance properties and agent card manifest
// API Calls:
//   - PATCH /agentRegistry/agentInstances/{agentInstanceId}
//   - PATCH /agentRegistry/agentCardManifests/{agentCardManifestId} (if manifest exists)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentinstance-update?view=graph-rest-beta
func (r *AgentInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentInstanceResourceModel
	var state AgentInstanceResourceModel

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

	requestBody, err := constructResource(ctx, &plan, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing agent instance for update",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		AgentRegistry().
		AgentInstances().
		ByAgentInstanceId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	if plan.AgentCardManifest != nil && state.AgentCardManifest != nil && !state.AgentCardManifest.ID.IsNull() {
		manifestId := state.AgentCardManifest.ID.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Updating agentCardManifest with ID: %s", manifestId))

		manifestBody, err := constructAgentCardManifest(ctx, plan.AgentCardManifest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing agent card manifest for update",
				fmt.Sprintf("Could not construct agent card manifest: %s", err.Error()),
			)
			return
		}

		_, err = r.client.
			AgentRegistry().
			AgentCardManifests().
			ByAgentCardManifestId(manifestId).
			Patch(ctx, manifestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, "Waiting for eventual consistency after update")
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

// Delete handles the Delete operation for agent instance resources.
//
// Operation: Deletes an agent instance
// API Calls:
//   - DELETE /agentRegistry/agentInstances/{agentInstanceId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentregistry-delete-agentinstances?view=graph-rest-beta
func (r *AgentInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object AgentInstanceResourceModel

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

	err := r.client.
		AgentRegistry().
		AgentInstances().
		ByAgentInstanceId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
