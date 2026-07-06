package graphBetaNetworkWebFilteringPolicyRule

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
)

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

	created, err := r.createWebFilteringPolicyRule(ctx, object.WebFilteringPolicyID.ValueString(), requestBody)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.id == nil {
		resp.Diagnostics.AddError("Error creating web filtering policy rule", "The API returned an invalid response without an id.")
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
		resp.Diagnostics.AddError("Error reading resource state after create", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

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

	rule, err := r.getWebFilteringPolicyRule(ctx, object.WebFilteringPolicyID.ValueString(), object.ID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

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

	if err := r.updateWebFilteringPolicyRule(ctx, state.WebFilteringPolicyID.ValueString(), state.ID.ValueString(), requestBody); err != nil {
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
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after update", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

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

	if err := r.deleteWebFilteringPolicyRuleWithPreconditionRetry(ctx, object.WebFilteringPolicyID.ValueString(), object.ID.ValueString()); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *NetworkWebFilteringPolicyRuleResource) deleteWebFilteringPolicyRuleWithPreconditionRetry(ctx context.Context, policyID, ruleID string) error {
	const (
		maxPreconditionRetries = 3
		preconditionRetryDelay = 5 * time.Second
	)

	var lastErr error
	for attempt := 0; attempt <= maxPreconditionRetries; attempt++ {
		if err := r.deleteWebFilteringPolicyRule(ctx, policyID, ruleID); err != nil {
			lastErr = err
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode != 412 && errorInfo.ErrorCode != "PreconditionFailed" {
				return err
			}

			// Tenant verification showed Microsoft Graph can return 412 when Terraform
			// deletes multiple rules from the same web filtering policy in parallel.
			// Retrying lets the policy/rule collection state settle before the next delete.
			if attempt < maxPreconditionRetries {
				tflog.Warn(ctx, "Retrying web filtering policy rule delete after Graph precondition failure", map[string]any{
					"policy_id": policyID,
					"rule_id":   ruleID,
					"attempt":   attempt + 1,
				})

				select {
				case <-time.After(preconditionRetryDelay):
					continue
				case <-ctx.Done():
					return fmt.Errorf("context cancelled while retrying web filtering policy rule delete: %w", ctx.Err())
				}
			}
		} else {
			return nil
		}
	}

	return lastErr
}
