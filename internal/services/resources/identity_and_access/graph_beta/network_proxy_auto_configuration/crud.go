// Package graphBetaNetworkProxyAutoConfiguration manages Global Secure Access custom proxy auto-configuration (PAC) files.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no explicitForwardProxyConfig or
//     proxyAutoConfigurationFiles request builder in networkaccess and no PAC file model in models/networkaccess.
//   - Learn: none. There is no proxyAutoConfigurationFile resource or API page; the endpoint and its name, content and
//     isEnabled shape follow the Entra admin centre's custom PAC file requests to
//     /networkAccess/explicitForwardProxyConfig/proxyAutoConfigurationFiles.
//   - Observed: checked by reading the SDK source on 2026-11-24; no live call, because the test tenant has no Global
//     Secure Access licence. The POST, GET, PATCH and DELETE go through custom_requests in raw JSON mode with pacResponse
//     and a map body. APIVersion is left empty so the requests use the beta client's own base URL, as the previous
//     transport did.
//
// Removal condition: the beta SDK ships a /networkAccess/explicitForwardProxyConfig/proxyAutoConfigurationFiles request
// builder with a PAC file model.
package graphBetaNetworkProxyAutoConfiguration

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

var errInvalidPACResponse = errors.New("graph omitted required custom PAC fields; existing state is retained")

// readbackPaths are the attributes the read-back after a write waits for.
var readbackPaths = crud.RootPaths("id", "name", "content", "is_enabled")

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
	body, err := constructResource(plan, nil)
	var remote *pacResponse
	if err == nil {
		_, err = customrequests.PostRequest(ctx, r.client.GetAdapter(), customrequests.PostRequestConfig{Endpoint: r.ResourcePath, RawBody: body, RawResponse: &remote}, nil, nil)
	}
	if err == nil && (remote == nil || remote.ID == nil || *remote.ID == "") {
		err = errInvalidPACResponse
	}
	if err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	// Retain the created ID if a subsequent read fails, so cleanup remains possible.
	plan.ID = types.StringValue(*remote.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicate := crud.AttributesMatchPredicate(req.Plan, crud.AttributeMatchOptions{Paths: readbackPaths})
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithConsistencyPredicate(predicate))
	readRequest := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readRequest, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
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
	// An empty or null body leaves remote nil, which MapRemoteStateToTerraform reports as an invalid response.
	var remote *pacResponse
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID.ValueString(),
		RawResponse:       &remote,
	}
	if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
		options := kiotaerrors.GraphErrorOptions{PreserveStateOnReadBadRequest: true, PreserveStateOnReadNotFound: false}
		kiotaerrors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions, options)
		return
	}
	if err := MapRemoteStateToTerraform(&state, remote); err != nil {
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
	body, err := constructResource(plan, &state)
	if err == nil && len(body) > 0 {
		patchConfig := customrequests.PatchRequestConfig{
			Endpoint:          r.ResourcePath,
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        plan.ID.ValueString(),
			RawBody:           body,
		}
		err = customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig)
	}
	if err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicate := crud.AttributesMatchPredicate(req.Plan, crud.AttributeMatchOptions{Paths: readbackPaths})
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName, crud.ReadWithConsistencyPredicate(predicate))
	readRequest := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readRequest, &crud.UpdateResponseContainer{UpdateResponse: resp}, opts); err != nil {
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
	deleteConfig := customrequests.DeleteRequestConfig{
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID.ValueString(),
	}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}
