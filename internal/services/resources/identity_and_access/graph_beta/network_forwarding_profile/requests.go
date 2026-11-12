package graphBetaNetworkForwardingProfile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

var (
	errInvalidProfileResponse = errors.New("graph returned a forwarding profile without an ID")
	errProfileNotCreated      = errors.New("forwarding profile not found after failed create")
)

// Variables so unit tests can shorten the wait.
var (
	createdLookupAttempts = 12
	createdLookupInterval = 5 * time.Second
)

// The v2 profile properties are missing from the generated SDK models, so the
// resource sends JSON matching the Entra portal requests.
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

// findCreated returns the custom profile matching the plan's name, traffic type
// and priority, or errProfileNotCreated when none exists.
func (r *NetworkForwardingProfileResource) findCreated(
	ctx context.Context,
	plan NetworkForwardingProfileResourceModel,
) (*sharedmodels.ForwardingProfileWire, error) {
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
	for i, p := range response.Value {
		if p.ID != nil && p.IsCustomProfile != nil && *p.IsCustomProfile &&
			p.Name != nil && *p.Name == plan.Name.ValueString() &&
			p.TrafficForwardingType != nil && *p.TrafficForwardingType == plan.TrafficForwardingType.ValueString() &&
			p.Priority != nil && *p.Priority == plan.Priority.ValueInt32() {
			return &response.Value[i], nil
		}
	}
	return nil, errProfileNotCreated
}

// waitForCreated polls for a profile created by a request that failed on the
// gateway; it can take a while to appear in the list.
func (r *NetworkForwardingProfileResource) waitForCreated(
	ctx context.Context,
	plan NetworkForwardingProfileResourceModel,
) *sharedmodels.ForwardingProfileWire {
	for range createdLookupAttempts {
		if created, err := r.findCreated(ctx, plan); err == nil && created != nil {
			return created
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
