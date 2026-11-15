package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	abstractions "github.com/microsoft/kiota-abstractions-go"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

const (
	forwardingPolicyLinkODataType = "#microsoft.graph.networkaccess.forwardingPolicyLink"
	forwardingPolicyODataType     = "#microsoft.graph.networkaccess.forwardingPolicy"
	// Links are eventually consistent; a single matching read can come from a fresh replica
	// followed by a stale one.
	stableReads = 3
)

// readbackInterval is a variable so unit tests can shorten it.
var readbackInterval = 2 * time.Second

var errWaitTimedOut = errors.New("timed out waiting for forwarding policy link change to propagate")

func profileTemplate(profileType string) string {
	if profileType == profileTypeDevice {
		return "/networkAccess/deviceForwardingProfiles/{forwardingProfileId}"
	}
	return "/networkAccess/forwardingProfiles/{forwardingProfileId}"
}

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) getProfile(
	ctx context.Context,
	profileType, profileID string,
) (*forwardingProfile, error) {
	var profile forwardingProfile
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		profileTemplate(profileType)+"?$expand=policies($expand=policy)",
		map[string]string{"forwardingProfileId": profileID},
		nil,
		&profile,
	); err != nil {
		return nil, fmt.Errorf("get forwarding profile %s: %w", profileID, err)
	}
	return &profile, nil
}

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) patchPolicies(
	ctx context.Context,
	profileType, profileID string,
	delta map[string]any,
) error {
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.PATCH,
		profileTemplate(profileType),
		map[string]string{"forwardingProfileId": profileID},
		map[string]any{"policies@delta": []map[string]any{delta}},
		nil,
	); err != nil {
		return fmt.Errorf("patch forwarding profile %s policies: %w", profileID, err)
	}
	return nil
}

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) addLink(
	ctx context.Context,
	plan *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel,
) error {
	link := map[string]any{
		"@odata.type": forwardingPolicyLinkODataType,
		"state":       plan.State.ValueString(),
		"configType":  plan.ConfigType.ValueString(),
		"policy": map[string]any{
			"@odata.type": forwardingPolicyODataType,
			"id":          plan.ForwardingPolicyID.ValueString(),
		},
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		link["priority"] = plan.Priority.ValueInt64()
	}
	return r.patchPolicies(ctx, plan.ProfileType.ValueString(), plan.ForwardingProfileID.ValueString(), link)
}

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) removeLink(
	ctx context.Context,
	profileType, profileID, linkID string,
) error {
	return r.patchPolicies(ctx, profileType, profileID, map[string]any{
		"@removed":    map[string]any{"reason": "deleted"},
		"@odata.type": forwardingPolicyLinkODataType,
		"id":          linkID,
	})
}

// findLink prefers the link ID and falls back to the policy ID because Graph
// assigns a new link ID whenever a policy is removed and re-added.
func findLink(links []policyLink, linkID, policyID string) *policyLink {
	if linkID != "" {
		for i := range links {
			if strings.EqualFold(links[i].ID, linkID) {
				return &links[i]
			}
		}
	}
	if policyID != "" {
		for i := range links {
			if links[i].Policy != nil && strings.EqualFold(links[i].Policy.ID, policyID) {
				return &links[i]
			}
		}
	}
	return nil
}

// waitUntil polls check until it reports true on stableReads consecutive calls.
func waitUntil(ctx context.Context, check func(context.Context) (bool, error)) error {
	consecutive := 0
	for {
		ok, err := check(ctx)
		if err != nil {
			return err
		}
		if ok {
			consecutive++
			if consecutive >= stableReads {
				return nil
			}
		} else {
			consecutive = 0
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: %w", errWaitTimedOut, ctx.Err())
		case <-time.After(readbackInterval):
		}
	}
}

// waitForLinkRemoved returns once linkID is consistently absent from the profile.
func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) waitForLinkRemoved(
	ctx context.Context,
	profileType, profileID, linkID string,
) error {
	return waitUntil(ctx, func(ctx context.Context) (bool, error) {
		profile, err := r.getProfile(ctx, profileType, profileID)
		if err != nil {
			if isNotFound(ctx, err) {
				return true, nil
			}
			return false, err
		}
		return findLink(profile.Policies, linkID, "") == nil, nil
	})
}
