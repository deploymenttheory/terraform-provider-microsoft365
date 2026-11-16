package graphBetaNetworkForwardingProfile

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

var (
	errInvalidProfileResponse = stderrors.New("graph returned a forwarding profile without an ID")
	errProfileNotCreated      = stderrors.New("forwarding profile not found after failed create")
)

// createdLookupAttemptsFor polls longest after errors that can hide a create
// (transport errors, which GraphError reports as 400 "URLError", 503 or 504;
// gateway failures; non-Graph 500s). 400/409 get a few checks (a retried POST
// fails as a duplicate); anything else gets one.
func createdLookupAttemptsFor(info kiotaerrors.GraphErrorInfo) int {
	switch {
	case info.ErrorCode == "URLError",
		info.StatusCode == 502, info.StatusCode == 503, info.StatusCode == 504,
		info.StatusCode == 500 && !info.IsODataError:
		return createdLookupAttempts
	case info.StatusCode == 400, info.StatusCode == 409:
		return duplicateLookupAttempts
	default:
		return 1
	}
}

// claimedProfiles holds IDs already created or adopted in this process, so two
// resources with the same name and priority cannot both adopt one profile.
var claimedProfiles sync.Map

func claimProfile(id string) { claimedProfiles.Store(strings.ToLower(id), struct{}{}) }

func profileClaimed(id string) bool {
	_, ok := claimedProfiles.Load(strings.ToLower(id))
	return ok
}

// Variables so unit tests can shorten the wait.
var (
	createdLookupAttempts   = 12
	duplicateLookupAttempts = 3
	createdLookupInterval   = 5 * time.Second
)

// The SDK lacks the v2 properties, so requests are raw JSON in the portal's shape.
func requestBody(
	ctx context.Context,
	plan NetworkForwardingProfileResourceModel,
	state *NetworkForwardingProfileResourceModel,
) (map[string]any, diag.Diagnostics) {
	body := map[string]any{}
	if state == nil {
		body["isCustomProfile"] = true
		body["trafficForwardingType"] = plan.TrafficForwardingType.ValueString()
	}
	if state == nil || !plan.Name.Equal(state.Name) {
		body["name"] = plan.Name.ValueString()
	}
	if !plan.Description.IsUnknown() &&
		(state == nil || !plan.Description.Equal(state.Description)) {
		if plan.Description.IsNull() {
			if state != nil {
				body["description"] = nil
			}
		} else {
			body["description"] = plan.Description.ValueString()
		}
	}
	if state == nil || !plan.Priority.Equal(state.Priority) {
		body["priority"] = plan.Priority.ValueInt32()
	}
	var prior *sharedmodels.ForwardingProfileSettings
	if state != nil {
		s := state.settings()
		prior = &s
	}
	diags := sharedmodels.SettingsPatch(ctx, body, plan.settings(), prior)
	return body, diags
}

func (r *NetworkForwardingProfileResource) create(
	ctx context.Context,
	body map[string]any,
) (*sharedmodels.ForwardingProfileWire, error) {
	var response sharedmodels.ForwardingProfileWire
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.POST,
		r.ResourcePath,
		nil,
		body,
		&response,
	); err != nil {
		return nil, fmt.Errorf("create forwarding profile: %w", err)
	}
	if response.ID == nil || *response.ID == "" {
		return nil, errInvalidProfileResponse
	}
	return &response, nil
}

func (r *NetworkForwardingProfileResource) list(
	ctx context.Context,
) ([]sharedmodels.ForwardingProfileWire, error) {
	var response struct {
		Value []sharedmodels.ForwardingProfileWire `json:"value"`
	}
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		r.ResourcePath,
		nil,
		nil,
		&response,
	); err != nil {
		return nil, fmt.Errorf("list forwarding profiles: %w", err)
	}
	return response.Value, nil
}

// customProfileIDs returns the IDs of the custom profiles that exist now.
func (r *NetworkForwardingProfileResource) customProfileIDs(ctx context.Context) (map[string]bool, error) {
	profiles, err := r.list(ctx)
	if err != nil {
		return nil, err
	}
	ids := make(map[string]bool, len(profiles))
	for _, p := range profiles {
		if p.ID != nil && p.IsCustomProfile != nil && *p.IsCustomProfile {
			ids[*p.ID] = true
		}
	}
	return ids, nil
}

// findCreated returns the custom profile matching the plan's name, traffic type
// and priority that is not in existing, or errProfileNotCreated when none is.
func (r *NetworkForwardingProfileResource) findCreated(
	ctx context.Context,
	plan NetworkForwardingProfileResourceModel,
	existing map[string]bool,
) (*sharedmodels.ForwardingProfileWire, error) {
	profiles, err := r.list(ctx)
	if err != nil {
		return nil, err
	}
	for i, p := range profiles {
		if p.ID != nil && !existing[*p.ID] && !profileClaimed(*p.ID) && p.IsCustomProfile != nil && *p.IsCustomProfile &&
			p.Name != nil && *p.Name == plan.Name.ValueString() &&
			p.TrafficForwardingType != nil && *p.TrafficForwardingType == plan.TrafficForwardingType.ValueString() &&
			p.Priority != nil && *p.Priority == plan.Priority.ValueInt32() {
			return &profiles[i], nil
		}
	}
	return nil, errProfileNotCreated
}

// waitForCreated looks for a profile created by a failed request; it can appear late.
func (r *NetworkForwardingProfileResource) waitForCreated(
	ctx context.Context,
	plan NetworkForwardingProfileResourceModel,
	existing map[string]bool,
	attempts int,
) *sharedmodels.ForwardingProfileWire {
	for attempt := 1; attempt <= attempts; attempt++ {
		created, err := r.findCreated(ctx, plan, existing)
		if err == nil && created != nil {
			return created
		}
		// Stop on lasting errors such as 403; transport errors (400 "URLError") can clear.
		if err != nil && !stderrors.Is(err, errProfileNotCreated) {
			info := kiotaerrors.GraphError(ctx, err)
			if info.StatusCode >= 400 && info.StatusCode < 500 && info.StatusCode != 429 &&
				info.ErrorCode != "URLError" {
				return nil
			}
		}
		if attempt == attempts {
			return nil
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(createdLookupInterval):
		}
	}
	return nil
}

func (r *NetworkForwardingProfileResource) get(
	ctx context.Context,
	id string,
) (*sharedmodels.ForwardingProfileWire, error) {
	var response sharedmodels.ForwardingProfileWire
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		r.ResourcePath+"/{profileId}",
		map[string]string{"profileId": id},
		nil,
		&response,
	); err != nil {
		return nil, fmt.Errorf("get forwarding profile: %w", err)
	}
	return &response, nil
}

func (r *NetworkForwardingProfileResource) patch(
	ctx context.Context,
	id string,
	body map[string]any,
) error {
	if len(body) == 0 {
		return nil
	}
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.PATCH,
		r.ResourcePath+"/{profileId}",
		map[string]string{"profileId": id},
		body,
		nil,
	); err != nil {
		return fmt.Errorf("update forwarding profile: %w", err)
	}
	return nil
}

func (r *NetworkForwardingProfileResource) delete(ctx context.Context, id string) error {
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.DELETE,
		r.ResourcePath+"/{profileId}",
		map[string]string{"profileId": id},
		nil,
		nil,
	); err != nil {
		return fmt.Errorf("delete forwarding profile: %w", err)
	}
	return nil
}
