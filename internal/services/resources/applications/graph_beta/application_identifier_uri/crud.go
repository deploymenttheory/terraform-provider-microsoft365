package graphBetaApplicationIdentifierUri

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

// Create handles the Create operation for Application Identifier URI resources.
//
// Operation: Adds an identifier URI to an application
// API Calls:
//   - GET /applications/{applicationId} (to retrieve existing URIs)
//   - PATCH /applications/{applicationId} (to update with new URI)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
func (r *ApplicationIdentifierUriResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ApplicationIdentifierUriResourceModel

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

	application, err := r.client.
		Applications().
		ByApplicationId(applicationID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	existingUris := application.GetIdentifierUris()

	for _, uri := range existingUris {
		if uri == object.IdentifierUri.ValueString() {
			resp.Diagnostics.AddError(
				"Identifier URI already exists",
				fmt.Sprintf("The identifier URI %s already exists on the application", object.IdentifierUri.ValueString()),
			)
			return
		}
	}

	requestBody, err := constructResource(ctx, &object, existingUris)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing identifier URI resource",
			fmt.Sprintf("Could not construct resource: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Adding identifier URI to application_id: %s", applicationID))

	_, err = r.client.
		Applications().
		ByApplicationId(applicationID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Waiting 10 seconds for eventual consistency after create")
	time.Sleep(10 * time.Second)

	object.Id = types.StringValue(fmt.Sprintf("%s/%s", object.ApplicationID.ValueString(), object.IdentifierUri.ValueString()))

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

// Read handles the Read operation for Application Identifier URI resources.
//
// Operation: Retrieves application to verify identifier URI exists
// API Calls:
//   - GET /applications/{applicationId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta
func (r *ApplicationIdentifierUriResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ApplicationIdentifierUriResourceModel

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

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	applicationID := object.ApplicationID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with identifier_uri: %s", ResourceName, object.IdentifierUri.ValueString()))

	application, err := r.client.
		Applications().
		ByApplicationId(applicationID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, application)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Application Identifier URI resources.
//
// Operation: Not supported - identifier URIs require recreation
// Note: Both application_id and identifier_uri are ForceNew, so updates will not be called
func (r *ApplicationIdentifierUriResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This should never be called since both attributes are ForceNew
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Application identifier URI resources cannot be updated. Both application_id and identifier_uri require recreation.",
	)
}

// Delete handles the Delete operation for Application Identifier URI resources.
//
// Operation: Removes an identifier URI from an application
// API Calls:
//   - GET /applications/{applicationId} (to retrieve existing URIs)
//   - PATCH /applications/{applicationId} (to update without the URI)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
func (r *ApplicationIdentifierUriResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object ApplicationIdentifierUriResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	applicationID := object.ApplicationID.ValueString()
	identifierUri := object.IdentifierUri.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Removing identifier URI %s from application: %s", identifierUri, applicationID))

	application, err := r.client.
		Applications().
		ByApplicationId(applicationID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	existingUris := application.GetIdentifierUris()

	requestBody, err := constructDeleteResource(ctx, identifierUri, existingUris)
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

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed identifier URI: %s", identifierUri))
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
