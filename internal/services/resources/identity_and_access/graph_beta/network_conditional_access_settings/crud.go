package graphBetaNetworkConditionalAccessSettings

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	kiotahttp "github.com/microsoft/kiota-http-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	builders "github.com/microsoftgraph/msgraph-beta-sdk-go/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// Create adopts the existing singleton and PATCHes only explicitly configured settings.
func (r *NetworkConditionalAccessSettingsResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan, config NetworkConditionalAccessSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
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
	if config.SignalingStatus.IsNull() {
		request, err := r.client.NetworkAccess().
			Settings().
			ConditionalAccess().
			ToGetRequestInformation(ctx, nil)
		if err != nil {

			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationCreate,
				r.ReadPermissions,
			)
			return

		}
		result, err := r.client.GetAdapter().
			SendPrimitive(ctx, request, "[]byte", abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
		remote, err := decodeSettingsResponse(result, err)
		if err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationCreate,
				r.ReadPermissions,
			)
			return
		}
		mapRemoteState(&plan, remote, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		body, err := constructResource(plan.SignalingStatus)
		if err != nil {
			resp.Diagnostics.AddError("Invalid signaling status", err.Error())
			return
		}

		if _, err = r.client.NetworkAccess().Settings().ConditionalAccess().Patch(ctx, body,
			&builders.SettingsConditionalAccessRequestBuilderPatchRequestConfiguration{
				// Avoid Kiota replaying an uncompressed body with a stale gzip header on retry.
				Options: []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)},
			}); err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationCreate,
				r.WritePermissions,
			)
			return
		}
		plan.ID = types.StringValue(singletonID)
	}
	// Save the fixed identity before readback so a successful PATCH is not lost on a GET failure.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), singletonID)...)
	}
	if config.SignalingStatus.IsNull() || resp.Diagnostics.HasError() {
		return
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = consistencyPredicate(plan.SignalingStatus)
	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Error reading settings after apply", err.Error())
	}
}

func (r *NetworkConditionalAccessSettingsResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkConditionalAccessSettingsResourceModel
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
	request, err := r.client.NetworkAccess().
		Settings().
		ConditionalAccess().
		ToGetRequestInformation(ctx, nil)
	if err != nil {

		// A missing singleton endpoint can mean licensing or availability, not deletion.
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
	result, err := r.client.GetAdapter().
		SendPrimitive(ctx, request, "[]byte", abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
	remote, err := decodeSettingsResponse(result, err)
	if err != nil {
		// A missing singleton endpoint can mean licensing or availability, not deletion.
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
	mapRemoteState(&state, remote, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), singletonID)...)
	}
}

func (r *NetworkConditionalAccessSettingsResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, config, state NetworkConditionalAccessSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
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
	if config.SignalingStatus.IsNull() || plan.SignalingStatus.Equal(state.SignalingStatus) {
		// Omitted settings must never be written back from an earlier refresh.
		plan.SignalingStatus = state.SignalingStatus
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}
	body, err := constructResource(plan.SignalingStatus)
	if err != nil {
		resp.Diagnostics.AddError("Invalid signaling status", err.Error())
		return
	}

	if _, err = r.client.NetworkAccess().Settings().ConditionalAccess().Patch(ctx, body,
		&builders.SettingsConditionalAccessRequestBuilderPatchRequestConfiguration{
			// Avoid Kiota replaying an uncompressed body with a stale gzip header on retry.
			Options: []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)},
		}); err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationUpdate,
			r.WritePermissions,
		)
		return
	}
	plan.ID = types.StringValue(singletonID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = consistencyPredicate(plan.SignalingStatus)
	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Error reading settings after update", err.Error())
	}
}

// Delete releases Terraform ownership without changing the tenant configuration.
func (r *NetworkConditionalAccessSettingsResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(
		ctx,
		"Removing conditional access settings from Terraform state; remote settings are preserved",
	)
	resp.State.RemoveResource(ctx)
}

var errMissingSignalingStatus = stderrors.New(
	"signaling_status must be a known enabled or disabled value before writing",
)

var errInvalidSettingsResponse = stderrors.New("invalid conditional access settings response")
