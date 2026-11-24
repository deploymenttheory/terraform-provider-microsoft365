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

// policyLinkRequest is the raw JSON policy link body sent to /networkAccess/filteringProfiles/{id}/policies.
type policyLinkRequest struct {
	ODataType    string           `json:"@odata.type,omitempty"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
	Priority     *int64           `json:"priority,omitempty"`
	State        string           `json:"state,omitempty"`
	LoggingState string           `json:"loggingState,omitempty"`
	Policy       *policyReference `json:"policy,omitempty"`
}

// policyReference references the existing policy a new link points to.
type policyReference struct {
	ID        string `json:"id,omitempty"`
	ODataType string `json:"@odata.type,omitempty"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
}
