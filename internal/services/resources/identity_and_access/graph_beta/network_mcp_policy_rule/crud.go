// Package graphBetaNetworkMCPPolicyRule manages the rules of Global Secure Access MCP policies.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no mcpPolicyRule model or MCP condition types in
//     models/networkaccess, and no mcpPolicies/{id}/policyRules request builder in networkaccess.
//   - Learn: none. There is no networkaccess-mcppolicyrule resource page.
//   - Observed: the request and response shape for /networkaccess/mcpPolicies/{id}/policyRules were taken from the
//     Entra admin centre Global Secure Access MCP policy blade. The package's earlier comments record that a live
//     POST/PATCH followed by GET preserved the order and duplicates of every condition values array. The SDK was checked
//     by reading its source on 2026-11-24 (no live call: the test tenant has no Global Secure Access licence). Every
//     call therefore goes through the custom_requests Config API in raw JSON mode, with APIVersion left empty so it
//     uses the beta client's own base URL.
//
// Removal condition: the beta SDK ships the networkaccess mcpPolicyRule model with its matching condition types and a
// /networkAccess/mcpPolicies/{id}/policyRules request builder.
package graphBetaNetworkMCPPolicyRule

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
)

// readbackMismatchDetail explains a successful write whose readback differs from the configuration.
const readbackMismatchDetail = "The API returned success but did not persist all requested attributes. Check the resource before retrying."

func (r *NetworkMCPPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkMCPPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructResource(ctx, &object, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing MCP resource", err.Error())
		return
	}
	var created *mcpPolicyRuleResponse
	postConfig := customrequests.PostRequestConfig{Endpoint: fmt.Sprintf(mcpPolicyRulesEndpoint, object.MCPPolicyID.ValueString()), RawBody: body, RawResponse: &created}
	_, err = customrequests.PostRequest(ctx, r.client.GetAdapter(), postConfig, nil, nil)
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
	object.Status = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, MCPPolicyRuleIdentity{ID: object.ID.ValueString(), MCPPolicyID: object.MCPPolicyID.ValueString()})...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	desired := object
	var remote *mcpPolicyRuleResponse
	getConfig := customrequests.GetRequestConfig{Endpoint: fmt.Sprintf(mcpPolicyRuleEndpoint, object.MCPPolicyID.ValueString(), object.ID.ValueString()), RawResponse: &remote}
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
	// The diff body always holds @odata.type, so more than one property means a requested value was not persisted.
	if diff, err := constructResource(ctx, &desired, &object); err != nil || len(diff) > 1 {
		resp.Diagnostics.AddError("MCP readback does not match configuration", readbackMismatchDetail)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkMCPPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkMCPPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	var remote *mcpPolicyRuleResponse
	getConfig := customrequests.GetRequestConfig{Endpoint: fmt.Sprintf(mcpPolicyRuleEndpoint, object.MCPPolicyID.ValueString(), object.ID.ValueString()), RawResponse: &remote}
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
		resp.Diagnostics.Append(resp.Identity.Set(ctx, MCPPolicyRuleIdentity{ID: object.ID.ValueString(), MCPPolicyID: object.MCPPolicyID.ValueString()})...)
	}
}

func (r *NetworkMCPPolicyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state NetworkMCPPolicyRuleResourceModel
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
	body, err := constructResource(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing MCP update", err.Error())
		return
	}
	// The body always holds @odata.type. Nothing else means only Terraform-side values such as timeouts changed.
	if len(body) > 1 {
		var patched *mcpPolicyRuleResponse
		endpoint := fmt.Sprintf(mcpPolicyRuleEndpoint, state.MCPPolicyID.ValueString(), state.ID.ValueString())
		patchConfig := customrequests.PatchRequestConfig{Endpoint: endpoint, RawBody: body, RawResponse: &patched}
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
	var remote *mcpPolicyRuleResponse
	getConfig := customrequests.GetRequestConfig{Endpoint: fmt.Sprintf(mcpPolicyRuleEndpoint, object.MCPPolicyID.ValueString(), object.ID.ValueString()), RawResponse: &remote}
	_, err = customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
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
	if diff, err := constructResource(ctx, &plan, &object); err != nil || len(diff) > 1 {
		resp.Diagnostics.AddError("MCP readback does not match configuration", readbackMismatchDetail)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, MCPPolicyRuleIdentity{ID: object.ID.ValueString(), MCPPolicyID: object.MCPPolicyID.ValueString()})...)
	}
}

func (r *NetworkMCPPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkMCPPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	deleteConfig := customrequests.DeleteRequestConfig{Endpoint: fmt.Sprintf(mcpPolicyRuleEndpoint, object.MCPPolicyID.ValueString(), object.ID.ValueString())}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	resp.State.RemoveResource(ctx)
}
