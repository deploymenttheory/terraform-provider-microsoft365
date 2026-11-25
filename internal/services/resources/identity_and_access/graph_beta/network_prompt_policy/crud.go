// Package graphBetaNetworkPromptPolicy manages Global Secure Access prompt policies.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no promptPolicy model in models/networkaccess and
//     no promptPolicies request builder in networkaccess.
//   - Learn: none. There is no networkaccess-promptpolicy resource page; the Global Secure Access API overview does not
//     list prompt policies.
//   - Observed: the request and response shape for /networkaccess/promptPolicies were taken from the Entra admin centre
//     Global Secure Access prompt policy blade (checked by reading the SDK source on 2026-11-24; no live call, because
//     the test tenant has no Global Secure Access licence). Every call therefore goes through the custom_requests
//     Config API in raw JSON mode, with APIVersion left empty so it uses the beta client's own base URL.
//
// Removal condition: the beta SDK ships a networkaccess promptPolicy model and a /networkAccess/promptPolicies request
// builder.
package graphBetaNetworkPromptPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// readbackMatch is what the read-back after a write waits for: a populated id and the written fields.
var readbackMatch = crud.AttributeMatchOptions{Paths: crud.RootPaths("name", "description", "default_action"), RequireKnown: crud.RootPaths("id")}

func (r *NetworkPromptPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkPromptPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	var created *promptPolicyResponse
	postConfig := customrequests.PostRequestConfig{Endpoint: promptPoliciesEndpoint, RawBody: constructResource(ctx, &object, nil), RawResponse: &created}
	_, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), postConfig, nil, nil)
	if err == nil && created == nil {
		err = errEmptyResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.ID == nil || *created.ID == "" {
		resp.Diagnostics.AddError("Invalid create response", "The API returned no resource id. Check Graph before retrying creation.")
		return
	}
	// Persist the response body identity before readback. The API Location header may contain a placeholder.
	object.ID = types.StringValue(*created.ID)
	object.Version = types.StringNull()
	object.LastModifiedDateTime = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithConsistencyPredicate(crud.AttributesMatchPredicate(resp.State, readbackMatch)))
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after create", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

func (r *NetworkPromptPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkPromptPolicyResourceModel
	operation := constants.TfOperationRead
	if ctxOperation, ok := ctx.Value("retry_operation").(string); ok {
		operation = ctxOperation
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
	var remote *promptPolicyResponse
	getConfig := customrequests.GetRequestConfig{Endpoint: promptPoliciesEndpoint + "/" + object.ID.ValueString(), RawResponse: &remote}
	_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && remote == nil {
		err = errEmptyResponse
	}
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions, errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid prompt resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...)
	}
}

func (r *NetworkPromptPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkPromptPolicyResourceModel
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
	// An empty body means only Terraform-side values such as timeouts changed, so no PATCH is sent.
	if body := constructResource(ctx, &plan, &state); len(body) > 0 {
		// The response is decoded only so a malformed body still fails the update; the readback below sets state.
		var patched *promptPolicyResponse
		patchConfig := customrequests.PatchRequestConfig{Endpoint: promptPoliciesEndpoint + "/" + state.ID.ValueString(), RawBody: body, RawResponse: &patched}
		if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}
	// ReadWithRetry starts from the previous known state and only replaces it after a
	// successful readback that matches the planned values.
	resp.State = req.State
	plan.ID = state.ID
	readState := req.State
	resp.Diagnostics.Append(readState.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: readState, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName, crud.ReadWithConsistencyPredicate(crud.AttributesMatchPredicate(readState, readbackMatch)))
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.UpdateResponseContainer{UpdateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after update", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

func (r *NetworkPromptPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkPromptPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	deleteConfig := customrequests.DeleteRequestConfig{Endpoint: promptPoliciesEndpoint + "/" + object.ID.ValueString()}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	resp.State.RemoveResource(ctx)
}
