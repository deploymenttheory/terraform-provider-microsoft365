// REF: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-forwardingpolicylink?view=graph-rest-beta
package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel represents the schema for the network forwarding profile private access policy link resource.
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

// policyRef represents the policy referenced by a forwarding policy link.
type policyRef struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// policyLink represents a forwarding profile policy link returned by Microsoft Graph.
type policyLink struct {
	ID         string     `json:"id"`
	State      string     `json:"state,omitempty"`
	Priority   *int64     `json:"priority,omitempty"`
	ConfigType string     `json:"configType,omitempty"`
	Policy     *policyRef `json:"policy,omitempty"`
}

// forwardingProfile represents a forwarding profile with its expanded policy links.
type forwardingProfile struct {
	ID       string       `json:"id"`
	Policies []policyLink `json:"policies"`
}

// removedLink is a policies@delta entry that removes a link. It is a struct because Graph needs @removed before
// @odata.type ("Invalid request parameters" otherwise) and json.Marshal sorts map keys.
type removedLink struct {
	Removed   removedReason `json:"@removed"`
	ODataType string        `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData key.
	ID        string        `json:"id"`
}

// removedReason represents the removal reason in a policy link delta request.
type removedReason struct {
	Reason string `json:"reason"`
}
