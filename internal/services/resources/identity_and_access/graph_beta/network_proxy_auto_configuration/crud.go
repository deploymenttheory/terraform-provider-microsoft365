package graphBetaNetworkProxyAutoConfiguration

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

func (r *NetworkProxyAutoConfigurationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkProxyAutoConfigurationResourceModel
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
	body, err := pacBody(plan, nil)
	if err != nil {
		err := err
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}
	var createResponse pacResponse
	createErr := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.POST,
		r.ResourcePath,
		nil,
		body,
		&createResponse,
	)
	remote, err := decodeCreateResponse(createResponse, createErr)
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
	// Retain the created ID if a subsequent read fails, so cleanup remains possible.
	plan.ID = types.StringValue(*remote.ID)
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
			"Error reading custom PAC after initial apply",
			err.Error(),
		)
	}
}

func (r *NetworkProxyAutoConfigurationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkProxyAutoConfigurationResourceModel
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
	var readResponse pacResponse
	readErr := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		r.ResourcePath+"/{pacId}",
		map[string]string{"pacId": state.ID.ValueString()},
		nil,
		&readResponse,
	)
	remote, err := &readResponse, readErr
	if err != nil {
		err = fmt.Errorf("get custom PAC: %w", err)

		errors.HandleKiotaGraphErrorWithOptions(
			ctx,
			err,
			resp,
			constants.TfOperationRead,
			r.ReadPermissions,
			errors.GraphErrorOptions{
				PreserveStateOnReadBadRequest: true,
				PreserveStateOnReadNotFound:   false,
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

func (r *NetworkProxyAutoConfigurationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkProxyAutoConfigurationResourceModel
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
	body, err := pacBody(plan, &state)
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
			r.ResourcePath+"/{pacId}",
			map[string]string{"pacId": plan.ID.ValueString()},
			body,
			nil,
		)
	}
	if err := updateErr; err != nil {
		err = fmt.Errorf("update custom PAC: %w", err)

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
		resp.Diagnostics.AddError("Error reading custom PAC after update", err.Error())
	}
}

func (r *NetworkProxyAutoConfigurationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkProxyAutoConfigurationResourceModel
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
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.DELETE,
		r.ResourcePath+"/{pacId}",
		map[string]string{"pacId": state.ID.ValueString()},
		nil,
		nil,
	); err != nil {
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
	plan NetworkProxyAutoConfigurationResourceModel,
	operation string,
) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkProxyAutoConfigurationResourceModel
		if state.Get(ctx, &actual).HasError() {
			return false
		}
		return actual.ID.Equal(plan.ID) && actual.Name.Equal(plan.Name) &&
			actual.Content.Equal(plan.Content) &&
			actual.IsEnabled.Equal(plan.IsEnabled)
	}
	return opts
}

var (
	errInvalidPACResponse = stderrors.New(
		"graph omitted required custom PAC fields; existing state is retained",
	)
	errInvalidPACConfig = stderrors.New(
		"name, content and is_enabled must be explicitly configured; name and content must be nonempty",
	)
)
