// REF: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-policy?view=graph-rest-beta
package graphBetaNetworkFilteringProfilePolicyLink

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkFilteringProfilePolicyLinkResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	FilteringProfileID types.String   `tfsdk:"filtering_profile_id"`
	PolicyID           types.String   `tfsdk:"policy_id"`
	PolicyLinkID       types.String   `tfsdk:"policy_link_id"`
	PolicyType         types.String   `tfsdk:"policy_type"`
	State              types.String   `tfsdk:"state"`
	Priority           types.Int64    `tfsdk:"priority"`
	LoggingState       types.String   `tfsdk:"logging_state"`
	Version            types.String   `tfsdk:"version"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}
