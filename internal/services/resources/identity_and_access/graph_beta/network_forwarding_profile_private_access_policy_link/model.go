package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	ForwardingProfileID types.String   `tfsdk:"forwarding_profile_id"`
	ProfileType         types.String   `tfsdk:"profile_type"`
	ForwardingPolicyID  types.String   `tfsdk:"forwarding_policy_id"`
	ConfigType          types.String   `tfsdk:"config_type"`
	Priority            types.Int64    `tfsdk:"priority"`
	State               types.String   `tfsdk:"state"`
	PolicyLinkID        types.String   `tfsdk:"policy_link_id"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

type policyRef struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type policyLink struct {
	ID         string     `json:"id"`
	State      string     `json:"state,omitempty"`
	Priority   *int64     `json:"priority,omitempty"`
	ConfigType string     `json:"configType,omitempty"`
	Policy     *policyRef `json:"policy,omitempty"`
}

type forwardingProfile struct {
	ID       string       `json:"id"`
	Policies []policyLink `json:"policies"`
}
