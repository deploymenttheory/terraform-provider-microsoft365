package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
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
	unlock := lockProfile(profileType, profileID)
	defer unlock()
	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var profile forwardingProfile
	err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		profileTemplate(profileType)+"?$expand=policies($expand=policy)",
		map[string]string{"forwardingProfileId": profileID},
		nil,
		&profile,
	)
	if err != nil {
		err = fmt.Errorf("get forwarding profile %s: %w", profileID, err)
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}

	existing := findLink(profile.Policies, "", plan.ForwardingPolicyID.ValueString())
	switch {
	case existing != nil && existing.ConfigType != "" &&
		!strings.EqualFold(existing.ConfigType, plan.ConfigType.ValueString()):
		// Replacing could leave the policy unlinked, so leave it in place.
		resp.Diagnostics.AddError(
			"Forwarding policy is already linked with a different config type",
			fmt.Sprintf("Policy %s is linked to profile %s as %q (link %s), not %q. Check config_type.",
				plan.ForwardingPolicyID.ValueString(), profileID, existing.ConfigType, existing.ID,
				plan.ConfigType.ValueString()),
		)
		return
	case existing != nil && linkMatchesPlan(existing, &plan):
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
		if err := customrequests.JSONRequest(
			ctx,
			r.client.GetAdapter(),
			abstractions.PATCH,
			profileTemplate(profileType),
			map[string]string{"forwardingProfileId": profileID},
			constructPolicyDelta(constructRemovedLink(existing.ID)),
			nil,
		); err != nil {
			err = fmt.Errorf("patch forwarding profile %s policies: %w", profileID, err)

			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
		consecutive := 0
		if err := crud.PollUntil(ctx, readbackInterval, func(ctx context.Context) (bool, error) {
			var currentProfile forwardingProfile
			readErr := customrequests.JSONRequest(
				ctx,
				r.client.GetAdapter(),
				abstractions.GET,
				profileTemplate(profileType)+"?$expand=policies($expand=policy)",
				map[string]string{"forwardingProfileId": profileID},
				nil,
				&currentProfile,
			)
			absent := false
			if readErr != nil {
				readErr = fmt.Errorf("get forwarding profile %s: %w", profileID, readErr)
				switch errors.GraphError(ctx, readErr).StatusCode {
				case 404:
					absent = true
				case 401, 403:
					return false, &crud.FatalPollError{Err: readErr}
				default:
					consecutive = 0
					return false, nil
				}
			} else {
				absent = findLink(currentProfile.Policies, existing.ID, "") == nil
			}
			if absent {
				consecutive++
			} else {
				consecutive = 0
			}
			return consecutive >= stableReads, nil
		}); err != nil {
			if ctx.Err() != nil {
				err = fmt.Errorf("%w: %w", errWaitTimedOut, ctx.Err())
			}

			resp.Diagnostics.AddError("Error waiting for existing forwarding policy link removal",
				fmt.Sprintf("Link %s was removed; the next apply re-adds the policy. %s", existing.ID, err.Error()))
			return
		}
		if err := customrequests.JSONRequest(
			ctx,
			r.client.GetAdapter(),
			abstractions.PATCH,
			profileTemplate(profileType),
			map[string]string{"forwardingProfileId": profileID},
			constructPolicyDelta(constructLink(&plan)),
			nil,
		); err != nil {
			err = fmt.Errorf("patch forwarding profile %s policies: %w", profileID, err)

			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)

			configType := existing.ConfigType
			if configType == "" {
				configType = plan.ConfigType.ValueString()
			}
			restore := NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel{
				ProfileType:         types.StringValue(profileType),
				ForwardingProfileID: types.StringValue(profileID),
				ForwardingPolicyID:  types.StringValue(existing.Policy.ID),
				ConfigType:          types.StringValue(configType),
				State:               types.StringValue(strings.ToLower(existing.State)),
				Priority:            types.Int64PointerValue(existing.Priority),
			}
			if err := customrequests.JSONRequest(
				ctx,
				r.client.GetAdapter(),
				abstractions.PATCH,
				profileTemplate(profileType),
				map[string]string{"forwardingProfileId": profileID},
				constructPolicyDelta(constructLink(&restore)),
				nil,
			); err != nil {
				err = fmt.Errorf("patch forwarding profile %s policies: %w", profileID, err)
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
		if err := customrequests.JSONRequest(
			ctx,
			r.client.GetAdapter(),
			abstractions.PATCH,
			profileTemplate(profileType),
			map[string]string{"forwardingProfileId": profileID},
			constructPolicyDelta(constructLink(&plan)),
			nil,
		); err != nil {
			err = fmt.Errorf("patch forwarding profile %s policies: %w", profileID, err)

			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}

	// PATCH returns the profile, not the link; Read resolves the link by policy ID.
	planned := plan
	plan.ID = types.StringNull()
	plan.PolicyLinkID = types.StringNull()
	if plan.Priority.IsUnknown() {
		plan.Priority = types.Int64Null()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		readbackOptions(&planned, constants.TfOperationCreate),
	); err != nil {
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

	var profile forwardingProfile
	err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		profileTemplate(state.ProfileType.ValueString())+"?$expand=policies($expand=policy)",
		map[string]string{"forwardingProfileId": state.ForwardingProfileID.ValueString()},
		nil,
		&profile,
	)
	if err != nil {
		err = fmt.Errorf("get forwarding profile %s: %w", state.ForwardingProfileID.ValueString(), err)
	}
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions,
			errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}

	link := findLink(profile.Policies, state.PolicyLinkID.ValueString(), state.ForwardingPolicyID.ValueString())
	if link == nil {
		tflog.Warn(ctx, "Forwarding policy link not found on profile, removing from state", map[string]any{
			"forwarding_profile_id": state.ForwardingProfileID.ValueString(),
			"policy_link_id":        state.PolicyLinkID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	mapLinkToState(&state, link)
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
	unlock := lockProfile(profileType, profileID)
	defer unlock()
	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var profile forwardingProfile
	err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		profileTemplate(profileType)+"?$expand=policies($expand=policy)",
		map[string]string{"forwardingProfileId": profileID},
		nil,
		&profile,
	)
	if err != nil {
		err = fmt.Errorf("get forwarding profile %s: %w", profileID, err)
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
	link := findLink(profile.Policies, state.PolicyLinkID.ValueString(), policyID)
	if link == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.PATCH,
		profileTemplate(profileType),
		map[string]string{"forwardingProfileId": profileID},
		constructPolicyDelta(constructRemovedLink(link.ID)),
		nil,
	); err != nil {
		err = fmt.Errorf("patch forwarding profile %s policies: %w", profileID, err)

		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		if !resp.Diagnostics.HasError() {
			resp.State.RemoveResource(ctx)
		}
		return
	}
	consecutive := 0
	if err := crud.PollUntil(ctx, readbackInterval, func(ctx context.Context) (bool, error) {
		var currentProfile forwardingProfile
		readErr := customrequests.JSONRequest(
			ctx,
			r.client.GetAdapter(),
			abstractions.GET,
			profileTemplate(profileType)+"?$expand=policies($expand=policy)",
			map[string]string{"forwardingProfileId": profileID},
			nil,
			&currentProfile,
		)
		absent := false
		if readErr != nil {
			readErr = fmt.Errorf("get forwarding profile %s: %w", profileID, readErr)
			switch errors.GraphError(ctx, readErr).StatusCode {
			case 404:
				absent = true
			case 401, 403:
				return false, &crud.FatalPollError{Err: readErr}
			default:
				consecutive = 0
				return false, nil
			}
		} else {
			absent = findLink(currentProfile.Policies, link.ID, "") == nil
		}
		if absent {
			consecutive++
		} else {
			consecutive = 0
		}
		return consecutive >= stableReads, nil
	}); err != nil {
		if ctx.Err() != nil {
			err = fmt.Errorf("%w: %w", errWaitTimedOut, ctx.Err())
		}

		resp.Diagnostics.AddError("Error waiting for forwarding policy link removal", err.Error())
		return
	}
	resp.State.RemoveResource(ctx)
}

// linkMatchesPlan checks whether the policy link matches the planned configuration.
func linkMatchesPlan(link *policyLink, plan *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel) bool {
	if !strings.EqualFold(link.State, plan.State.ValueString()) {
		return false
	}
	if link.ConfigType != "" && !strings.EqualFold(link.ConfigType, plan.ConfigType.ValueString()) {
		return false
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		return link.Priority != nil && *link.Priority == plan.Priority.ValueInt64()
	}
	return true
}

// mapLinkToState keeps configured casing for IDs that Graph echoes back case-insensitively.
func mapLinkToState(state *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel, link *policyLink) {
	if !strings.EqualFold(state.PolicyLinkID.ValueString(), link.ID) {
		state.PolicyLinkID = types.StringValue(link.ID)
	}
	state.ID = compositeID(state.ForwardingProfileID.ValueString(), state.PolicyLinkID.ValueString())
	if link.Policy != nil && !strings.EqualFold(state.ForwardingPolicyID.ValueString(), link.Policy.ID) {
		state.ForwardingPolicyID = types.StringValue(link.Policy.ID)
	}
	if link.ConfigType != "" {
		state.ConfigType = types.StringValue(link.ConfigType)
	}
	if link.State != "" {
		state.State = types.StringValue(strings.ToLower(link.State))
	}
	if link.Priority != nil {
		state.Priority = types.Int64Value(*link.Priority)
	} else {
		state.Priority = types.Int64Null()
	}
}

// readbackOptions configures shared read retries to verify the preceding write.
func readbackOptions(
	plan *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel,
	operation string,
) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.RetryInterval = readbackInterval
	opts.MaxRetries = 60
	consecutive := 0
	opts.ConsistencyPredicate = func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel
		if state.Raw.IsNull() || state.Get(ctx, &actual).HasError() || !matchesPlan(&actual, plan) {
			consecutive = 0
			return false
		}
		consecutive++
		return consecutive >= stableReads
	}
	return opts
}

// matchesPlan checks whether remote state matches the planned settings.
func matchesPlan(actual, plan *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel) bool {
	if actual.PolicyLinkID.ValueString() == "" ||
		!strings.EqualFold(actual.ForwardingPolicyID.ValueString(), plan.ForwardingPolicyID.ValueString()) ||
		!actual.State.Equal(plan.State) ||
		!actual.ConfigType.Equal(plan.ConfigType) {
		return false
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		return actual.Priority.Equal(plan.Priority)
	}
	return true
}

// readbackInterval is a variable so unit tests can shorten it.
var readbackInterval = 2 * time.Second

var errWaitTimedOut = stderrors.New("timed out waiting for forwarding policy link change to propagate")

// profileLocks serialises writes to one profile: each is a read-modify-write,
// and Terraform runs sibling links in parallel.
var profileLocks sync.Map

// lockProfile serializes changes to policy links on the same forwarding profile.
func lockProfile(profileType, profileID string) (unlock func()) {
	key := profileType + "/" + strings.ToLower(profileID)
	mu, _ := profileLocks.LoadOrStore(key, &sync.Mutex{})
	m, _ := mu.(*sync.Mutex)
	m.Lock()
	return m.Unlock
}
