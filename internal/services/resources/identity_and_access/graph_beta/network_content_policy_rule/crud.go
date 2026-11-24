// Package graphBetaNetworkContentPolicyRule manages the rules of Global Secure Access content (file) policies.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no fileRule, filePolicyFqdnDestination,
//     filePolicyUrlDestination or filePolicyWebCategoryDestination model in models/networkaccess, and no
//     filePolicies/{id}/policyRules request builder in networkaccess.
//   - Learn: none. There is no networkaccess-filepolicy or networkaccess-filerule resource page.
//   - Observed: the request and response shape for /networkaccess/filePolicies/{id}/policyRules were taken from the
//     Entra admin centre Global Secure Access content policy blade (checked by reading the SDK source on 2026-11-24;
//     no live call, because the test tenant has no Global Secure Access licence). Every call therefore goes through
//     the custom_requests Config API in raw JSON mode, with APIVersion left empty so it uses the beta client's own
//     base URL.
//
// Removal condition: the beta SDK ships the networkaccess fileRule model with its file policy destination types and a
// /networkAccess/filePolicies/{id}/policyRules request builder.
package graphBetaNetworkContentPolicyRule

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

func (r *NetworkContentPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkContentPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing resource for Create Method", err.Error())
		return
	}
	var created contentPolicyRuleResponse
	postConfig := customrequests.PostRequestConfig{Endpoint: fmt.Sprintf(contentPolicyRulesEndpoint, object.ContentPolicyID.ValueString()), RawBody: body, RawResponse: &created}
	if _, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), postConfig, nil, nil); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.ID == nil {
		resp.Diagnostics.AddError("Error creating content policy rule", "The API returned an invalid response without an id.")
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

func (r *NetworkContentPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkContentPolicyRuleResourceModel
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
	var remote contentPolicyRuleResponse
	endpoint := fmt.Sprintf(contentPolicyRuleEndpoint, object.ContentPolicyID.ValueString(), object.ID.ValueString())
	getConfig := customrequests.GetRequestConfig{Endpoint: endpoint, RawResponse: &remote}
	if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	MapRemoteStateToTerraform(ctx, &object, &remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkContentPolicyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkContentPolicyRuleResourceModel
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
	body, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing resource for Update Method", err.Error())
		return
	}
	patchConfig := customrequests.PatchRequestConfig{Endpoint: fmt.Sprintf(contentPolicyRuleEndpoint, state.ContentPolicyID.ValueString(), state.ID.ValueString()), RawBody: body}
	logFields := map[string]any{"policy_id": state.ContentPolicyID.ValueString(), "rule_id": state.ID.ValueString()}
	retryOpts := crud.PreconditionRetryOptions{Operation: constants.TfOperationUpdate, ResourceTypeName: ResourceName, LogFields: logFields}
	patchRule := func(ctx context.Context) error {
		return customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig)
	}
	if err := crud.WriteWithPreconditionRetry(ctx, patchRule, retryOpts); err != nil {
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

func (r *NetworkContentPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkContentPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	deleteConfig := customrequests.DeleteRequestConfig{Endpoint: fmt.Sprintf(contentPolicyRuleEndpoint, object.ContentPolicyID.ValueString(), object.ID.ValueString())}
	logFields := map[string]any{"policy_id": object.ContentPolicyID.ValueString(), "rule_id": object.ID.ValueString()}
	retryOpts := crud.PreconditionRetryOptions{Operation: constants.TfOperationDelete, ResourceTypeName: ResourceName, LogFields: logFields}
	deleteRule := func(ctx context.Context) error {
		return customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig)
	}
	if err := crud.WriteWithPreconditionRetry(ctx, deleteRule, retryOpts); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}
