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

// CustomHeaderModel is one header that Global Secure Access adds to matching requests.
type CustomHeaderModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// webFilteringRuleRequest is the raw JSON web filtering rule sent to /networkaccess/webFilteringPolicies/{id}/policyRules.
type webFilteringRuleRequest struct {
	ODataType   string  `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	// Action is a webFilteringRuleAction for allow rules, which always carry headerSettings, or a
	// webFilteringBlockAction for block rules, which Graph rejects with any headerSettings.
	Action             any                            `json:"action"`
	Priority           *int64                         `json:"priority,omitempty"`
	Settings           webFilteringRuleSettings       `json:"settings"`
	MatchingConditions webFilteringMatchingConditions `json:"matchingConditions"`
}

// webFilteringRule is the raw JSON web filtering rule returned by Graph.
type webFilteringRule struct {
	ID                 *string                        `json:"id"`
	Name               *string                        `json:"name"`
	Description        *string                        `json:"description"`
	Priority           *int64                         `json:"priority"`
	Action             *webFilteringRuleAction        `json:"action"`
	Settings           webFilteringRuleSettings       `json:"settings"`
	MatchingConditions webFilteringMatchingConditions `json:"matchingConditions"`
}

// webFilteringRuleAction is a rule action. Allow actions send headerSettings as null when there are no custom headers.
type webFilteringRuleAction struct {
	ODataType      string                      `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
	HeaderSettings *webFilteringHeaderSettings `json:"headerSettings"`
}

// webFilteringBlockAction is a block action, sent without headerSettings.
type webFilteringBlockAction struct {
	ODataType string `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
}

// webFilteringHeaderSettings holds the header modifications of an allow action.
type webFilteringHeaderSettings struct {
	Modifications []webFilteringHeaderModification `json:"modifications"`
}

// webFilteringHeaderModification is one header modification; Terraform writes only headerModificationAdd entries.
type webFilteringHeaderModification struct {
	ODataType   string `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
	HeaderName  string `json:"headerName"`
	HeaderValue string `json:"headerValue"`
}

// webFilteringRuleSettings holds the rule status.
type webFilteringRuleSettings struct {
	Status *string `json:"status,omitempty"`
}

// webFilteringMatchingConditions holds the rule's destinations and optional sources.
type webFilteringMatchingConditions struct {
	Destinations *webFilteringDestinations `json:"destinations,omitempty"`
	Sources      *webFilteringSources      `json:"sources,omitempty"`
}

// webFilteringDestinations holds the URL/FQDN and web category targets and the comma-separated HTTP methods.
type webFilteringDestinations struct {
	Targets           []webFilteringDestinationTarget `json:"targets"`
	HTTPRequestMethod *string                         `json:"httpRequestMethod,omitempty"`
}

// webFilteringDestinationTarget is one URL/FQDN or web category destination.
type webFilteringDestinationTarget struct {
	ODataType string   `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
	Values    []string `json:"values"`
}

// webFilteringSources holds the comma-separated session types.
type webFilteringSources struct {
	SessionType *string `json:"sessionType,omitempty"`
}
