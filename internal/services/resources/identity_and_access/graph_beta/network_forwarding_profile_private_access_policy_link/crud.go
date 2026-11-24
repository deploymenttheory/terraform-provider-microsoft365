// Package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink links Private Access forwarding policies to
// forwarding profiles.
//
// SDK deviation: missing property
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, models/networkaccess.ForwardingPolicyLink (through
//     PolicyLink) has no configType property, so the generated serializer cannot send it and the parser leaves it in
//     additional data. networkaccess.ForwardingProfile also has no policies@delta annotation, and a link removal must
//     write @removed before @odata.type, which a map in additional data cannot guarantee.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-forwardingpolicylink?view=graph-rest-beta
//     does not list configType either.
//   - Observed: the configType values (appPolicy, dnsPolicy) and the policies@delta add and remove shapes follow the
//     Entra admin centre's requests to /networkAccess/forwardingProfiles/{id}. Graph answers a removal whose @odata.type
//     comes before @removed with "Invalid request parameters". Checked against the SDK source on 2026-11-24; no live
//     call, because the test tenant has no Global Secure Access licence.
//
// Removal condition: networkaccess.ForwardingPolicyLink exposes configType and the SDK can send an ordered
// policies@delta collection on a forwarding profile PATCH.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no deviceForwardingProfiles request builder in
//     networkaccess, so the policy links of pre-logon (device) profiles cannot be read or changed through the SDK.
//   - Learn: none. There is no deviceForwardingProfiles resource or API page; the endpoint follows the Entra admin
//     centre's pre-logon forwarding profile requests to /networkAccess/deviceForwardingProfiles/{id}.
//   - Observed: checked by reading the SDK source on 2026-11-24; no live call, because the test tenant has no Global
//     Secure Access licence. Both profile types go through custom_requests in raw JSON mode with forwardingProfile and
//     map bodies. APIVersion is left empty so the requests use the beta client's own base URL, as the previous
//     transport did.
//
// Removal condition: the beta SDK ships a /networkAccess/deviceForwardingProfiles request builder.
package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// profileLinksLockScope serialises writes to one profile's links: each is a read-modify-write, and Terraform runs
// sibling links in parallel.
const profileLinksLockScope = "forwarding profile policy links"

// expandPolicyLinks is appended to the profile GET so the response carries each link and its policy.
const expandPolicyLinks = "?$expand=policies($expand=policy)"

// profileEndpoints maps profile_type to its forwarding profile collection.
var profileEndpoints = map[string]string{
	profileTypeUser:   "/networkAccess/forwardingProfiles",
	profileTypeDevice: "/networkAccess/deviceForwardingProfiles",
}

// readbackInterval is a variable so unit tests can shorten it.
var readbackInterval = 2 * time.Second

var (
	errWaitTimedOut         = stderrors.New("timed out waiting for forwarding policy link change to propagate")
	errEmptyProfileResponse = stderrors.New("graph returned an empty forwarding profile response")
)

// Create handles the Create operation for network forwarding profile private access policy link resources.
//
// Operation: Adds or replaces the Private Access policy link using policies@delta
// API Calls:
//   - GET /networkAccess/forwardingProfiles/{id}?$expand=policies($expand=policy)
//   - PATCH /networkAccess/forwardingProfiles/{id} (policies@delta)
//   - GET /networkAccess/deviceForwardingProfiles/{id}?$expand=policies($expand=policy) (device profiles)
//   - PATCH /networkAccess/deviceForwardingProfiles/{id} (device profiles; policies@delta)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-forwardingpolicylink?view=graph-rest-beta
func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Lock first so the timeout does not run while queued behind sibling links.
	profileType := plan.ProfileType.ValueString()
	profileID := plan.ForwardingProfileID.ValueString()
	unlock, err := crud.AcquireKeyedLock(ctx, profileLinksLockScope, profileType+"/"+profileID)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for forwarding profile lock", err.Error())
		return
	}
	defer unlock()
	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var profile *forwardingProfile
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          profileEndpoints[profileType],
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        profileID,
		EndpointSuffix:    expandPolicyLinks,
		RawResponse:       &profile,
	}
	_, err = customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && profile == nil {
		err = errEmptyProfileResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}

	// Every write is a policies@delta PATCH on the profile.
	patchConfig := customrequests.PatchRequestConfig{Endpoint: profileEndpoints[profileType], ResourceIDPattern: customrequests.ResourceIDPatternSegment, ResourceID: profileID}
	existing := findPolicyLinkByIDOrPolicyID(profile.Policies, "", plan.ForwardingPolicyID.ValueString())
	switch {
	case existing != nil && existing.ConfigType != "" && !strings.EqualFold(existing.ConfigType, plan.ConfigType.ValueString()):
		// Replacing could leave the policy unlinked, so leave it in place.
		resp.Diagnostics.AddError(
			"Forwarding policy is already linked with a different config type",
			fmt.Sprintf("Policy %s is linked to profile %s as %q (link %s), not %q. Check config_type.",
				plan.ForwardingPolicyID.ValueString(), profileID, existing.ConfigType, existing.ID,
				plan.ConfigType.ValueString()),
		)
		return
	case existing != nil && strings.EqualFold(existing.State, plan.State.ValueString()) &&
		(existing.ConfigType == "" || strings.EqualFold(existing.ConfigType, plan.ConfigType.ValueString())) &&
		(plan.Priority.IsNull() || plan.Priority.IsUnknown() || (existing.Priority != nil && *existing.Priority == plan.Priority.ValueInt64())):
		// New user profiles auto-link every private app policy; adopt it instead of failing.
		tflog.Warn(ctx, "Adopting existing forwarding policy link that already matches the configuration", map[string]any{
			"forwarding_profile_id": profileID,
			"policy_link_id":        existing.ID,
		})
	case existing != nil:
		tflog.Warn(ctx, "Replacing existing forwarding policy link whose state or priority differs", map[string]any{
			"forwarding_profile_id": profileID,
			"policy_link_id":        existing.ID,
		})
		// Graph rejects link updates and a remove+add in one delta, so this takes two writes.
		removal := map[string]any{"policies@delta": []any{removedLink{Removed: removedReason{Reason: "deleted"}, ODataType: forwardingPolicyLinkODataType, ID: existing.ID}}}
		patchConfig.RawBody = removal
		if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
		if err := r.waitForPolicyLinkRemoval(ctx, profileType, profileID, existing.ID); err != nil {
			resp.Diagnostics.AddError("Error waiting for existing forwarding policy link removal",
				fmt.Sprintf("Link %s was removed; the next apply re-adds the policy. %s", existing.ID, err.Error()))
			return
		}
		patchConfig.RawBody = constructResource(&plan)
		if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			// Re-add the link that was removed above: Microsoft does not recreate removed auto-links.
			configType := existing.ConfigType
			if configType == "" {
				configType = plan.ConfigType.ValueString()
			}
			previous := NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel{
				ProfileType:         types.StringValue(profileType),
				ForwardingProfileID: types.StringValue(profileID),
				ForwardingPolicyID:  types.StringValue(existing.Policy.ID),
				ConfigType:          types.StringValue(configType),
				State:               types.StringValue(strings.ToLower(existing.State)),
				Priority:            types.Int64PointerValue(existing.Priority),
			}
			patchConfig.RawBody = constructResource(&previous)
			if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
				resp.Diagnostics.AddError(
					"Previous forwarding policy link was removed and could not be restored",
					fmt.Sprintf("Create removed link %s from profile %s before its add failed, and re-adding it failed: %s",
						existing.ID, profileID, err.Error()),
				)
				return
			}
			resp.Diagnostics.AddWarning(
				"Previous forwarding policy link restored",
				fmt.Sprintf("The add failed, so the link that Create removed from profile %s was re-added with its "+
					"previous state and priority. Graph gives it a new link ID.", profileID),
			)
			return
		}
	default:
		patchConfig.RawBody = constructResource(&plan)
		if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}

	// PATCH returns the profile, not the link; Read resolves the link by policy ID.
	plan.ID = types.StringNull()
	plan.PolicyLinkID = types.StringNull()
	if plan.Priority.IsUnknown() {
		plan.Priority = types.Int64Null()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait until the link reflects the plan. An unset priority is chosen by Graph, and Graph may return the policy ID
	// in another case.
	caseInsensitive := true
	matches := crud.AttributesMatchPredicate(req.Plan, crud.AttributeMatchOptions{
		Paths:              crud.RootPaths("forwarding_policy_id", "state", "config_type", "priority"),
		RequireKnown:       crud.RootPaths("policy_link_id"),
		IgnoreNullExpected: true,
		Overrides:          []crud.AttributeMatchOverride{{Path: path.MatchRoot("forwarding_policy_id"), CaseInsensitiveStrings: &caseInsensitive}},
	})
	predicate := crud.ReadWithConsistencyPredicate(crud.RequireStableReads(matches))
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithRetryInterval(readbackInterval), crud.ReadWithMaxRetries(60), predicate)
	readRequest := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readRequest, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading forwarding policy link after create", err.Error())
	}
}

