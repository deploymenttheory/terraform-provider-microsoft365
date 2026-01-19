package graphBetaNamedLocation

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

// Create handles the Create operation for Named Location resources.
func (r *NamedLocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NamedLocationResourceModel

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
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		Identity().
		ConditionalAccess().
		NamedLocations().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if createdResource == nil || createdResource.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error extracting resource ID",
			"Created resource ID is missing from response",
		)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())
	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.MaxRetries = 60                 // Up from default 30
	opts.RetryInterval = 5 * time.Second // Up from default 2 seconds

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

// Read handles the Read operation for Named Location resources.
func (r *NamedLocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NamedLocationResourceModel

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

	tflog.Debug(ctx, "Making GET request to retrieve named location")

	remoteResource, err := r.client.
		Identity().
		ConditionalAccess().
		NamedLocations().
		ByNamedLocationId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, remoteResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Named Location resources.
func (r *NamedLocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NamedLocationResourceModel
	var state NamedLocationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Making PATCH request to update named location")

	_, err = r.client.
		Identity().
		ConditionalAccess().
		NamedLocations().
		ByNamedLocationId(state.ID.ValueString()).
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
	opts.MaxRetries = 60                 // Up from default 30
	opts.RetryInterval = 5 * time.Second // Up from default 2 seconds

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for Named Location resources.
//
// This function implements a specialized deletion workflow required by Microsoft Graph's
// Named Location API constraints. The complexity exists because:
//
// 1. TRUSTED IP NAMED LOCATIONS CANNOT BE DELETED DIRECTLY
//   - Microsoft Graph API will reject DELETE requests for IP Named Locations with isTrusted=true
//   - This is a security feature to prevent accidental deletion of trusted network locations
//   - The API requires isTrusted to be explicitly set to false before deletion is allowed
//
// 2. EVENTUAL CONSISTENCY CHALLENGES
//
//   - Microsoft Graph API exhibits eventual consistency behavior
//
//   - A PATCH request to set isTrusted=false may not immediately take effect
//
//   - Subsequent GET requests may still show isTrusted=true for a period of time
//
//   - Attempting DELETE before the change propagates will fail
//
//     3. DELETION WORKFLOW FOR TRUSTED IP LOCATIONS:
//     Step 1: GET resource and check if it's an ipNamedLocation with isTrusted=true
//     Step 2: If conditions met, PATCH to set isTrusted=false
//     Step 3: Poll with GET requests until isTrusted=false is confirmed (eventual consistency)
//     Step 4: Execute DELETE operation
//     Step 5: Remove from Terraform state
//
// 4. DELETION WORKFLOW FOR OTHER NAMED LOCATIONS:
//   - Country Named Locations and non-trusted IP locations can be deleted directly
//   - Skip steps 2-3 and proceed directly to DELETE operation
//
// This approach ensures reliable deletion across all Named Location types while handling
// the API's security constraints and eventual consistency behavior.
func (r *NamedLocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NamedLocationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if !r.handleTrustedIPLocation(ctx, state.ID.ValueString(), resp) {
		return
	}

	if !r.waitForConditionalAccessPolicyReferences(ctx, state.ID.ValueString(), resp) {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Making DELETE request for %s with ID: %s", ResourceName, state.ID.ValueString()))

	deleteOptions := crud.DefaultDeleteWithRetryOptions()
	deleteOptions.ResourceTypeName = ResourceName
	deleteOptions.ResourceID = state.ID.ValueString()
	deleteOptions.RetryInterval = 10 * time.Second
	deleteOptions.MaxRetries = 6

	err := crud.DeleteWithRetry(ctx, func(ctx context.Context) error {
		tflog.Debug(ctx, fmt.Sprintf("Executing DELETE call for named location %s", state.ID.ValueString()))
		deleteErr := r.client.
			Identity().
			ConditionalAccess().
			NamedLocations().
			ByNamedLocationId(state.ID.ValueString()).
			Delete(ctx, nil)

		if deleteErr != nil {
			errorInfo := errors.GraphError(ctx, deleteErr)
			tflog.Debug(ctx, fmt.Sprintf("DELETE call returned error: status=%d, category=%s, message=%s",
				errorInfo.StatusCode, errorInfo.Category, errorInfo.ErrorMessage))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("DELETE call succeeded for named location %s", state.ID.ValueString()))
		}
		return deleteErr
	}, deleteOptions)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		tflog.Error(ctx, fmt.Sprintf("DeleteWithRetry failed: status=%d, category=%s, code=%s, message=%s",
			errorInfo.StatusCode, errorInfo.Category, errorInfo.ErrorCode, errorInfo.ErrorMessage))
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted %s with ID: %s", ResourceName, state.ID.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
