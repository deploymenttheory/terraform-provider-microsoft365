package graphBetaNetworkContentPolicyRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// contentPolicyRulesEndpoint is the portal-backed Graph collection of a content policy's rules; %s is the policy ID.
const contentPolicyRulesEndpoint = "networkaccess/filePolicies/%s/policyRules"

// contentPolicyRuleEndpoint addresses one content policy rule; the first %s is the policy ID and the second the rule ID.
const contentPolicyRuleEndpoint = "networkaccess/filePolicies/%s/policyRules/%s"

// contentPolicyRuleDestinationObjectType is the element type of the destinations list.
var contentPolicyRuleDestinationObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type":   types.StringType,
	"values": types.SetType{ElemType: types.StringType},
}}

type NetworkContentPolicyRuleResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	ContentPolicyID  types.String   `tfsdk:"content_policy_id"`
	Name             types.String   `tfsdk:"name"`
	Description      types.String   `tfsdk:"description"`
	Action           types.String   `tfsdk:"action"`
	Priority         types.Int64    `tfsdk:"priority"`
	Status           types.String   `tfsdk:"status"`
	Activities       types.Set      `tfsdk:"activities"`
	ContentTypes     types.Set      `tfsdk:"content_types"`
	TextContentTypes types.Set      `tfsdk:"text_content_types"`
	Destinations     types.List     `tfsdk:"destinations"`
	SessionTypes     types.Set      `tfsdk:"session_types"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

type ContentPolicyRuleDestinationModel struct {
	Type   types.String `tfsdk:"type"`
	Values types.Set    `tfsdk:"values"`
}

// contentPolicyRuleResponse is the JSON shape of a file rule returned by Graph. The SDK has no model for it (see the SDK
// deviation in crud.go).
type contentPolicyRuleResponse struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Action      *string `json:"action"`
	Priority    *int64  `json:"priority"`
	Settings    *struct {
		Status *string `json:"status"`
	} `json:"settings"`
	MatchingConditions *struct {
		FileAttributes *struct {
			Activities       *string   `json:"activities"`
			ContentTypes     []*string `json:"contentTypes"`
			TextContentTypes []*string `json:"textContentTypes"`
		} `json:"fileAttributes"`
		Destinations []*contentPolicyRuleDestinationResponse `json:"destinations"`
		Sources      *struct {
			SessionType *string `json:"sessionType"`
		} `json:"sources"`
	} `json:"matchingConditions"`
}

// contentPolicyRuleDestinationResponse is one file policy destination; @odata.type selects web category, FQDN or URL.
type contentPolicyRuleDestinationResponse struct {
	ODataType *string   `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
	Values    []*string `json:"values"`
}
