package graphBetaNetworkForwardingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

func mapRemoteStateToTerraform(
	ctx context.Context,
	state *NetworkForwardingProfileResourceModel,
	remote *sharedmodels.ForwardingProfileWire,
) diag.Diagnostics {
	var diags diag.Diagnostics
	if remote.IsCustomProfile != nil && !*remote.IsCustomProfile {
		diags.AddError(
			"Not a custom forwarding profile",
			"This profile is a built-in default profile. Manage it with microsoft365_graph_beta_identity_and_access_network_default_forwarding_profile.",
		)
		return diags
	}
	state.ID = types.StringPointerValue(remote.ID)
	state.Name = types.StringPointerValue(remote.Name)
	state.Description = types.StringPointerValue(remote.Description)
	state.TrafficForwardingType = types.StringPointerValue(remote.TrafficForwardingType)
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
