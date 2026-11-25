package graphBetaNetworkCustomBlockPage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// readbackMatch is what the read-back after a write waits for: the singleton id, the written state and, when set, the configuration.
var readbackMatch = crud.AttributeMatchOptions{
	Paths: crud.RootPaths("id", "state", "configuration"),
	Overrides: []crud.AttributeMatchOverride{
		{Path: path.MatchRoot("configuration"), IgnoreNullExpected: new(true)},
		{Path: path.MatchRoot("configuration").AtName("body"), IgnoreNullExpected: new(false)},
	},
}

// Create adopts and PATCHes the existing singleton; there is no create API.
func (r *NetworkCustomBlockPageResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data NetworkCustomBlockPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		data.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructResource(&data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid custom block page", err.Error())
		return
	}
	_, err = r.client.NetworkAccess().Settings().CustomBlockPage().Patch(ctx, body, nil)
	if err != nil {
		grapherrors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}
	data.ID = types.StringValue(singletonID)
	// Keep the singleton tracked even if its post-write read fails.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithConsistencyPredicate(crud.AttributesMatchPredicate(resp.State, readbackMatch))),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading custom block page after apply", err.Error())
	}
}

func (r *NetworkCustomBlockPageResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data NetworkCustomBlockPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		data.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	remote, err := r.client.NetworkAccess().Settings().CustomBlockPage().Get(ctx, nil)
	if err != nil {
		// A singleton cannot be deleted. In particular, do not interpret 404 (for
		// example, feature unavailable) as permission to drop its management state.
		info := grapherrors.GraphError(ctx, err)
		resp.Diagnostics.AddError(
			"Error reading custom block page",
			fmt.Sprintf(
				"Could not read %s (HTTP %d, %s): %s. Required permission: %s.",
				ResourceName,
				info.StatusCode,
				info.ErrorCode,
				info.ErrorMessage,
				r.ReadPermissions,
			),
		)
		return
	}
	if err := mapRemoteState(&data, remote); err != nil {
		resp.Diagnostics.AddError("Invalid custom block page response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkCustomBlockPageResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data NetworkCustomBlockPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		data.Timeouts.Update,
		UpdateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructResource(&data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid custom block page", err.Error())
		return
	}
	_, err = r.client.NetworkAccess().Settings().CustomBlockPage().Patch(ctx, body, nil)
	if err != nil {
		grapherrors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationUpdate,
			r.WritePermissions,
		)
		return
	}
	data.ID = types.StringValue(singletonID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName, crud.ReadWithConsistencyPredicate(crud.AttributesMatchPredicate(resp.State, readbackMatch))),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading custom block page after update", err.Error())
	}
}

// Delete relinquishes management without issuing DELETE, PATCH, or a reset.
func (r *NetworkCustomBlockPageResource) Delete(
	ctx context.Context,
	_ resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	resp.State.RemoveResource(ctx)
}
