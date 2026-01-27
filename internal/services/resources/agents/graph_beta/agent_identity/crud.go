package graphBetaAgentIdentity

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

// Create handles the Create operation for agent identity resources.
//
// Operation: Creates a new agent identity service principal
// API Calls:
//   - POST /servicePrincipals/microsoft.graph.agentIdentity
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentidentity-post?view=graph-rest-beta
// Note: Uses custom request as SDK doesn't support the agentIdentity cast endpoint
func (r *AgentIdentityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AgentIdentityResourceModel

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
			"Error constructing agent identity",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Use custom request since the SDK doesn't support the cast endpoint for agentIdentity
	// POST /servicePrincipals/microsoft.graph.agentIdentity
	// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-post?view=graph-rest-beta
	config := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    "servicePrincipals/microsoft.graph.agentIdentity",
		RequestBody: requestBody,
	}

	createdResource, err := customrequests.PostRequest(
		ctx,
		r.client.GetAdapter(),
		config,
		CreateAgentIdentityResponseFactory(),
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// Extract the ID from the response
	if agentIdentity, ok := createdResource.(*AgentIdentityResponse); ok && agentIdentity.GetId() != nil {
		object.ID = types.StringValue(*agentIdentity.GetId())
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

	// Sponsors and owners are not immediately available after creation
	tflog.Debug(ctx, "Waiting 20 seconds for eventual consistency after create")
	time.Sleep(20 * time.Second)

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

// Read handles the Read operation for agent identity resources.
//
// Operation: Retrieves agent identity details including sponsors and owners
// API Calls:
//   - GET /servicePrincipals/{id}
//   - GET /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors
//   - GET /servicePrincipals/{id}/owners
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentidentity-get?view=graph-rest-beta
func (r *AgentIdentityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentIdentityResourceModel

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

	servicePrincipal, err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, servicePrincipal)

	// Fetch sponsors using custom request (SDK doesn't support the cast endpoint)
	// GET /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors
	// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-list-sponsors?view=graph-rest-beta
	sponsorConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          fmt.Sprintf("/servicePrincipals/%s/microsoft.graph.agentIdentity/sponsors", object.ID.ValueString()),
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

	// Fetch owners using SDK (owners endpoint works with standard API)
	// GET /servicePrincipals/{id}/owners
	// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-list-owners?view=graph-rest-beta
	owners, err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(object.ID.ValueString()).
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

// Update handles the Update operation for agent identity resources.
//
// Operation: Updates agent identity properties and manages sponsors and owners
// API Calls:
//   - PATCH /servicePrincipals/{id}
//   - DELETE /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors/{sponsorId}/$ref (for removed sponsors)
//   - POST /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors/$ref (for new sponsors)
//   - DELETE /servicePrincipals/{id}/owners/{ownerId}/$ref (for removed owners)
//   - POST /servicePrincipals/{id}/owners/$ref (for new owners)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentidentity-update?view=graph-rest-beta
func (r *AgentIdentityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentIdentityResourceModel
	var state AgentIdentityResourceModel

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
			"Error constructing agent identity for update",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		ServicePrincipals().
		ByServicePrincipalId(state.ID.ValueString()).
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

	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency after update")
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for agent identity resources.
//
// Operation: Deletes agent identity service principal with optional hard delete
// API Calls:
//   - DELETE /servicePrincipals/{id}
//   - GET /directory/deletedItems/microsoft.graph.servicePrincipal (if hard_delete is true)
//   - DELETE /directory/deletedItems/{id} (if hard_delete is true)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/agentidentity-delete?view=graph-rest-beta
// Note: When hard_delete is true, performs soft delete then permanent deletion from directory deleted items after eventual consistency delay
func (r *AgentIdentityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object AgentIdentityResourceModel

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

	agentIdentityId := object.ID.ValueString()

	softDeleteFunc := func(ctx context.Context) error {
		return r.client.
			ServicePrincipals().
			ByServicePrincipalId(agentIdentityId).
			Delete(ctx, nil)
	}

	deleteOpts := directory.DeleteOptions{
		MaxRetries:    10,
		RetryInterval: 5 * time.Second,
		ResourceType:  directory.ResourceTypeServicePrincipal,
		ResourceID:    agentIdentityId,
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
