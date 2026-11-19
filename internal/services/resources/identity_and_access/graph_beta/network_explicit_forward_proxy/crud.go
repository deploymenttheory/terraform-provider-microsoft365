package graphBetaNetworkExplicitForwardProxy

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func (r *NetworkExplicitForwardProxyResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkExplicitForwardProxyResourceModel
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
	body, err := proxyPatch(plan, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	var updateErr error
	if len(body) > 0 {
		updateErr = customrequests.JSONRequest(
			ctx,
			r.client.GetAdapter(),
			abstractions.PATCH,
			r.ResourcePath,
			nil,
			body,
			nil,
		)
	}
	if err := updateErr; err != nil {
		err = fmt.Errorf("patch explicit proxy: %w", err)

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
	plan.ID = types.StringValue("explicitForwardProxyConfig")
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
			"Error reading explicit forward proxy after initial apply",
			err.Error(),
		)
	}
}

func (r *NetworkExplicitForwardProxyResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkExplicitForwardProxyResourceModel
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
	var readResponse proxyResponse
	readErr := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		r.ResourcePath,
		nil,
		nil,
		&readResponse,
	)
	remote, err := &readResponse, readErr
	if err != nil {
		err = fmt.Errorf("get explicit proxy: %w", err)

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
	if err := mapResponse(&state, remote); err != nil {
		resp.Diagnostics.AddError("Invalid Graph response", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkExplicitForwardProxyResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkExplicitForwardProxyResourceModel
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
	body, err := proxyPatch(plan, &state)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}
	var updateErr error
	if len(body) > 0 {
		updateErr = customrequests.JSONRequest(
			ctx,
			r.client.GetAdapter(),
			abstractions.PATCH,
			r.ResourcePath,
			nil,
			body,
			nil,
		)
	}
	if err := updateErr; err != nil {
		err = fmt.Errorf("patch explicit proxy: %w", err)

		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationUpdate,
			r.WritePermissions,
		)
		return
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
		resp.Diagnostics.AddError("Error reading explicit forward proxy after update", err.Error())
	}
}

// Delete releases Terraform management without sending a Graph request.
func (r *NetworkExplicitForwardProxyResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	resp.State.RemoveResource(ctx)
}

func readbackOptions(
	plan NetworkExplicitForwardProxyResourceModel,
	operation string,
) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkExplicitForwardProxyResourceModel
		if state.Get(ctx, &actual).HasError() {
			return false
		}
		return actual.ID.Equal(plan.ID) &&
			actual.InternetAccess != nil && plan.InternetAccess != nil && actual.InternetAccess.IsEnabled.Equal(plan.InternetAccess.IsEnabled) &&
			actual.InternetAccess.IsSourceIPSessionAffinityEnabled.Equal(
				plan.InternetAccess.IsSourceIPSessionAffinityEnabled,
			) &&
			actual.InternetAccess.SourceIPSessionAffinityOptions.Equal(
				plan.InternetAccess.SourceIPSessionAffinityOptions,
			)
	}
	return opts
}

var (
	errInvalidProxyResponse = stderrors.New(
		"graph omitted required explicit proxy settings or returned unsupported session affinity options; existing state is retained",
	)
	errUnknownProxySettings = stderrors.New(
		"internet_access attributes must be known and non-null, with affinity enabled exactly when options is not none",
	)
)
