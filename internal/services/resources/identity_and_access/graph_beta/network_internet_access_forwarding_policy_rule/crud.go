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
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

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

	created, err := r.createRuleWithPreconditionRetry(ctx, object.ForwardingPolicyID.ValueString(), requestBody)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.id == nil {
		resp.Diagnostics.AddError("Error creating internet access forwarding policy rule", "The API returned an invalid response without an id.")
		return
	}

	object.ID = types.StringValue(*created.id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading internet access forwarding policy rule state after create", err.Error())
	}
}

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

	rule, err := r.getRule(ctx, object.ForwardingPolicyID.ValueString(), object.ID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

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

	if err := r.updateRuleWithPreconditionRetry(ctx, state.ForwardingPolicyID.ValueString(), state.ID.ValueString(), requestBody); err != nil {
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
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading internet access forwarding policy rule state after update", err.Error())
	}
}

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

	if err := r.deleteRuleWithPreconditionRetry(ctx, object.ForwardingPolicyID.ValueString(), object.ID.ValueString()); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Deleted %s with id %s", ResourceName, object.ID.ValueString()))
	resp.State.RemoveResource(ctx)
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) createRuleWithPreconditionRetry(ctx context.Context, policyID string, requestBody s.Parsable) (*internetAccessForwardingRuleResponse, error) {
	var created *internetAccessForwardingRuleResponse
	err := r.withPolicyRulePreconditionRetry(ctx, "create", policyID, "", func() error {
		result, err := r.createRule(ctx, policyID, requestBody)
		if err != nil {
			return err
		}
		created = result
		return nil
	})
	return created, err
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) updateRuleWithPreconditionRetry(ctx context.Context, policyID, ruleID string, requestBody s.Parsable) error {
	return r.withPolicyRulePreconditionRetry(ctx, "update", policyID, ruleID, func() error {
		return r.updateRule(ctx, policyID, ruleID, requestBody)
	})
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) deleteRuleWithPreconditionRetry(ctx context.Context, policyID, ruleID string) error {
	return r.withPolicyRulePreconditionRetry(ctx, "delete", policyID, ruleID, func() error {
		return r.deleteRule(ctx, policyID, ruleID)
	})
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) withPolicyRulePreconditionRetry(ctx context.Context, operation, policyID, ruleID string, fn func() error) error {
	const (
		maxPreconditionRetries = 3
		preconditionRetryDelay = 5 * time.Second
	)

	var lastErr error
	for attempt := 0; attempt <= maxPreconditionRetries; attempt++ {
		if err := fn(); err != nil {
			lastErr = err
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode != 412 && errorInfo.ErrorCode != "PreconditionFailed" {
				return err
			}
			if attempt < maxPreconditionRetries {
				tflog.Warn(ctx, "Retrying internet access forwarding policy rule operation after Graph precondition failure", map[string]any{
					"operation": operation,
					"policy_id": policyID,
					"rule_id":   ruleID,
					"attempt":   attempt + 1,
				})

				select {
				case <-time.After(preconditionRetryDelay):
					continue
				case <-ctx.Done():
					return fmt.Errorf("context cancelled while retrying internet access forwarding policy rule %s: %w", operation, ctx.Err())
				}
			}
		} else {
			return nil
		}
	}

	return lastErr
}
