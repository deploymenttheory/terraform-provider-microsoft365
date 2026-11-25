// Package graphBetaNetworkPromptPolicyRule manages the rules of Global Secure Access prompt policies.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no promptRule model or conversation scheme types
//     in models/networkaccess, and no promptPolicies/{id}/policyRules request builder in networkaccess.
//   - Learn: none. There is no networkaccess-promptrule resource page.
//   - Observed: the request and response shape for /networkaccess/promptPolicies/{id}/policyRules were taken from the
//     Entra admin centre Global Secure Access prompt policy blade. The response carries portal fields that Graph
//     metadata omits, notably action. The SDK was checked by reading its source on 2026-11-24 (no live call: the test
//     tenant has no Global Secure Access licence). Every call therefore goes through the custom_requests Config API in
//     raw JSON mode, with APIVersion left empty so it uses the beta client's own base URL.
//
// Removal condition: the beta SDK ships the networkaccess promptRule model with its conversation scheme types and a
// /networkAccess/promptPolicies/{id}/policyRules request builder.
package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// readbackMatch is what the read-back after a write waits for: a populated id and the written fields. Graph returns an unset custom scheme JSON path as an empty string.
var readbackMatch = crud.AttributeMatchOptions{
	Paths:        crud.RootPaths("name", "description", "action", "priority", "enabled", "prompt_logging", "scan_result", "conversation_schemes"),
	RequireKnown: crud.RootPaths("id"),
	Overrides:    []crud.AttributeMatchOverride{{Path: path.MatchRoot("conversation_schemes").AtAnyListIndex().AtName("json_path"), NullMatchesEmpty: new(true)}},
}

func (r *NetworkPromptPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkPromptPolicyRuleResourceModel
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
		resp.Diagnostics.AddError("Error constructing prompt resource", err.Error())
		return
	}
	var created *promptPolicyRuleResponse
	postConfig := customrequests.PostRequestConfig{Endpoint: fmt.Sprintf(promptPolicyRulesEndpoint, object.PromptPolicyID.ValueString()), RawBody: body, RawResponse: &created}
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
	// Persist the response body identity before readback. The API Location header may contain a placeholder.
	object.ID = types.StringValue(*created.ID)
	object.Status = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, PromptPolicyRuleIdentity{ID: object.ID.ValueString(), PromptPolicyID: object.PromptPolicyID.ValueString()})...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithConsistencyPredicate(promptPolicyRuleConsistencyPredicate(&object)))
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after create", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

func (r *NetworkPromptPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkPromptPolicyRuleResourceModel
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
	var remote *promptPolicyRuleResponse
	endpoint := fmt.Sprintf(promptPolicyRuleEndpoint, object.PromptPolicyID.ValueString(), object.ID.ValueString())
	getConfig := customrequests.GetRequestConfig{Endpoint: endpoint, RawResponse: &remote}
	_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && remote == nil {
		err = errEmptyResponse
	}
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions, errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		if remote.Settings != nil && remote.Settings.Status != nil {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("status"), *remote.Settings.Status)...)
		}
		resp.Diagnostics.AddError("Invalid prompt resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, PromptPolicyRuleIdentity{ID: object.ID.ValueString(), PromptPolicyID: object.PromptPolicyID.ValueString()})...)
	}
}

func (r *NetworkPromptPolicyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkPromptPolicyRuleResourceModel
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
		resp.Diagnostics.AddError("Error constructing prompt update", err.Error())
		return
	}
	// The body always holds @odata.type. Nothing else means only Terraform-side values such as timeouts changed.
	if len(body) > 1 {
		// The response is decoded only so a malformed body still fails the update; the readback below sets state.
		var patched *promptPolicyRuleResponse
		endpoint := fmt.Sprintf(promptPolicyRuleEndpoint, state.PromptPolicyID.ValueString(), state.ID.ValueString())
		patchConfig := customrequests.PatchRequestConfig{Endpoint: endpoint, RawBody: body, RawResponse: &patched}
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
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName, crud.ReadWithConsistencyPredicate(promptPolicyRuleConsistencyPredicate(&plan)))
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.UpdateResponseContainer{UpdateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after update", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

func (r *NetworkPromptPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkPromptPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	deleteConfig := customrequests.DeleteRequestConfig{Endpoint: fmt.Sprintf(promptPolicyRuleEndpoint, object.PromptPolicyID.ValueString(), object.ID.ValueString())}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	resp.State.RemoveResource(ctx)
}
