// Package graphBetaNetworkContentPolicy manages Global Secure Access content (file) policies.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no filePolicy model in
//     models/networkaccess and no filePolicies request builder in networkaccess.
//   - Learn: none. There is no networkaccess-filepolicy resource page; the Global Secure Access
//     API overview does not list file policies.
//   - Observed: the request and response shape for /networkaccess/filePolicies were taken from
//     the Entra admin centre Global Secure Access content policy blade. Every call therefore goes
//     through the custom_requests Config API in raw JSON mode, with APIVersion left empty so it
//     uses the beta client's own base URL.
//
// Removal condition: the beta SDK ships a networkaccess filePolicy model and a
// /networkAccess/filePolicies request builder.
package graphBetaNetworkContentPolicy

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

func (r *NetworkContentPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkContentPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	var created contentPolicyResponse
	postConfig := customrequests.PostRequestConfig{Endpoint: contentPoliciesEndpoint, RawBody: constructResource(ctx, &object, nil), RawResponse: &created}
	if _, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), postConfig, nil, nil); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.ID == nil {
		resp.Diagnostics.AddError("Error creating content policy", "The API returned an invalid response without an id.")
		return
	}
	object.ID = types.StringValue(*created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after create", fmt.Sprintf("Could not read resource state: %s", err))
	}
}

func (r *NetworkContentPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkContentPolicyResourceModel
	var identity sharedmodels.ResourceIdentity
	operation := constants.TfOperationRead
	if op, ok := ctx.Value("retry_operation").(string); ok {
		operation = op
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
	identity.ID = object.ID.ValueString()
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var remote contentPolicyResponse
	getConfig := customrequests.GetRequestConfig{Endpoint: contentPoliciesEndpoint + "/" + object.ID.ValueString(), RawResponse: &remote}
	if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	MapRemoteStateToTerraform(ctx, &object, &remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkContentPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkContentPolicyResourceModel
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
	// The PATCH answers with the updated policy. An empty response is reported as a failed update.
	var updated *contentPolicyResponse
	endpoint := contentPoliciesEndpoint + "/" + state.ID.ValueString()
	patchConfig := customrequests.PatchRequestConfig{Endpoint: endpoint, RawBody: constructResource(ctx, &plan, &state), RawResponse: &updated}
	err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig)
	if err == nil && updated == nil {
		err = errEmptyUpdateResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.UpdateResponseContainer{UpdateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after update", fmt.Sprintf("Could not read resource state: %s", err))
	}
}

func (r *NetworkContentPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkContentPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	deleteConfig := customrequests.DeleteRequestConfig{Endpoint: contentPoliciesEndpoint + "/" + object.ID.ValueString()}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}
