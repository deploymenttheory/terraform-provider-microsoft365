package graphBetaNetworkDeviceForwardingProfile

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// NetworkDeviceForwardingProfileResourceModel describes a pre-logon (device) forwarding profile.
type NetworkDeviceForwardingProfileResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	Name                  types.String   `tfsdk:"name"`
	Description           types.String   `tfsdk:"description"`
	TrafficForwardingType types.String   `tfsdk:"traffic_forwarding_type"`
	Priority              types.Int32    `tfsdk:"priority"`
	State                 types.String   `tfsdk:"state"`
	ClientFallbackAction  types.String   `tfsdk:"client_fallback_action"`
	LocalNetworkSettings  types.Object   `tfsdk:"local_network_settings"`
	AssignmentRules       types.Object   `tfsdk:"assignment_rules"`
	Version               types.String   `tfsdk:"version"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

func (m NetworkDeviceForwardingProfileResourceModel) settings() sharedmodels.ForwardingProfileSettings {
	return sharedmodels.ForwardingProfileSettings{
		State:                m.State,
		ClientFallbackAction: m.ClientFallbackAction,
		// Device profiles reject agent-only source matching, so it is never sent.
		ForwardingSourceMatching: types.ObjectNull(sharedmodels.ForwardingSourceMatchingAttrTypes),
		LocalNetworkSettings:     m.LocalNetworkSettings,
		AssignmentRules:          m.AssignmentRules,
	}
}
