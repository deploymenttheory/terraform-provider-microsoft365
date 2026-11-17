package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	abstractions "github.com/microsoft/kiota-abstractions-go"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

const (
	forwardingPolicyLinkODataType = "#microsoft.graph.networkaccess.forwardingPolicyLink"
	forwardingPolicyODataType     = "#microsoft.graph.networkaccess.forwardingPolicy"
	// A matching read can be followed by a stale one.
	stableReads = 3
)

// readbackInterval is a variable so unit tests can shorten it.
var readbackInterval = 2 * time.Second

var errWaitTimedOut = errors.New("timed out waiting for forwarding policy link change to propagate")

// profileLocks serialises writes to one profile: each is a read-modify-write,
// and Terraform runs sibling links in parallel.
var profileLocks sync.Map

func lockProfile(profileType, profileID string) (unlock func()) {
	key := profileType + "/" + strings.ToLower(profileID)
	mu, _ := profileLocks.LoadOrStore(key, &sync.Mutex{})
	m, _ := mu.(*sync.Mutex)
	m.Lock()
	return m.Unlock
}

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
	delta any,
) error {
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.PATCH,
		profileTemplate(profileType),
		map[string]string{"forwardingProfileId": profileID},
		map[string]any{"policies@delta": []any{delta}},
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
	return r.patchPolicies(ctx, profileType, profileID, removedLink{
		Removed:   removedReason{Reason: "deleted"},
		ODataType: forwardingPolicyLinkODataType,
		ID:        linkID,
	})
}

// removedLink is a struct because Graph needs @removed before @odata.type
// ("Invalid request parameters" otherwise) and json.Marshal sorts map keys.
type removedLink struct {
	Removed   removedReason `json:"@removed"`
	ODataType string        `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData key.
	ID        string        `json:"id"`
}

type removedReason struct {
	Reason string `json:"reason"`
}

// findLink falls back to the policy ID because a re-added link gets a new ID.
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
			switch kiotaerrors.GraphError(ctx, err).StatusCode {
			case 404:
				return true, nil
			case 401, 403:
				return false, err
			}
			// The removal succeeded; keep polling through transient errors.
			return false, nil
		}
		return findLink(profile.Policies, linkID, "") == nil, nil
	})
}
