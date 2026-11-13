package graphBetaNetworkDefaultForwardingProfile

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

var errDefaultProfileNotFound = errors.New(
	"default forwarding profile not found; the tenant must be onboarded to Global Secure Access and use traffic forwarding profiles v2",
)

// findDefault returns the built-in profile for a traffic type. Its ID differs
// per tenant, so it is looked up instead of configured.
func (r *NetworkDefaultForwardingProfileResource) findDefault(
	ctx context.Context,
	trafficType string,
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
		if p.TrafficForwardingType != nil && *p.TrafficForwardingType == trafficType &&
			p.IsCustomProfile != nil && !*p.IsCustomProfile {
			return &response.Value[i], nil
		}
	}
	return nil, fmt.Errorf("%w: %q", errDefaultProfileNotFound, trafficType)
}

func (r *NetworkDefaultForwardingProfileResource) get(
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

func (r *NetworkDefaultForwardingProfileResource) patch(
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

func mapRemoteStateToTerraform(
	ctx context.Context,
	state *NetworkDefaultForwardingProfileResourceModel,
	remote *sharedmodels.ForwardingProfileWire,
) diag.Diagnostics {
	var diags diag.Diagnostics
	if remote.IsCustomProfile == nil || *remote.IsCustomProfile {
		diags.AddError(
			"Not a default forwarding profile",
			"This profile is a custom profile. Manage it with microsoft365_graph_beta_identity_and_access_network_forwarding_profile.",
		)
		return diags
	}
	if !state.TrafficForwardingType.IsNull() && !state.TrafficForwardingType.IsUnknown() &&
		remote.TrafficForwardingType != nil && state.TrafficForwardingType.ValueString() != *remote.TrafficForwardingType {
		diags.AddError(
			"Traffic type mismatch",
			fmt.Sprintf(
				"Profile %s is the default %q profile, not %q.",
				*remote.ID,
				*remote.TrafficForwardingType,
				state.TrafficForwardingType.ValueString(),
			),
		)
		return diags
	}
	state.ID = types.StringPointerValue(remote.ID)
	state.TrafficForwardingType = types.StringPointerValue(remote.TrafficForwardingType)
	state.Name = types.StringPointerValue(remote.Name)
	state.Description = types.StringPointerValue(remote.Description)
	state.Priority = types.Int32PointerValue(remote.Priority)
	state.Version = types.StringPointerValue(remote.Version)
	settings, d := sharedmodels.SettingsFromWire(ctx, remote)
	diags.Append(d...)
	state.State = settings.State
	state.ClientFallbackAction = settings.ClientFallbackAction
	state.ForwardingSourceMatching = settings.ForwardingSourceMatching
	state.LocalNetworkSettings = settings.LocalNetworkSettings
	state.AssignmentRules = settings.AssignmentRules
	return diags
}
