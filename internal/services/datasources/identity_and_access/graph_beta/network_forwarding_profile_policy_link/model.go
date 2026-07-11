package graphBetaNetworkForwardingProfilePolicyLink

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkForwardingProfilePolicyLinkDataSourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	ForwardingProfileID   types.String   `tfsdk:"forwarding_profile_id"`
	ForwardingProfileName types.String   `tfsdk:"forwarding_profile_name"`
	TrafficForwardingType types.String   `tfsdk:"traffic_forwarding_type"`
	PolicyName            types.String   `tfsdk:"policy_name"`
	PolicyLinkID          types.String   `tfsdk:"policy_link_id"`
	Priority              types.Int64    `tfsdk:"priority"`
	State                 types.String   `tfsdk:"state"`
	Version               types.String   `tfsdk:"version"`
	PolicyID              types.String   `tfsdk:"policy_id"`
	PolicyDescription     types.String   `tfsdk:"policy_description"`
	PolicyVersion         types.String   `tfsdk:"policy_version"`
	PrivateAccessAppID    types.String   `tfsdk:"private_access_app_id"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
