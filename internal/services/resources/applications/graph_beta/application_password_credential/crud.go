package graphBetaApplicationPasswordCredential

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for application password credential resources.
//
// Operation: Adds a password credential to an application
// API Calls:
//   - POST /applications/{id}/addPassword
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-addpassword?view=graph-rest-beta
// Note: Password secret is only returned during creation and cannot be retrieved later
func (r *ApplicationPasswordCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ApplicationPasswordCredentialResourceModel

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

	applicationID := object.ApplicationID.ValueString()

	requestBody, err := constructAddPasswordRequest(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing password credential request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling addPassword for application_id: %s", applicationID))

	createdCredential, err := r.client.
		Applications().
		ByApplicationId(applicationID).
		AddPassword().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// Map response to state directly within create since the values are
	// not retrievable after creation.
	MapRemoteResourceStateToTerraform(ctx, &object, createdCredential)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform a Read with retry to verify the credential exists and handle eventual consistency
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

// Read handles the Read operation for application password credential resources.
//
// Operation: Fetches the parent application and reads the password credential from its collection
// API Calls:
//   - GET /applications/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta
// Note: The secret_text value cannot be retrieved from the API after creation, so it's preserved from state
func (r *ApplicationPasswordCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ApplicationPasswordCredentialResourceModel
	var identity sharedmodels.ResourceIdentity

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	// Determine the operation context (read vs create)
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

	applicationID := object.ApplicationID.ValueString()
	keyID := object.KeyID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with key_id: %s from application: %s (operation: %s)", ResourceName, keyID, applicationID, operation))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Preserve secret_text from state since it cannot be retrieved from API
	secretTextFromState := object.SecretText

	application, err := r.client.
		Applications().
		ByApplicationId(applicationID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	// Find the matching password credential by key_id
	passwordCredentials := application.GetPasswordCredentials()
	if passwordCredentials == nil {
		tflog.Error(ctx, fmt.Sprintf("No password credentials array found on application %s", applicationID))
		resp.Diagnostics.AddError(
			"Password credential not found",
			fmt.Sprintf("No password credentials found on application %s", applicationID),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found %d password credentials on application, looking for key_id: %s", len(passwordCredentials), keyID))

	var foundCredential graphmodels.PasswordCredentialable
	for i, cred := range passwordCredentials {
		credKeyID := ""
		if cred.GetKeyId() != nil {
			credKeyID = cred.GetKeyId().String()
		}
		tflog.Debug(ctx, fmt.Sprintf("Password credential [%d]: key_id=%s, display_name=%s", i, credKeyID, *cred.GetDisplayName()))

		if cred.GetKeyId() != nil && cred.GetKeyId().String() == keyID {
			foundCredential = cred
			tflog.Debug(ctx, fmt.Sprintf("Found matching password credential at index %d", i))
			break
		}
	}

	if foundCredential == nil {
		tflog.Error(ctx, fmt.Sprintf("Password credential with key_id %s not found on application %s after checking %d credentials", keyID, applicationID, len(passwordCredentials)))
		resp.Diagnostics.AddError(
			"Password credential not found",
			fmt.Sprintf("Password credential with key_id %s not found on application %s", keyID, applicationID),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, foundCredential)

	// Restore secret_text from state (cannot be retrieved from API)
	object.SecretText = secretTextFromState

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity.ID = object.ApplicationID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for application password credential resources.
//
// Operation: Updates password by removing old credential and adding new one
// API Calls:
//   - POST /applications/{id}/removePassword
//   - POST /applications/{id}/addPassword
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-removepassword?view=graph-rest-beta
// Note: Password credentials cannot be updated directly; changes require delete and recreate
func (r *ApplicationPasswordCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationPasswordCredentialResourceModel
	var state ApplicationPasswordCredentialResourceModel

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

	applicationID := state.ApplicationID.ValueString()
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

	err = r.client.
		Applications().
		ByApplicationId(applicationID).
		RemovePassword().
		Post(ctx, removeRequestBody, nil)

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

	createdCredential, err := r.client.
		Applications().
		ByApplicationId(applicationID).
		AddPassword().
		Post(ctx, addRequestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// Map response to state directly since the values are not retrievable after creation
	MapRemoteResourceStateToTerraform(ctx, &plan, createdCredential)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform a Read with retry to verify the credential exists and handle eventual consistency
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

// Delete handles the Delete operation for application password credential resources.
//
// Operation: Removes a password credential from an application
// API Calls:
//   - POST /applications/{id}/removePassword
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-removepassword?view=graph-rest-beta
func (r *ApplicationPasswordCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationPasswordCredentialResourceModel

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

	applicationID := data.ApplicationID.ValueString()
	keyID := data.KeyID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Deleting password credential with key_id: %s from application: %s", keyID, applicationID))

	removeRequestBody, err := constructRemovePasswordRequest(ctx, keyID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing remove password request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	err = r.client.
		Applications().
		ByApplicationId(applicationID).
		RemovePassword().
		Post(ctx, removeRequestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted password credential with key_id: %s", keyID))
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
