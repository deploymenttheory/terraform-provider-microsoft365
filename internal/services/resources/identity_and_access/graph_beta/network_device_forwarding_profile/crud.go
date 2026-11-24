// Package graphBetaNetworkDeviceForwardingProfile manages Global Secure Access pre-logon (device) forwarding profiles.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no deviceForwardingProfiles request builder in
//     networkaccess and no device forwarding profile model in models/networkaccess. The forwarding profile v2 properties
//     (clientFallbackAction, localNetworkSettings, assignmentRules) are missing from networkaccess.ForwardingProfile too.
//   - Learn: none. There is no deviceForwardingProfiles resource or API page; the endpoint and its request shapes follow
//     the Entra admin centre's pre-logon forwarding profile requests to /networkaccess/deviceForwardingProfiles (see
//     https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding).
//   - Observed: checked by reading the SDK source on 2026-11-24; no live call, because the test tenant has no Global
//     Secure Access licence. The list, POST, GET, PATCH and DELETE go through custom_requests in raw JSON mode with
//     sharedmodels.ForwardingProfileWire. APIVersion is left empty so the requests use the beta client's own base URL, as
//     the previous transport did.
//
// Removal condition: the beta SDK ships a /networkAccess/deviceForwardingProfiles request builder whose model exposes
// clientFallbackAction, localNetworkSettings and assignmentRules.
package graphBetaNetworkDeviceForwardingProfile

import (
	"context"
	stderrors "errors"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

var (
	errInvalidProfileResponse = stderrors.New("graph returned a forwarding profile without an ID")
	errEmptyProfileResponse   = stderrors.New("graph returned an empty forwarding profile response")
	errProfileNotCreated      = stderrors.New("forwarding profile not found after failed create")
)

// claimedProfiles holds the lower-cased IDs of profiles already created, adopted or deleted in this provider process,
// so two resources with the same name and priority cannot both adopt one profile, and a stale list after a delete is
// never adopted on a re-create. It is a record of claims, not a lock: entries are never released.
var claimedProfiles sync.Map

// Variables so unit tests can shorten the post-create lookup.
var (
	createdLookupAttempts   = 12
	duplicateLookupAttempts = 3
	createdLookupInterval   = 5 * time.Second
)

// readbackPaths are the attributes the read-back after a write waits for.
var readbackPaths = crud.RootPaths("name", "description", "priority", "state", "client_fallback_action", "local_network_settings", "assignment_rules")

// Create handles the Create operation for network device forwarding profile resources.
//
// Operation: Creates a new network device forwarding profile
// API Calls:
//   - GET /networkaccess/deviceForwardingProfiles (custom profiles that exist before the POST)
//   - POST /networkaccess/deviceForwardingProfiles
//   - GET /networkaccess/deviceForwardingProfiles (only after a failed POST, to adopt the profile it created)
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
func (r *NetworkDeviceForwardingProfileResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkDeviceForwardingProfileResourceModel
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

	body, diags := constructResource(ctx, &plan, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// POST can return 504 after creating the profile, and a retried POST then fails
	// as a duplicate. List first so only a profile created by this request is adopted.
	var before struct {
		Value []sharedmodels.ForwardingProfileWire `json:"value"`
	}
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternNone,
		RawResponse:       &before,
	}
	_, listErr := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if listErr != nil {
		tflog.Warn(ctx, "Could not list forwarding profiles before create", map[string]any{"error": listErr.Error()})
	}
	existing := make(map[string]bool, len(before.Value))
	for _, p := range before.Value {
		if p.ID != nil && p.IsCustomProfile != nil && *p.IsCustomProfile {
			existing[*p.ID] = true
		}
	}

	var remote *sharedmodels.ForwardingProfileWire
	_, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), customrequests.PostRequestConfig{Endpoint: r.ResourcePath, RawBody: body, RawResponse: &remote}, nil, nil)
	if err == nil && (remote == nil || remote.ID == nil || *remote.ID == "") {
		err = errInvalidProfileResponse
	}
	// Without the list a match may be pre-existing, so report the error.
	if err != nil && listErr == nil {
		// Poll longest after errors that can hide a create (transport errors, which GraphError reports as 400
		// "URLError", 503 or 504; gateway failures; non-Graph 500s). 400/409 get a few checks (a retried POST fails
		// as a duplicate); anything else gets one.
		info := errors.GraphError(ctx, err)
		attempts := 1
		switch {
		case info.ErrorCode == "URLError", info.StatusCode == 502, info.StatusCode == 503, info.StatusCode == 504, info.StatusCode == 500 && !info.IsODataError:
			attempts = createdLookupAttempts
		case info.StatusCode == 400, info.StatusCode == 409:
			attempts = duplicateLookupAttempts
		}
		// The profile a failed POST created can appear late. Stop on lasting errors such as 403; transport errors
		// (400 "URLError") can clear.
		var created *sharedmodels.ForwardingProfileWire
		attempt := 0
		_ = crud.PollUntil(ctx, createdLookupInterval, func(ctx context.Context) (bool, error) {
			attempt++
			var after struct {
				Value []sharedmodels.ForwardingProfileWire `json:"value"`
			}
			getConfig := customrequests.GetRequestConfig{
				Endpoint:          r.ResourcePath,
				ResourceIDPattern: customrequests.ResourceIDPatternNone,
				RawResponse:       &after,
			}
			if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
				lookup := errors.GraphError(ctx, err)
				if attempt >= attempts || (lookup.StatusCode >= 400 && lookup.StatusCode < 500 && lookup.StatusCode != 429 && lookup.ErrorCode != "URLError") {
					return false, &crud.FatalPollError{Err: err}
				}
				return false, nil
			}
			for i, p := range after.Value {
				if p.ID == nil || existing[*p.ID] || p.IsCustomProfile == nil || !*p.IsCustomProfile || p.Name == nil || *p.Name != plan.Name.ValueString() {
					continue
				}
				if p.TrafficForwardingType == nil || *p.TrafficForwardingType != plan.TrafficForwardingType.ValueString() {
					continue
				}
				if p.Priority == nil || *p.Priority != plan.Priority.ValueInt32() {
					continue
				}
				// Claim atomically, so a profile another resource has claimed is skipped.
				if _, claimed := claimedProfiles.LoadOrStore(strings.ToLower(*p.ID), struct{}{}); claimed {
					continue
				}
				created = &after.Value[i]
				return true, nil
			}
			if attempt >= attempts {
				return false, &crud.FatalPollError{Err: errProfileNotCreated}
			}
			return false, errProfileNotCreated
		})
		if created != nil {
			tflog.Warn(ctx, "Forwarding profile create returned an error but the profile exists; adopting it", map[string]any{"status": info.StatusCode, "id": *created.ID})
			remote, err = created, nil
		}
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	claimedProfiles.Store(strings.ToLower(*remote.ID), struct{}{})
	// Persist the ID before reading back so a failed read still allows cleanup.
	plan.ID = types.StringPointerValue(remote.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Assignments apply asynchronously; wait until GET reflects the plan.
	predicate := crud.RequireStableReads(crud.AttributesMatchPredicate(req.Plan, crud.AttributeMatchOptions{Paths: readbackPaths}))
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithConsistencyPredicate(predicate))
	readRequest := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readRequest, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading forwarding profile after create", err.Error())
	}
}

