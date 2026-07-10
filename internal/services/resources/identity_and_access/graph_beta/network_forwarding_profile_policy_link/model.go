package graphBetaNetworkForwardingProfilePolicyLink

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkForwardingProfilePolicyLinkResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	ForwardingProfileID   types.String   `tfsdk:"forwarding_profile_id"`
	PolicyLinkID          types.String   `tfsdk:"policy_link_id"`
	State                 types.String   `tfsdk:"state"`
	Priority              types.Int64    `tfsdk:"priority"`
	Version               types.String   `tfsdk:"version"`
	PolicyID              types.String   `tfsdk:"policy_id"`
	PolicyName            types.String   `tfsdk:"policy_name"`
	PolicyDescription     types.String   `tfsdk:"policy_description"`
	TrafficForwardingType types.String   `tfsdk:"traffic_forwarding_type"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
