package graphBetaAgentIdentityBlueprintCertificateCredential

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for certificate credentials.
//
// Key credentials don't have a dedicated POST endpoint - they're managed as a property
// of the application object via PATCH. This means we must:
//
//  1. Fetch existing keyCredentials first (unless replacing all), because PATCH replaces
//     the entire keyCredentials array, not just append to it.
//
//  2. PATCH the application with the combined credentials using the agentIdentityBlueprint
//     cast endpoint (required for agent blueprints, standard endpoint returns 400).
//
//  3. Fetch the application again to discover the API-generated keyId, since the API
//     generates the keyId - we cannot set it ourselves. We match by displayName.
//
// 4. Call Read with retry to handle eventual consistency and populate full state.
func (r *AgentIdentityBlueprintCertificateCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AgentIdentityBlueprintCertificateCredentialResourceModel

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

	var existingCredentials []graphmodels.KeyCredentialable
	if !object.ReplaceExistingCertificates.ValueBool() {
		tflog.Debug(ctx, "Fetching existing key credentials to preserve them")

		existingApp, err := r.client.
			Applications().
			ByApplicationId(blueprintID).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
			return
		}

		if existingApp != nil {
			existingCredentials = existingApp.GetKeyCredentials()
			tflog.Debug(ctx, fmt.Sprintf("Found %d existing key credentials", len(existingCredentials)))
		}
	}

	requestBody, err := constructResource(ctx, &object, existingCredentials)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing certificate credential request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Patching application with certificate for blueprint_id: %s", blueprintID))

	config := customrequests.PatchRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "applications",
		ResourceID:        blueprintID,
		ResourceIDPattern: "/{id}",
		EndpointSuffix:    "/microsoft.graph.agentIdentityBlueprint",
		RequestBody:       requestBody,
	}

	err = customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), config)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency after create")
	time.Sleep(15 * time.Second)

	updatedApp, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}

	newCredentials := updatedApp.GetKeyCredentials()
	keyID := FindKeyCredentialByDisplayName(newCredentials, object.DisplayName.ValueString())

	if keyID == nil {
		resp.Diagnostics.AddError(
			"Error finding new certificate credential",
			"Could not find the newly created certificate credential in the application",
		)
		return
	}

	object.KeyID = convert.GraphToFrameworkUUID(keyID)
	tflog.Debug(ctx, fmt.Sprintf("Found new certificate with key_id: %s", keyID.String()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
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

// Read handles the Read operation.
// Uses GET /applications/{id} to retrieve the key credentials.
func (r *AgentIdentityBlueprintCertificateCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentIdentityBlueprintCertificateCredentialResourceModel

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

	blueprintID := object.BlueprintID.ValueString()
	keyID := object.KeyID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with blueprint_id: %s and key_id: %s", ResourceName, blueprintID, keyID))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	application, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if application == nil {
		resp.Diagnostics.AddError(
			"Error reading certificate credential",
			fmt.Sprintf("Received nil application when reading blueprint_id: %s", blueprintID),
		)
		return
	}

	// Check if the credential still exists
	credential := FindKeyCredentialByKeyID(application.GetKeyCredentials(), keyID)
	if credential == nil {
		if operation == constants.TfOperationRead {
			tflog.Warn(ctx, fmt.Sprintf("Credential with keyId %s not found, removing from state", keyID))
			resp.State.RemoveResource(ctx)
			return
		}
		// During Create/Update, return error to allow retry
		resp.Diagnostics.AddError(
			"Credential not found",
			fmt.Sprintf("Credential with keyId %s not yet available, retry may be needed", keyID),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, credential)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
// Certificate credentials cannot be updated in-place. This performs a delete and recreate.
func (r *AgentIdentityBlueprintCertificateCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentIdentityBlueprintCertificateCredentialResourceModel
	var state AgentIdentityBlueprintCertificateCredentialResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Updating certificate credential - removing old key_id: %s", oldKeyID))

	existingApp, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}

	existingCredentials := existingApp.GetKeyCredentials()

	deleteBody, err := constructDeleteResource(ctx, oldKeyID, existingCredentials)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing delete request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	deleteConfig := customrequests.PatchRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "applications",
		ResourceID:        blueprintID,
		ResourceIDPattern: "/{id}",
		EndpointSuffix:    "/microsoft.graph.agentIdentityBlueprint",
		RequestBody:       deleteBody,
	}

	err = customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Successfully removed old certificate credential")

	time.Sleep(5 * time.Second)

	existingApp, err = r.client.Applications().ByApplicationId(blueprintID).Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}

	existingCredentials = existingApp.GetKeyCredentials()

	addBody, err := constructResource(ctx, &plan, existingCredentials)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing add request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	addConfig := customrequests.PatchRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "applications",
		ResourceID:        blueprintID,
		ResourceIDPattern: "/{id}",
		EndpointSuffix:    "/microsoft.graph.agentIdentityBlueprint",
		RequestBody:       addBody,
	}

	err = customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), addConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Successfully added new certificate credential")
	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency after create")
	time.Sleep(15 * time.Second)

	updatedApp, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}

	newCredentials := updatedApp.GetKeyCredentials()
	keyID := FindKeyCredentialByDisplayName(newCredentials, plan.DisplayName.ValueString())

	if keyID == nil {
		resp.Diagnostics.AddError(
			"Error finding new certificate credential",
			"Could not find the newly created certificate credential in the application",
		)
		return
	}

	plan.KeyID = convert.GraphToFrameworkUUID(keyID)
	tflog.Debug(ctx, fmt.Sprintf("Found new certificate with key_id: %s", keyID.String()))

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

// Delete handles the Delete operation.
// Uses PATCH /applications/{id}/microsoft.graph.agentIdentityBlueprint to remove the certificate
// from the keyCredentials property while preserving other credentials.
func (r *AgentIdentityBlueprintCertificateCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AgentIdentityBlueprintCertificateCredentialResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Deleting certificate credential with key_id: %s from blueprint: %s", keyID, blueprintID))

	// Get existing credentials
	existingApp, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
		return
	}

	existingCredentials := existingApp.GetKeyCredentials()

	// Construct request body without the credential being deleted
	requestBody, err := constructDeleteResource(ctx, keyID, existingCredentials)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing delete request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	config := customrequests.PatchRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "applications",
		ResourceID:        blueprintID,
		ResourceIDPattern: "/{id}",
		EndpointSuffix:    "/microsoft.graph.agentIdentityBlueprint",
		RequestBody:       requestBody,
	}

	err = customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), config)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted certificate credential with key_id: %s", keyID))
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
