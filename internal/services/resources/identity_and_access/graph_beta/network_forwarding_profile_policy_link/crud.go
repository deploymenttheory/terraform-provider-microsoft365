package graphBetaNetworkForwardingProfilePolicyLink

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *NetworkForwardingProfilePolicyLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkForwardingProfilePolicyLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if err := r.patchPolicyLinkState(ctx, object.ForwardingProfileID.ValueString(), object.PolicyLinkID.ValueString(), constructResource(&object)); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = compositeID(object.ForwardingProfileID.ValueString(), object.PolicyLinkID.ValueString())
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
		resp.Diagnostics.AddError("Error reading forwarding policy link state after create", err.Error())
	}
}

func (r *NetworkForwardingProfilePolicyLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkForwardingProfilePolicyLinkResourceModel
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

	link, err := r.getPolicyLink(ctx, object.ForwardingProfileID.ValueString(), object.PolicyLinkID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	if link == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, object.ForwardingProfileID.ValueString(), link)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkForwardingProfilePolicyLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkForwardingProfilePolicyLinkResourceModel
	var state NetworkForwardingProfilePolicyLinkResourceModel
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

	if err := r.patchPolicyLinkState(ctx, state.ForwardingProfileID.ValueString(), state.PolicyLinkID.ValueString(), constructResource(&plan)); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ID = state.ID
	plan.ForwardingProfileID = state.ForwardingProfileID
	plan.PolicyLinkID = state.PolicyLinkID
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
		resp.Diagnostics.AddError("Error reading forwarding policy link state after update", err.Error())
	}
}

func (r *NetworkForwardingProfilePolicyLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkForwardingProfilePolicyLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Removing %s from Terraform state without deleting Microsoft-managed forwarding policy link", ResourceName), map[string]any{
		"forwarding_profile_id": object.ForwardingProfileID.ValueString(),
		"policy_link_id":        object.PolicyLinkID.ValueString(),
	})
	resp.State.RemoveResource(ctx)
}
