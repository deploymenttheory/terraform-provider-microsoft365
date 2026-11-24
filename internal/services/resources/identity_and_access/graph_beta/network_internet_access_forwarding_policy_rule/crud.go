package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

// Create handles the Create operation for network internet access forwarding policy rule resources.
//
// Operation: Creates a new network internet access forwarding policy rule
// API Calls:
//   - POST /networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-manage-internet-access-profile
func (r *NetworkInternetAccessForwardingPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkInternetAccessForwardingPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object, false)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing internet access forwarding policy rule", err.Error())
		return
	}

	// Graph can return 412 when Terraform writes several rules of the same policy in parallel; retrying lets the
	// policy settle.
	var created networkaccess.PolicyRuleable
	post := func(ctx context.Context) error {
		var err error
		created, err = r.client.NetworkAccess().ForwardingPolicies().ByForwardingPolicyId(object.ForwardingPolicyID.ValueString()).PolicyRules().Post(ctx, requestBody, nil)
		return err
	}
	retryOptions := crud.PreconditionRetryOptions{Operation: constants.TfOperationCreate, ResourceTypeName: ResourceName, LogFields: map[string]any{"policy_id": object.ForwardingPolicyID.ValueString()}}
	if err := crud.WriteWithPreconditionRetry(ctx, post, retryOptions); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created == nil || created.GetId() == nil {
		resp.Diagnostics.AddError("Error creating internet access forwarding policy rule", "The API returned an invalid response without an id.")
		return
	}

	object.ID = types.StringValue(*created.GetId())
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading internet access forwarding policy rule state after create", err.Error())
	}
}

// Read handles the Read operation for network internet access forwarding policy rule resources.
//
// Operation: Retrieves the network internet access forwarding policy rule by ID
// API Calls:
//   - GET /networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-manage-internet-access-profile
func (r *NetworkInternetAccessForwardingPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkInternetAccessForwardingPolicyRuleResourceModel
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

	remoteResource, err := r.client.NetworkAccess().ForwardingPolicies().ByForwardingPolicyId(object.ForwardingPolicyID.ValueString()).PolicyRules().ByPolicyRuleId(object.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	if remoteResource == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	// The SDK picks the model from @odata.type; any other rule type cannot be mapped to this resource.
	rule, ok := remoteResource.(networkaccess.InternetAccessForwardingRuleable)
	if !ok {
		resp.Diagnostics.AddError("Error reading internet access forwarding policy rule", fmt.Sprintf("The API returned %T, expected an internet access forwarding rule.", remoteResource))
		return
	}

	MapRemoteStateToTerraform(ctx, &object, rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

// Update handles the Update operation for network internet access forwarding policy rule resources.
//
// Operation: Updates the network internet access forwarding policy rule
// API Calls:
//   - PATCH /networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-manage-internet-access-profile
func (r *NetworkInternetAccessForwardingPolicyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkInternetAccessForwardingPolicyRuleResourceModel
	var state NetworkInternetAccessForwardingPolicyRuleResourceModel
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

	plan.ID = state.ID
	requestBody, err := constructResource(ctx, &plan, true)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing internet access forwarding policy rule update", err.Error())
		return
	}

	patch := func(ctx context.Context) error {
		_, err := r.client.NetworkAccess().ForwardingPolicies().ByForwardingPolicyId(state.ForwardingPolicyID.ValueString()).PolicyRules().ByPolicyRuleId(state.ID.ValueString()).Patch(ctx, requestBody, nil)
		return err
	}
	retryOptions := crud.PreconditionRetryOptions{Operation: constants.TfOperationUpdate, ResourceTypeName: ResourceName, LogFields: map[string]any{"policy_id": state.ForwardingPolicyID.ValueString(), "rule_id": state.ID.ValueString()}}
	if err := crud.WriteWithPreconditionRetry(ctx, patch, retryOptions); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ForwardingPolicyID = state.ForwardingPolicyID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading internet access forwarding policy rule state after update", err.Error())
	}
}

// Delete handles the Delete operation for network internet access forwarding policy rule resources.
//
// Operation: Deletes the network internet access forwarding policy rule
// API Calls:
//   - DELETE /networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-manage-internet-access-profile
func (r *NetworkInternetAccessForwardingPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkInternetAccessForwardingPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	remove := func(ctx context.Context) error {
		return r.client.NetworkAccess().ForwardingPolicies().ByForwardingPolicyId(object.ForwardingPolicyID.ValueString()).PolicyRules().ByPolicyRuleId(object.ID.ValueString()).Delete(ctx, nil)
	}
	retryOptions := crud.PreconditionRetryOptions{Operation: constants.TfOperationDelete, ResourceTypeName: ResourceName, LogFields: map[string]any{"policy_id": object.ForwardingPolicyID.ValueString(), "rule_id": object.ID.ValueString()}}
	if err := crud.WriteWithPreconditionRetry(ctx, remove, retryOptions); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Deleted %s with id %s", ResourceName, object.ID.ValueString()))
	resp.State.RemoveResource(ctx)
}