// Read handles the Read operation for network device forwarding profile resources.
//
// Operation: Retrieves the network device forwarding profile by ID
// API Calls:
//   - GET /networkaccess/deviceForwardingProfiles/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
func (r *NetworkDeviceForwardingProfileResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkDeviceForwardingProfileResourceModel
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

	var remote *sharedmodels.ForwardingProfileWire
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID.ValueString(),
		RawResponse:       &remote,
	}
	_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && remote == nil {
		err = errEmptyProfileResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions)
		return
	}
	resp.Diagnostics.Append(MapRemoteStateToTerraform(ctx, &state, remote)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles the Update operation for network device forwarding profile resources.
//
// Operation: Updates the network device forwarding profile
// API Calls:
//   - PATCH /networkaccess/deviceForwardingProfiles/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
func (r *NetworkDeviceForwardingProfileResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkDeviceForwardingProfileResourceModel
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

	body, diags := constructResource(ctx, &plan, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(body) > 0 {
		patchConfig := customrequests.PatchRequestConfig{
			Endpoint:          r.ResourcePath,
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        state.ID.ValueString(),
			RawBody:           body,
		}
		if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicate := crud.RequireStableReads(crud.AttributesMatchPredicate(req.Plan, crud.AttributeMatchOptions{Paths: readbackPaths}))
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName, crud.ReadWithConsistencyPredicate(predicate))
	readRequest := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readRequest, &crud.UpdateResponseContainer{UpdateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading forwarding profile after update", err.Error())
	}
}

// Delete handles the Delete operation for network device forwarding profile resources.
//
// Operation: Deletes the network device forwarding profile
// API Calls:
//   - DELETE /networkaccess/deviceForwardingProfiles/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
func (r *NetworkDeviceForwardingProfileResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkDeviceForwardingProfileResourceModel
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
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	// A stale list could still show it; never adopt it on a re-create.
	claimedProfiles.Store(strings.ToLower(state.ID.ValueString()), struct{}{})
	resp.State.RemoveResource(ctx)
}