// Read handles the Read operation for network forwarding profile private access policy link resources.
//
// Operation: Retrieves the linked policy from the selected user or device forwarding profile
// API Calls:
//   - GET /networkAccess/forwardingProfiles/{id}?$expand=policies($expand=policy)
//   - GET /networkAccess/deviceForwardingProfiles/{id}?$expand=policies($expand=policy) (device profiles)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-forwardingpolicylink?view=graph-rest-beta
func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operation := constants.TfOperationRead
	if op, ok := ctx.Value("retry_operation").(string); ok && op != "" {
		operation = op
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var profile *forwardingProfile
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          profileEndpoints[state.ProfileType.ValueString()],
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ForwardingProfileID.ValueString(),
		EndpointSuffix:    expandPolicyLinks,
		RawResponse:       &profile,
	}
	_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && profile == nil {
		err = errEmptyProfileResponse
	}
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions, errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}

	link := findPolicyLinkByIDOrPolicyID(profile.Policies, state.PolicyLinkID.ValueString(), state.ForwardingPolicyID.ValueString())
	if link == nil {
		tflog.Warn(ctx, "Forwarding policy link not found on profile, removing from state", map[string]any{
			"forwarding_profile_id": state.ForwardingProfileID.ValueString(),
			"policy_link_id":        state.PolicyLinkID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	MapRemoteStateToTerraform(&state, link)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles the Update operation for network forwarding profile private access policy link resources.
//
// Operation: Stores plan-only changes; changes to link settings require replacement
// API Calls:
//   - None (replacement required for link settings)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-forwardingpolicylink?view=graph-rest-beta
func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles the Delete operation for network forwarding profile private access policy link resources.
//
// Operation: Removes the owned policy link using policies@delta while retaining its profile and policy
// API Calls:
//   - GET /networkAccess/forwardingProfiles/{id}?$expand=policies($expand=policy)
//   - PATCH /networkAccess/forwardingProfiles/{id} (policies@delta)
//   - GET /networkAccess/deviceForwardingProfiles/{id}?$expand=policies($expand=policy) (device profiles)
//   - PATCH /networkAccess/deviceForwardingProfiles/{id} (device profiles; policies@delta)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-forwardingpolicylink?view=graph-rest-beta
func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	profileType := state.ProfileType.ValueString()
	profileID := state.ForwardingProfileID.ValueString()
	unlock, err := crud.AcquireKeyedLock(ctx, profileLinksLockScope, profileType+"/"+profileID)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for forwarding profile lock", err.Error())
		return
	}
	defer unlock()
	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var profile *forwardingProfile
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          profileEndpoints[profileType],
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        profileID,
		EndpointSuffix:    expandPolicyLinks,
		RawResponse:       &profile,
	}
	_, err = customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && profile == nil {
		err = errEmptyProfileResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
		if !resp.Diagnostics.HasError() {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	// Match by policy only when the link ID was never stored (failed read-back).
	// Otherwise the policy's link may belong to a create_before_destroy replacement.
	policyID := ""
	if state.PolicyLinkID.ValueString() == "" {
		policyID = state.ForwardingPolicyID.ValueString()
	}
	link := findPolicyLinkByIDOrPolicyID(profile.Policies, state.PolicyLinkID.ValueString(), policyID)
	if link == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	removal := map[string]any{"policies@delta": []any{removedLink{Removed: removedReason{Reason: "deleted"}, ODataType: forwardingPolicyLinkODataType, ID: link.ID}}}
	patchConfig := customrequests.PatchRequestConfig{
		Endpoint:          profileEndpoints[profileType],
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        profileID,
		RawBody:           removal,
	}
	if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		if !resp.Diagnostics.HasError() {
			resp.State.RemoveResource(ctx)
		}
		return
	}
	if err := r.waitForPolicyLinkRemoval(ctx, profileType, profileID, link.ID); err != nil {
		resp.Diagnostics.AddError("Error waiting for forwarding policy link removal", err.Error())
		return
	}
	resp.State.RemoveResource(ctx)
}

// waitForPolicyLinkRemoval returns once linkID is absent from the profile on crud.StableReads reads in a row, because a
// read that no longer shows the link can be followed by a stale one that does. A 404 for the profile counts as absent;
// 401 and 403 stop the wait, and other errors are retried because the removal itself succeeded.
func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) waitForPolicyLinkRemoval(ctx context.Context, profileType, profileID, linkID string) error {
	absent := 0
	err := crud.PollUntil(ctx, readbackInterval, func(ctx context.Context) (bool, error) {
		var profile *forwardingProfile
		getConfig := customrequests.GetRequestConfig{
			Endpoint:          profileEndpoints[profileType],
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        profileID,
			EndpointSuffix:    expandPolicyLinks,
			RawResponse:       &profile,
		}
		_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
		removed := err == nil && profile != nil && findPolicyLinkByIDOrPolicyID(profile.Policies, linkID, "") == nil
		if err != nil {
			switch errors.GraphError(ctx, err).StatusCode {
			case 404:
				removed = true
			case 401, 403:
				return false, &crud.FatalPollError{Err: err}
			}
		}
		if !removed {
			absent = 0
			return false, nil
		}
		absent++
		return absent >= crud.StableReads, nil
	})
	if err != nil && ctx.Err() != nil {
		return fmt.Errorf("%w: %w", errWaitTimedOut, ctx.Err())
	}
	if err != nil {
		return fmt.Errorf("read forwarding profile %s: %w", profileID, err)
	}
	return nil
}
