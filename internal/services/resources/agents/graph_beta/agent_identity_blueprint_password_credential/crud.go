package graphBetaAgentIdentityBlueprintPasswordCredential

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
// Uses POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/addPassword
// Maps the response directly to state since there is no GET endpoint for password credentials.
func (r *AgentIdentityBlueprintPasswordCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AgentIdentityBlueprintPasswordCredentialResourceModel

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

	requestBody, err := constructAddPasswordRequest(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing password credential request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling addPassword for blueprint_id: %s", blueprintID))

	// Use custom request with agentIdentityBlueprint cast endpoint
	endpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/addPassword", blueprintID)

	config := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    endpoint,
		RequestBody: requestBody,
	}

	result, err := customrequests.PostRequest(
		ctx,
		r.client.GetAdapter(),
		config,
		graphmodels.CreatePasswordCredentialFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	createdCredential, ok := result.(graphmodels.PasswordCredentialable)
	if !ok {
		resp.Diagnostics.AddError(
			"Error creating password credential",
			fmt.Sprintf("Failed to cast response to PasswordCredentialable, got type: %T", result),
		)
		return
	}

	// Map response to state directly within create since, the values are
	// not retrievable after creation.
	MapRemoteResourceStateToTerraform(ctx, &object, createdCredential)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation.
// Since there is no GET endpoint for individual password credentials, this function
// simply returns the current state. The state is preserved from Create/Update operations.
func (r *AgentIdentityBlueprintPasswordCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentIdentityBlueprintPasswordCredentialResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

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

	// No API call - just return current state
	// The password credential API does not have a GET endpoint for individual credentials.
	// The secret_text cannot be retrieved after creation, so we preserve the state as-is.

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
// Since password credentials cannot be updated in-place, this performs:
// 1. POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/removePassword (delete old)
// 2. POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/addPassword (create new)
func (r *AgentIdentityBlueprintPasswordCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentIdentityBlueprintPasswordCredentialResourceModel
	var state AgentIdentityBlueprintPasswordCredentialResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Updating password credential - removing old key_id: %s", oldKeyID))

	// Step 1: Remove the old password credential
	removeRequestBody, err := constructRemovePasswordRequest(ctx, oldKeyID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing remove password request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	removeEndpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/removePassword", blueprintID)

	removeConfig := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    removeEndpoint,
		RequestBody: removeRequestBody,
	}

	err = customrequests.PostRequestNoContent(ctx, r.client.GetAdapter(), removeConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Successfully removed old password credential")

	// Step 2: Add the new password credential
	addRequestBody, err := constructAddPasswordRequest(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing add password request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	addEndpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/addPassword", blueprintID)

	addConfig := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    addEndpoint,
		RequestBody: addRequestBody,
	}

	result, err := customrequests.PostRequest(
		ctx,
		r.client.GetAdapter(),
		addConfig,
		graphmodels.CreatePasswordCredentialFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	createdCredential, ok := result.(graphmodels.PasswordCredentialable)
	if !ok {
		resp.Diagnostics.AddError(
			"Error creating new password credential",
			fmt.Sprintf("Failed to cast response to PasswordCredentialable, got type: %T", result),
		)
		return
	}

	// Map response to state
	MapRemoteResourceStateToTerraform(ctx, &plan, createdCredential)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation.
// Uses POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/removePassword
func (r *AgentIdentityBlueprintPasswordCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AgentIdentityBlueprintPasswordCredentialResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Deleting password credential with key_id: %s from blueprint: %s", keyID, blueprintID))

	removeRequestBody, err := constructRemovePasswordRequest(ctx, keyID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing remove password request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	endpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/removePassword", blueprintID)

	config := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    endpoint,
		RequestBody: removeRequestBody,
	}

	err = customrequests.PostRequestNoContent(ctx, r.client.GetAdapter(), config)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted password credential with key_id: %s", keyID))
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
