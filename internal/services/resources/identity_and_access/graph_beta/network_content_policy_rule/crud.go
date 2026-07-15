package graphBetaNetworkContentPolicyRule

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
	s "github.com/microsoft/kiota-abstractions-go/serialization"
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
	created, err := r.createContentPolicyRule(ctx, object.ContentPolicyID.ValueString(), body)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created == nil || created.id == nil {
		resp.Diagnostics.AddError("Error creating content policy rule", "The API returned an invalid response without an id.")
		return
	}
	object.ID = types.StringValue(*created.id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
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
	rule, err := r.getContentPolicyRule(ctx, object.ContentPolicyID.ValueString(), object.ID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	MapRemoteStateToTerraform(ctx, &object, rule)
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
	if err := r.updateContentPolicyRuleWithPreconditionRetry(ctx, state.ContentPolicyID.ValueString(), state.ID.ValueString(), body); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
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
	if err := r.deleteContentPolicyRuleWithPreconditionRetry(ctx, object.ContentPolicyID.ValueString(), object.ID.ValueString()); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *NetworkContentPolicyRuleResource) updateContentPolicyRuleWithPreconditionRetry(ctx context.Context, policyID, ruleID string, body s.Parsable) error {
	return r.withContentPolicyRulePreconditionRetry(ctx, "update", policyID, ruleID, func() error { return r.updateContentPolicyRule(ctx, policyID, ruleID, body) })
}

func (r *NetworkContentPolicyRuleResource) deleteContentPolicyRuleWithPreconditionRetry(ctx context.Context, policyID, ruleID string) error {
	return r.withContentPolicyRulePreconditionRetry(ctx, "delete", policyID, ruleID, func() error { return r.deleteContentPolicyRule(ctx, policyID, ruleID) })
}

func (r *NetworkContentPolicyRuleResource) withContentPolicyRulePreconditionRetry(ctx context.Context, operation, policyID, ruleID string, fn func() error) error {
	const maxRetries = 3
	const retryDelay = 5 * time.Second
	attempt := 0
	return crud.PollUntil(ctx, retryDelay, func(ctx context.Context) (bool, error) {
		err := fn()
		if err == nil {
			return true, nil
		}

		info := errors.GraphError(ctx, err)
		if info.StatusCode != 412 && info.ErrorCode != "PreconditionFailed" {
			return false, &crud.FatalPollError{Err: err}
		}
		if attempt >= maxRetries {
			return false, &crud.FatalPollError{Err: err}
		}

		attempt++
		tflog.Warn(ctx, "Retrying content policy rule operation after Graph precondition failure", map[string]any{
			"operation": operation,
			"policy_id": policyID,
			"rule_id":   ruleID,
			"attempt":   attempt,
		})
		return false, err
	})
}
