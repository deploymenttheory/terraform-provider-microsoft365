// Package graphBetaNetworkWebFilteringPolicyRule manages the rules of Global Secure Access web filtering policies.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no webFilteringPolicies request builder in
//     networkaccess and no webFilteringRule, webFilteringActionAllow, webFilteringActionBlock,
//     webFilteringUrlDestination, webFilteringWebCategoryDestination or headerModificationAdd model in
//     models/networkaccess.
//   - Learn: none. Learn documents only the generic filteringRule hierarchy
//     (https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringrule?view=graph-rest-beta), not
//     /networkaccess/webFilteringPolicies/{id}/policyRules or the webFiltering* discriminators.
//   - Observed: the endpoint and payload follow the Entra admin centre's Global Secure Access web content filtering
//     requests. Graph rejects headerSettings on block actions, even as null, and keeps existing headers on an allow
//     action PATCH that omits it. Checked against the SDK source on 2026-11-24; no live call, because the test tenant
//     has no Global Secure Access licence. Requests go through custom_requests in raw JSON mode with APIVersion left
//     empty, so they use the beta client's own base URL as the previous transport did.
//
// Removal condition: the beta SDK ships a /networkAccess/webFilteringPolicies/{id}/policyRules request builder and
// the webFiltering* rule, action, destination and header modification models.
package graphBetaNetworkWebFilteringPolicyRule

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// webFilteringPoliciesEndpoint is the portal-observed policy collection; note the lower-case "networkaccess" segment.
const webFilteringPoliciesEndpoint = "/networkaccess/webFilteringPolicies"

// Create handles the Create operation for network web filtering policy rule resources.
//
// Operation: Creates a new rule in a web filtering policy
// API Calls:
//   - POST /networkaccess/webFilteringPolicies/{webFilteringPolicyId}/policyRules
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-configure-web-content-filtering
func (r *NetworkWebFilteringPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkWebFilteringPolicyRuleResourceModel

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
		resp.Diagnostics.AddError("Error constructing resource for Create Method", fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()))
		return
	}

	var created webFilteringRule
	rulesEndpoint := webFilteringPoliciesEndpoint + "/" + object.WebFilteringPolicyID.ValueString() + "/policyRules"
	postConfig := customrequests.PostRequestConfig{
		Endpoint:    rulesEndpoint,
		RawBody:     requestBody,
		RawResponse: &created,
	}
	if _, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), postConfig, nil, nil); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.ID == nil {
		resp.Diagnostics.AddError("Error creating web filtering policy rule", "The API returned an invalid response without an id.")
		return
	}

	object.ID = types.StringValue(*created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after create", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

// Read handles the Read operation for network web filtering policy rule resources.
//
// Operation: Retrieves the web filtering policy rule by ID
// API Calls:
//   - GET /networkaccess/webFilteringPolicies/{webFilteringPolicyId}/policyRules/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-configure-web-content-filtering
func (r *NetworkWebFilteringPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkWebFilteringPolicyRuleResourceModel
	var identity sharedmodels.ResourceIdentity

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

	identity.ID = object.ID.ValueString()
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var rule webFilteringRule
	rulesEndpoint := webFilteringPoliciesEndpoint + "/" + object.WebFilteringPolicyID.ValueString() + "/policyRules"
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          rulesEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        object.ID.ValueString(),
		RawResponse:       &rule,
	}
	if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, &rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

// Update handles the Update operation for network web filtering policy rule resources.
//
// Operation: Updates the web filtering policy rule, retrying while Graph reports a policy precondition failure
// API Calls:
//   - PATCH /networkaccess/webFilteringPolicies/{webFilteringPolicyId}/policyRules/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-configure-web-content-filtering
func (r *NetworkWebFilteringPolicyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkWebFilteringPolicyRuleResourceModel
	var state NetworkWebFilteringPolicyRuleResourceModel

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
		resp.Diagnostics.AddError("Error constructing resource for Update Method", fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()))
		return
	}

	// Graph can return 412 when Terraform writes several rules of the same policy in parallel; retrying lets the
	// policy settle.
	rulesEndpoint := webFilteringPoliciesEndpoint + "/" + state.WebFilteringPolicyID.ValueString() + "/policyRules"
	patchConfig := customrequests.PatchRequestConfig{
		Endpoint:          rulesEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID.ValueString(),
		RawBody:           requestBody,
	}
	retryOptions := crud.PreconditionRetryOptions{
		Operation:        constants.TfOperationUpdate,
		ResourceTypeName: ResourceName,
		LogFields:        map[string]any{"policy_id": state.WebFilteringPolicyID.ValueString(), "rule_id": state.ID.ValueString()},
	}
	patch := func(ctx context.Context) error {
		return customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig)
	}
	if err := crud.WriteWithPreconditionRetry(ctx, patch, retryOptions); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after update", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

// Delete handles the Delete operation for network web filtering policy rule resources.
//
// Operation: Deletes the web filtering policy rule, retrying while Graph reports a policy precondition failure
// API Calls:
//   - DELETE /networkaccess/webFilteringPolicies/{webFilteringPolicyId}/policyRules/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-configure-web-content-filtering
func (r *NetworkWebFilteringPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkWebFilteringPolicyRuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	rulesEndpoint := webFilteringPoliciesEndpoint + "/" + object.WebFilteringPolicyID.ValueString() + "/policyRules"
	deleteConfig := customrequests.DeleteRequestConfig{
		Endpoint:          rulesEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        object.ID.ValueString(),
	}
	retryOptions := crud.PreconditionRetryOptions{
		Operation:        constants.TfOperationDelete,
		ResourceTypeName: ResourceName,
		LogFields:        map[string]any{"policy_id": object.WebFilteringPolicyID.ValueString(), "rule_id": object.ID.ValueString()},
	}
	remove := func(ctx context.Context) error {
		return customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig)
	}
	if err := crud.WriteWithPreconditionRetry(ctx, remove, retryOptions); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)
}
