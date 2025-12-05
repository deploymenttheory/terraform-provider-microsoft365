package graphBetaAgentIdentityBlueprintKeyCredential

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
// Uses POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/addKey
// Then calls Read with retry for stating.
func (r *AgentIdentityBlueprintKeyCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AgentIdentityBlueprintKeyCredentialResourceModel

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

	blueprintID := object.BlueprintID.ValueString()

	requestBody, err := constructAddKeyRequest(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing key credential request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling addKey for blueprint_id: %s", blueprintID))

	endpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/addKey", blueprintID)

	config := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    endpoint,
		RequestBody: requestBody,
	}

	result, err := customrequests.PostRequest(
		ctx,
		r.client.GetAdapter(),
		config,
		graphmodels.CreateKeyCredentialFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	createdCredential, ok := result.(graphmodels.KeyCredentialable)
	if !ok {
		resp.Diagnostics.AddError(
			"Error creating key credential",
			fmt.Sprintf("Failed to cast response to KeyCredentialable, got type: %T", result),
		)
		return
	}

	// Set the key_id from the response for Read to use
	if createdCredential.GetKeyId() != nil {
		object.KeyID = types.StringValue(createdCredential.GetKeyId().String())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait for eventual consistency
	tflog.Debug(ctx, "Waiting 10 seconds for eventual consistency after create")
	time.Sleep(10 * time.Second)

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
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
// Uses GET /applications/{id}/microsoft.graph.agentIdentityBlueprint?$select=keyCredentials
// to fetch key credentials and find the matching one by keyId.
func (r *AgentIdentityBlueprintKeyCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentIdentityBlueprintKeyCredentialResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with key_id: %s", ResourceName, object.KeyID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	blueprintID := object.BlueprintID.ValueString()
	keyID := object.KeyID.ValueString()

	// Fetch application with keyCredentials using custom request with OData $select
	getConfig := customrequests.GetRequestConfig{
		APIVersion: customrequests.GraphAPIBeta,
		Endpoint:   fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint", blueprintID),
		QueryParameters: map[string]string{
			"$select": "keyCredentials",
		},
	}

	result, err := customrequests.GetRequest(
		ctx,
		r.client.GetAdapter(),
		getConfig,
		graphmodels.CreateApplicationFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	application, ok := result.(graphmodels.Applicationable)
	if !ok {
		resp.Diagnostics.AddError(
			"Error reading key credentials",
			fmt.Sprintf("Failed to cast response to Applicationable, got type: %T", result),
		)
		return
	}

	// Map the response to find our key credential
	if err := MapRemoteResourceStateToTerraform(ctx, &object, application, keyID); err != nil {
		resp.Diagnostics.AddError(
			"Error mapping key credential state",
			fmt.Sprintf("Could not map key credential: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
// Since key credentials cannot be updated in-place, this performs:
// 1. POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/removeKey (delete old)
// 2. POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/addKey (create new)
// 3. Calls Read with retry for stating
func (r *AgentIdentityBlueprintKeyCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentIdentityBlueprintKeyCredentialResourceModel
	var state AgentIdentityBlueprintKeyCredentialResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

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

	blueprintID := state.BlueprintID.ValueString()
	oldKeyID := state.KeyID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Updating key credential - removing old key_id: %s", oldKeyID))

	// Step 1: Remove the old key credential
	removeRequestBody, err := constructRemoveKeyRequest(ctx, oldKeyID, plan.Proof.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing remove key request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	removeEndpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/removeKey", blueprintID)

	removeConfig := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    removeEndpoint,
		RequestBody: removeRequestBody,
	}

	err = customrequests.PostRequestNoContent(ctx, r.client.GetAdapter(), removeConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update (remove)", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Successfully removed old key credential")

	// Step 2: Add the new key credential
	addRequestBody, err := constructAddKeyRequest(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing add key request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	addEndpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/addKey", blueprintID)

	addConfig := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    addEndpoint,
		RequestBody: addRequestBody,
	}

	result, err := customrequests.PostRequest(
		ctx,
		r.client.GetAdapter(),
		addConfig,
		graphmodels.CreateKeyCredentialFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update (add)", r.WritePermissions)
		return
	}

	createdCredential, ok := result.(graphmodels.KeyCredentialable)
	if !ok {
		resp.Diagnostics.AddError(
			"Error creating new key credential",
			fmt.Sprintf("Failed to cast response to KeyCredentialable, got type: %T", result),
		)
		return
	}

	// Set the new key_id from the response for Read to use
	if createdCredential.GetKeyId() != nil {
		plan.KeyID = types.StringValue(createdCredential.GetKeyId().String())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait for eventual consistency
	tflog.Debug(ctx, "Waiting 10 seconds for eventual consistency after update")
	time.Sleep(10 * time.Second)

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
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

// Delete handles the Delete operation.
// Uses POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/removeKey
func (r *AgentIdentityBlueprintKeyCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AgentIdentityBlueprintKeyCredentialResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete method for: %s", ResourceName))

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
	keyID := data.KeyID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Deleting key credential with key_id: %s from blueprint: %s", keyID, blueprintID))

	removeRequestBody, err := constructRemoveKeyRequest(ctx, keyID, data.Proof.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing remove key request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	endpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/removeKey", blueprintID)

	config := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    endpoint,
		RequestBody: removeRequestBody,
	}

	err = customrequests.PostRequestNoContent(ctx, r.client.GetAdapter(), config)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted key credential with key_id: %s", keyID))
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
