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

// Create handles the Create operation for network default forwarding profile resources.
//
// Operation: Finds and adopts the built-in profile, then applies configured differences
// API Calls:
//   - GET /networkaccess/forwardingProfiles (finds the built-in profile)
//   - PATCH /networkaccess/forwardingProfiles/{id} (only when settings differ)
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
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

// Read handles the Read operation for network default forwarding profile resources.
//
// Operation: Retrieves the built-in forwarding profile by ID or imported traffic type
// API Calls:
//   - GET /networkaccess/forwardingProfiles/{id}
//   - GET /networkaccess/forwardingProfiles (when imported by traffic type)
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
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

// Update handles the Update operation for network default forwarding profile resources.
//
// Operation: Updates the network default forwarding profile
// API Calls:
//   - PATCH /networkaccess/forwardingProfiles/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
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

// Delete handles the Delete operation for network default forwarding profile resources.
//
// Operation: Removes Terraform state while retaining the built-in forwarding profile
// API Calls:
//   - None (state only)
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
func (r *NetworkDefaultForwardingProfileResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	resp.State.RemoveResource(ctx)
}

// readbackOptions configures shared read retries to verify the preceding write.
func readbackOptions(
	plan NetworkDefaultForwardingProfileResourceModel,
	operation string,
) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = crud.RequireStableReads(func(ctx context.Context, s tfsdk.State) bool {
		var actual NetworkDefaultForwardingProfileResourceModel
		if s.Get(ctx, &actual).HasError() {
			return false
		}
		return matches(plan.State, actual.State) &&
			matches(plan.ClientFallbackAction, actual.ClientFallbackAction) &&
			matches(plan.ForwardingSourceMatching, actual.ForwardingSourceMatching) &&
			matches(plan.LocalNetworkSettings, actual.LocalNetworkSettings) &&
			matches(plan.AssignmentRules, actual.AssignmentRules)
	})
	return opts
}

// matches checks whether remote state matches the configured settings.
func matches(planned, actual attr.Value) bool {
	return planned.IsUnknown() || planned.Equal(actual)
}
