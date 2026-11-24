// Package graphBetaNetworkMCPPolicy manages Global Secure Access MCP policies.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no mcpPolicy model in models/networkaccess and
//     no mcpPolicies request builder in networkaccess.
//   - Learn: none. There is no networkaccess-mcppolicy resource page; the Global Secure Access API overview does not
//     list MCP policies.
//   - Observed: the request and response shape for /networkaccess/mcpPolicies were taken from the Entra admin centre
//     Global Secure Access MCP policy blade (checked by reading the SDK source on 2026-11-24; no live call, because the
//     test tenant has no Global Secure Access licence). Every call therefore goes through the custom_requests Config
//     API in raw JSON mode, with APIVersion left empty so it uses the beta client's own base URL.
//
// Removal condition: the beta SDK ships a networkaccess mcpPolicy model and a /networkAccess/mcpPolicies request builder.
package graphBetaNetworkMCPPolicy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// readbackMismatchDetail explains a successful write whose readback differs from the configuration.
const readbackMismatchDetail = "The API returned success but did not persist all requested attributes. Check the resource before retrying."

func (r *NetworkMCPPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkMCPPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	var created *mcpPolicyResponse
	postConfig := customrequests.PostRequestConfig{Endpoint: mcpPoliciesEndpoint, RawBody: constructResource(ctx, &object, nil), RawResponse: &created}
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
	// Persist the response body identity before readback. The observed Location header is only the Graph root.
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
	desired := object
	var remote *mcpPolicyResponse
	getConfig := customrequests.GetRequestConfig{Endpoint: mcpPoliciesEndpoint + "/" + object.ID.ValueString(), RawResponse: &remote}
	_, err = customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && remote == nil {
		err = errEmptyResponse
	}
	// This is a create readback, not a refresh: even a 404 must retain the created id.
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid MCP resource response", err.Error())
		return
	}
	if len(constructResource(ctx, &desired, &object)) > 0 {
		resp.Diagnostics.AddError("MCP readback does not match configuration", readbackMismatchDetail)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkMCPPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkMCPPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	var remote *mcpPolicyResponse
	getConfig := customrequests.GetRequestConfig{Endpoint: mcpPoliciesEndpoint + "/" + object.ID.ValueString(), RawResponse: &remote}
	_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && remote == nil {
		err = errEmptyResponse
	}
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions, errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid MCP resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...)
	}
}

func (r *NetworkMCPPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state NetworkMCPPolicyResourceModel
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
		var patched *mcpPolicyResponse
		patchConfig := customrequests.PatchRequestConfig{Endpoint: mcpPoliciesEndpoint + "/" + state.ID.ValueString(), RawBody: body, RawResponse: &patched}
		err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig)
		if err == nil && patched == nil {
			err = errEmptyResponse
		}
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}
	// Keep the previous known state if readback fails after a successful PATCH.
	resp.State = req.State
	object := plan
	object.ID = state.ID
	var remote *mcpPolicyResponse
	getConfig := customrequests.GetRequestConfig{Endpoint: mcpPoliciesEndpoint + "/" + object.ID.ValueString(), RawResponse: &remote}
	_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && remote == nil {
		err = errEmptyResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid MCP resource response", err.Error())
		return
	}
	if len(constructResource(ctx, &plan, &object)) > 0 {
		resp.Diagnostics.AddError("MCP readback does not match configuration", readbackMismatchDetail)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...)
	}
}

func (r *NetworkMCPPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkMCPPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	deleteConfig := customrequests.DeleteRequestConfig{Endpoint: mcpPoliciesEndpoint + "/" + object.ID.ValueString()}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	resp.State.RemoveResource(ctx)
}
