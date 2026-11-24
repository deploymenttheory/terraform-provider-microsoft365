// Package graphBetaNetworkExplicitForwardProxy manages the Global Secure Access explicit forward proxy configuration.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no explicitForwardProxyConfig request builder in
//     networkaccess and no explicit forward proxy model in models/networkaccess.
//   - Learn: none. There is no explicitForwardProxyConfig resource or API page; the endpoint and its internetAccess shape
//     follow the Entra admin centre's explicit forward proxy requests to /networkAccess/explicitForwardProxyConfig.
//   - Observed: checked by reading the SDK source on 2026-11-24; no live call, because the test tenant has no Global
//     Secure Access licence. The GET and PATCH go through custom_requests in raw JSON mode with proxyResponse and a map
//     body. APIVersion is left empty so the requests use the beta client's own base URL, as the previous transport did.
//
// Removal condition: the beta SDK ships a /networkAccess/explicitForwardProxyConfig request builder with a model for its
// internetAccess settings.
package graphBetaNetworkExplicitForwardProxy

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

var errInvalidProxyResponse = errors.New("graph omitted required explicit proxy settings or returned unsupported session affinity options; existing state is retained")

// readbackPaths are the attributes the read-back after a write waits for.
var readbackPaths = crud.RootPaths("id", "internet_access")

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
	body, err := constructResource(plan, nil)
	if err == nil && len(body) > 0 {
		patchConfig := customrequests.PatchRequestConfig{
			Endpoint:          r.ResourcePath,
			ResourceIDPattern: customrequests.ResourceIDPatternNone,
			RawBody:           body,
		}
		err = customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig)
	}
	if err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	// Record adoption before readback so an unavailable replica does not lose management state.
	plan.ID = types.StringValue("explicitForwardProxyConfig")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicate := crud.AttributesMatchPredicate(req.Plan, crud.AttributeMatchOptions{Paths: readbackPaths})
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithConsistencyPredicate(predicate))
	readRequest := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readRequest, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
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
	// An empty or null body leaves remote nil, which MapRemoteStateToTerraform reports as an invalid response.
	var remote *proxyResponse
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternNone,
		RawResponse:       &remote,
	}
	if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
		options := kiotaerrors.GraphErrorOptions{PreserveStateOnReadBadRequest: true, PreserveStateOnReadNotFound: true}
		kiotaerrors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions, options)
		return
	}
	if err := MapRemoteStateToTerraform(&state, remote); err != nil {
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
	body, err := constructResource(plan, &state)
	if err == nil && len(body) > 0 {
		patchConfig := customrequests.PatchRequestConfig{
			Endpoint:          r.ResourcePath,
			ResourceIDPattern: customrequests.ResourceIDPatternNone,
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
