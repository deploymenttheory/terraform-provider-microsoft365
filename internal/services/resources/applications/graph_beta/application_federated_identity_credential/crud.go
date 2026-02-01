package graphBetaApplicationFederatedIdentityCredential

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

// Create handles the Create operation for application federated identity credential resources.
//
// Operation: Creates a new federated identity credential for an application
// API Calls:
//   - POST /applications/{applicationId}/federatedIdentityCredentials
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-post-federatedidentitycredentials?view=graph-rest-beta
func (r *ApplicationFederatedIdentityCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ApplicationFederatedIdentityCredentialResourceModel

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

	applicationID := object.ApplicationID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Creating federated identity credential for application: %s", applicationID))

	// Use standard Kiota SDK
	createdCredential, err := r.client.
		Applications().
		ByApplicationId(applicationID).
		FederatedIdentityCredentials().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if createdCredential == nil || createdCredential.GetId() == nil {
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

// Read handles the Read operation for application federated identity credential resources.
//
// Operation: Retrieves a federated identity credential by ID
// API Calls:
//   - GET /applications/{applicationId}/federatedIdentityCredentials/{credentialId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-get?view=graph-rest-beta
func (r *ApplicationFederatedIdentityCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ApplicationFederatedIdentityCredentialResourceModel

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

	applicationID := object.ApplicationID.ValueString()
	credentialID := object.ID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Fetching federated identity credential - application_id: %s, credential_id: %s", applicationID, credentialID))

	// Use standard Kiota SDK
	credential, err := r.client.
		Applications().
		ByApplicationId(applicationID).
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

// Update handles the Update operation for application federated identity credential resources.
//
// Operation: Updates an existing federated identity credential
// API Calls:
//   - PATCH /applications/{applicationId}/federatedIdentityCredentials/{credentialId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-update?view=graph-rest-beta
func (r *ApplicationFederatedIdentityCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationFederatedIdentityCredentialResourceModel
	var state ApplicationFederatedIdentityCredentialResourceModel

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

	applicationID := state.ApplicationID.ValueString()
	credentialID := state.ID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Updating federated identity credential - application_id: %s, credential_id: %s", applicationID, credentialID))

	_, err = r.client.
		Applications().
		ByApplicationId(applicationID).
		FederatedIdentityCredentials().
		ByFederatedIdentityCredentialId(credentialID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency after update")
	time.Sleep(15 * time.Second)

	plan.ID = state.ID
	plan.ApplicationID = state.ApplicationID

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

// Delete handles the Delete operation for application federated identity credential resources.
//
// Operation: Deletes a federated identity credential
// API Calls:
//   - DELETE /applications/{applicationId}/federatedIdentityCredentials/{credentialId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-delete?view=graph-rest-beta
func (r *ApplicationFederatedIdentityCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationFederatedIdentityCredentialResourceModel

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

	applicationID := data.ApplicationID.ValueString()
	credentialID := data.ID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Deleting federated identity credential - application_id: %s, credential_id: %s", applicationID, credentialID))

	// Use standard Kiota SDK
	err := r.client.
		Applications().
		ByApplicationId(applicationID).
		FederatedIdentityCredentials().
		ByFederatedIdentityCredentialId(credentialID).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
