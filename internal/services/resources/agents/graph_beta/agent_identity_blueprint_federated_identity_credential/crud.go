package graphBetaAgentIdentityBlueprintFederatedIdentityCredential

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for agent identity blueprint federated identity credential resources.
//
// Operation: Creates a new federated identity credential for an agent identity blueprint
// API Calls:
//   - POST /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials
//
// Reference: https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-post?view=graph-rest-beta
// Note: Uses custom request as SDK doesn't support the agentIdentityBlueprint cast endpoint
func (r *AgentIdentityBlueprintFederatedIdentityCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AgentIdentityBlueprintFederatedIdentityCredentialResourceModel

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
			"Error constructing federated identity credential",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	blueprintID := object.BlueprintID.ValueString()

	// Use custom request with agentIdentityBlueprint cast endpoint
	endpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials", blueprintID)

	config := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    endpoint,
		RequestBody: requestBody,
	}

	result, err := customrequests.PostRequest(
		ctx,
		r.client.GetAdapter(),
		config,
		graphmodels.CreateFederatedIdentityCredentialFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	createdCredential, ok := result.(graphmodels.FederatedIdentityCredentialable)
	if !ok || createdCredential == nil || createdCredential.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error creating federated identity credential",
			"The API returned an invalid response",
		)
		return
	}

	object.ID = types.StringValue(*createdCredential.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Waiting 10 seconds for eventual consistency after create")
	time.Sleep(10 * time.Second)

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

// Read handles the Read operation for agent identity blueprint federated identity credential resources.
//
// Operation: Retrieves a federated identity credential by ID
// API Calls:
//   - GET /applications/{blueprintId}/federatedIdentityCredentials/{credentialId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-get?view=graph-rest-beta
// Note: Standard SDK endpoint works without cast for read operations
func (r *AgentIdentityBlueprintFederatedIdentityCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentIdentityBlueprintFederatedIdentityCredentialResourceModel

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

	blueprintID := object.BlueprintID.ValueString()
	credentialID := object.ID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Fetching federated identity credential - blueprint_id: %s, credential_id: %s", blueprintID, credentialID))

	// Use standard Kiota SDK - both standard and cast endpoints work
	credential, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		FederatedIdentityCredentials().
		ByFederatedIdentityCredentialId(credentialID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if credential == nil {
		resp.Diagnostics.AddError(
			"Error reading federated identity credential",
			fmt.Sprintf("Received nil credential for ID: %s", credentialID),
		)
		return
	}

	MapRemoteStateFromSDK(ctx, &object, credential)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for agent identity blueprint federated identity credential resources.
//
// Operation: Updates an existing federated identity credential
// API Calls:
//   - PATCH /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-update?view=graph-rest-beta
// Note: Uses custom request as SDK doesn't support the agentIdentityBlueprint cast endpoint
func (r *AgentIdentityBlueprintFederatedIdentityCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentIdentityBlueprintFederatedIdentityCredentialResourceModel
	var state AgentIdentityBlueprintFederatedIdentityCredentialResourceModel

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

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing federated identity credential",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	blueprintID := state.BlueprintID.ValueString()
	credentialID := state.ID.ValueString()

	// Use custom request with agentIdentityBlueprint cast endpoint
	config := customrequests.PatchRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          fmt.Sprintf("/applications/%s/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials", blueprintID),
		ResourceID:        credentialID,
		ResourceIDPattern: "/{id}",
		RequestBody:       requestBody,
	}

	err = customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), config)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency after update")
	time.Sleep(15 * time.Second)

	plan.ID = state.ID
	plan.BlueprintID = state.BlueprintID

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

// Delete handles the Delete operation for agent identity blueprint federated identity credential resources.
//
// Operation: Deletes a federated identity credential
// API Calls:
//   - DELETE /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-delete?view=graph-rest-beta
// Note: Uses custom request as SDK doesn't support the agentIdentityBlueprint cast endpoint
func (r *AgentIdentityBlueprintFederatedIdentityCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AgentIdentityBlueprintFederatedIdentityCredentialResourceModel

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

	blueprintID := data.BlueprintID.ValueString()
	credentialID := data.ID.ValueString()

	// Use custom request with agentIdentityBlueprint cast endpoint
	config := customrequests.DeleteRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          fmt.Sprintf("/applications/%s/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials", blueprintID),
		ResourceID:        credentialID,
		ResourceIDPattern: "/{id}",
	}

	err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), config)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
