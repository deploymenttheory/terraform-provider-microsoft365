package graphBetaNetworkDefaultForwardingProfile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// mapRemoteStateToTerraform maps the remote Microsoft Graph response to Terraform state.
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

// decodeFindDefaultResponse validates and maps the Microsoft Graph response.
func decodeFindDefaultResponse(response struct {
	Value []sharedmodels.ForwardingProfileWire `json:"value"`
}, err error, trafficType string,
) (*sharedmodels.ForwardingProfileWire, error) {
	if err != nil {
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
