package graphBetaApplicationCertificateCredential

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for application certificate credential resources.
//
// Operation: Adds a certificate credential to an application
// API Calls:
//   - GET /applications/{id} (to fetch existing keyCredentials if not replacing)
//   - PATCH /applications/{id} (to add certificate to keyCredentials array)
//   - GET /applications/{id} (to retrieve API-generated keyId)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-addkey?view=graph-rest-beta
// Note: This resource uses PATCH on the keyCredentials array rather than the addKey action
// to avoid the proof-of-possession requirement which complicates Terraform management.
func (r *ApplicationCertificateCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ApplicationCertificateCredentialResourceModel

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

	// Fetch existing credentials if not replacing all
	var existingCredentials []interface{}
	if !object.ReplaceExistingCertificates.ValueBool() {
		tflog.Debug(ctx, "Fetching existing key credentials to preserve them")

		existingApp, err := r.client.
			Applications().
			ByApplicationId(applicationID).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
			return
		}

		if existingApp != nil && existingApp.GetKeyCredentials() != nil {
			for _, cred := range existingApp.GetKeyCredentials() {
				existingCredentials = append(existingCredentials, cred)
			}
			tflog.Debug(ctx, fmt.Sprintf("Found %d existing key credentials", len(existingCredentials)))
		}
	}

	// Construct request body with new credential
	requestBody, err := constructResource(ctx, &object, existingCredentials)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing certificate credential request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Patching application with certificate for application_id: %s", applicationID))

	// Patch the application with the new keyCredentials array
	_, err = r.client.
		Applications().
		ByApplicationId(applicationID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency after create")
	time.Sleep(15 * time.Second)

	// Read back to get the API-generated key_id
	updatedApp, err := r.client.
		Applications().
		ByApplicationId(applicationID).
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
			"Could not find the newly created certificate credential in the application. This may indicate the certificate was not added successfully.",
		)
		return
	}

	object.KeyID = convert.GraphToFrameworkUUID(keyID)
	tflog.Debug(ctx, fmt.Sprintf("Found new certificate with key_id: %s", keyID.String()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read with retry for eventual consistency
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

// Read handles the Read operation for application certificate credential resources.
//
// Operation: Retrieves certificate credential by reading application keyCredentials array
// API Calls:
//   - GET /applications/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta
func (r *ApplicationCertificateCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ApplicationCertificateCredentialResourceModel

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

	applicationID := object.ApplicationID.ValueString()
	keyID := object.KeyID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with application_id: %s and key_id: %s", ResourceName, applicationID, keyID))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	application, err := r.client.
		Applications().
		ByApplicationId(applicationID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if application == nil {
		resp.Diagnostics.AddError(
			"Error reading certificate credential",
			fmt.Sprintf("Received nil application when reading application_id: %s", applicationID),
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

// Update handles the Update operation for application certificate credential resources.
//
// Operation: Updates certificate by removing old credential and adding new one
// API Calls:
//   - GET /applications/{id} (to fetch existing keyCredentials)
//   - PATCH /applications/{id} (to remove old certificate)
//   - PATCH /applications/{id} (to add new certificate)
//   - GET /applications/{id} (to retrieve API-generated keyId)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-addkey?view=graph-rest-beta
// Note: Certificate credentials cannot be updated directly; changes require delete and recreate.
// This is enforced at the schema level with RequiresReplace plan modifiers.
func (r *ApplicationCertificateCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationCertificateCredentialResourceModel
	var state ApplicationCertificateCredentialResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Updating certificate credential - removing old key_id: %s", oldKeyID))

	// Get existing credentials
	existingApp, err := r.client.
		Applications().
		ByApplicationId(applicationID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}

	existingCredentials := existingApp.GetKeyCredentials()

	// Step 1: Remove old credential
	var credsToKeep []interface{}
	for _, cred := range existingCredentials {
		if cred.GetKeyId() != nil && cred.GetKeyId().String() != oldKeyID {
			credsToKeep = append(credsToKeep, cred)
		}
	}

	deleteBody, err := constructDeleteResource(ctx, oldKeyID, existingCredentials)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing delete request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Applications().
		ByApplicationId(applicationID).
		Patch(ctx, deleteBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Successfully removed old certificate credential, waiting 5 seconds")
	time.Sleep(5 * time.Second)

	// Step 2: Add new credential
	addBody, err := constructResource(ctx, &plan, credsToKeep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing add request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Applications().
		ByApplicationId(applicationID).
		Patch(ctx, addBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Successfully added new certificate credential")
	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency after update")
	time.Sleep(15 * time.Second)

	// Read back to get the new key_id
	updatedApp, err := r.client.
		Applications().
		ByApplicationId(applicationID).
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

	// Read with retry for eventual consistency
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

// Delete handles the Delete operation for application certificate credential resources.
//
// Operation: Removes certificate credential from application keyCredentials array
// API Calls:
//   - GET /applications/{id} (to fetch existing keyCredentials)
//   - PATCH /applications/{id} (to remove certificate while preserving others)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-removekey?view=graph-rest-beta
func (r *ApplicationCertificateCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationCertificateCredentialResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Deleting certificate credential with key_id: %s from application: %s", keyID, applicationID))

	// Get existing credentials
	existingApp, err := r.client.
		Applications().
		ByApplicationId(applicationID).
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

	_, err = r.client.
		Applications().
		ByApplicationId(applicationID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted certificate credential with key_id: %s", keyID))
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
