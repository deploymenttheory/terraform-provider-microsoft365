// SDK deviation: unknown enum value
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, networkaccess.SettingsConditionalAccessRequestBuilder.Get and
//     models/networkaccess.ConditionalAccessSettings — a signalingStatus value the generated Status enum does not know is
//     dropped to nil, so an unsupported value reads as a missing one. A JSON null response body also panics in the adapter's
//     backing-store parse instead of returning an error.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/networkaccess-conditionalaccesssettings-get?view=graph-rest-beta
//   - Observed: 2026-11-24, GET /beta/networkAccess/settings/conditionalAccess through the generated builder against httpmock,
//     not the live API (the test tenant has no Global Secure Access licence). Body {"signalingStatus":"bogus"} returned a
//     nil status; body null panicked with a nil pointer dereference. The GET therefore goes through
//     custom_requests.GetRequestByResourceId and signalingStatus is decoded as a string, so an unknown value or a null body
//     is reported and the existing state is kept. APIVersion is left empty so the GET uses the beta client's own base URL, as
//     the generated builder does. The PATCH uses the generated builder.
//
// Removal condition: the beta SDK returns an error instead of panicking on a JSON null body and surfaces an unknown
// signalingStatus instead of dropping it, so the GET can use SettingsConditionalAccessRequestBuilder.Get.
package graphBetaNetworkConditionalAccessSettings

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	kiotahttp "github.com/microsoft/kiota-http-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// readbackMatch is what the read-back after a write waits for: the singleton id and the written signaling status.
var readbackMatch = crud.AttributeMatchOptions{Paths: crud.RootPaths("id", "signaling_status")}

var errInvalidSettingsResponse = stderrors.New("invalid conditional access settings response")

// Create adopts the existing singleton and PATCHes only explicitly configured settings.
func (r *NetworkConditionalAccessSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config NetworkConditionalAccessSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	if config.SignalingStatus.IsNull() {
		getConfig := customrequests.GetRequestConfig{Endpoint: r.ResourcePath, ResourceIDPattern: customrequests.ResourceIDPatternNone}
		raw, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
		var remote *conditionalAccessSettingsResponse
		if err == nil {
			if len(raw) == 0 {
				err = errInvalidSettingsResponse
			} else if decodeErr := json.Unmarshal(raw, &remote); decodeErr != nil {
				err = fmt.Errorf("%w: %w", errInvalidSettingsResponse, decodeErr)
			} else if remote == nil {
				err = errInvalidSettingsResponse
			}
		}
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
			return
		}
		MapRemoteStateToTerraform(&plan, remote, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		body, err := constructResource(plan.SignalingStatus)
		if err != nil {
			resp.Diagnostics.AddError("Invalid signaling status", err.Error())
			return
		}
		// Avoid Kiota replaying an uncompressed body with a stale gzip header on retry.
		options := []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)}
		patchConfig := &networkaccess.SettingsConditionalAccessRequestBuilderPatchRequestConfiguration{Options: options}
		if _, err = r.client.NetworkAccess().Settings().ConditionalAccess().Patch(ctx, body, patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
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
	predicate := crud.ReadWithConsistencyPredicate(crud.AttributesMatchPredicate(resp.State, readbackMatch))
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, predicate)
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading settings after apply", err.Error())
	}
}

func (r *NetworkConditionalAccessSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkConditionalAccessSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	getConfig := customrequests.GetRequestConfig{Endpoint: r.ResourcePath, ResourceIDPattern: customrequests.ResourceIDPatternNone}
	raw, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	var remote *conditionalAccessSettingsResponse
	if err == nil {
		if len(raw) == 0 {
			err = errInvalidSettingsResponse
		} else if decodeErr := json.Unmarshal(raw, &remote); decodeErr != nil {
			err = fmt.Errorf("%w: %w", errInvalidSettingsResponse, decodeErr)
		} else if remote == nil {
			err = errInvalidSettingsResponse
		}
	}
	if err != nil {
		// A missing singleton endpoint can mean licensing or availability, not deletion.
		readOptions := errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true, PreserveStateOnReadNotFound: true}
		errors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions, readOptions)
		return
	}
	MapRemoteStateToTerraform(&state, remote, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), singletonID)...)
	}
}

func (r *NetworkConditionalAccessSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, config, state NetworkConditionalAccessSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
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
	// Avoid Kiota replaying an uncompressed body with a stale gzip header on retry.
	options := []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)}
	patchConfig := &networkaccess.SettingsConditionalAccessRequestBuilderPatchRequestConfiguration{Options: options}
	if _, err = r.client.NetworkAccess().Settings().ConditionalAccess().Patch(ctx, body, patchConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}
	plan.ID = types.StringValue(singletonID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicate := crud.ReadWithConsistencyPredicate(crud.AttributesMatchPredicate(resp.State, readbackMatch))
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName, predicate)
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.UpdateResponseContainer{UpdateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading settings after update", err.Error())
	}
}

// Delete releases Terraform ownership without changing the tenant configuration.
func (r *NetworkConditionalAccessSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Removing conditional access settings from Terraform state; remote settings are preserved")
	resp.State.RemoveResource(ctx)
}
