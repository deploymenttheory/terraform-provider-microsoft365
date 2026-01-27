package graphBetaAgentIdentityBlueprintIdentifierUri

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
)

// Create handles the Create operation for agent identity blueprint identifier URI resources.
//
// Operation: Adds an identifier URI and optional scope to an agent identity blueprint
// API Calls:
//   - GET /applications/{id} (to fetch existing identifierUris)
//   - PATCH /applications/{id} (to add identifier URI to array)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
func (r *AgentIdentityBlueprintIdentifierUriResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AgentIdentityBlueprintIdentifierUriResourceModel

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

	// First, get the current application to retrieve existing identifier URIs and api configuration
	// Note: We don't use $expand here as 'api' is a complex type, not a navigation property
	// Instead, we ensure all properties are returned by not using $select
	requestConfig := &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{},
	}

	application, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		Get(ctx, requestConfig)

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

	tflog.Debug(ctx, fmt.Sprintf("Adding identifier URI to blueprint_id: %s", blueprintID))

	_, err = r.client.
		Applications().
		ByApplicationId(blueprintID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Waiting 10 seconds for eventual consistency after create")
	time.Sleep(10 * time.Second)

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

// Read handles the Read operation for agent identity blueprint identifier URI resources.
//
// Operation: Retrieves application to verify identifier URI and scope exist
// API Calls:
//   - GET /applications/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta
func (r *AgentIdentityBlueprintIdentifierUriResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentIdentityBlueprintIdentifierUriResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with identifier_uri: %s", ResourceName, object.IdentifierUri.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	blueprintID := object.BlueprintID.ValueString()

	// Get the application with api property fully populated including oauth2PermissionScopes
	// Note: We don't use $expand or $select to ensure all properties including nested ones are returned
	requestConfig := &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{},
	}

	application, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, application)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for agent identity blueprint identifier URI resources.
//
// Operation: Updates scope configuration (identifier_uri changes trigger replacement)
// API Calls:
//   - GET /applications/{id} (to fetch existing identifierUris)
//   - PATCH /applications/{id} (to update scope configuration)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
// Note: Changes to identifier_uri attribute trigger resource replacement
func (r *AgentIdentityBlueprintIdentifierUriResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentIdentityBlueprintIdentifierUriResourceModel
	var state AgentIdentityBlueprintIdentifierUriResourceModel

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

	requestConfig := &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{},
	}

	application, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	existingUris := application.GetIdentifierUris()

	requestBody, err := constructResource(ctx, &plan, existingUris)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing identifier URI resource",
			fmt.Sprintf("Could not construct resource: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating identifier URI configuration for blueprint_id: %s", blueprintID))

	_, err = r.client.
		Applications().
		ByApplicationId(blueprintID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

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

// Delete handles the Delete operation for agent identity blueprint identifier URI resources.
//
// Operation: Removes identifier URI from application identifierUris array
// API Calls:
//   - GET /applications/{id} (to fetch existing identifierUris)
//   - PATCH /applications/{id} (to remove identifier URI from array)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
func (r *AgentIdentityBlueprintIdentifierUriResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AgentIdentityBlueprintIdentifierUriResourceModel

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
	identifierUri := data.IdentifierUri.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Removing identifier URI %s from blueprint: %s", identifierUri, blueprintID))

	requestConfig := &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{},
	}

	application, err := r.client.
		Applications().
		ByApplicationId(blueprintID).
		Get(ctx, requestConfig)

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
		ByApplicationId(blueprintID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed identifier URI: %s", identifierUri))
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
