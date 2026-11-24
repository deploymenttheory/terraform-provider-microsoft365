// Package graphBetaNetworkDefaultForwardingProfile manages the built-in Global Secure Access forwarding profiles.
//
// SDK deviation: missing property
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, models/networkaccess.ForwardingProfile (through
//     ForwardingProfileBase and Profile) has no clientFallbackAction, localNetworkSettings, forwardingSourceMatching or
//     assignmentRules property, so the generated serializer cannot send them and the parser leaves them in additional data.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-forwardingprofile?view=graph-rest-beta
//     does not list these properties either; they are not in the beta $metadata.
//   - Observed: the property shapes follow the Entra admin centre's forwarding profile v2 requests to
//     /networkaccess/forwardingProfiles (see https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding).
//     Checked by reading the SDK source on 2026-11-24; no live call, because the test tenant has no Global Secure Access
//     licence. The list, GET and PATCH go through custom_requests in raw JSON mode with sharedmodels.ForwardingProfileWire.
//     APIVersion is left empty so the requests use the beta client's own base URL, as the previous transport did.
//
// Removal condition: networkaccess.ForwardingProfile exposes clientFallbackAction, localNetworkSettings,
// forwardingSourceMatching and assignmentRules.
package graphBetaNetworkDefaultForwardingProfile

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

var (
	errDefaultProfileNotFound = stderrors.New("default forwarding profile not found; the tenant must be onboarded to Global Secure Access and use traffic forwarding profiles v2")
	errEmptyProfileResponse   = stderrors.New("graph returned an empty forwarding profile response")
)

// readbackPaths are the settings the read-back after a write waits for.
var readbackPaths = crud.RootPaths("state", "client_fallback_action", "forwarding_source_matching", "local_network_settings", "assignment_rules")

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

	// The built-in profile's ID varies by tenant, so find it by traffic type.
	var profiles struct {
		Value []sharedmodels.ForwardingProfileWire `json:"value"`
	}
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternNone,
		RawResponse:       &profiles,
	}
	_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	var remote *sharedmodels.ForwardingProfileWire
	if err == nil {
		err = fmt.Errorf("%w: %q", errDefaultProfileNotFound, plan.TrafficForwardingType.ValueString())
		for i, p := range profiles.Value {
			if p.TrafficForwardingType != nil && *p.TrafficForwardingType == plan.TrafficForwardingType.ValueString() && p.IsCustomProfile != nil && !*p.IsCustomProfile {
				remote, err = &profiles.Value[i], nil
				break
			}
		}
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	var current NetworkDefaultForwardingProfileResourceModel
	current.TrafficForwardingType = plan.TrafficForwardingType
	resp.Diagnostics.Append(MapRemoteStateToTerraform(ctx, &current, remote)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Send only configured values that differ from the adopted profile.
	body, diags := constructResource(ctx, plan, current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(body) > 0 {
		patchConfig := customrequests.PatchRequestConfig{
			Endpoint:          r.ResourcePath,
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        *remote.ID,
			RawBody:           body,
		}
		if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}
	plan.ID = types.StringPointerValue(remote.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicate := crud.RequireStableReads(crud.AttributesMatchPredicate(req.Plan, crud.AttributeMatchOptions{Paths: readbackPaths}))
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithConsistencyPredicate(predicate))
	readRequest := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readRequest, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
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

	var remote *sharedmodels.ForwardingProfileWire
	var err error
	if state.ID.IsNull() || state.ID.IsUnknown() {
		// Imported by traffic type: the built-in profile's ID varies by tenant.
		var profiles struct {
			Value []sharedmodels.ForwardingProfileWire `json:"value"`
		}
		getConfig := customrequests.GetRequestConfig{
			Endpoint:          r.ResourcePath,
			ResourceIDPattern: customrequests.ResourceIDPatternNone,
			RawResponse:       &profiles,
		}
		_, err = customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
		if err == nil {
			err = fmt.Errorf("%w: %q", errDefaultProfileNotFound, state.TrafficForwardingType.ValueString())
			for i, p := range profiles.Value {
				if p.TrafficForwardingType != nil && *p.TrafficForwardingType == state.TrafficForwardingType.ValueString() && p.IsCustomProfile != nil && !*p.IsCustomProfile {
					remote, err = &profiles.Value[i], nil
					break
				}
			}
		}
	} else {
		getConfig := customrequests.GetRequestConfig{
			Endpoint:          r.ResourcePath,
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        state.ID.ValueString(),
			RawResponse:       &remote,
		}
		_, err = customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
		if err == nil && remote == nil {
			err = errEmptyProfileResponse
		}
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

	body, diags := constructResource(ctx, plan, state)
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
