package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	profileType := plan.ProfileType.ValueString()
	profileID := plan.ForwardingProfileID.ValueString()
	unlock := lockProfile(profileType, profileID)
	defer unlock()
	profile, err := r.getProfile(ctx, profileType, profileID)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}

	existing := findLink(profile.Policies, "", plan.ForwardingPolicyID.ValueString())
	switch {
	case existing != nil && existing.ConfigType != "" &&
		!strings.EqualFold(existing.ConfigType, plan.ConfigType.ValueString()):
		// Replacing could leave the policy unlinked, so leave it in place.
		resp.Diagnostics.AddError(
			"Forwarding policy is already linked with a different config type",
			fmt.Sprintf("Policy %s is linked to profile %s as %q (link %s), not %q. Check config_type.",
				plan.ForwardingPolicyID.ValueString(), profileID, existing.ConfigType, existing.ID,
				plan.ConfigType.ValueString()),
		)
		return
	case existing != nil && linkMatchesPlan(existing, &plan):
		// New user profiles auto-link every private app policy; adopt it instead of failing.
		tflog.Warn(ctx, "Adopting existing forwarding policy link that already matches the configuration", map[string]any{
			"forwarding_profile_id": profileID,
			"policy_link_id":        existing.ID,
		})
	case existing != nil:
		tflog.Warn(ctx, "Replacing existing forwarding policy link whose state or priority differs", map[string]any{
			"forwarding_profile_id": profileID,
			"policy_link_id":        existing.ID,
		})
		// Graph rejects link updates and a remove+add in one delta, so this takes two writes.
		if err := r.removeLink(ctx, profileType, profileID, existing.ID); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
		if err := r.waitForLinkRemoved(ctx, profileType, profileID, existing.ID); err != nil {
			resp.Diagnostics.AddError("Error waiting for existing forwarding policy link removal", err.Error())
			return
		}
		if err := r.addLink(ctx, &plan); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			r.restoreLink(ctx, profileType, profileID, existing, resp)
			return
		}
	default:
		if err := r.addLink(ctx, &plan); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}

	// PATCH returns the profile, not the link; Read resolves the link by policy ID.
	planned := plan
	plan.ID = types.StringNull()
	plan.PolicyLinkID = types.StringNull()
	if plan.Priority.IsUnknown() {
		plan.Priority = types.Int64Null()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		readbackOptions(&planned, constants.TfOperationCreate),
	); err != nil {
		resp.Diagnostics.AddError("Error reading forwarding policy link after create", err.Error())
	}
}

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operation := constants.TfOperationRead
	if op, ok := ctx.Value("retry_operation").(string); ok && op != "" {
		operation = op
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	profile, err := r.getProfile(ctx, state.ProfileType.ValueString(), state.ForwardingProfileID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions,
			errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}

	link := findLink(profile.Policies, state.PolicyLinkID.ValueString(), state.ForwardingPolicyID.ValueString())
	if link == nil {
		tflog.Warn(ctx, "Forwarding policy link not found on profile, removing from state", map[string]any{
			"forwarding_profile_id": state.ForwardingProfileID.ValueString(),
			"policy_link_id":        state.PolicyLinkID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	mapLinkToState(&state, link)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update only stores plan-only changes such as timeouts; everything else forces replacement.
func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	profileType := state.ProfileType.ValueString()
	profileID := state.ForwardingProfileID.ValueString()
	unlock := lockProfile(profileType, profileID)
	defer unlock()
	profile, err := r.getProfile(ctx, profileType, profileID)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
		if !resp.Diagnostics.HasError() {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	link := findLink(profile.Policies, state.PolicyLinkID.ValueString(), state.ForwardingPolicyID.ValueString())
	if link == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := r.removeLink(ctx, profileType, profileID, link.ID); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		if !resp.Diagnostics.HasError() {
			resp.State.RemoveResource(ctx)
		}
		return
	}
	if err := r.waitForLinkRemoved(ctx, profileType, profileID, link.ID); err != nil {
		resp.Diagnostics.AddError("Error waiting for forwarding policy link removal", err.Error())
		return
	}
	resp.State.RemoveResource(ctx)
}

// restoreLink re-adds the link Create removed when the new add fails.
// Microsoft does not recreate removed auto-links.
func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) restoreLink(
	ctx context.Context,
	profileType, profileID string,
	previous *policyLink,
	resp *resource.CreateResponse,
) {
	restore := NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel{
		ProfileType:         types.StringValue(profileType),
		ForwardingProfileID: types.StringValue(profileID),
		ForwardingPolicyID:  types.StringValue(previous.Policy.ID),
		ConfigType:          types.StringValue(previous.ConfigType),
		State:               types.StringValue(strings.ToLower(previous.State)),
		Priority:            types.Int64PointerValue(previous.Priority),
	}
	if err := r.addLink(ctx, &restore); err != nil {
		resp.Diagnostics.AddError(
			"Previous forwarding policy link was removed and could not be restored",
			fmt.Sprintf("Create removed link %s from profile %s before its add failed, and re-adding it failed: %s",
				previous.ID, profileID, err.Error()),
		)
		return
	}
	resp.Diagnostics.AddWarning(
		"Previous forwarding policy link restored",
		fmt.Sprintf("The add failed, so the link that Create removed from profile %s was re-added with its "+
			"previous state and priority. Graph gives it a new link ID.", profileID),
	)
}

func linkMatchesPlan(link *policyLink, plan *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel) bool {
	if !strings.EqualFold(link.State, plan.State.ValueString()) {
		return false
	}
	if link.ConfigType != "" && !strings.EqualFold(link.ConfigType, plan.ConfigType.ValueString()) {
		return false
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		return link.Priority != nil && *link.Priority == plan.Priority.ValueInt64()
	}
	return true
}

// mapLinkToState keeps configured casing for IDs that Graph echoes back case-insensitively.
func mapLinkToState(state *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel, link *policyLink) {
	if !strings.EqualFold(state.PolicyLinkID.ValueString(), link.ID) {
		state.PolicyLinkID = types.StringValue(link.ID)
	}
	state.ID = compositeID(state.ForwardingProfileID.ValueString(), state.PolicyLinkID.ValueString())
	if link.Policy != nil && !strings.EqualFold(state.ForwardingPolicyID.ValueString(), link.Policy.ID) {
		state.ForwardingPolicyID = types.StringValue(link.Policy.ID)
	}
	if link.ConfigType != "" {
		state.ConfigType = types.StringValue(link.ConfigType)
	}
	if link.State != "" {
		state.State = types.StringValue(strings.ToLower(link.State))
	}
	if link.Priority != nil {
		state.Priority = types.Int64Value(*link.Priority)
	} else {
		state.Priority = types.Int64Null()
	}
}

func readbackOptions(
	plan *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel,
	operation string,
) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.RetryInterval = readbackInterval
	opts.MaxRetries = 60
	consecutive := 0
	opts.ConsistencyPredicate = func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
		if state.Raw.IsNull() || state.Get(ctx, &actual).HasError() || !matchesPlan(&actual, plan) {
			consecutive = 0
			return false
		}
		consecutive++
		return consecutive >= stableReads
	}
	return opts
}

func matchesPlan(actual, plan *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel) bool {
	if actual.PolicyLinkID.ValueString() == "" ||
		!strings.EqualFold(actual.ForwardingPolicyID.ValueString(), plan.ForwardingPolicyID.ValueString()) ||
		!actual.State.Equal(plan.State) ||
		!actual.ConfigType.Equal(plan.ConfigType) {
		return false
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		return actual.Priority.Equal(plan.Priority)
	}
	return true
}
