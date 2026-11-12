package graphBetaNetworkDefaultForwardingProfile

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// NetworkDefaultForwardingProfileResourceModel describes a built-in forwarding profile.
type NetworkDefaultForwardingProfileResourceModel struct {
	ID                       types.String   `tfsdk:"id"`
	TrafficForwardingType    types.String   `tfsdk:"traffic_forwarding_type"`
	Name                     types.String   `tfsdk:"name"`
	Description              types.String   `tfsdk:"description"`
	Priority                 types.Int32    `tfsdk:"priority"`
	State                    types.String   `tfsdk:"state"`
	ClientFallbackAction     types.String   `tfsdk:"client_fallback_action"`
	ForwardingSourceMatching types.Object   `tfsdk:"forwarding_source_matching"`
	LocalNetworkSettings     types.Object   `tfsdk:"local_network_settings"`
	AssignmentRules          types.Object   `tfsdk:"assignment_rules"`
	Version                  types.String   `tfsdk:"version"`
	Timeouts                 timeouts.Value `tfsdk:"timeouts"`
}

func (m NetworkDefaultForwardingProfileResourceModel) settings() sharedmodels.ForwardingProfileSettings {
	return sharedmodels.ForwardingProfileSettings{
		State:                    m.State,
		ClientFallbackAction:     m.ClientFallbackAction,
		ForwardingSourceMatching: m.ForwardingSourceMatching,
		LocalNetworkSettings:     m.LocalNetworkSettings,
		AssignmentRules:          m.AssignmentRules,
	}
}
