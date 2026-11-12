package graphBetaNetworkDefaultForwardingProfile

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

func (r *NetworkDefaultForwardingProfileResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkDefaultForwardingProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	remote, err := r.findDefault(ctx, plan.TrafficForwardingType.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	var current NetworkDefaultForwardingProfileResourceModel
	current.TrafficForwardingType = plan.TrafficForwardingType
	resp.Diagnostics.Append(mapRemoteStateToTerraform(ctx, &current, remote)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Send only configured values that differ from the adopted profile.
	body := map[string]any{}
	prior := current.settings()
	resp.Diagnostics.Append(sharedmodels.SettingsPatch(ctx, body, plan.settings(), &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.patch(ctx, *remote.ID, body); err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}
	plan.ID = types.StringPointerValue(remote.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		readbackOptions(plan, constants.TfOperationCreate),
	); err != nil {
		resp.Diagnostics.AddError(
			"Error reading default forwarding profile after initial apply",
			err.Error(),
		)
	}
}

func (r *NetworkDefaultForwardingProfileResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkDefaultForwardingProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		state.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	var (
		remote *sharedmodels.ForwardingProfileWire
		err    error
	)
	if state.ID.IsNull() || state.ID.IsUnknown() {
		// Imported by traffic type.
		remote, err = r.findDefault(ctx, state.TrafficForwardingType.ValueString())
	} else {
		remote, err = r.get(ctx, state.ID.ValueString())
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions)
		return
	}
	resp.Diagnostics.Append(mapRemoteStateToTerraform(ctx, &state, remote)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkDefaultForwardingProfileResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkDefaultForwardingProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Update,
		UpdateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	body := map[string]any{}
	prior := state.settings()
	resp.Diagnostics.Append(sharedmodels.SettingsPatch(ctx, body, plan.settings(), &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.patch(ctx, state.ID.ValueString(), body); err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationUpdate,
			r.WritePermissions,
		)
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		readbackOptions(plan, constants.TfOperationUpdate),
	); err != nil {
		resp.Diagnostics.AddError(
			"Error reading default forwarding profile after update",
			err.Error(),
		)
	}
}

// Delete only removes state. Disabling or resetting a default profile on
// destroy would change traffic handling for the whole tenant.
func (r *NetworkDefaultForwardingProfileResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	resp.State.RemoveResource(ctx)
}

func readbackOptions(
	plan NetworkDefaultForwardingProfileResourceModel,
	operation string,
) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = func(ctx context.Context, s tfsdk.State) bool {
		var actual NetworkDefaultForwardingProfileResourceModel
		if s.Get(ctx, &actual).HasError() {
			return false
		}
		return matches(plan.State, actual.State) &&
			matches(plan.ClientFallbackAction, actual.ClientFallbackAction) &&
			matches(plan.ForwardingSourceMatching, actual.ForwardingSourceMatching) &&
			matches(plan.LocalNetworkSettings, actual.LocalNetworkSettings) &&
			matches(plan.AssignmentRules, actual.AssignmentRules)
	}
	return opts
}

func matches(planned, actual attr.Value) bool {
	return planned.IsUnknown() || planned.Equal(actual)
}
