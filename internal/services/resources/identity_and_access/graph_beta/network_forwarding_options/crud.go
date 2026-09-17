package graphBetaNetworkForwardingOptions

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

var (
	errMissingDNSState = stderrors.New("skip_dns_lookup_state must be explicitly configured")
	errInvalidDNSState = stderrors.New("skip_dns_lookup_state must be enabled or disabled")
)

func (r *NetworkForwardingOptionsResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkForwardingOptionsResourceModel
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
	if err := r.patch(ctx, plan); err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}
	// Record adoption before readback so an unavailable replica does not lose management state.
	plan.ID = types.StringValue(singletonID)
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
			"Error reading forwarding options after initial apply",
			err.Error(),
		)
	}
}

func (r *NetworkForwardingOptionsResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkForwardingOptionsResourceModel
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
	remote, err := r.client.NetworkAccess().Settings().ForwardingOptions().Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(
			ctx,
			err,
			resp,
			constants.TfOperationRead,
			r.ReadPermissions,
			errors.GraphErrorOptions{
				PreserveStateOnReadBadRequest: true,
				PreserveStateOnReadNotFound:   true,
			},
		)
		return
	}
	if remote == nil || remote.GetSkipDnsLookupState() == nil {
		resp.Diagnostics.AddError(
			"Invalid forwarding options response",
			"Graph returned no skipDnsLookupState; existing state is retained.",
		)
		return
	}
	value := remote.GetSkipDnsLookupState().String()
	if value != "enabled" && value != "disabled" {
		resp.Diagnostics.AddError(
			"Unsupported DNS forwarding state",
			fmt.Sprintf("Graph returned skipDnsLookupState %q; existing state is retained.", value),
		)
		return
	}
	state.ID = types.StringValue(singletonID)
	state.SkipDNSLookupState = types.StringValue(value)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkForwardingOptionsResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkForwardingOptionsResourceModel
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
	if !plan.SkipDNSLookupState.Equal(state.SkipDNSLookupState) {
		if err := r.patch(ctx, plan); err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationUpdate,
				r.WritePermissions,
			)
			return
		}
	}
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
		resp.Diagnostics.AddError("Error reading forwarding options after update", err.Error())
	}
}

// Delete releases Terraform management without sending a Graph request.
func (r *NetworkForwardingOptionsResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	resp.State.RemoveResource(ctx)
}

func (r *NetworkForwardingOptionsResource) patch(
	ctx context.Context,
	plan NetworkForwardingOptionsResourceModel,
) error {
	if plan.SkipDNSLookupState.IsNull() || plan.SkipDNSLookupState.IsUnknown() {
		return errMissingDNSState
	}
	value := plan.SkipDNSLookupState.ValueString()
	if value != "enabled" && value != "disabled" {
		return errInvalidDNSState
	}
	parsed, err := models.ParseStatus(value)
	if err != nil {
		return fmt.Errorf("parse DNS forwarding status: %w", err)
	}
	// A fresh model prevents Graph additional properties from being echoed back into shared settings.
	body := models.NewForwardingOptions()
	body.SetSkipDnsLookupState(parsed.(*models.Status))
	_, err = r.client.NetworkAccess().Settings().ForwardingOptions().Patch(ctx, body, nil)
	if err != nil {
		return fmt.Errorf("patch forwarding options: %w", err)
	}
	return nil
}

func readbackOptions(
	plan NetworkForwardingOptionsResourceModel,
	operation string,
) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkForwardingOptionsResourceModel
		if state.Get(ctx, &actual).HasError() {
			return false
		}
		return actual.ID.ValueString() == singletonID &&
			actual.SkipDNSLookupState.Equal(plan.SkipDNSLookupState)
	}
	return opts
}
