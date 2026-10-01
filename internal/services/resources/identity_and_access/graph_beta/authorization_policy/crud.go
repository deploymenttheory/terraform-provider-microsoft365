package graphBetaAuthorizationPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create configures the authorization policy.
//
// The authorizationPolicy resource is a singleton that always exists in the tenant; there is no documented
// POST endpoint to create it. This operation therefore uses a PATCH request to apply the desired
// configuration, mirroring the Update behaviour.
//
// API Calls:
//   - PATCH /policies/authorizationPolicy/authorizationPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/authorizationpolicy-update?view=graph-rest-beta
func (r *AuthorizationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AuthorizationPolicyResourceModel

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

	_, err = r.client.
		Policies().
		AuthorizationPolicy().
		ByAuthorizationPolicyId(singletonID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// This is a singleton resource — set the static ID immediately so the subsequent Read
	// can locate the correct state entry.
	object.ID = types.StringValue(singletonID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = consistencyPredicate(&object)

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

}

// Read retrieves the current state of the authorization policy from the Graph API.
//
// API Calls:
//   - GET /policies/authorizationPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/authorizationpolicy-get?view=graph-rest-beta
func (r *AuthorizationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AuthorizationPolicyResourceModel

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

	remoteResource, err := r.client.
		Policies().
		AuthorizationPolicy().
		Get(ctx, nil)

	if err != nil {
		// A singleton cannot be deleted by this resource. Keep state on unavailable or malformed reads.
		diagnosticResponse := &resource.UpdateResponse{}
		errors.HandleKiotaGraphError(ctx, err, diagnosticResponse, operation, r.ReadPermissions)
		resp.Diagnostics.Append(diagnosticResponse.Diagnostics...)
		return
	}

	if err := mapRemoteState(ctx, &object, remoteResource); err != nil {
		resp.Diagnostics.AddError("Invalid authorization policy response", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update applies changes to the authorization policy.
//
// API Calls:
//   - PATCH /policies/authorizationPolicy/authorizationPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/authorizationpolicy-update?view=graph-rest-beta
func (r *AuthorizationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AuthorizationPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	_, err = r.client.
		Policies().
		AuthorizationPolicy().
		ByAuthorizationPolicyId(singletonID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(singletonID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = consistencyPredicate(&plan)

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

}

// Delete releases Terraform management without changing the tenant-wide settings.
func (r *AuthorizationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}
