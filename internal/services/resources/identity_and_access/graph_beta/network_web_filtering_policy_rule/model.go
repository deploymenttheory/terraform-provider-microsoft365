package graphBetaNetworkWebFilteringPolicyRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkWebFilteringPolicyRuleResourceModel represents a web content
// filtering rule inside a Global Secure Access web filtering policy.
type NetworkWebFilteringPolicyRuleResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	WebFilteringPolicyID types.String   `tfsdk:"web_filtering_policy_id"`
	Name                 types.String   `tfsdk:"name"`
	Description          types.String   `tfsdk:"description"`
	Priority             types.Int64    `tfsdk:"priority"`
	Action               types.String   `tfsdk:"action"`
	Status               types.String   `tfsdk:"status"`
	UrlsOrFqdns          types.Set      `tfsdk:"urls_or_fqdns"`
	WebCategories        types.Set      `tfsdk:"web_categories"`
	HTTPMethods          types.Set      `tfsdk:"http_methods"`
	SessionTypes         types.Set      `tfsdk:"session_types"`
	CustomHeaders        types.List     `tfsdk:"custom_headers"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type customHeaderModel struct {
	HeaderName  types.String `tfsdk:"header_name"`
	HeaderValue types.String `tfsdk:"header_value"`
}
