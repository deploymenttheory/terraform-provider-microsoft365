package graphBetaNetworkForwardingProfile

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

func (r *NetworkForwardingProfileResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkForwardingProfileResourceModel
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

	body, diags := requestBody(ctx, plan, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// POST can return 504 after creating the profile, and a retried POST then fails
	// as a duplicate. List first so only a profile created by this request is adopted.
	existing, listErr := r.customProfileIDs(ctx)
	if listErr != nil {
		tflog.Warn(ctx, "Could not list forwarding profiles before create", map[string]any{
			"error": listErr.Error(),
		})
	}
	remote, err := r.create(ctx, body)
	// Without the list a match may be pre-existing, so report the error.
	if err != nil && listErr == nil {
		info := errors.GraphError(ctx, err)
		if created := r.waitForCreated(ctx, plan, existing, createdLookupAttemptsFor(info)); created != nil {
			tflog.Warn(
				ctx,
				"Forwarding profile create returned an error but the profile exists; adopting it",
				map[string]any{
					"status": info.StatusCode,
					"id":     *created.ID,
				},
			)
			remote, err = created, nil
		}
	}
	if err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}
	claimProfile(*remote.ID)
	// Persist the ID before reading back so a failed read still allows cleanup.
	plan.ID = types.StringPointerValue(remote.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Assignments apply asynchronously; wait until GET reflects the plan.
	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		readbackOptions(plan, constants.TfOperationCreate),
	); err != nil {
		resp.Diagnostics.AddError("Error reading forwarding profile after create", err.Error())
	}
}

func (r *NetworkForwardingProfileResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkForwardingProfileResourceModel
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

	remote, err := r.get(ctx, state.ID.ValueString())
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

func (r *NetworkForwardingProfileResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkForwardingProfileResourceModel
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

	body, diags := requestBody(ctx, plan, &state)
	resp.Diagnostics.Append(diags...)
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
		resp.Diagnostics.AddError("Error reading forwarding profile after update", err.Error())
	}
}

func (r *NetworkForwardingProfileResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkForwardingProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		state.Timeouts.Delete,
		DeleteTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	if err := r.delete(ctx, state.ID.ValueString()); err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationDelete,
			r.WritePermissions,
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func readbackOptions(
	plan NetworkForwardingProfileResourceModel,
	operation string,
) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = sharedmodels.RequireStableReads(func(ctx context.Context, s tfsdk.State) bool {
		var actual NetworkForwardingProfileResourceModel
		if s.Get(ctx, &actual).HasError() {
			return false
		}
		return matches(plan.Name, actual.Name) &&
			matches(plan.Description, actual.Description) &&
			matches(plan.Priority, actual.Priority) &&
			matches(plan.State, actual.State) &&
			matches(plan.ClientFallbackAction, actual.ClientFallbackAction) &&
			matches(plan.ForwardingSourceMatching, actual.ForwardingSourceMatching) &&
			matches(plan.LocalNetworkSettings, actual.LocalNetworkSettings) &&
			matches(plan.AssignmentRules, actual.AssignmentRules)
	})
	return opts
}

// matches reports whether a configured value has reached Graph. Unknown planned
// values are computed by Graph and are not compared.
func matches(planned, actual attr.Value) bool {
	return planned.IsUnknown() || planned.Equal(actual)
}
