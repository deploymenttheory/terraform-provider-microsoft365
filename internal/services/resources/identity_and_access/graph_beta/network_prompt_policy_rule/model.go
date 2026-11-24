package graphBetaNetworkPromptPolicyRule

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// promptPolicyRulesEndpoint is the portal-backed Graph collection of the rules of a prompt policy; %s is the policy ID.
const promptPolicyRulesEndpoint = "networkaccess/promptPolicies/%s/policyRules"

// promptPolicyRuleEndpoint addresses one prompt policy rule; the first %s is the policy ID and the second the rule ID.
const promptPolicyRuleEndpoint = "networkaccess/promptPolicies/%s/policyRules/%s"

// conversationSchemeObjectType is the element type of conversation_schemes.
var conversationSchemeObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type":        types.StringType,
	"url":         types.StringType,
	"json_path":   types.StringType,
	"scheme_name": types.StringType,
}}

type NetworkPromptPolicyRuleResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	PromptPolicyID      types.String   `tfsdk:"prompt_policy_id"`
	Name                types.String   `tfsdk:"name"`
	Description         types.String   `tfsdk:"description"`
	Action              types.String   `tfsdk:"action"`
	Priority            types.Int32    `tfsdk:"priority"`
	Enabled             types.Bool     `tfsdk:"enabled"`
	Status              types.String   `tfsdk:"status"`
	ConversationSchemes types.List     `tfsdk:"conversation_schemes"`
	PromptLogging       types.String   `tfsdk:"prompt_logging"`
	ScanResult          types.String   `tfsdk:"scan_result"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

type ConversationSchemeModel struct {
	Type       types.String `tfsdk:"type"`
	URL        types.String `tfsdk:"url"`
	JSONPath   types.String `tfsdk:"json_path"`
	SchemeName types.String `tfsdk:"scheme_name"`
}

type PromptPolicyRuleIdentity struct {
	ID             string `tfsdk:"id"`
	PromptPolicyID string `tfsdk:"prompt_policy_id"`
}

// promptPolicyRuleResponse is the JSON shape of a prompt policy rule returned by Graph. The SDK has no model for it (see
// the SDK deviation in crud.go); decoding it directly also keeps portal fields omitted from Graph metadata, notably action.
type promptPolicyRuleResponse struct {
	ID            *string `json:"id"`
	ODataType     *string `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Action        *string `json:"action"`
	PromptLogging *string `json:"promptLogging"`
	Priority      *int64  `json:"priority"`
	Settings      *struct {
		Status *string `json:"status"`
	} `json:"settings"`
	Conditions *struct {
		ScanResult *string `json:"scanResult"`
		// ConversationSchemes stays raw so a present null can be told apart from a missing property.
		ConversationSchemes json.RawMessage `json:"conversationSchemes"`
	} `json:"matchingConditions"`
}

// conversationSchemeResponse is one entry of matchingConditions.conversationSchemes.
type conversationSchemeResponse struct {
	ODataType  *string `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
	URL        *string `json:"url"`
	JSONPath   *string `json:"jsonPath"`
	SchemeName *string `json:"schemeName"`
}
